package handler

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kevwan/go-stash/stash/es"
	"github.com/kevwan/go-stash/stash/filter"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessageHandler struct {
	writer         *es.Writer
	indexer        *es.Index
	filters        []filter.FilterFunc
	processedCount int64     // 成功处理的消息数量
	errorCount     int64     // 错误处理的消息数量
	startTime      time.Time // 启动时间
}

func NewHandler(writer *es.Writer, indexer *es.Index) *MessageHandler {
	handler := &MessageHandler{
		writer:    writer,
		indexer:   indexer,
		startTime: time.Now(),
	}

	// 启动统计输出协程
	go handler.startStatsReporter()

	return handler
}

func (mh *MessageHandler) AddFilters(filters ...filter.FilterFunc) {
	for _, f := range filters {
		mh.filters = append(mh.filters, f)
	}
}

func (mh *MessageHandler) Consume(key, val string) error {
	// 双重日志确保能看到消费记录
	consumeMsg := fmt.Sprintf("✅ KAFKA CONSUMED - Key: %s, Size: %d bytes", key, len(val))
	fmt.Printf("%s\n", consumeMsg)
	logx.Infof("Processing message: key=%s, size=%d, data=%s", key, len(val), val[:min(len(val), 200)])

	var m map[string]interface{}
	if err := jsoniter.Unmarshal([]byte(val), &m); err != nil {
		atomic.AddInt64(&mh.errorCount, 1)
		errorMsg := fmt.Sprintf("❌ ERROR - Failed to unmarshal JSON: %v", err)
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		logx.Errorf("Failed to unmarshal JSON: %v, data: %s", err, val[:min(len(val), 100)])
		return err
	}

	index := mh.indexer.GetIndex(m)
	logx.Infof("Using ES index: %s", index)

	for _, proc := range mh.filters {
		if m = proc(m); m == nil {
			logx.Info("Message dropped by filter")
			return nil
		}
	}

	// 最后一道保险：确保@timestamp字段一定存在
	mh.ensureTimestamp(m)

	bs, err := jsoniter.Marshal(m)
	if err != nil {
		atomic.AddInt64(&mh.errorCount, 1)
		errorMsg := fmt.Sprintf("❌ ERROR - Failed to marshal processed data: %v", err)
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		logx.Errorf("Failed to marshal processed data: %v", err)
		return err
	}

	// 记录最终要写入ES的数据
	logx.Infof("Writing to ES index: %s, data: %s", index, string(bs)[:min(len(bs), 200)])

	if err := mh.writer.Write(index, string(bs)); err != nil {
		atomic.AddInt64(&mh.errorCount, 1)
		// 双重错误输出，确保能看到
		errorMsg := fmt.Sprintf("❌ CRITICAL: Failed to write to Elasticsearch: %v, index: %s (Total errors: %d)",
			err, index, atomic.LoadInt64(&mh.errorCount))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		logx.Errorf("Failed to write to Elasticsearch: %v, index: %s, data: %s", err, index, string(bs)[:min(len(bs), 200)])
		return err
	}

	// 增加成功计数
	atomic.AddInt64(&mh.processedCount, 1)

	// 双重成功日志确保能看到
	successMsg := fmt.Sprintf("✅ SUCCESS - Kafka message processed and written to ES: index=%s (Total: %d)",
		index, atomic.LoadInt64(&mh.processedCount))
	fmt.Printf("%s\n", successMsg)
	logx.Info("Successfully written to Elasticsearch")
	return nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ensureTimestamp 确保数据中有@timestamp字段
func (mh *MessageHandler) ensureTimestamp(m map[string]interface{}) {
	const timestampKey = "@timestamp"

	// 如果已经有@timestamp字段，直接返回
	if _, exists := m[timestampKey]; exists {
		logx.Infof("@timestamp already exists: %v", m[timestampKey])
		return
	}

	logx.Info("@timestamp field missing, attempting to generate one")

	var t time.Time
	var timestampSet bool

	// 尝试从常见的时间字段中获取时间戳
	timeFields := []string{"created_at", "timestamp", "Timestamp", "time", "event_time", "log_time"}

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
		logx.Info("No timestamp field found, using current time")
	} else {
		logx.Infof("Found timestamp field, parsed time: %v", t)
	}

	// 格式化为Elasticsearch标准格式
	formattedTimestamp := t.Format("2006-01-02T15:04:05.000Z")

	m[timestampKey] = formattedTimestamp
	logx.Infof("Set @timestamp to: %s", formattedTimestamp)
}

// startStatsReporter 启动统计报告协程
func (mh *MessageHandler) startStatsReporter() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒报告一次统计信息
	defer ticker.Stop()

	// 立即更新健康检查文件
	mh.updateHealthCheck()

	for range ticker.C {
		processed := atomic.LoadInt64(&mh.processedCount)
		errors := atomic.LoadInt64(&mh.errorCount)
		uptime := time.Since(mh.startTime)

		rate := float64(0)
		if uptime.Seconds() > 0 {
			rate = float64(processed) / uptime.Seconds()
		}

		statsMsg := fmt.Sprintf("📊 KAFKA CONSUME STATS - Processed: %d, Errors: %d, Rate: %.2f msg/s, Uptime: %v",
			processed, errors, rate, uptime.Truncate(time.Second))

		fmt.Printf("%s\n", statsMsg)
		logx.Infof("Handler stats - Processed: %d, Errors: %d, Rate: %.2f msg/s, Uptime: %v",
			processed, errors, rate, uptime.Truncate(time.Second))

		// 更新健康检查文件
		mh.updateHealthCheck()
	}
}

// updateHealthCheck 更新健康检查文件
func (mh *MessageHandler) updateHealthCheck() {
	healthFile := "/tmp/go-stash-health"
	processed := atomic.LoadInt64(&mh.processedCount)
	errors := atomic.LoadInt64(&mh.errorCount)

	healthData := fmt.Sprintf("processed=%d,errors=%d,timestamp=%d",
		processed, errors, time.Now().Unix())

	// 忽略错误，健康检查是非关键功能
	_ = os.WriteFile(healthFile, []byte(healthData), 0644)
}
