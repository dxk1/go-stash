package main

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/kevwan/go-stash/stash/config"
	"github.com/kevwan/go-stash/stash/es"
	"github.com/kevwan/go-stash/stash/filter"
	"github.com/kevwan/go-stash/stash/handler"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
)

//var configFile = flag.String("f", "etc/config.yaml", "Specify the config file")

func toKqConf(c config.KafkaConf) []kq.KqConf {
	var ret []kq.KqConf

	for _, topic := range c.Topics {
		ret = append(ret, kq.KqConf{
			ServiceConf: c.ServiceConf,
			Brokers:     c.Brokers,
			Group:       c.Group,
			Topic:       topic,
			Offset:      c.Offset,
			Conns:       c.Conns,
			Consumers:   c.Consumers,
			Processors:  c.Processors,
			MinBytes:    c.MinBytes,
			MaxBytes:    c.MaxBytes,
			Username:    c.Username,
			Password:    c.Password,
		})
	}

	return ret
}

func main() {
	flag.Parse()

	var c config.Config
	//conf.MustLoad(*configFile, &c)
	config.InitConf(&c)
	proc.SetTimeToForceQuit(c.GracePeriod)

	cs, _ := json.Marshal(c)
	logx.Infof("read config file succ: %s\n", string(cs))
	group := service.NewServiceGroup()
	defer group.Stop()

	for i, processor := range c.Clusters {
		logx.Infof("Setting up cluster %d with ES hosts: %v", i+1, processor.Output.ElasticSearch.Hosts)

		client, err := elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL(processor.Output.ElasticSearch.Hosts...),
			elastic.SetBasicAuth(processor.Output.ElasticSearch.Username, processor.Output.ElasticSearch.Password),
		)
		if err != nil {
			logx.Errorf("Failed to create ES client: %v", err)
			logx.Must(err)
		}
		logx.Infof("ES client created successfully")

		logx.Infof("Creating filters for cluster %d", i+1)
		filters := filter.CreateFilters(processor)
		logx.Infof("Created %d filters", len(filters))

		logx.Infof("Creating ES writer...")
		writer, err := es.NewWriter(processor.Output.ElasticSearch)
		if err != nil {
			logx.Errorf("Failed to create ES writer: %v", err)
			logx.Must(err)
		}
		logx.Infof("ES writer created successfully")

		var loc *time.Location
		if len(processor.Output.ElasticSearch.TimeZone) > 0 {
			logx.Infof("Using timezone: %s", processor.Output.ElasticSearch.TimeZone)
			loc, err = time.LoadLocation(processor.Output.ElasticSearch.TimeZone)
			if err != nil {
				logx.Errorf("Failed to load timezone: %v", err)
				logx.Must(err)
			}
		} else {
			logx.Info("Using local timezone")
			loc = time.Local
		}

		logx.Infof("Creating ES indexer with pattern: %s", processor.Output.ElasticSearch.Index)
		indexer := es.NewIndex(client, processor.Output.ElasticSearch.Index, loc)

		logx.Info("Creating message handler and adding filters...")
		handle := handler.NewHandler(writer, indexer)
		handle.AddFilters(filters...)
		handle.AddFilters(filter.AddUriFieldFilter("url", "uri"))

		kafkaConfigs := toKqConf(processor.Input.Kafka)
		logx.Infof("Setting up %d Kafka consumers", len(kafkaConfigs))
		for j, k := range kafkaConfigs {
			logx.Infof("Adding Kafka consumer %d: topic=%s, group=%s", j+1, k.Topic, k.Group)
			group.Add(kq.MustNewQueue(k, handle))
		}
	}

	fmt.Println("🚀 go-stash is starting all services...")
	logx.Info("All Kafka consumers and ES writers configured, starting services...")

	// 启动服务
	group.Start()

	fmt.Println("✅ go-stash started successfully and consuming Kafka messages")
	logx.Info("go-stash service started successfully")
}
