package OpenLLM

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-io/requests"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

var _ LLM = (*OpenAI)(nil)

// ============================================================================
// OpenAI SDK客户端封装 / OpenAI SDK Client Wrapper
// ============================================================================

// openAIClient 封装OpenAI SDK客户端，直接使用SDK原生类型
// openAIClient wraps OpenAI SDK client, using SDK native types directly
type OpenAI struct {
	options []Option
	client  *openai.Client
}

// newOpenAIClient 创建OpenAI SDK客户端
// newOpenAIClient creates a new OpenAI SDK client
func CreateOpenAI(opts ...Option) *OpenAI {

	options := newOptions(opts)
	client := openai.NewClient(
		option.WithBaseURL(options.URL),
		option.WithAPIKey(options.APIKey),
		option.WithHTTPClient(requests.New().HTTPClient(options.HTTPClientOptions...)),
	)

	return &OpenAI{
		options: opts,
		client:  &client,
	}
}

// ============================================================================
// 底层SDK调用方法 / Native SDK Call Methods
// ============================================================================

// chatCompletion 调用OpenAI Chat Completion API（非流式）
// chatCompletion calls OpenAI Chat Completion API (non-streaming)
// 使用SDK原生类型
func (o *OpenAI) ChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	return o.client.Chat.Completions.New(ctx, params)
}

// chatCompletionStream 调用OpenAI Chat Completion API（流式）
// chatCompletionStream calls OpenAI Chat Completion API (streaming)
// 返回Stream接口
func (o *OpenAI) ChatCompletionStream(ctx context.Context, params openai.ChatCompletionNewParams, streamOutput StreamOutput) (*openai.ChatCompletion, error) {

	// 创建流式请求
	stream := o.client.Chat.Completions.NewStreaming(ctx, params, option.WithJSONSet("stream", true))
	// 累积流式数据
	var content strings.Builder

	completion := &openai.ChatCompletion{
		ID:      "stream-" + requests.GenId(),
		Object:  "chat.completion_stream",
		Model:   params.Model,
		Created: time.Now().Unix(),
		Choices: []openai.ChatCompletionChoice{},
		Usage:   openai.CompletionUsage{},
	}

	// 工具调用的临时存储（因为stream中tool_call是增量的）
	tools := make(map[int64]*openai.ChatCompletionMessageToolCallUnion)

	// 用于跟踪是否已初始化 choice
	choiceInitialized := false

	// 处理流式响应
	for stream.Next() {
		chunk := stream.Current()
		// 处理每个choice的delta
		for _, choice := range chunk.Choices {
			// 初始化 choice（只在第一次遇到时）
			if !choiceInitialized {
				completion.Choices = append(completion.Choices, openai.ChatCompletionChoice{
					Index: choice.Index,
					Message: openai.ChatCompletionMessage{
						Role:    "assistant",
						Content: "",
					},
				})
				choiceInitialized = true
			}

			delta := choice.Delta
			// 累积content并实时输出
			if delta.Content != "" {
				content.WriteString(delta.Content)
				if streamOutput != nil {
					streamOutput(delta.Content)
				}
			}

			// 累积tool_calls（增量式）
			for _, delta := range delta.ToolCalls {
				tool, ok := tools[delta.Index]
				if !ok {
					tool = &openai.ChatCompletionMessageToolCallUnion{
						ID:   delta.ID,
						Type: "function",
						Function: openai.ChatCompletionMessageFunctionToolCallFunction{
							Name:      delta.Function.Name,
							Arguments: delta.Function.Arguments,
						},
					}
				} else {
					tool.Function.Arguments += delta.Function.Arguments
				}
				tools[delta.Index] = tool
			}

			// 保存finish_reason（最后一个chunk才有）
			if choice.FinishReason != "" {
				completion.Choices[0].FinishReason = choice.FinishReason
			}
		}

		// 保存usage（最后一个chunk才有）
		if chunk.Usage.CompletionTokens > 0 || chunk.Usage.PromptTokens > 0 {
			completion.Usage = chunk.Usage
		}
	}

	// 检查stream是否有错误
	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("agent:stream error: %w", err)
	}

	// 设置最终的内容和工具调用
	if len(completion.Choices) > 0 {
		completion.Choices[0].Message.Content = content.String()

		// 将map转换为slice（按index排序）
		if len(tools) > 0 {
			completion.Choices[0].Message.ToolCalls = make([]openai.ChatCompletionMessageToolCallUnion, 0, len(tools))
			for _, tool := range tools {
				if tool != nil {
					completion.Choices[0].Message.ToolCalls = append(completion.Choices[0].Message.ToolCalls, *tool)
				}
			}
		}
	}

	return completion, nil

}

// Completion 执行单次对话完成（非流式）
// Completion performs a single conversation completion (non-streaming)
func (o *OpenAI) Completion(ctx context.Context, input *Input, opts ...Option) (*Output, error) {
	// 1. 适配：Union类型 → SDK原生类型
	params, err := o.GenerateOpenAIChatCompletionNewParams(input, opts...)
	if err != nil {
		return nil, NewLLMError(ProviderOpenAI, "CONVERT_ERROR", "转换请求参数失败", err)
	}

	// 2. 记录开始时间
	startTime := time.Now()

	// 3. 调用底层SDK（使用原生类型）
	completion, err := o.ChatCompletion(ctx, params)
	if err != nil {
		return nil, NewLLMError(ProviderOpenAI, "API_ERROR", "OpenAI API调用失败", err)
	}

	// 5. 检查响应
	if len(completion.Choices) == 0 {
		return nil, NewLLMError(ProviderOpenAI, "EMPTY_RESPONSE", "OpenAI返回空响应", nil)
	}

	// 6. 适配：SDK原生类型 → Union类型
	return fromOpenAIResponse(completion, time.Since(startTime)), nil
}

// CompletionStream 执行单次对话完成（流式）
// CompletionStream performs a single conversation completion (streaming)
func (o *OpenAI) CompletionStream(ctx context.Context, input *Input, streamOutput StreamOutput, opts ...Option) (*Output, error) {
	// 1. 适配：Union类型 → SDK原生类型
	params, err := o.GenerateOpenAIChatCompletionNewParams(input, opts...)
	if err != nil {
		return nil, NewLLMError(ProviderOpenAI, "CONVERT_ERROR", "转换请求参数失败", err)
	}

	// 2. 记录开始时间
	startTime := time.Now()

	// 3. 调用底层SDK（使用原生类型）
	completion, err := o.ChatCompletionStream(ctx, params, streamOutput)
	if err != nil {
		return nil, NewLLMError(ProviderOpenAI, "API_ERROR", "OpenAI API调用失败", err)
	}
	return fromOpenAIResponse(completion, time.Since(startTime)), nil
}

// ============================================================================
// 适配逻辑 / Adapter Logic
// 将Union类型转换为SDK原生类型，或将SDK原生类型转换为Union类型
// ============================================================================

// toOpenAIParams 将Union请求转换为OpenAI SDK原生参数
// toOpenAIParams converts Union request to OpenAI SDK native parameters
func (o *OpenAI) GenerateOpenAIChatCompletionNewParams(input *Input, opts ...Option) (openai.ChatCompletionNewParams, error) {
	// 1. 转换消息
	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(input.Messages))
	for _, msg := range input.Messages {
		openaiMsg, err := toOpenAIMessage(msg)
		if err != nil {
			return openai.ChatCompletionNewParams{}, err
		}
		messages = append(messages, openaiMsg)
	}

	options := newOptions(o.options, opts...)

	// 2. 构建基础参数
	params := openai.ChatCompletionNewParams{
		Messages:            messages,
		Model:               input.Model,
		Temperature:         openai.Float(options.Temperature),
		MaxCompletionTokens: openai.Int(options.MaxTokens),
		TopP:                openai.Float(options.TopP),
	}

	// Gemini模型不支持Seed参数
	if !strings.HasPrefix(strings.ToLower(input.Model), "gemini") {
		params.Seed = openai.Int(options.Seed)
	}

	if input.ToolChoice != nil {
		params.ToolChoice = toOpenAIToolChoice(*input.ToolChoice)
	}

	for _, tool := range input.Tools {
		openaiTool, err := tool.ToOpenAITool()
		if err != nil {
			return openai.ChatCompletionNewParams{}, fmt.Errorf("转换工具 %s 失败: %w", tool.Name, err)
		}
		params.Tools = append(params.Tools, openaiTool)
	}

	return params, nil
}

// fromOpenAIResponse 将OpenAI SDK响应转换为Union响应
// fromOpenAIResponse converts OpenAI SDK response to Union response
func fromOpenAIResponse(completion *openai.ChatCompletion, duration time.Duration) *Output {
	choice := completion.Choices[0]
	message := choice.Message

	// 构建响应
	output := &Output{
		Content:      message.Content,
		FinishReason: choice.FinishReason,
		Cost:         duration,
		RawResponse:  completion,
	}

	// 转换工具调用
	for _, tc := range message.ToolCalls {
		var args map[string]any
		json.Unmarshal([]byte(tc.Function.Arguments), &args)
		output.ToolCalls = append(output.ToolCalls, ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: args,
		})
	}

	// 添加Token使用情况
	output.TokenUsage = TokenUsage{
		InputTokens:  int64(completion.Usage.PromptTokens),
		OutputTokens: int64(completion.Usage.CompletionTokens),
		TotalTokens:  int64(completion.Usage.PromptTokens + completion.Usage.CompletionTokens),
	}

	return output
}

// toOpenAIMessage 将Message消息转换为OpenAI消息
// toOpenAIMessage converts Message message to OpenAI message
func toOpenAIMessage(msg Message) (openai.ChatCompletionMessageParamUnion, error) {
	switch msg.Role {
	case RoleSystem:
		return openai.SystemMessage(msg.Content), nil

	case RoleUser:
		return openai.UserMessage(msg.Content), nil

	case RoleAssistant:
		// Assistant消息需要手动构建，因为可能包含工具调用
		var assistant openai.ChatCompletionAssistantMessageParam
		assistant.Content.OfString = param.NewOpt(msg.Content)

		// 转换工具调用
		if len(msg.ToolCalls) > 0 {
			toolCalls := make([]openai.ChatCompletionMessageToolCallUnionParam, 0, len(msg.ToolCalls))
			for _, tc := range msg.ToolCalls {
				argsJSON, _ := json.Marshal(tc.Arguments)
				var funcCall openai.ChatCompletionMessageFunctionToolCallParam
				funcCall.ID = tc.ID
				funcCall.Type = "function"
				funcCall.Function.Name = tc.Name
				funcCall.Function.Arguments = string(argsJSON)
				toolCalls = append(toolCalls, openai.ChatCompletionMessageToolCallUnionParam{
					OfFunction: &funcCall,
				})
			}
			assistant.ToolCalls = toolCalls
		}

		return openai.ChatCompletionMessageParamUnion{OfAssistant: &assistant}, nil

	case RoleTool:
		return openai.ToolMessage(msg.Content, msg.ToolCallID), nil

	default:
		return openai.ChatCompletionMessageParamUnion{}, fmt.Errorf("不支持的消息角色: %s", msg.Role)
	}
}

// toOpenAIToolChoice 将Tool工具选择转换为OpenAI格式
// toOpenAIToolChoice converts Union tool choice to OpenAI format
func toOpenAIToolChoice(choice ToolChoiceOption) openai.ChatCompletionToolChoiceOptionUnionParam {
	switch choice.Type {
	case ToolChoiceAuto:
		var result openai.ChatCompletionToolChoiceOptionUnionParam
		result.OfAuto = param.NewOpt("auto")
		return result
	case ToolChoiceRequired:
		var result openai.ChatCompletionToolChoiceOptionUnionParam
		result.OfAuto = param.NewOpt("required")
		return result
	case ToolChoiceNone:
		var result openai.ChatCompletionToolChoiceOptionUnionParam
		result.OfAuto = param.NewOpt("none")
		return result
	case ToolChoiceSpecific:
		var namedChoice openai.ChatCompletionNamedToolChoiceParam
		namedChoice.Type = "function"
		namedChoice.Function.Name = choice.ToolName

		var result openai.ChatCompletionToolChoiceOptionUnionParam
		result.OfFunctionToolChoice = &namedChoice
		return result
	default:
		var result openai.ChatCompletionToolChoiceOptionUnionParam
		result.OfAuto = param.NewOpt("auto")
		return result
	}
}
