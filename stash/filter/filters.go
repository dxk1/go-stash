package filter

import "github.com/kevwan/go-stash/stash/config"

const (
	filterDrop         = "drop"
	filterRemoveFields = "remove_field"
	filterTransfer     = "transfer"
	filterTimestamp    = "timestamp"
	filterAnalysis     = "analysis"
	filterAdd          = "add"
	opAnd              = "and"
	opOr               = "or"
	typeContains       = "contains"
	typeMatch          = "match"
)

type FilterFunc func(map[string]interface{}) map[string]interface{}

func CreateFilters(p config.Cluster) []FilterFunc {
	var filters []FilterFunc

	for _, f := range p.Filters {
		switch f.Action {
		case filterDrop:
			filters = append(filters, DropFilter(f.Conditions))
		case filterRemoveFields:
			filters = append(filters, RemoveFieldFilter(f.Fields))
		case filterTransfer:
			filters = append(filters, TransferFilter(f.Field, f.Target))
		case filterTimestamp:
			filters = append(filters, TimestampFilter(f.Field))
		case filterAnalysis:
			filters = append(filters, AnalysisFilter())
		case filterAdd:
			filters = append(filters, AddFilter(f.Field, f.Target, f.Match))
		}
	}

	// 添加专门处理created_at字段的filter
	filters = append(filters, CreatedAtTimestampFilter())

	// 在所有filters的最后添加EnsureTimestampFilter作为最后的保险
	// 确保@timestamp字段总是存在，如果前面的filters没有设置的话
	filters = append(filters, EnsureTimestampFilter())

	return filters
}
