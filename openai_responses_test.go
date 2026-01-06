package OpenLLM

import (
	"context"
	"testing"
)

// Test_OpenAIResponses_Completion 测试 Responses API 非流式调用
func Test_OpenAIResponses_Completion(t *testing.T) {
	t.Skip("需要配置真实的 API Key 才能运行")

	client := CreateOpenAIResponses(
		APIKey("your-api-key"),
	)

	input := &Input{
		Model: "gpt-4o",
		Messages: []Message{
			UserMessage("你好，请简单介绍一下自己"),
		},
	}

	output, err := client.Completion(context.Background(), input)
	if err != nil {
		t.Fatalf("Completion 失败: %v", err)
	}

	t.Logf("Content: %s", output.Content)
	t.Logf("Token Usage: %+v", output.TokenUsage)
	t.Logf("Cost: %v", output.Cost)
}

// Test_OpenAIResponses_CompletionStream 测试 Responses API 流式调用
func Test_OpenAIResponses_CompletionStream(t *testing.T) {
	t.Skip("需要配置真实的 API Key 才能运行")

	client := CreateOpenAIResponses(
		APIKey("your-api-key"),
	)

	input := &Input{
		Model: "gpt-4o",
		Messages: []Message{
			UserMessage("请用 100 字介绍 Go 语言"),
		},
	}

	output, err := client.CompletionStream(context.Background(), input, func(content string) {
		t.Logf("Stream: %s", content)
	})
	if err != nil {
		t.Fatalf("CompletionStream 失败: %v", err)
	}

	t.Logf("Final Content: %s", output.Content)
	t.Logf("Token Usage: %+v", output.TokenUsage)
}
