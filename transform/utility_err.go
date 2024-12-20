package transform

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brexhq/substation/v2/config"
	"github.com/brexhq/substation/v2/message"

	iconfig "github.com/brexhq/substation/v2/internal/config"
)

type utilityErrConfig struct {
	// Message is the error message to return.
	Message string `json:"message"`

	ID string `json:"id"`
}

func (c *utilityErrConfig) Decode(in interface{}) error {
	return iconfig.Decode(in, c)
}

func newUtilityErr(_ context.Context, cfg config.Config) (*utilityErr, error) {
	conf := utilityErrConfig{}
	if err := conf.Decode(cfg.Settings); err != nil {
		return nil, fmt.Errorf("transform utility_err: %v", err)
	}

	if conf.ID == "" {
		conf.ID = "utility_err"
	}

	tf := utilityErr{
		conf: conf,
	}

	return &tf, nil
}

type utilityErr struct {
	conf utilityErrConfig
}

func (tf *utilityErr) Transform(_ context.Context, msg *message.Message) ([]*message.Message, error) {
	if msg.HasFlag(message.IsControl) {
		return []*message.Message{msg}, nil
	}

	return []*message.Message{msg}, fmt.Errorf("%s", tf.conf.Message)
}

func (tf *utilityErr) String() string {
	b, _ := json.Marshal(tf.conf)
	return string(b)
}
