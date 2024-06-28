/**
 * Package filter
 * @Author: tbb
 * @Date: 2024/6/26 15:13
 */
package filter

func AddFilter(inField, outField string, match map[string]string) FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		if len(match) > 0 {
			for i, o := range match {
				v, ok := m[i]
				if !ok {
					continue
				}
				m[o] = v
			}
		}

		val, ok := m[inField]
		if !ok {
			return m
		}

		m[outField] = val

		return m
	}
}
