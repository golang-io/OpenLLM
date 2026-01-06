package OpenLLM

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/golang-io/requests"
)

var input = &Input{
	Model:    "deepseek-v3.1-terminus",
	Messages: []Message{UserMessage("查询北京明天上午10点的天气?")},
	Tools:    Tools,
}

func Test_CreateOpenAI(t *testing.T) {
	api := CreateOpenAI(
		URL(os.Getenv("OpenLLM_LKEAP_BASE_URL")),
		APIKey(os.Getenv("OpenLLM_LKEAP_API_KEY")),
		HTTPClientOptions(requests.Trace(1024000)),
	)
	output, err := api.CompletionStream(context.Background(), input, func(content string) {
		log.Printf("%s", content)
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", output)
	log.Print("-------------")
	output2, err := api.CompletionStream(context.Background(), input, func(content string) {
		log.Printf("%s", content)
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", output2)

}

func Test_CreateAzure(t *testing.T) {
	api := CreateAzure(
		URL(os.Getenv("OpenLLM_AZURE_URL")),
		APIKey(os.Getenv("OpenLLM_AZURE_API_KEY")),
		HTTPClientOptions(requests.Trace(1024000)),
	)
	input.Model = "gpt-5-mini"
	output, err := api.Completion(context.Background(), input)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", output)
	log.Print("-------------")
	output2, err := api.CompletionStream(context.Background(), input, func(content string) {
		log.Printf("%s", content)
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", output2)
}

func TestVenusCompletion(t *testing.T) {
	venus := CreateOpenAI(
		HTTPClientOptions(requests.Trace(1024000)),
	)
	output, err := venus.CompletionStream(context.Background(), &Input{
		Model:    Qwen3VL235BA22BThinking,
		Messages: []Message{UserMessage("Hello, how are you?")},
	}, func(content string) {
		log.Printf("%s", content)
	})
	log.Printf("%#v, err=%v", output, err)
}
