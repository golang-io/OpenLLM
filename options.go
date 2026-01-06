package OpenLLM

import (
	"os"

	"github.com/golang-io/requests"
)

// ============================================================================
// 配置选项 / Configuration Options
// ============================================================================

// Options LLM 客户端配置选项
// Options defines configuration options for LLM clients
type Options struct {
	Provider          string            `json:"provider,omitempty"`            // 提供商类型 / Provider type
	URL               string            `json:"url,omitempty"`                 // API基础URL / API base URL
	APIKey            string            `json:"api_key,omitempty"`             // API密钥 / API key
	Temperature       float64           `json:"temperature,omitempty"`         // 默认温度 / Default temperature
	MaxTokens         int64             `json:"max_tokens,omitempty"`          // 默认最大token数 / Default max tokens
	TopP              float64           `json:"top_p,omitempty"`               // 默认TopP / Default top-p
	JSONSet           map[string]any    `json:"json_set,omitempty"`            // 扩展配置（提供商特定）/ Extended config (provider-specific)
	Seed              int64             `json:"seed,omitempty"`                // 随机种子 / Random seed
	HTTPClientOptions []requests.Option `json:"http_client_options,omitempty"` // HTTP客户端配置 / HTTP client options
}

// Option 配置函数类型
// Option is a function type for configuring Options
type Option func(*Options)

// newOptions 创建并初始化配置选项
// newOptions creates and initializes configuration options
func newOptions(opts []Option, extends ...Option) *Options {
	options := &Options{
		URL:         os.Getenv("OpenLLM_BASE_URL"),
		APIKey:      os.Getenv("OpenLLM_API_KEY"),
		Temperature: 0.2,
		MaxTokens:   128 * 1000,
		Seed:        88,
		JSONSet:     make(map[string]any),
	}
	for _, o := range opts {
		o(options)
	}
	for _, o := range extends {
		o(options)
	}
	return options
}

// ============================================================================
// 配置函数 / Configuration Functions
// ============================================================================

// URL 设置 API 基础 URL
// URL sets the API base URL
func URL(url string) Option {
	return func(options *Options) {
		options.URL = url
	}
}

// APIKey 设置 API 密钥
// APIKey sets the API key
func APIKey(apiKey string) Option {
	return func(options *Options) {
		options.APIKey = apiKey
	}
}

// Seed 设置随机种子（用于结果可复现）
// Seed sets the random seed (for reproducible results)
func Seed(seed int64) Option {
	return func(options *Options) {
		options.Seed = seed
	}
}

// Temperature 设置温度参数（控制随机性，范围 0.0-2.0）
// Temperature sets the temperature parameter (controls randomness, range 0.0-2.0)
func Temperature(temperature float64) Option {
	return func(options *Options) {
		options.Temperature = temperature
	}
}

// MaxTokens 设置最大输出 token 数
// MaxTokens sets the maximum number of output tokens
func MaxTokens(maxTokens int64) Option {
	return func(options *Options) {
		options.MaxTokens = maxTokens
	}
}

// TopP 设置 Top-P 采样参数（核采样，范围 0.0-1.0）
// TopP sets the Top-P sampling parameter (nucleus sampling, range 0.0-1.0)
func TopP(topP float64) Option {
	return func(options *Options) {
		options.TopP = topP
	}
}

// JSONSet 设置扩展配置（提供商特定参数）
// JSONSet sets extended configuration (provider-specific parameters)
func JSONSet(jsonSet map[string]any) Option {
	return func(options *Options) {
		options.JSONSet = jsonSet
	}
}

// HTTPClientOptions 设置 HTTP 客户端配置（如超时、代理等）
// HTTPClientOptions sets HTTP client configuration (such as timeout, proxy, etc.)
func HTTPClientOptions(httpClientOptions ...requests.Option) Option {
	return func(options *Options) {
		options.HTTPClientOptions = append(options.HTTPClientOptions, httpClientOptions...)
	}
}
