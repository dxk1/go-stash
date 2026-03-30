/**
 * Package filter
 * @Author: tbb
 * @Date: 2026/03/25
 */
package filter

import (
	"time"
)

const timestampKey = "@timestamp"

// EnsureTimestampFilter 确保每个文档都有@timestamp字段
// 如果已经存在@timestamp字段则保持不变，否则尝试从其他时间字段推导，最后使用当前时间
func EnsureTimestampFilter() FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		// 检查是否已经有@timestamp字段
		if _, exists := m[timestampKey]; exists {
			// 如果已经存在@timestamp字段，直接返回
			return m
		}

		var t time.Time
		var timestampSet bool

		// 尝试从常见的时间字段中获取时间戳
		timeFields := []string{"created_at", "timestamp", "Timestamp", "time", "@timestamp", "event_time", "log_time"}

		for _, field := range timeFields {
			if val, exists := m[field]; exists {
				// 尝试解析不同类型的时间戳
				if timeStr, ok := val.(string); ok {
					// 支持多种时间格式
					formats := []string{
						"2006-01-02T15:04:05.000Z",
						"2006-01-02T15:04:05Z",
						"2006-01-02T15:04:05.999999999Z07:00", // RFC3339Nano
						"2006-01-02T15:04:05Z07:00",           // RFC3339
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

					if timestampSet {
						break
					}
				} else if timestamp, ok := val.(float64); ok {
					// Unix时间戳（秒）
					t = time.Unix(int64(timestamp), 0)
					timestampSet = true
					break
				} else if timestamp, ok := val.(int64); ok {
					// Unix时间戳（秒）
					t = time.Unix(timestamp, 0)
					timestampSet = true
					break
				} else if timeObj, ok := val.(time.Time); ok {
					// 直接的time.Time对象
					t = timeObj
					timestampSet = true
					break
				}
			}
		}

		// 如果没有找到任何时间字段，使用当前时间
		if !timestampSet {
			t = time.Now()
		}

		// 格式化为Elasticsearch标准格式
		formattedTimestamp := t.Format("2006-01-02T15:04:05.000Z")

		m[timestampKey] = formattedTimestamp
		return m
	}
}
