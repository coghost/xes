package xes

import (
	"context"
)

const (
	ReceiverEs int = iota + 1
	ReceiverStdout
	ReceiverPp
	ReceiverDump
)

type Option struct {
	// the log Receiver: is not used by now
	Receiver int `json:"receiver" yaml:"receiver" redis:"receiver" structs:"receiver"`

	Disabled bool `json:"disabled" yaml:"disabled" redis:"disabled" structs:"disabled"`

	// log Level
	Level int `json:"level" yaml:"level" redis:"level" structs:"level"`

	// Retries
	Retries int `json:"retries" yaml:"retries" redis:"retries" structs:"retries"`

	// Timeout
	Timeout int `json:"timeout" yaml:"timeout" redis:"timeout" structs:"timeout"`

	// index name for es
	IndexName   string `json:"index_name" yaml:"index_name" redis:"index_name" structs:"index_name"`
	IndexSuffix string `json:"index_suffix" yaml:"index_suffix" redis:"index_suffix" structs:"index_suffix"`

	// Urls
	Urls []string `json:"urls" yaml:"urls" redis:"urls" structs:"urls"`

	// WithFuncName
	WithFuncName bool `json:"with_func_name" yaml:"with_func_name" redis:"with_func_name" structs:"with_func_name"`
	WithSysinfo  bool `json:"with_sysinfo" yaml:"with_sysinfo" redis:"with_sysinfo" structs:"with_sysinfo"`

	// ctx
	ctx context.Context

	// basic auth
	Username string `json:"username" yaml:"username" redis:"username" structs:"username"`
	Password string `json:"password" yaml:"password" redis:"password" structs:"password"`
}

type EsOption func(o *Option)

func newDefaultOption() *Option {
	return &Option{
		Level:    DebugLevel,
		Receiver: ReceiverEs,
		// set daily suffix
		IndexSuffix: "2006-01-02",
		// retry times
		Retries: 3,
		Timeout: 10,
		ctx:     context.Background(),
		Urls:    []string{},
	}
}

func bindOptions(opt *Option, opts ...EsOption) {
	for _, f := range opts {
		f(opt)
	}
}

func WithUrls(urls ...string) EsOption {
	return func(o *Option) {
		o.Urls = append(o.Urls, urls...)
	}
}

func WithIndexName(indexName string) EsOption {
	return func(o *Option) {
		o.IndexName = indexName
	}
}

func WithReceiver(i int) EsOption {
	return func(o *Option) {
		o.Receiver = i
	}
}

func WithLevel(level int) EsOption {
	return func(o *Option) {
		o.Level = level
	}
}

func WithUsername(username string) EsOption {
	return func(o *Option) {
		o.Username = username
	}
}

func WithPassword(password string) EsOption {
	return func(o *Option) {
		o.Password = password
	}
}
