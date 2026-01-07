package OpenLLM

import (
	"context"
	"time"
)

// ============================================================================
// LLM接口定义 / LLM Interface Definition
// ============================================================================

// LLM 大语言模型统一接口 - 适配所有提供商
// LLM is a unified interface for all language model providers
type LLM interface {
	// Complete 执行单次对话完成（非流式）
	// Complete performs a single conversation completion (non-streaming)
	Completion(context.Context, *Input, ...Option) (*Output, error)

	// CompleteStream 执行单次对话完成（流式）
	// CompleteStream performs a single conversation completion (streaming)
	// callback函数在收到每个chunk时被调用，用于实时输出
	// The callback function is called for each chunk received for real-time output
	CompletionStream(context.Context, *Input, StreamOutput, ...Option) (*Output, error)
}

// ============================================================================
// 流式输出回调 / Streaming Output Function
// ============================================================================

// StreamOutput 流式输出回调函数
// StreamOutput is called for each chunk of streaming output
// 参数:
//   - content: 本次输出的文本片段 / The text chunk for this output
type StreamOutput func(content string)

// ============================================================================
// 提供商信息 / Provider Information
// ============================================================================

// ProviderInfo 提供商信息
// ProviderInfo contains information about the LLM provider
type ProviderInfo struct {
	Type    ProviderType `json:"type"`     // 提供商类型 / Provider type
	Name    string       `json:"name"`     // 提供商名称 / Provider name
	Version string       `json:"version"`  // API版本 / API version
	Model   string       `json:"model"`    // 当前使用的模型 / Current model
	BaseURL string       `json:"base_url"` // API基础URL / API base URL
	// 能力标识 / Capability flags
	Capabilities ProviderCapabilities `json:"capabilities"`
}

// ProviderCapabilities 提供商能力
// ProviderCapabilities represents the capabilities of a provider
type ProviderCapabilities struct {
	ToolCall      bool `json:"tool_call"`      // 是否支持工具调用 / Support tool calls
	Thinking      bool `json:"thinking"`       // 是否支持思考 / Support thinking
	Streaming     bool `json:"streaming"`      // 是否支持流式输出 / Support streaming
	Temperature   bool `json:"temperature"`    // 是否支持温度参数 / Support temperature
	TopP          bool `json:"top_p"`          // 是否支持TopP参数 / Support top-p
	Seed          bool `json:"seed"`           // 是否支持随机种子 / Support seed
	SystemMessage bool `json:"system_message"` // 是否支持系统消息 / Support system message
}

// ============================================================================
// 配置结构 / Configuration Structure
// ============================================================================

// Config LLM配置 - 用于创建LLM实例
// Config is used to create an LLM instance
type Config struct {
	Provider ProviderType `json:"provider"` // 提供商类型 / Provider type
	BaseURL  string       `json:"base_url"` // API基础URL / API base URL
	APIKey   string       `json:"api_key"`  // API密钥 / API key
	Model    string       `json:"model"`    // 模型名称 / Model name

}

// ProviderType 提供商类型
// ProviderType represents the type of LLM provider
type ProviderType string

const (
	ProviderOpenAI  ProviderType = "openai"  // OpenAI
	ProviderAzure   ProviderType = "azure"   // Azure OpenAI
	ProviderClaude  ProviderType = "claude"  // Anthropic Claude
	ProviderHunyuan ProviderType = "hunyuan" // 腾讯混元 / Tencent Hunyuan
	ProviderGemini  ProviderType = "gemini"  // Google Gemini
	ProviderCustom  ProviderType = "custom"  // 自定义提供商 / Custom provider
)

// ============================================================================
// 错误定义 / Error Definitions
// ============================================================================

// LLMError LLM调用错误
// LLMError represents an error from LLM call
type LLMError struct {
	Provider ProviderType `json:"provider"` // 提供商 / Provider
	Code     string       `json:"code"`     // 错误码 / Error code
	Message  string       `json:"message"`  // 错误消息 / Error message
	Cause    error        `json:"cause"`    // 原始错误 / Original error
}

// Error 实现error接口
// Error implements the error interface
func (e *LLMError) Error() string {
	if e.Cause != nil {
		return string(e.Provider) + ": " + e.Message + " - " + e.Cause.Error()
	}
	return string(e.Provider) + ": " + e.Message
}

// Unwrap 实现errors.Unwrap接口
// Unwrap implements the errors.Unwrap interface
func (e *LLMError) Unwrap() error {
	return e.Cause
}

// ============================================================================
// 构造错误的辅助函数 / Helper Functions for Error Construction
// ============================================================================

// NewLLMError 创建LLM错误
// NewLLMError creates a new LLM error
func NewLLMError(provider ProviderType, code, message string, cause error) *LLMError {
	return &LLMError{
		Provider: provider,
		Code:     code,
		Message:  message,
		Cause:    cause,
	}
}

// ============================================================================
// 统一消息结构 / Union Message Structure
// ============================================================================

// MessageRole 消息角色
// MessageRole represents the role of a message sender
type MessageRole string

const (
	RoleSystem    MessageRole = "system"    // 系统消息 / System message
	RoleUser      MessageRole = "user"      // 用户消息 / User message
	RoleAssistant MessageRole = "assistant" // 助手消息 / Assistant message
	RoleTool      MessageRole = "tool"      // 工具结果 / Tool result
)

// Message 统一消息格式 - 适配所有LLM提供商
// Message represents a unified message format across all LLM providers
type Message struct {
	Role       MessageRole    `json:"role"`                   // 消息角色 / Message role
	Content    string         `json:"content"`                // 文本内容 / Text content
	ToolCalls  []ToolCall     `json:"tool_calls,omitempty"`   // 工具调用列表（仅assistant） / Tool calls (assistant only)
	ToolCallID string         `json:"tool_call_id,omitempty"` // 工具调用ID（仅tool） / Tool call ID (tool only)
	Name       string         `json:"name,omitempty"`         // 消息名称（可选） / Message name (optional)
	Metadata   map[string]any `json:"metadata,omitempty"`     // 元数据（扩展用） / Metadata (for extension)
}

// ============================================================================
// 统一工具调用结构 / Union Tool Call Structure
// ============================================================================

// ToolCall 统一工具调用格式
// ToolCall represents a unified tool call format
type ToolCall struct {
	ID        string         `json:"id"`        // 工具调用ID / Tool call ID
	Name      string         `json:"name"`      // 工具名称 / Tool name
	Arguments map[string]any `json:"arguments"` // 工具参数 / Tool arguments
}

// Tool 工具定义
// Tool represents a function tool definition for LLM function calling
type Tool struct {
	Name        string       `json:"name"`               // 工具名称，必须唯一 / Tool name, must be unique
	Description string       `json:"description"`        // 工具描述，用于LLM理解工具用途 / Tool description for LLM understanding
	Parameters  *JSONSchema  `json:"parameters"`         // JSON Schema格式的参数定义 / Parameters in JSON Schema format
	Metadata    ToolMetadata `json:"metadata,omitempty"` // 工具元数据，用于扩展 / Tool metadata for extension
}

// ToolMetadata 工具元数据
// ToolMetadata contains additional metadata for a tool
type ToolMetadata struct {
	Category   string         `json:"category,omitempty"`   // 工具分类 / Tool category
	Version    string         `json:"version,omitempty"`    // 工具版本 / Tool version
	Author     string         `json:"author,omitempty"`     // 工具作者 / Tool author
	Tags       []string       `json:"tags,omitempty"`       // 标签 / Tags
	Extensions map[string]any `json:"extensions,omitempty"` // 扩展字段 / Extension fields
}

// JSONSchema JSON Schema定义
// JSONSchema represents a JSON Schema definition for tool parameters
type JSONSchema struct {
	Type        string                 `json:"type"`                  // 类型：object, string, number, integer, boolean, array / Schema type
	Properties  map[string]*JSONSchema `json:"properties,omitempty"`  // 属性定义（仅object类型）/ Property definitions (object type only)
	Required    []string               `json:"required,omitempty"`    // 必填字段列表（仅object类型）/ Required field list (object type only)
	Items       *JSONSchema            `json:"items,omitempty"`       // 数组元素定义（仅array类型）/ Array item definition (array type only)
	Enum        []any                  `json:"enum,omitempty"`        // 枚举值 / Enum values
	Default     any                    `json:"default,omitempty"`     // 默认值 / Default value
	Description string                 `json:"description,omitempty"` // 字段描述 / Field description
	Minimum     *float64               `json:"minimum,omitempty"`     // 最小值（number/integer）/ Minimum value (number/integer)
	Maximum     *float64               `json:"maximum,omitempty"`     // 最大值（number/integer）/ Maximum value (number/integer)
}

// ============================================================================
// 工具选择策略 / Tool Choice Strategy
// ============================================================================

// ToolChoiceType 工具选择类型
// ToolChoiceType represents the type of tool choice
type ToolChoiceType string

const (
	ToolChoiceAuto     ToolChoiceType = "auto"     // 自动选择 / Auto select
	ToolChoiceRequired ToolChoiceType = "required" // 必须调用工具 / Must call tool
	ToolChoiceNone     ToolChoiceType = "none"     // 不调用工具 / Don't call tool
	ToolChoiceSpecific ToolChoiceType = "specific" // 指定特定工具 / Specific tool
)

// ToolChoiceOption 工具选择选项
// ToolChoiceOption represents tool choice configuration
type ToolChoiceOption struct {
	Type     ToolChoiceType `json:"type"`                // 选择类型 / Choice type
	ToolName string         `json:"tool_name,omitempty"` // 指定工具名（当Type=Specific时） / Specific tool name (when Type=Specific)
}

// Auto 返回自动选择的工具选项
// Auto returns auto tool choice option
func Auto() ToolChoiceOption {
	return ToolChoiceOption{Type: ToolChoiceAuto}
}

// Required 返回必须调用工具的选项
// Required returns required tool choice option
func Required() ToolChoiceOption {
	return ToolChoiceOption{Type: ToolChoiceRequired}
}

// None 返回不调用工具的选项
// None returns none tool choice option
func None() ToolChoiceOption {
	return ToolChoiceOption{Type: ToolChoiceNone}
}

// Specific 返回指定工具的选项
// Specific returns specific tool choice option
func Specific(toolName string) ToolChoiceOption {
	return ToolChoiceOption{Type: ToolChoiceSpecific, ToolName: toolName}
}

// ============================================================================
// 统一请求结构 / Union Request Structure
// ============================================================================

// Input 统一LLM请求格式 - 适配所有LLM提供商
// Input represents a unified LLM request format across all providers
type Input struct {
	Model      string            `json:"model"`                 // 模型名称 / Model name
	Messages   []Message         `json:"messages"`              // 消息列表 / Message list
	Tools      []Tool            `json:"tools,omitempty"`       // 工具列表 / Tool list
	ToolChoice *ToolChoiceOption `json:"tool_choice,omitempty"` // 工具选择策略 / Tool choice strategy
	Stream     bool              `json:"stream"`                // 是否流式 / Whether streaming

}

// ============================================================================
// 统一响应结构 / Union Response Structure
// ============================================================================

// FinishReason 结束原因
// FinishReason represents the reason why the generation stopped
type FinishReason string

const (
	FinishReasonStop          FinishReason = "STOP"           // 正常结束 / Normal stop
	FinishReasonLength        FinishReason = "length"         // 达到长度限制 / Length limit reached
	FinishReasonToolCalls     FinishReason = "tool_calls"     // 需要调用工具 / Tool calls required
	FinishReasonContentFilter FinishReason = "content_filter" // 内容过滤 / Content filtered
	FinishReasonError         FinishReason = "error"          // 错误 / Error
)

// FromFinishReason 将结束原因转换为Union格式
// FromFinishReason converts finish reason to Union format
func FromFinishReason(reason string) FinishReason {
	switch reason {
	case "stop", "STOP":
		return FinishReasonStop
	case "length":
		return FinishReasonLength
	case "tool_calls":
		return FinishReasonToolCalls
	case "content_filter":
		return FinishReasonContentFilter
	default:
		return FinishReasonStop
	}
}

// TokenUsage Token使用情况
// TokenUsage represents token consumption statistics
type TokenUsage struct {
	InputTokens    int64 `json:"input_tokens"`    // 输入token数 / Input tokens
	ThinkingTokens int64 `json:"thinking_tokens"` // 思考token数（仅支持思考模型） / Thinking tokens (only for thinking models)
	OutputTokens   int64 `json:"output_tokens"`   // 输出token数 / Output tokens
	TotalTokens    int64 `json:"total_tokens"`    // 总token数 / Total tokens
}

// Output 统一LLM响应格式 - 适配所有LLM提供商
// Output represents a unified LLM response format across all providers
type Output struct {
	StartAt      time.Time     `json:"start_at"`      // 开始时间 / Start time
	Content      string        `json:"content"`       // 文本内容 / Text content
	Thinking     string        `json:"thinking"`      // 思考内容（仅支持思考模型，如Gemini） / Thinking content (only for thinking models like Gemini)
	ToolCalls    []ToolCall    `json:"tool_calls"`    // 工具调用列表 / Tool call list
	FinishReason string        `json:"finish_reason"` // 结束原因 / Finish reason
	TokenUsage   TokenUsage    `json:"token_usage"`   // Token使用情况 / Token usage
	Cost         time.Duration `json:"cost"`          // 调用耗时 / Call duration
	RawResponse  any           `json:"raw_response"`  // 原始响应（调试用） / Raw response (for debugging)
	// 扩展字段（提供商特定数据）/ Extended fields (provider-specific data)
	Extra map[string]any `json:"extra,omitempty"`
}

// ============================================================================
// 消息构造器 / Message Builders
// ============================================================================

// SystemMessage 创建系统消息
// SystemMessage creates a system message
func SystemMessage(content string) Message {
	return Message{
		Role:    RoleSystem,
		Content: content,
	}
}

// UserMessage 创建用户消息
// UserMessage creates a user message
func UserMessage(content string) Message {
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}

// AssistantMessage 创建助手消息（纯文本）
// AssistantMessage creates an assistant message (text only)
func AssistantMessage(content string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

// AssistantMessageWithTools 创建助手消息（包含工具调用）
// AssistantMessageWithTools creates an assistant message with tool calls
func AssistantMessageWithTools(content string, toolCalls []ToolCall) Message {
	return Message{
		Role:      RoleAssistant,
		Content:   content,
		ToolCalls: toolCalls,
	}
}

// ToolMessage 创建工具结果消息
// ToolMessage creates a tool result message
func ToolMessage(content string, toolCallID string) Message {
	return Message{
		Role:       RoleTool,
		Content:    content,
		ToolCallID: toolCallID,
	}
}
