package OpenLLM

import (
	"fmt"

	"github.com/openai/openai-go/v3"
)

// ============================================================================
// Tool 转换方法 / Tool Conversion Methods
// ============================================================================

// ToOpenAITool 将Tool转换为OpenAI ChatCompletionToolUnionParam
// ToOpenAITool converts Tool to OpenAI ChatCompletionToolUnionParam
func (t *Tool) ToOpenAITool() (openai.ChatCompletionToolUnionParam, error) {
	if err := t.Validate(); err != nil {
		return openai.ChatCompletionToolUnionParam{}, fmt.Errorf("工具验证失败: %w", err)
	}

	params := t.Parameters.ToOpenAIParameters()

	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        t.Name,
		Description: openai.String(t.Description),
		Parameters:  params,
	}), nil
}

// Validate 验证Tool定义是否有效
// Validate checks if the Tool definition is valid
func (t *Tool) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("工具名称不能为空")
	}
	if t.Description == "" {
		return fmt.Errorf("工具描述不能为空")
	}
	if t.Parameters == nil {
		return fmt.Errorf("工具参数定义不能为空")
	}
	if t.Parameters.Type != "object" {
		return fmt.Errorf("工具参数类型必须是object，当前为: %s", t.Parameters.Type)
	}
	// 验证参数Schema
	return t.Parameters.Validate()
}

// ============================================================================
// JSONSchema 转换方法 / JSONSchema Conversion Methods
// ============================================================================

// ToOpenAIParameters 将JSONSchema转换为OpenAI FunctionParameters
// ToOpenAIParameters converts JSONSchema to OpenAI FunctionParameters
func (s *JSONSchema) ToOpenAIParameters() openai.FunctionParameters {
	params := openai.FunctionParameters{
		"type": s.Type,
	}

	if s.Description != "" {
		params["description"] = s.Description
	}

	if len(s.Properties) > 0 {
		properties := make(map[string]any)
		for name, prop := range s.Properties {
			properties[name] = prop.toOpenAIProperty()
		}
		params["properties"] = properties
	}

	if len(s.Required) > 0 {
		params["required"] = s.Required
	}

	if s.Items != nil {
		params["items"] = s.Items.toOpenAIProperty()
	}

	if len(s.Enum) > 0 {
		params["enum"] = s.Enum
	}

	if s.Default != nil {
		params["default"] = s.Default
	}

	if s.Minimum != nil {
		params["minimum"] = *s.Minimum
	}

	if s.Maximum != nil {
		params["maximum"] = *s.Maximum
	}

	return params
}

// toOpenAIProperty 将JSONSchema转换为OpenAI属性格式（内部方法）
// toOpenAIProperty converts JSONSchema to OpenAI property format (internal method)
func (s *JSONSchema) toOpenAIProperty() map[string]any {
	prop := map[string]any{
		"type": s.Type,
	}

	if s.Description != "" {
		prop["description"] = s.Description
	}

	if len(s.Properties) > 0 {
		properties := make(map[string]any)
		for name, p := range s.Properties {
			properties[name] = p.toOpenAIProperty()
		}
		prop["properties"] = properties
	}

	if len(s.Required) > 0 {
		prop["required"] = s.Required
	}

	if s.Items != nil {
		prop["items"] = s.Items.toOpenAIProperty()
	}

	if len(s.Enum) > 0 {
		prop["enum"] = s.Enum
	}

	if s.Default != nil {
		prop["default"] = s.Default
	}

	if s.Minimum != nil {
		prop["minimum"] = *s.Minimum
	}

	if s.Maximum != nil {
		prop["maximum"] = *s.Maximum
	}

	return prop
}

// Validate 验证JSONSchema是否有效
// Validate checks if the JSONSchema is valid
func (s *JSONSchema) Validate() error {
	if s.Type == "" {
		return fmt.Errorf("JSONSchema类型不能为空")
	}

	// 验证必填字段是否在properties中定义
	if s.Type == "object" && len(s.Required) > 0 {
		for _, required := range s.Required {
			if _, exists := s.Properties[required]; !exists {
				return fmt.Errorf("必填字段 %s 未在properties中定义", required)
			}
		}
	}

	// 递归验证嵌套的Schema
	for name, prop := range s.Properties {
		if err := prop.Validate(); err != nil {
			return fmt.Errorf("属性 %s 验证失败: %w", name, err)
		}
	}

	if s.Items != nil {
		if err := s.Items.Validate(); err != nil {
			return fmt.Errorf("数组items验证失败: %w", err)
		}
	}

	return nil
}
