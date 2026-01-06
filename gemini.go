package OpenLLM

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/golang-io/requests"
	"google.golang.org/genai"
)

type Gemini struct {
	options []Option
	client  *genai.Client
	once    sync.Once
}

func CreateGemini(ctx context.Context, opts ...Option) *Gemini {
	options := newOptions(opts)
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     options.APIKey,
		Backend:    genai.BackendGeminiAPI,
		HTTPClient: requests.New().HTTPClient(options.HTTPClientOptions...),
	})
	if err != nil {
		panic(err)
	}
	return &Gemini{
		options: opts,
		client:  client,
	}
}

// extractThinkingContent 从响应中提取 thinking 内容
// extractThinkingContent extracts thinking content from the response
func extractThinkingContent(candidates []*genai.Candidate) string {
	if len(candidates) == 0 || candidates[0].Content == nil {
		return ""
	}

	var thinkingParts []string
	for _, part := range candidates[0].Content.Parts {
		if part.Thought && part.Text != "" {
			thinkingParts = append(thinkingParts, part.Text)
		}
	}

	if len(thinkingParts) == 0 {
		return ""
	}
	return thinkingParts[0] // 通常只有一个 thinking part，如果有多个则合并
}

func (g *Gemini) Completion(ctx context.Context, input *Input, opts ...Option) (*Output, error) {
	// options := newOptions(opts, g.options...)
	output := &Output{StartAt: time.Now()}
	result, err := g.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(input.Messages[0].Content),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
				ThinkingLevel:   genai.ThinkingLevelLow,
			},
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						{
							Name:        "get_weather",
							Description: "Get the weather of a city",
							Parameters: &genai.Schema{
								Type: "object",
								Properties: map[string]*genai.Schema{
									"city": {
										Type: "string",
									},
								},
								Required: []string{"city"},
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	// 提取 thinking 内容 / Extract thinking content
	thinking := extractThinkingContent(result.Candidates)

	if len(result.Candidates) > 0 {
		for _, part := range result.Candidates[0].Content.Parts {
			if part.FunctionCall != nil {
				output.ToolCalls = append(output.ToolCalls, ToolCall{
					ID:        part.FunctionCall.ID,
					Name:      part.FunctionCall.Name,
					Arguments: part.FunctionCall.Args,
				})
			}
		}
	}

	output.Content = result.Text()
	output.Thinking = thinking
	output.FinishReason = string(result.Candidates[0].FinishReason)
	output.TokenUsage = TokenUsage{
		InputTokens:    int64(result.UsageMetadata.PromptTokenCount),
		ThinkingTokens: int64(result.UsageMetadata.ThoughtsTokenCount),
		OutputTokens:   int64(result.UsageMetadata.CandidatesTokenCount + result.UsageMetadata.ThoughtsTokenCount),
		TotalTokens:    int64(result.UsageMetadata.TotalTokenCount),
	}
	return output, nil
}

func Int32(value int) *int32 {
	v := int32(value)
	return &v
}

func (g *Gemini) CompletionStream(ctx context.Context, input *Input, streamOutput StreamOutput, opts ...Option) (*Output, error) {
	// options := newOptions(opts)

	response := g.client.Models.GenerateContentStream(ctx, "gemini-2.5-flash",
		genai.Text(input.Messages[0].Content),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true, // 启用 thinking 内容 / Enable thinking content
			},
		},
	)

	output := &Output{StartAt: time.Now()}

	for chunk, err := range response {
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(chunk.Candidates) > 0 && chunk.Candidates[0].Content != nil {
			// 遍历所有 parts，区分 thinking 内容和普通内容
			// Iterate through all parts, distinguish between thinking content and normal content
			for _, part := range chunk.Candidates[0].Content.Parts {
				if part.Text == "" {
					continue
				}
				if part.Thought {
					// thinking 内容累积到 Thinking 字段 / Accumulate thinking content to Thinking field
					output.Thinking += part.Text
					streamOutput(part.Text)
				} else {
					// 普通内容累积到 Content 字段并流式输出 / Accumulate normal content to Content field and stream output
					output.Content += part.Text
					streamOutput(part.Text)
				}
			}
		}
		if len(chunk.Candidates) > 0 {
			output.FinishReason = string(chunk.Candidates[0].FinishReason)
		}
		output.TokenUsage.InputTokens += int64(chunk.UsageMetadata.PromptTokenCount)
		output.TokenUsage.ThinkingTokens += int64(chunk.UsageMetadata.ThoughtsTokenCount)
		output.TokenUsage.OutputTokens += int64(chunk.UsageMetadata.CandidatesTokenCount + chunk.UsageMetadata.ThoughtsTokenCount)
		output.TokenUsage.TotalTokens += int64(chunk.UsageMetadata.TotalTokenCount)
	}
	output.Cost = time.Since(output.StartAt)
	return output, nil
}
