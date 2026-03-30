/**
 * Package filter
 * @Author: tbb
 * @Date: 2026/03/25
 */
package filter

import (
	"time"
)

// CreatedAtTimestampFilter 专门处理created_at字段的时间戳转换
func CreatedAtTimestampFilter() FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		var t time.Time
		var timestampSet bool

		// 检查是否已经有@timestamp字段
		if _, exists := m["@timestamp"]; exists {
			// 如果已经存在@timestamp字段，直接返回
			return m
		}

		// 尝试从created_at字段获取时间戳
		if val, exists := m["created_at"]; exists {
			if timeStr, ok := val.(string); ok {
				// 支持多种时间格式
				formats := []string{
					"2006-01-02T15:04:05Z", // "2026-03-25T06:47:07Z"
					"2006-01-02T15:04:05.000Z",
					"2006-01-02T15:04:05.999999999Z07:00",
					"2006-01-02T15:04:05Z07:00",
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05",
					time.RFC3339,
					time.RFC3339Nano,
				}

				for _, format := range formats {
					if parsed, err := time.Parse(format, timeStr); err == nil {
						t = parsed
						timestampSet = true
						break
					}
				}
			}
		}

		// 如果没有成功解析created_at，使用当前时间
		if !timestampSet {
			t = time.Now()
		}

		// 格式化为Elasticsearch标准格式
		formattedTimestamp := t.Format("2006-01-02T15:04:05.000Z")

		// 设置@timestamp字段
		m["@timestamp"] = formattedTimestamp

		// 为了兼容旧版本，也设置@timestamp0（如果需要的话）
		m["@timestamp0"] = t

		return m
	}
}
