/**
 * Package filter
 * @Author: tbb
 * @Date: 2023/12/5 17:01
 */
package filter

import (
	"fmt"
	"time"
)

func TimestampFilter(field string) FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		val, ok := m[field]
		if !ok {
			return m
		}

		timestamp, ok := val.(float64)
		if !ok {
			return m
		}
		t := time.Unix(int64(timestamp), 0)
		// 将 time.Time 类型转换为格式化的日期和时间字符串
		// 注意：这里设置了毫秒为 343，因为您提供的示例中包含了毫秒
		formattedTime := fmt.Sprintf("%s.%03d", t.Format("2006-01-02 15:04:05"), 343)

		m["@timestamp"] = formattedTime
		return m
	}
}
