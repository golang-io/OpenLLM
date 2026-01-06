package OpenLLM

import (
	"testing"
)

func TestTool_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tool    *Tool
		wantErr bool
	}{
		{
			name: "valid tool",
			tool: &Tool{
				Name:        "get_weather",
				Description: "获取城市天气",
				Parameters: &JSONSchema{
					Type: "object",
					Properties: map[string]*JSONSchema{
						"city": {
							Type:        "string",
							Description: "城市名称",
						},
					},
					Required: []string{"city"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			tool: &Tool{
				Description: "test",
				Parameters: &JSONSchema{
					Type: "object",
				},
			},
			wantErr: true,
		},
		{
			name: "missing description",
			tool: &Tool{
				Name: "test",
				Parameters: &JSONSchema{
					Type: "object",
				},
			},
			wantErr: true,
		},
		{
			name: "parameters not object type",
			tool: &Tool{
				Name:        "test",
				Description: "test",
				Parameters: &JSONSchema{
					Type: "string",
				},
			},
			wantErr: true,
		},
		{
			name: "required field not in properties",
			tool: &Tool{
				Name:        "test",
				Description: "test",
				Parameters: &JSONSchema{
					Type: "object",
					Properties: map[string]*JSONSchema{
						"city": {
							Type: "string",
						},
					},
					Required: []string{"country"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tool.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tool.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTool_ToOpenAITool(t *testing.T) {
	tool := &Tool{
		Name:        "get_weather",
		Description: "获取城市天气",
		Parameters: &JSONSchema{
			Type: "object",
			Properties: map[string]*JSONSchema{
				"city": {
					Type:        "string",
					Description: "城市名称",
				},
				"unit": {
					Type:        "string",
					Description: "温度单位",
					Enum:        []any{"celsius", "fahrenheit"},
					Default:     "celsius",
				},
			},
			Required: []string{"city"},
		},
		Metadata: ToolMetadata{
			Category: "weather",
			Version:  "1.0.0",
			Tags:     []string{"weather", "api"},
		},
	}

	openaiTool, err := tool.ToOpenAITool()
	if err != nil {
		t.Fatalf("ToOpenAITool() error = %v", err)
	}

	if openaiTool.OfFunction == nil {
		t.Fatal("Expected OfFunction to be non-nil")
	}

	funcDef := openaiTool.OfFunction.Function
	if funcDef.Name != "get_weather" {
		t.Errorf("Expected name 'get_weather', got %s", funcDef.Name)
	}

	if funcDef.Description.Or("") != "获取城市天气" {
		t.Errorf("Expected description '获取城市天气', got %s", funcDef.Description.Or(""))
	}

	params := funcDef.Parameters
	if params["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", params["type"])
	}

	properties, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be map[string]any")
	}

	if len(properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(properties))
	}

	required, ok := params["required"].([]string)
	if !ok {
		t.Fatal("Expected required to be []string")
	}

	if len(required) != 1 || required[0] != "city" {
		t.Errorf("Expected required ['city'], got %v", required)
	}
}

func TestJSONSchema_Validate(t *testing.T) {
	tests := []struct {
		name    string
		schema  *JSONSchema
		wantErr bool
	}{
		{
			name: "valid object schema",
			schema: &JSONSchema{
				Type: "object",
				Properties: map[string]*JSONSchema{
					"name": {
						Type: "string",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid array schema",
			schema: &JSONSchema{
				Type: "array",
				Items: &JSONSchema{
					Type: "string",
				},
			},
			wantErr: false,
		},
		{
			name: "nested object schema",
			schema: &JSONSchema{
				Type: "object",
				Properties: map[string]*JSONSchema{
					"address": {
						Type: "object",
						Properties: map[string]*JSONSchema{
							"city": {
								Type: "string",
							},
						},
						Required: []string{"city"},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty type",
			schema:  &JSONSchema{},
			wantErr: true,
		},
		{
			name: "required field not in properties",
			schema: &JSONSchema{
				Type:     "object",
				Required: []string{"name"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONSchema.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToolCall_ID(t *testing.T) {
	// 测试ToolCall.ID字段（确保从Id改为ID）
	tc := ToolCall{
		ID:        "call_123",
		Name:      "test",
		Arguments: map[string]any{"key": "value"},
	}

	if tc.ID != "call_123" {
		t.Errorf("Expected ID 'call_123', got %s", tc.ID)
	}
}
