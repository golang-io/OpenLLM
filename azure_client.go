package OpenLLM

import (
	"context"
	"log"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
)

var _ LLM = (*Azure)(nil)

// ============================================================================
// Azure OpenAI SDK客户端封装 / Azure OpenAI SDK Client Wrapper
// ============================================================================

// azureClient 封装Azure OpenAI SDK客户端
// Azure OpenAI兼容OpenAI协议，直接复用OpenAI的client
// azureClient wraps Azure OpenAI SDK client
// Azure OpenAI is compatible with OpenAI protocol, reuses OpenAI's client
type Azure struct {
	client  *OpenAI // 复用OpenAI client
	options []Option
}

// newAzureClient 创建Azure OpenAI SDK客户端
// newAzureClient creates a new Azure OpenAI SDK client
func CreateAzure(opts ...Option) *Azure {
	// Azure OpenAI兼容OpenAI协议，直接使用OpenAI client
	return &Azure{
		client:  CreateOpenAI(opts...),
		options: opts,
	}
}

// // chatCompletion 调用Azure OpenAI Chat Completion API（非流式）
// // chatCompletion calls Azure OpenAI Chat Completion API (non-streaming)
// func (c *Azure) chatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
// 	return c.client.chatCompletion(ctx, params)
// }

// // chatCompletionStream 调用Azure OpenAI Chat Completion API（流式）
// // chatCompletionStream calls Azure OpenAI Chat Completion API (streaming)
// func (c *Azure) chatCompletionStream(ctx context.Context, params openai.ChatCompletionNewParams) openAIStream {
// 	return c.client.chatCompletionStream(ctx, params)
// }

// ============================================================================
// LLM接口实现 / LLM Interface Implementation
// ============================================================================

// Completion 执行单次对话完成（非流式）
// Completion performs a single conversation completion (non-streaming)
func (a *Azure) Completion(ctx context.Context, input *Input, opts ...Option) (*Output, error) {

	// 2. 适配：Union类型 → SDK原生类型
	params, err := a.GenerateAzureChatCompletionNewParams(input, opts...)
	if err != nil {
		return nil, NewLLMError(ProviderAzure, "CONVERT_ERROR", "转换请求参数失败", err)
	}

	// 3. 记录开始时间
	startTime := time.Now()

	// 4. 调用底层SDK（使用原生类型）
	completion, err := a.client.ChatCompletion(ctx, params)
	if err != nil {
		return nil, NewLLMError(ProviderAzure, "API_ERROR", "Azure OpenAI API调用失败", err)
	}

	// 5. 计算耗时
	duration := time.Since(startTime)

	// 6. 检查响应
	if len(completion.Choices) == 0 {
		return nil, NewLLMError(ProviderAzure, "EMPTY_RESPONSE", "Azure OpenAI返回空响应", nil)
	}

	// 7. 适配：SDK原生类型 → Union类型
	return fromAzureResponse(completion, duration), nil
}

// CompletionStream 执行单次对话完成（流式）
// CompletionStream performs a single conversation completion (streaming)
func (a *Azure) CompletionStream(ctx context.Context, input *Input, streamOutput StreamOutput, opts ...Option) (*Output, error) {
	params, err := a.GenerateAzureChatCompletionNewParams(input, opts...)
	if err != nil {
		return nil, NewLLMError(ProviderAzure, "CONVERT_ERROR", "转换请求参数失败", err)
	}

	// 3. 记录开始时间
	startTime := time.Now()

	// 4. 调用底层SDK（使用原生类型）
	completion, err := a.client.ChatCompletionStream(ctx, params, streamOutput)
	if err != nil {
		return nil, NewLLMError(ProviderAzure, "API_ERROR", "Azure OpenAI API调用失败", err)
	}

	// 检查是否有响应内容 / Check if response has content
	if len(completion.Choices) == 0 {
		return nil, NewLLMError(ProviderAzure, "EMPTY_RESPONSE", "Azure OpenAI API返回空响应", nil)
	}

	return fromAzureResponse(completion, time.Since(startTime)), nil
}

// Provider 获取提供商信息
// Provider returns the provider information
func (p *Azure) Provider() ProviderInfo {
	return ProviderInfo{
		Type:    ProviderAzure,
		Name:    "Azure OpenAI",
		Version: "v1",
		// Model:   p.config.Model,
		// BaseURL: p.config.BaseURL,
		Capabilities: ProviderCapabilities{
			ToolCall:      true,
			Streaming:     true,
			Temperature:   true,
			TopP:          true,
			Seed:          true,
			SystemMessage: true,
		},
	}
}

// ============================================================================
// 适配逻辑 / Adapter Logic
// Azure OpenAI兼容OpenAI协议，复用OpenAI的适配逻辑，但需要过滤不支持的参数
// ============================================================================

// toAzureParams 将Union请求转换为Azure OpenAI参数
// Azure兼容OpenAI协议，但需要过滤不支持的参数
// toAzureParams converts Union request to Azure OpenAI parameters
// Azure is compatible with OpenAI protocol, but needs to filter unsupported parameters
func (a *Azure) GenerateAzureChatCompletionNewParams(input *Input, opts ...Option) (openai.ChatCompletionNewParams, error) {

	// 复用OpenAI的适配逻辑
	params, err := a.client.GenerateOpenAIChatCompletionNewParams(input, opts...)
	if err != nil {
		return openai.ChatCompletionNewParams{}, err
	}
	// 复制配置以避免修改原始配置

	if params.Temperature.Valid() {
		log.Printf("[Azure] 警告: Azure OpenAI不支持Temperature参数，将被忽略")
		params.Temperature = param.Opt[float64]{}
	}
	if params.TopP.Valid() {
		log.Printf("[Azure] 警告: Azure OpenAI不支持TopP参数，将被忽略")
		params.TopP = param.Opt[float64]{}
	}
	if params.Seed.Valid() {
		log.Printf("[Azure] 警告: Azure OpenAI不支持Seed参数，将被忽略")
		params.Seed = param.Opt[int64]{}
	}

	return params, nil
}

// fromAzureResponse 将Azure OpenAI响应转换为Union响应
// Azure兼容OpenAI协议，直接复用OpenAI的转换逻辑
// fromAzureResponse converts Azure OpenAI response to Union response
// Azure is compatible with OpenAI protocol, reuses OpenAI's conversion logic
func fromAzureResponse(completion *openai.ChatCompletion, duration time.Duration) *Output {
	return fromOpenAIResponse(completion, duration)
}
