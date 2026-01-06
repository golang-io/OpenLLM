package OpenLLM

// ============================================================================
// 测试用 Mock 工具 / Mock Tools for Testing
// ============================================================================

// Tools 预定义的测试工具列表，用于单元测试和示例
// Tools is a predefined list of test tools for unit testing and examples
var Tools = []Tool{
	{
		Name:        "query_weather",
		Description: "指定一个时间，查询一个城市这个时间的天气情况",
		Parameters: &JSONSchema{
			Type: "object",
			Properties: map[string]*JSONSchema{
				"city": {
					Type:        "string",
					Description: "城市信息：例如：beijing，guangzhou",
				},
				"time": {
					Type:        "string",
					Description: "时间：例如：2006-01-02 15:04:05",
				},
			},
			Required: []string{"city"},
		},
	},
	{
		Name:        "current_time",
		Description: "查询当前时间",
		Parameters: &JSONSchema{
			Type: "object",
			Properties: map[string]*JSONSchema{
				"tz": {
					Type:        "string",
					Description: "时区：例如：Asia/Shanghai",
				},
			},
			Required: []string{"tz"},
		},
	},
}
