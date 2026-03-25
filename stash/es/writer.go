package es

import (
	"context"

	"github.com/kevwan/go-stash/stash/config"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	Writer struct {
		docType  string
		client   *elastic.Client
		inserter *executors.ChunkExecutor
	}

	valueWithIndex struct {
		index string
		val   string
	}
)

func NewWriter(c config.ElasticSearchConf) (*Writer, error) {
	logx.Infof("Creating ES client for hosts: %v, compress: %v, maxChunkBytes: %d",
		c.Hosts, c.Compress, c.MaxChunkBytes)

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(c.Hosts...),
		elastic.SetGzip(c.Compress),
		elastic.SetBasicAuth(c.Username, c.Password),
	)
	if err != nil {
		logx.Errorf("Failed to create ES client: %v", err)
		return nil, err
	}

	logx.Info("ES client created successfully")
	writer := Writer{
		docType: c.DocType,
		client:  client,
	}
	writer.inserter = executors.NewChunkExecutor(writer.execute, executors.WithChunkBytes(c.MaxChunkBytes))
	logx.Infof("ES writer initialized with docType: %s, maxChunkBytes: %d", c.DocType, c.MaxChunkBytes)
	return &writer, nil
}

func (w *Writer) Write(index, val string) error {
	logx.Infof("Adding document to queue for index: %s, size: %d bytes", index, len(val))
	return w.inserter.Add(valueWithIndex{
		index: index,
		val:   val,
	}, len(val))
}

func (w *Writer) execute(vals []interface{}) {
	logx.Infof("Executing bulk write with %d documents", len(vals))

	var bulk = w.client.Bulk()
	for _, val := range vals {
		pair := val.(valueWithIndex)
		req := elastic.NewBulkIndexRequest().Index(pair.index)

		// 注意：在Elasticsearch 7.x+ 版本中，_type参数已被废弃
		// 在Elasticsearch 8.x 版本中，_type参数被完全移除
		// 如果您使用的是ES 6.x，可以取消下面注释来启用类型设置
		// if len(w.docType) > 0 {
		//     req = req.Type(w.docType)
		// }

		req = req.Doc(pair.val)
		bulk.Add(req)
	}

	logx.Infof("Sending bulk request to Elasticsearch...")
	resp, err := bulk.Do(context.Background())
	if err != nil {
		logx.Errorf("Elasticsearch bulk operation failed: %v", err)
		return
	}

	logx.Infof("Bulk operation completed. Errors: %v, Took: %d ms", resp.Errors, resp.Took)

	// bulk error in docs will report in response items
	if !resp.Errors {
		logx.Infof("All %d documents successfully indexed", len(vals))
		return
	}

	logx.Errorf("Bulk operation reported errors, checking individual items...")
	errorCount := 0
	for _, imap := range resp.Items {
		for action, item := range imap {
			if item.Error == nil {
				continue
			}

			errorCount++
			logx.Errorf("ES indexing error [%s]: Index=%s, ID=%s, Status=%d, Error=%v",
				action, item.Index, item.Id, item.Status, item.Error)
		}
	}

	logx.Errorf("Total errors in bulk operation: %d out of %d documents", errorCount, len(vals))
}
