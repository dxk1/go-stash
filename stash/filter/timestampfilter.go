/**
 * Package filter
 * @Author: tbb
 * @Date: 2023/12/5 17:01
 */
package filter

import (
	"time"
)

func TimestampFilter(field string) FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		// 尝试从指定字段获取timestamp
		val, ok := m[field]
		var t time.Time
		var timestampSet bool

		if ok {
			// 如果字段存在，尝试解析timestamp
			if timestamp, isFloat := val.(float64); isFloat {
				// 支持秒级时间戳
				t = time.Unix(int64(timestamp), 0)
				timestampSet = true
			} else if timestampStr, isString := val.(string); isString {
				// 支持字符串格式的时间戳
				// 尝试多种时间格式
				formats := []string{
					"2006-01-02T15:04:05.000Z",
					"2006-01-02T15:04:05Z",
					"2006-01-02 15:04:05",
					time.RFC3339,
					time.RFC3339Nano,
				}

				for _, format := range formats {
					if parsed, err := time.Parse(format, timestampStr); err == nil {
						t = parsed
						timestampSet = true
						break
					}
				}
			} else if timestampInt, isInt := val.(int64); isInt {
				// 支持int64类型的时间戳
				t = time.Unix(timestampInt, 0)
				timestampSet = true
			}
		}

		// 如果没有成功解析到timestamp，使用当前时间
		if !timestampSet {
			t = time.Now()
		}

		// 格式化为Elasticsearch标准的时间格式
		formattedTimestamp := t.Format("2006-01-02T15:04:05.000Z")

		// 设置@timestamp字段（字符串格式，符合ES标准）
		m["@timestamp"] = formattedTimestamp

		// 保留原始时间对象（如果需要的话）
		m["@timestamp_obj"] = t

		return m
	}
}
