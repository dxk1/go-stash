package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"

	"github.com/zeromicro/go-zero/core/service"
)

type (
	Condition struct {
		Key   string `yaml:"Key"`
		Value string `yaml:"Value"`
		Type  string `json:",default=match,options=match|contains"  yaml:"Type"`
		Op    string `json:",default=and,options=and|or"  yaml:"Op"`
	}

	ElasticSearchConf struct {
		Hosts         []string `yaml:"Hosts"`
		Index         string   `yaml:"Index"`
		DocType       string   `json:",default=doc" yaml:"DocType"`
		TimeZone      string   `json:",optional" yaml:"TimeZone"`
		MaxChunkBytes int      `json:",default=15728640" yaml:"MaxChunkBytes"` // default 15M
		Compress      bool     `json:",default=false" yaml:"Compress"`
		Username      string   `json:",optional" yaml:"Username"`
		Password      string   `json:",optional" yaml:"Password"`
	}

	Filter struct {
		Action     string            `json:",options=drop|remove_field|transfer|timestamp|add" yaml:"Action"`
		Conditions []Condition       `json:",optional" yaml:"Conditions"`
		Fields     []string          `json:",optional" yaml:"Fields"`
		Field      string            `json:",optional" yaml:"Field"`
		Target     string            `json:",optional" yaml:"Target"`
		Match      map[string]string `json:",optional" yaml:"Match"`
	}

	KafkaConf struct {
		service.ServiceConf
		Brokers    []string `yaml:"Brokers"`
		Group      string   `yaml:"Group"`
		Topics     []string `yaml:"Topics"`
		Offset     string   `json:",options=first|last,default=last" yaml:"Offset"`
		Conns      int      `json:",default=1" yaml:"Conns"`
		Consumers  int      `json:",default=8" yaml:"Consumers"`
		Processors int      `json:",default=8" yaml:"Processors"`
		MinBytes   int      `json:",default=10240" yaml:"MinBytes"`    // 10K
		MaxBytes   int      `json:",default=10485760" yaml:"MaxBytes"` // 10M
		Username   string   `json:",optional" yaml:"Username"`
		Password   string   `json:",optional" yaml:"Password"`
	}

	Cluster struct {
		Input struct {
			Kafka KafkaConf `yaml:"Kafka"`
		} `yaml:"Input"`
		Filters []Filter `json:",optional" yaml:"Filters"`
		Output  struct {
			ElasticSearch ElasticSearchConf `yaml:"ElasticSearch"`
		} `yaml:"Output"`
	}

	Config struct {
		Clusters    []Cluster     `yaml:"Clusters"`
		GracePeriod time.Duration `json:",default=10s" yaml:"GracePeriod"`
	}
)

// InitConf ...
func InitConf(conf *Config) {
	f, err := ioutil.ReadFile("/usr/local/services/scf_stash/scf_stash.yaml")
	if err != nil {
		logx.Errorf("read config file failed: %s\n", err)
	}

	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		logx.Errorf("parse config failed: %s\n", err)
	}

	if conf.GracePeriod == 0 {
		conf.GracePeriod = 10 * time.Second
	}

	if len(conf.Clusters) > 0 {
		for i := range conf.Clusters {
			c := &conf.Clusters[i] // 使用索引获取cluster的引用
			if c.Input.Kafka.Conns == 0 {
				c.Input.Kafka.Conns = 1
			}

			if c.Input.Kafka.Consumers == 0 {
				c.Input.Kafka.Consumers = 8
			}

			if c.Input.Kafka.Processors == 0 {
				c.Input.Kafka.Processors = 8
			}

			if c.Input.Kafka.MinBytes == 0 {
				c.Input.Kafka.MinBytes = 10240
			}

			if c.Input.Kafka.MaxBytes == 0 {
				c.Input.Kafka.MaxBytes = 10485760
			}

			if c.Output.ElasticSearch.DocType == "" {
				c.Output.ElasticSearch.DocType = "doc"
			}

			if c.Output.ElasticSearch.MaxChunkBytes == 0 {
				c.Output.ElasticSearch.MaxChunkBytes = 15728640
			}

			if len(c.Filters) > 0 {
				for k := range c.Filters {
					fi := &c.Filters[k] // 使用索引获取Filters的引用
					if len(fi.Conditions) > 0 {
						for j := range fi.Conditions {
							ci := fi.Conditions[j]
							if ci.Type == "" {
								ci.Type = "match"
							}

							if ci.Op == "" {
								ci.Op = "and"
							}
						}
					}
				}
			}
		}

	}
}
