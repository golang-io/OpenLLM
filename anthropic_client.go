package OpenLLM

import (
	"context"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// 确保 Anthropic 实现了 LLM 接口
// Ensure Anthropic implements the LLM interface
var _ LLM = (*Anthropic)(nil)

// ============================================================================
// Anthropic Claude SDK客户端封装 / Anthropic Claude SDK Client Wrapper
// ============================================================================

// Anthropic 封装 Anthropic Claude SDK 客户端
// Anthropic wraps the Anthropic Claude SDK client
type Anthropic struct {
	client  anthropic.Client // Anthropic SDK 客户端 / Anthropic SDK client
	options []Option         // 配置选项 / Configuration options
}

// CreateAnthropic 创建 Anthropic Claude 客户端
// CreateAnthropic creates a new Anthropic Claude client
func CreateAnthropic(opts ...Option) *Anthropic {
	options := newOptions(opts)
	return &Anthropic{
		client: anthropic.NewClient(
			option.WithAPIKey(options.APIKey), // defaults to os.LookupEnv("ANTHROPIC_API_KEY")
		),
		options: opts,
	}
}

// ============================================================================
// LLM接口实现 / LLM Interface Implementation
// ============================================================================

// Completion 执行单次对话完成（非流式）
// Completion performs a single conversation completion (non-streaming)
func (a *Anthropic) Completion(ctx context.Context, input *Input, opts ...Option) (*Output, error) {
	options := newOptions(opts, a.options...)

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: int64(options.MaxTokens),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(input.Messages[0].Content)),
		},
		Model: anthropic.Model(input.Model),
	})
	if err != nil {
		return nil, err
	}
	return &Output{
		Content: message.Content[0].Text,
	}, nil
}

// CompletionStream 执行单次对话完成（流式）
// CompletionStream performs a single conversation completion (streaming)
func (a *Anthropic) CompletionStream(ctx context.Context, input *Input, streamOutput StreamOutput, opts ...Option) (*Output, error) {
	options := newOptions(opts, a.options...)

	stream := a.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		MaxTokens: int64(options.MaxTokens),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(input.Messages[0].Content)),
		},
		Model: anthropic.Model(input.Model),
	})

	message := anthropic.Message{}
	for stream.Next() {
		event := stream.Current()
		err := message.Accumulate(event)
		if err != nil {
			panic(err)
		}

		switch eventVariant := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				streamOutput(deltaVariant.Text)
			}

		}
	}

	if stream.Err() != nil {
		return nil, stream.Err()
	}
	return &Output{
		Content: message.Content[0].Text,
	}, nil
}

// Provider 获取提供商信息
// Provider returns the provider information
func (a *Anthropic) Provider() ProviderInfo {
	return ProviderInfo{
		Type: "anthropic",
	}
}
