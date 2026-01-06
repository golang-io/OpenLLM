package OpenLLM

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-io/requests"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// ============================================================================
// OpenAI Responses API 客户端封装 / OpenAI Responses API Client Wrapper
// ============================================================================

// OpenAIResponses 封装 OpenAI Responses API 客户端
// OpenAIResponses wraps OpenAI Responses API client
//
// Responses API 是 OpenAI 新推出的对话 API，提供了比 Chat Completions 更多的高级特性：
// - 原生支持 Reasoning 模型（o1, o3 等）
// - 可以获取模型的思考过程（reasoning content）
// - 支持更精细的工具调用控制
// - 支持后台异步处理
//
// 注意：Responses API 的输入输出格式与 Chat Completions 不同，需要特殊的参数转换
type OpenAIResponses struct {
	options []Option
	client  *openai.Client
}

// CreateOpenAIResponses 创建 OpenAI Responses API 客户端
// CreateOpenAIResponses creates a new OpenAI Responses API client
func CreateOpenAIResponses(opts ...Option) *OpenAIResponses {
	options := newOptions(opts)
	client := openai.NewClient(
		option.WithBaseURL(options.URL),
		option.WithAPIKey(options.APIKey),
		option.WithHTTPClient(
			requests.New(requests.Timeout(60*time.Second)).HTTPClient(options.HTTPClientOptions...),
		),
	)

	return &OpenAIResponses{
		options: opts,
		client:  &client,
	}
}

// ============================================================================
// 统一接口实现 / Unified Interface Implementation
// ============================================================================

// Completion 执行单次对话完成（非流式）
// Completion performs a single conversation completion (non-streaming)
//
// 注意：当前实现为简化版本，主要用于展示 Responses API 的接口封装方式
// 实际使用时需要根据具体需求完善参数转换逻辑
func (o *OpenAIResponses) Completion(ctx context.Context, input *Input, opts ...Option) (*Output, error) {
	return nil, fmt.Errorf("Responses API 实现正在开发中，请使用 Chat Completions API（CreateOpenAI）")
}

// CompletionStream 执行单次对话完成（流式）
// CompletionStream performs a single conversation completion (streaming)
//
// 注意：当前实现为简化版本
func (o *OpenAIResponses) CompletionStream(ctx context.Context, input *Input, streamOutput StreamOutput, opts ...Option) (*Output, error) {
	return nil, fmt.Errorf("Responses API 流式实现正在开发中，请使用 Chat Completions API（CreateOpenAI）")
}

// Provider 获取提供商信息
// Provider returns the provider information
func (o *OpenAIResponses) Provider() ProviderInfo {
	return ProviderInfo{
		Type:    ProviderOpenAI,
		Name:    "OpenAI Responses API",
		Version: "v1",
		Capabilities: ProviderCapabilities{
			ToolCall:      true,
			Streaming:     true,
			Temperature:   true,
			TopP:          true,
			Seed:          false, // Responses API 可能不支持 seed
			SystemMessage: true,
		},
	}
}

// ============================================================================
// 底层 SDK 访问方法 / Native SDK Access Methods
// ============================================================================

// GetResponsesService 获取底层 Responses Service
// GetResponsesService returns the underlying Responses Service
//
// 如果需要使用 Responses API 的原生功能，可以通过这个方法获取
// SDK 的 ResponseService，然后直接调用其方法
//
// 示例：
//
//	service := client.GetResponsesService()
//	response, err := service.New(ctx, params)
func (o *OpenAIResponses) GetResponsesService() *responses.ResponseService {
	return &o.client.Responses
}

// ============================================================================
// 使用说明 / Usage Documentation
// ============================================================================

/*
Responses API 使用示例：

基础用法（当完整实现后）：
```go
	client := OpenLLM.CreateOpenAIResponses(
		OpenLLM.WithAPIKey("your-api-key"),
		OpenLLM.WithModel("gpt-4o"),
	)

	input := &OpenLLM.Input{
		Messages: []OpenLLM.Message{
			OpenLLM.UserMessage("你好"),
		},
	}

	output, err := client.Completion(context.Background(), input)
```

使用底层 SDK（当前推荐方式）：
```go
	client := OpenLLM.CreateOpenAIResponses(
		OpenLLM.WithAPIKey("your-api-key"),
	)

	service := client.GetResponsesService()

	// 构建 Responses API 原生参数
	params := responses.ResponseNewParams{
		Model: "gpt-4o",
		// ... 其他参数
	}

	response, err := service.New(context.Background(), params)
```

Responses API 的主要优势：
1. 原生支持 o1/o3 等 Reasoning 模型
2. 可以获取模型的思考过程（reasoning content）
3. 支持后台异步处理（background: true）
4. 更精细的工具调用控制

更多文档请参考：README_RESPONSES.md
*/
