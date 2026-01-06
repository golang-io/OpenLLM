# OpenLLM

> 统一的大语言模型（LLM）Go 语言接口库，支持 OpenAI、Gemini、DeepSeek、千问等多种模型

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

## 📖 目录

- [概述](#概述)
- [核心特性](#核心特性)
- [支持的模型](#支持的模型)
- [快速开始](#快速开始)
- [详细功能说明](#详细功能说明)
  - [基本对话](#1-基本对话)
  - [流式输出](#2-流式输出)
  - [工具调用完整指南](#3-工具调用完整指南)
  - [Thinking 模式](#4-thinking-模式推理过程)
- [支持的提供商详解](#支持的提供商详解)
  - [OpenAI](#openai)
  - [Gemini](#gemini)  
  - [中国大模型](#中国大模型)
  - [Azure OpenAI](#azure-openai)
  - [Anthropic Claude](#anthropic-claude)
- [API 文档](#api-文档)
- [高级主题](#高级主题)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)
- [贡献指南](#贡献指南)

---

## 概述

OpenLLM 提供了统一的 Go 语言接口来访问各种大语言模型，屏蔽不同提供商之间的 API 差异。

### 为什么选择 OpenLLM？

OpenLLM提供兼容OpenAI接口的统一接入点

- ✅ **Google Gemini** 提供 OpenAI 兼容端点
- ✅ **DeepSeek、通义千问、Kimi** 等中国模型完全兼容
- ✅ **Ollama、vLLM** 等推理平台原生支持
- ✅ **Mistral AI、智谱** 等国际厂商兼容

**OpenLLM 的价值**

- 🎯 **统一接口**：一套代码支持所有主流 LLM
- 🔄 **轻松切换**：更换模型只需修改配置
- 🛠️ **工具调用**：标准化的 Function Calling 支持
- 📡 **流式输出**：统一的流式响应处理
- 🧠 **Thinking 支持**：内置推理过程提取
- 🌍 **中国模型**：完美支持 DeepSeek、千问、Kimi 等
- 📝 **类型安全**：完整的 Go 类型定义
- 🚀 **高性能**：最小化适配开销

---

## 核心特性

### ✅ 已实现

| 功能 | 说明 | 支持的模型 |
|------|------|---------|
| **基本对话** | 标准的文本对话 | 所有模型 |
| **流式输出** | 实时流式响应 | 所有模型 |
| **工具调用** | Function Calling | OpenAI/Gemini/DeepSeek/千问 |
| **Thinking** | 推理过程提取 | Gemini/o1/o3/DeepSeek |
| **多模态** | 图像理解（规划中） | GPT-4o/Gemini |
| **OpenAI Chat API** | 完整支持 | ✅ |
| **OpenAI Responses API** | o1/o3 推理模型 | ✅ |
| **Gemini 原生 SDK** | 高级特性支持 | ✅ |
| **Gemini OpenAI 兼容** | 快速集成 | ✅ |
| **中国模型** | DeepSeek/千问/Kimi | ✅ |

### 🚧 规划中

- [ ] Azure OpenAI 完整支持
- [ ] Claude (Anthropic) 优化
- [ ] 腾讯混元
- [ ] 本地模型 (Ollama)
- [ ] 批量请求
- [ ] 嵌入模型

---

## 支持的模型

### 🌟 推荐模型

| 模型 | 适用场景 | 特点 | API 端点 |
|------|---------|------|---------|
| **GPT-4o** | 通用对话、多模态 | 最新旗舰，速度快 | OpenAI |
| **o3-mini** | 复杂推理、编程 | 强推理能力 | Responses API |
| **Gemini 2.0 Flash** | 超大上下文 | 2M tokens，多模态 | Gemini 原生 |
| **DeepSeek V3** | 性价比、中文 | 国产最强开源 | OpenAI 兼容 |
| **通义千问** | 中文场景 | 阿里云生态 | OpenAI 兼容 |

### 📋 完整模型列表

#### OpenAI 系列

**Chat Completions API** (标准对话):
```
- gpt-4o / gpt-4o-mini         # 多模态旗舰（文本+图像）
- gpt-4-turbo / gpt-4          # GPT-4 系列
- gpt-3.5-turbo                # 经典高性价比
```

**Responses API** (推理模型):
```
- o3 / o3-mini                 # 最新推理模型（2025）
- o1-preview / o1-mini         # o1 系列
- gpt-4.1 / gpt-4.1-mini       # 新一代
- gpt-4.5                      # Orion（2025年2月）
```

#### Google Gemini

```
- gemini-2.0-flash-exp         # 最新实验版
- gemini-1.5-pro               # 旗舰模型（2M 上下文）
- gemini-1.5-flash             # 快速模型（1M 上下文）
- gemini-1.5-flash-8b          # 超轻量版
```

**两种使用方式**:
1. **原生 SDK** (`CreateGemini`) - 推荐，支持 Thinking 等高级特性
2. **OpenAI 兼容** (`CreateOpenAI`) - 快速集成现有代码

#### 中国大模型

| 厂商 | 模型 | OpenAI 兼容端点 | 特点 |
|------|------|----------------|------|
| **DeepSeek** | deepseek-v3 | `https://api.deepseek.com/v1` | 开源最强，推理能力优秀 |
| **阿里云** | qwen-plus/qwen-max | `https://dashscope.aliyuncs.com/compatible-mode/v1` | 阿里云生态 |
| **月之暗面** | moonshot-v1 | `https://api.moonshot.cn/v1` | Kimi，200K 上下文 |
| **智谱AI** | glm-4 | `https://open.bigmodel.cn/api/paas/v4` | ChatGLM 系列 |

---

## 快速开始

### 安装

```bash
go get github.com/golang-io/OpenLLM
```

### 1 分钟上手

```go
package main

import (
    "context"
    "fmt"
    "github.com/golang-io/OpenLLM"
)

func main() {
    // 创建 LLM 客户端
    llm := OpenLLM.CreateOpenAI(
        OpenLLM.APIKey("your-api-key"),
    )

    // 调用模型
    output, err := llm.Completion(context.Background(), &OpenLLM.Input{
        Model: "gpt-4o",
        Messages: []OpenLLM.Message{
            OpenLLM.UserMessage("你好，介绍一下Go语言"),
        },
    })
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output.Content)
}
```

---

## 详细功能说明

### 1. 基本对话

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.APIKey("sk-xxx"),
)

output, err := llm.Completion(ctx, &OpenLLM.Input{
    Model: "gpt-4o",
    Messages: []OpenLLM.Message{
        OpenLLM.SystemMessage("你是一个编程助手"),
        OpenLLM.UserMessage("如何实现快速排序？"),
    },
})

fmt.Println(output.Content)
```

### 2. 流式输出

#### 基本流式输出

```go
output, err := llm.CompletionStream(
    ctx,
    input,
    func(content string) {
        fmt.Print(content)  // 实时输出每个字符
    },
)

// 流式结束后，output 包含完整内容
fmt.Println("\n\n完整内容:", output.Content)
```

#### 流式输出的工具调用处理

**重要**：工具调用的参数是分片传输的，必须等待流式结束后才能获取完整信息。

```go
output, err := llm.CompletionStream(ctx, input, func(content string) {
    fmt.Print(content)  // ✅ 只有普通文本会实时输出
})

// ✅ 流式结束后检查工具调用
if len(output.ToolCalls) > 0 {
    fmt.Println("\n收到工具调用请求:")
    for _, call := range output.ToolCalls {
        fmt.Printf("- %s(%v)\n", call.Name, call.Arguments)
    }
    
    // 此时参数已完整，可以执行工具
    results := executeTools(output.ToolCalls)
}
```

**流式处理流程**：

```
LLM Stream Start
    ↓
Content Chunk → 实时输出给用户
    ↓
Tool Call Chunk → 内部累积（不输出）
    ↓
Stream End
    ↓
返回 Output {
    Content: "累积的文本",
    ToolCalls: [完整的工具调用]
}
```

### 3. 工具调用完整指南

#### 3.1 定义工具

**简单工具**：

```go
weatherTool := &OpenLLM.Tool{
    Name:        "get_weather",
    Description: "获取指定城市的天气信息",
    Parameters: &OpenLLM.JSONSchema{
        Type: "object",
        Properties: map[string]*OpenLLM.JSONSchema{
            "city": {
                Type:        "string",
                Description: "城市名称，如北京、上海",
            },
        },
        Required: []string{"city"},
    },
}
```

**复杂工具（带枚举和默认值）**：

```go
searchTool := &OpenLLM.Tool{
    Name:        "search_web",
    Description: "在互联网上搜索信息",
    Parameters: &OpenLLM.JSONSchema{
        Type: "object",
        Properties: map[string]*OpenLLM.JSONSchema{
            "query": {
                Type:        "string",
                Description: "搜索关键词",
            },
            "limit": {
                Type:        "integer",
                Description: "返回结果数量",
                Default:     10,
                Minimum:     float64Ptr(1),
                Maximum:     float64Ptr(100),
            },
            "lang": {
                Type:        "string",
                Description: "搜索语言",
                Enum:        []any{"zh", "en", "ja"},
                Default:     "zh",
            },
        },
        Required: []string{"query"},
    },
    Metadata: OpenLLM.ToolMetadata{
        Category: "search",
        Version:  "1.0.0",
        Tags:     []string{"web", "search"},
    },
}

func float64Ptr(f float64) *float64 { return &f }
```

**嵌套对象工具**：

```go
databaseTool := &OpenLLM.Tool{
    Name:        "query_database",
    Description: "查询数据库",
    Parameters: &OpenLLM.JSONSchema{
        Type: "object",
        Properties: map[string]*OpenLLM.JSONSchema{
            "connection": {
                Type:        "object",
                Description: "数据库连接信息",
                Properties: map[string]*OpenLLM.JSONSchema{
                    "host": {
                        Type:        "string",
                        Description: "数据库主机",
                    },
                    "port": {
                        Type:        "integer",
                        Description: "端口号",
                        Default:     5432,
                    },
                },
                Required: []string{"host"},
            },
            "query": {
                Type:        "string",
                Description: "SQL查询语句",
            },
        },
        Required: []string{"connection", "query"},
    },
}
```

#### 3.2 工具验证

```go
if err := weatherTool.Validate(); err != nil {
    log.Fatal("工具定义错误:", err)
}
```

验证规则：
- ✅ Name 不能为空
- ✅ Description 不能为空
- ✅ Parameters 必须是 object 类型
- ✅ Required 字段必须在 Properties 中定义
- ✅ 递归验证所有嵌套 Schema

#### 3.3 调用 LLM

```go
output, err := llm.Completion(ctx, &OpenLLM.Input{
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("北京天气怎么样？"),
    },
    Tools: []OpenLLM.Tool{*weatherTool},
    // 可选：工具选择策略
    ToolChoice: &OpenLLM.ToolChoiceOption{
        Type: OpenLLM.ToolChoiceAuto, // auto | required | none
    },
})
```

#### 3.4 处理工具调用

**单次工具调用**：

```go
if len(output.ToolCalls) > 0 {
    for _, call := range output.ToolCalls {
        fmt.Printf("工具: %s\n", call.Name)
        fmt.Printf("参数: %v\n", call.Arguments)
        
        // 执行工具（由业务层实现）
        result := executeWeatherAPI(call.Arguments["city"].(string))
        
        // 构建下一轮对话
        input.Messages = append(input.Messages,
            OpenLLM.AssistantMessageWithTools("", output.ToolCalls),
            OpenLLM.ToolMessage(result, call.ID),
        )
    }
    
    // 继续对话获取最终答案
    finalOutput, _ := llm.Completion(ctx, input)
    fmt.Println("最终答案:", finalOutput.Content)
}
```

**并行工具调用**：

LLM 可能同时返回多个工具调用，你可以选择并行或顺序执行：

```go
if len(output.ToolCalls) > 1 {
    // 判断是否可以并行执行
    if canExecuteParallel(output.ToolCalls) {
        results := executeToolsParallel(ctx, output.ToolCalls)
    } else {
        results := executeToolsSequential(ctx, output.ToolCalls)
    }
}

// 并行执行示例
func executeToolsParallel(ctx context.Context, calls []OpenLLM.ToolCall) map[string]string {
    results := make(map[string]string)
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, call := range calls {
        wg.Add(1)
        go func(tc OpenLLM.ToolCall) {
            defer wg.Done()
            result := executeTool(ctx, tc)
            mu.Lock()
            results[tc.ID] = result
            mu.Unlock()
        }(call)
    }

    wg.Wait()
    return results
}

// 判断是否可以并行执行
func canExecuteParallel(calls []OpenLLM.ToolCall) bool {
    // 简单策略：只读工具可以并行
    for _, call := range calls {
        if isWriteTool(call.Name) {
            return false  // 有写操作，顺序执行
        }
    }
    return true
}
```

#### 3.5 完整对话流程

```go
func chatWithTools(llm OpenLLM.LLM, userQuery string) string {
    input := &OpenLLM.Input{
        Messages: []OpenLLM.Message{
            OpenLLM.UserMessage(userQuery),
        },
        Tools: []OpenLLM.Tool{weatherTool, timeTool},
    }

    // 第一次调用
    output, _ := llm.Completion(ctx, input)
    
    // 循环处理工具调用（支持多轮）
    for len(output.ToolCalls) > 0 {
        fmt.Println("执行工具调用...")
        
        // 添加助手消息（包含工具调用）
        input.Messages = append(input.Messages,
            OpenLLM.AssistantMessageWithTools(output.Content, output.ToolCalls))
        
        // 执行所有工具
        for _, call := range output.ToolCalls {
            result := executeTool(ctx, call)
            input.Messages = append(input.Messages,
                OpenLLM.ToolMessage(result, call.ID))
        }
        
        // 再次调用 LLM
        output, _ = llm.Completion(ctx, input)
    }
    
    return output.Content
}
```

### 4. Thinking 模式（推理过程）

某些模型支持输出推理过程，帮助理解模型的思考方式。

#### 4.1 支持 Thinking 的模型

- ✅ **Gemini 系列**：原生支持 `IncludeThoughts`
- ✅ **o1/o3 系列**：通过 Responses API 原生支持
- ✅ **DeepSeek**：部分模型支持（通过 Responses API）

#### 4.2 Gemini Thinking

```go
llm := OpenLLM.CreateGemini(
    context.Background(),
    OpenLLM.APIKey("your-key"),
)

output, err := llm.Completion(ctx, &OpenLLM.Input{
    Model: "gemini-1.5-flash",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("解释量子纠缠的原理"),
    },
})

// 推理过程
if output.Thinking != "" {
    fmt.Println("思考过程:", output.Thinking)
}

// 最终答案
fmt.Println("答案:", output.Content)
```

#### 4.3 Streaming 模式下区分 Thinking

Gemini 在流式输出中可以区分推理内容和普通内容：

```go
output, err := llm.CompletionStream(ctx, input, func(chunk OpenLLM.StreamChunk) {
    switch chunk.Type {
    case OpenLLM.StreamContentTypeReasoning:
        fmt.Print("[思考] ", chunk.Content)
    case OpenLLM.StreamContentTypeNormal:
        fmt.Print(chunk.Content)
    }
})

// 流式结束后
fmt.Println("\n完整思考过程:", output.Thinking)
fmt.Println("完整答案:", output.Content)
```

---

## 支持的提供商详解

### OpenAI

#### Chat Completions API（标准对话）

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.APIKey("sk-xxx"),
    // 可选参数
    OpenLLM.Temperature(0.7),
    OpenLLM.MaxTokens(2000),
    OpenLLM.TopP(0.9),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "gpt-4o",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("你的问题"),
    },
})
```

#### Responses API（推理模型）

用于 o1/o3 等推理模型：

```go
llm := OpenLLM.CreateOpenAIResponses(
    OpenLLM.APIKey("sk-xxx"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "o3-mini",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("证明费马大定理"),
    },
})

// o1/o3 原生支持推理
fmt.Println("推理过程:", output.Thinking)
fmt.Println("证明:", output.Content)
```

**Responses API vs Chat Completions**:

| 特性 | Chat Completions | Responses API |
|------|------------------|---------------|
| 端点 | `/v1/chat/completions` | `/v1/responses` |
| 推理支持 | 部分模型 | 原生支持 |
| 适用模型 | GPT-3.5/4/4o | o1/o3/gpt-4.1 |
| 工具调用 | Function calling | 增强工具系统 |
| 状态管理 | 无状态 | 支持 Response ID |

### Gemini

Google Gemini 提供两种使用方式：

#### 方式 1：原生 SDK（推荐）

```go
llm := OpenLLM.CreateGemini(
    context.Background(),
    OpenLLM.APIKey("your-gemini-api-key"),
)

// 支持高级特性
output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "gemini-1.5-flash",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("解释相对论"),
    },
})

// 获取 Thinking
fmt.Println("思考:", output.Thinking)
fmt.Println("答案:", output.Content)
```

**优势**：
- ✅ 支持 Thinking（IncludeThoughts）
- ✅ 支持多模态（图像、音频、视频）
- ✅ 超大上下文（最高 2M tokens）
- ✅ 原生特性完整支持

#### 方式 2：OpenAI 兼容端点

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.URL("https://generativelanguage.googleapis.com/v1beta/openai/"),
    OpenLLM.APIKey("your-gemini-api-key"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "gemini-1.5-flash",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("你的问题"),
    },
})
```

**优势**：
- ✅ 快速集成现有 OpenAI 代码
- ✅ 兼容 OpenAI SDK 工具
- ✅ 无需修改业务逻辑

**选择建议**：
- 🎯 **需要 Thinking**：使用原生 SDK
- 🎯 **需要多模态**：使用原生 SDK
- 🎯 **快速迁移**：使用 OpenAI 兼容端点

### 中国大模型

所有中国大模型都通过 OpenAI 兼容端点使用：

#### DeepSeek

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.URL("https://api.deepseek.com/v1"),
    OpenLLM.APIKey("sk-xxx"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "deepseek-v3",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("你的问题"),
    },
})
```

**特点**：
- ✅ 国产最强开源模型
- ✅ 推理能力优秀（支持 Thinking）
- ✅ 性价比极高
- ✅ 完全兼容 OpenAI API

#### 通义千问（阿里云）

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.URL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
    OpenLLM.APIKey("sk-xxx"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "qwen-plus",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("你的问题"),
    },
})
```

**可用模型**：
- `qwen-max` - 旗舰模型
- `qwen-plus` - 性价比模型
- `qwen-turbo` - 快速模型

#### Kimi（月之暗面）

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.URL("https://api.moonshot.cn/v1"),
    OpenLLM.APIKey("sk-xxx"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "moonshot-v1-8k",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("你的问题"),
    },
})
```

**特点**：
- ✅ 200K 超长上下文
- ✅ 中文理解能力强

#### 智谱 AI

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.URL("https://open.bigmodel.cn/api/paas/v4"),
    OpenLLM.APIKey("your-api-key"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "glm-4",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("你的问题"),
    },
})
```

#### 中国模型对比

| 模型 | 上下文 | 特点 | 价格（相对） |
|------|--------|------|-------------|
| **DeepSeek V3** | 64K | 推理能力强 | ⭐ 极低 |
| **通义千问 Max** | 8K | 阿里云生态 | ⭐⭐ 中等 |
| **Kimi** | 200K | 超长上下文 | ⭐⭐⭐ 较高 |
| **ChatGLM-4** | 128K | 多模态 | ⭐⭐ 中等 |

### Azure OpenAI

```go
llm := OpenLLM.CreateAzure(
    OpenLLM.URL("https://your-resource.openai.azure.com/"),
    OpenLLM.APIKey("your-azure-key"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "gpt-4",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("你的问题"),
    },
})
```

**注意**：Azure 不支持某些参数，会自动过滤：
- ⚠️ `Temperature` - 自动忽略
- ⚠️ `TopP` - 自动忽略
- ⚠️ `Seed` - 自动忽略

### Anthropic Claude

```go
llm := OpenLLM.CreateAnthropic(
    OpenLLM.APIKey("sk-ant-xxx"),
)

output, _ := llm.Completion(ctx, &OpenLLM.Input{
    Model: "claude-3-5-sonnet-20241022",
    Messages: []OpenLLM.Message{
        OpenLLM.UserMessage("Hello"),
    },
})
```

**特点**：
- ✅ 长上下文（200K）
- ✅ 强大的代码理解
- ⚠️ API 格式有差异（自动适配）

---

## API 文档

### 核心接口

```go
// LLM 统一接口
type LLM interface {
    // 非流式调用
    Completion(ctx context.Context, input *Input, opts ...Option) (*Output, error)
    
    // 流式调用
    CompletionStream(ctx context.Context, input *Input, streamOutput StreamOutput, opts ...Option) (*Output, error)
    
    // 获取提供商信息
    Provider() ProviderInfo
}
```

### Input 结构

```go
type Input struct {
    Model      string            // 模型名称
    Messages   []Message         // 消息列表
    Tools      []Tool            // 工具列表（可选）
    ToolChoice *ToolChoiceOption // 工具选择策略（可选）
    Stream     bool              // 是否流式（自动设置）
}
```

### Output 结构

```go
type Output struct {
    StartAt      time.Time     // 开始时间
    Content      string        // 文本内容
    Thinking     string        // 推理内容（支持的模型）
    ToolCalls    []ToolCall    // 工具调用列表
    FinishReason string        // 结束原因：stop/tool_calls/length
    TokenUsage   TokenUsage    // Token 使用情况
    Cost         time.Duration // 调用耗时
    RawResponse  any           // 原始响应（调试用）
    Extra        map[string]any // 扩展字段
}
```

### Message 类型

```go
// 系统消息
OpenLLM.SystemMessage("你是一个助手")

// 用户消息
OpenLLM.UserMessage("你好")

// 助手消息
OpenLLM.AssistantMessage("你好！")

// 助手消息（含工具调用）
OpenLLM.AssistantMessageWithTools("", toolCalls)

// 工具结果消息
OpenLLM.ToolMessage("结果", toolCallID)
```

### Tool 结构

```go
type Tool struct {
    Name        string       // 工具名称
    Description string       // 工具描述
    Parameters  *JSONSchema  // 参数 Schema
    Metadata    ToolMetadata // 元数据（可选）
}

type JSONSchema struct {
    Type        string                  // object/string/number/integer/boolean/array
    Properties  map[string]*JSONSchema  // 属性定义
    Required    []string                // 必填字段
    Items       *JSONSchema             // 数组元素（array 类型）
    Enum        []any                   // 枚举值
    Default     any                     // 默认值
    Description string                  // 描述
    Minimum     *float64                // 最小值（number/integer）
    Maximum     *float64                // 最大值（number/integer）
}
```

### 工具选择策略

```go
// 自动选择
ToolChoice: &OpenLLM.ToolChoiceOption{
    Type: OpenLLM.ToolChoiceAuto,
}

// 必须调用工具
ToolChoice: &OpenLLM.ToolChoiceOption{
    Type: OpenLLM.ToolChoiceRequired,
}

// 不调用工具
ToolChoice: &OpenLLM.ToolChoiceOption{
    Type: OpenLLM.ToolChoiceNone,
}

// 指定特定工具
ToolChoice: &OpenLLM.ToolChoiceOption{
    Type:     OpenLLM.ToolChoiceSpecific,
    ToolName: "get_weather",
}
```

### 配置选项

```go
// 创建客户端时的配置选项
OpenLLM.URL("https://api.openai.com/v1")           // API 端点
OpenLLM.APIKey("sk-xxx")                            // API 密钥
OpenLLM.Temperature(0.7)                            // 温度参数（0.0-2.0）
OpenLLM.MaxTokens(2000)                             // 最大输出 tokens
OpenLLM.TopP(0.9)                                   // Top-P 采样（0.0-1.0）
OpenLLM.Seed(42)                                    // 随机种子（可复现）
OpenLLM.HTTPClientOptions(requests.Timeout(30))    // HTTP 配置

// 注意：Model 在 Input 中指定，不是创建客户端时的选项
input := &OpenLLM.Input{
    Model: "gpt-4o",  // ✅ 正确：在这里指定模型
    Messages: []OpenLLM.Message{...},
}
```

---

## 高级主题

### 1. 错误处理

```go
output, err := llm.Completion(ctx, input)
if err != nil {
    // 类型断言获取详细错误
    if llmErr, ok := err.(*OpenLLM.LLMError); ok {
        log.Printf("提供商: %s", llmErr.Provider)
        log.Printf("错误码: %s", llmErr.Code)
        log.Printf("错误信息: %s", llmErr.Message)
        log.Printf("原始错误: %v", llmErr.Err)
    }
    return err
}
```

### 2. 超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

output, err := llm.Completion(ctx, input)
```

### 3. 重试策略

```go
import "github.com/golang-io/requests"

llm := OpenLLM.CreateOpenAI(
    OpenLLM.APIKey("sk-xxx"),
    OpenLLM.HTTPClientOptions(
        requests.Retry(3),                // 重试 3 次
        requests.RetryDelay(time.Second), // 重试间隔
    ),
)
```

### 4. 调试模式

```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.APIKey("sk-xxx"),
    OpenLLM.HTTPClientOptions(
        requests.Trace(1024000), // 打印请求/响应详情
    ),
)
```

### 5. 批量处理

```go
queries := []string{"问题1", "问题2", "问题3"}
results := make([]string, len(queries))

var wg sync.WaitGroup
for i, query := range queries {
    wg.Add(1)
    go func(idx int, q string) {
        defer wg.Done()
        output, _ := llm.Completion(ctx, &OpenLLM.Input{
            Messages: []OpenLLM.Message{
                OpenLLM.UserMessage(q),
            },
        })
        results[idx] = output.Content
    }(i, query)
}

wg.Wait()
```

---

## 最佳实践

### 1. 选择合适的模型

```go
// 复杂推理任务（数学、逻辑、编程）
model := "o3-mini"  // 或 "deepseek-v3"

// 日常对话
model := "gpt-4o-mini"  // 或 "qwen-plus"

// 超大上下文（长文档处理）
model := "gemini-1.5-pro"  // 2M tokens

// 性价比场景
model := "deepseek-v3"  // 或 "qwen-turbo"
```

### 2. 参数调优

```go
// 创意写作（高随机性）
OpenLLM.Temperature(1.2)
OpenLLM.TopP(0.95)

// 精确任务（低随机性）
OpenLLM.Temperature(0.2)
OpenLLM.TopP(0.5)

// 可复现结果
OpenLLM.Seed(42)
OpenLLM.Temperature(0.0)
```

### 3. 工具调用最佳实践

```go
// ✅ 清晰的描述
tool := &OpenLLM.Tool{
    Name:        "get_weather",
    Description: "获取指定城市的实时天气信息，包括温度、湿度、风速等详细数据",
    // ...
}

// ❌ 模糊的描述
tool := &OpenLLM.Tool{
    Name:        "get_weather",
    Description: "天气",  // 太简略
    // ...
}

// ✅ 验证工具定义
if err := tool.Validate(); err != nil {
    log.Fatal(err)
}

// ✅ 判断是否可以并行执行
if canExecuteParallel(output.ToolCalls) {
    results := executeToolsParallel(ctx, output.ToolCalls)
}
```

### 4. 流式输出最佳实践

```go
// ✅ 推荐：简单直观
output, err := llm.CompletionStream(ctx, input, func(content string) {
    fmt.Print(content)
})

// 流式结束后处理工具调用
if len(output.ToolCalls) > 0 {
    executeTools(output.ToolCalls)
}

// ❌ 错误：尝试在流式过程中解析工具参数
output, err := llm.CompletionStream(ctx, input, func(content string) {
    // 不要这样做！工具参数不完整
    // parseToolCall(content)
})
```

### 5. Token 使用优化

```go
// 限制输出长度
OpenLLM.MaxTokens(500)

// 监控 Token 使用
output, _ := llm.Completion(ctx, input)
log.Printf("Prompt Tokens: %d", output.TokenUsage.PromptTokens)
log.Printf("Completion Tokens: %d", output.TokenUsage.CompletionTokens)
log.Printf("Total Tokens: %d", output.TokenUsage.TotalTokens)
```

---

## 常见问题

### Q: 为什么要设计统一接口？

A: 不同 LLM 提供商的 API 差异很大，统一接口可以：
- 🔄 轻松切换模型进行对比
- 📝 业务代码不依赖具体提供商
- 🧪 方便 A/B 测试
- 🚀 快速集成新模型

### Q: 哪些模型支持工具调用？

A: 
- ✅ **OpenAI**: GPT-4/GPT-3.5 系列
- ✅ **Gemini**: 1.5/2.0 系列
- ✅ **DeepSeek**: V3
- ✅ **通义千问**: qwen-plus/qwen-max
- ⚠️ **o1/o3**: Responses API 支持但能力较弱

### Q: Thinking 内容在哪里获取？

A: 
- **Gemini**: `Output.Thinking` 字段（原生支持）
- **o1/o3**: `Output.Thinking` 字段（Responses API）
- **其他模型**: 不支持

### Q: 流式模式下工具调用如何处理？

A: **必须等待流式结束**才能获取完整的 `ToolCalls`。工具参数是分片传输的，流式过程中是不完整的 JSON。

```go
output, _ := llm.CompletionStream(ctx, input, streamFunc)

// ✅ 流式结束后处理
if len(output.ToolCalls) > 0 {
    executeTools(output.ToolCalls)
}
```

### Q: 如何选择 Gemini 使用方式？

A:
- **原生 SDK** (`CreateGemini`): 需要 Thinking、多模态等高级特性
- **OpenAI 兼容** (`CreateOpenAI`): 快速集成，兼容现有代码

### Q: 性能如何？

A: 
- 适配器只做必要的数据转换，几乎无性能损失
- 流式输出直接转发，无缓冲
- 使用 Go 原生类型，零拷贝

### Q: 如何处理并发请求？

A: 
```go
// 创建单个客户端，多个 goroutine 共享
llm := OpenLLM.CreateOpenAI(...)

// 并发调用（客户端是并发安全的）
for i := 0; i < 10; i++ {
    go func() {
        output, _ := llm.Completion(ctx, input)
    }()
}
```

### Q: 支持代理吗？

A:
```go
llm := OpenLLM.CreateOpenAI(
    OpenLLM.APIKey("sk-xxx"),
    OpenLLM.HTTPClientOptions(
        requests.Proxy("http://proxy.example.com:8080"),
    ),
)
```

### Q: 如何添加新的提供商？

A: 实现 `LLM` 接口的两个方法：
1. `Completion(ctx, input, opts) (*Output, error)`
2. `CompletionStream(ctx, input, streamOutput, opts) (*Output, error)`

参考 `openai_client.go` 和 `gemini.go` 的实现。

---

## 贡献指南

欢迎贡献代码、提交 Issue 或改进文档！

### 开发计划

- [ ] 完善 Azure OpenAI 支持
- [ ] 优化 Claude 适配
- [ ] 支持腾讯混元
- [ ] 支持本地模型 (Ollama)
- [ ] 支持批量请求 API
- [ ] 支持嵌入模型（Embeddings）
- [ ] 支持图像生成（DALL-E）
- [ ] 添加更多测试用例

### 提交 PR 前

1. ✅ 运行所有测试：`go test -v ./...`
2. ✅ 检查代码格式：`go fmt ./...`
3. ✅ 添加注释（中英文）
4. ✅ 更新 README.md（如有 API 变更）
5. ✅ 添加单元测试

### 代码规范

- 🔤 **命名**: 使用 Go 惯用命名（ID 而非 Id）
- 📝 **注释**: 中英文双语注释
- 🧪 **测试**: 每个新功能都应有测试
- 📚 **文档**: 更新 README.md 和代码注释

---

## 许可证

MIT License

---

## 相关链接

- [OpenAI API 文档](https://platform.openai.com/docs/api-reference)
- [Google Gemini 文档](https://ai.google.dev/docs)
- [DeepSeek 文档](https://platform.deepseek.com/api-docs/)
- [阿里云百炼](https://help.aliyun.com/zh/model-studio/)
- [Anthropic Claude 文档](https://docs.anthropic.com/)

---

<p align="center">
  <sub>Built with ❤️ for the Go & AI community</sub>
</p>
