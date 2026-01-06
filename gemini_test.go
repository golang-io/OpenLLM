package OpenLLM

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-io/requests"
)

func Test_CreateGemini(t *testing.T) {
	gemini := CreateGemini(context.Background())
	resp, err := gemini.Completion(context.Background(), &Input{
		Messages: []Message{UserMessage("please call the tool get_weather to get the weather of a city, the city is Beijing")},
	}, HTTPClientOptions(requests.Trace(1024000)))

	t.Logf("%#v, %v", resp, err)
}

func Test_Gemini_CompletionStream(t *testing.T) {
	gemini := CreateGemini(context.Background(), HTTPClientOptions(requests.Trace(1024000)))
	output, err := gemini.CompletionStream(context.Background(), &Input{
		Messages: []Message{UserMessage("hello! how are you? please answer in Chinese. 你的输出多一点，我要测试stream模式")},
	}, func(content string) {
		fmt.Printf("%s", content)
		time.Sleep(1 * time.Second)
	})

	t.Logf("%#v, err=%v", output, err)
}

func Test_Gemini_Completion_FromOpenAI(t *testing.T) {
	gemini := CreateOpenAI(
		URL(os.Getenv("OpenLLM_GOOGLE_BASE_URL")),
		APIKey(os.Getenv("OpenLLM_GOOGLE_API_KEY")),
		HTTPClientOptions(requests.Trace(1024000)),
	)
	output, err := gemini.CompletionStream(context.Background(), &Input{
		Model:    Gemini25Flash,
		Messages: []Message{UserMessage("查询北京现在的天气. 你的输出多一点，我要测试stream模式")},
		Tools:    Tools,
	}, func(content string) {
		fmt.Printf("%s", content)
	})

	t.Logf("%#v, err=%v", output, err)
}
