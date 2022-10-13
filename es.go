package xes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	ErrIndexRequired   = errors.New("index name is required")
	ErrServersRequired = errors.New("server urls required")
	ErrServerIsEmpty   = errors.New("server url is empty")
)

// EsLogger an elastic writer can be used for zerolog
type EsLogger struct {
	option *Option
	client *elastic.Client

	discard bool
}

func NewEsLogger(esOpts ...EsOption) (*EsLogger, error) {
	opt := newDefaultOption()
	bindOptions(opt, esOpts...)
	if opt.IndexSuffix != "" {
		today := time.Now().Format(opt.IndexSuffix)
		opt.IndexName = fmt.Sprintf("%s-%s", opt.IndexName, today)
	}

	el := &EsLogger{option: opt}

	if opt.Receiver != ReceiverEs {
		return el, nil
	}

	if opt.IndexName == "" {
		return nil, ErrIndexRequired
	}

	if len(opt.Urls) == 0 {
		return nil, ErrServersRequired
	}
	for _, v := range opt.Urls {
		if v == "" {
			return nil, ErrServerIsEmpty
		}
	}

	httpClient := &http.Client{Timeout: time.Duration(opt.Timeout) * 1000 * time.Millisecond}
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(opt.Urls...),
		elastic.SetHttpClient(httpClient),
		elastic.SetRetrier(newMyRetrier(opt.Retries)),
		elastic.SetBasicAuth(opt.Username, opt.Password),
	)

	if err != nil {
		return nil, err
	}

	el.client = client
	return el, nil
}

// Write implements io.Writer, so we can use EsLogger as io.Writer
func (l *EsLogger) Write(p []byte) (n int, err error) {
	str := cast.ToString(p)
	err = l.write(str)
	return len(p), err
}

func (l *EsLogger) write(body interface{}) error {
	body = l.purify(body)

	if l.discard {
		return nil
	}

	_, err := l.client.Index().
		Index(l.option.IndexName).BodyJson(body).
		Do(l.option.ctx)

	if err != nil {
		dump.P(body)
		fmt.Fprintf(os.Stderr, "add log failed with error: %v\n", err)
		return err
	}

	_, err = l.client.Flush().Index(l.option.IndexName).Do(l.option.ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "flush log failed with error: %v\n", err)
		return err
	}

	return err
}

func (l *EsLogger) purify(raw interface{}) interface{} {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(raw.(string)), &dat); err != nil {
		return raw
	}

	if level, b := dat["level"]; b {
		l.checkLevel(level.(string))
	}

	if b, e := json.Marshal(dat); e == nil {
		return string(b)
	}
	return raw
}

func (l *EsLogger) checkLevel(level string) {
	li := stringLevel(level)
	l.discard = li < l.option.Level
}
