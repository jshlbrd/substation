package transform

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brexhq/substation/v2/config"
	"github.com/brexhq/substation/v2/message"

	iconfig "github.com/brexhq/substation/v2/internal/config"
)

type objectCopyConfig struct {
	ID     string         `json:"id"`
	Object iconfig.Object `json:"object"`
}

func (c *objectCopyConfig) Decode(in interface{}) error {
	return iconfig.Decode(in, c)
}

func newObjectCopy(_ context.Context, cfg config.Config) (*objectCopy, error) {
	conf := objectCopyConfig{}
	if err := conf.Decode(cfg.Settings); err != nil {
		return nil, fmt.Errorf("transform object_copy: %v", err)
	}

	if conf.ID == "" {
		conf.ID = "object_copy"
	}

	tf := objectCopy{
		conf:            conf,
		hasObjectKey:    conf.Object.SourceKey != "" && conf.Object.TargetKey == "",
		hasObjectSetKey: conf.Object.SourceKey == "" && conf.Object.TargetKey != "",
	}

	return &tf, nil
}

type objectCopy struct {
	conf            objectCopyConfig
	hasObjectKey    bool
	hasObjectSetKey bool
}

func (tf *objectCopy) Transform(ctx context.Context, msg *message.Message) ([]*message.Message, error) {
	if msg.HasFlag(message.IsControl) {
		return []*message.Message{msg}, nil
	}

	if tf.hasObjectKey {
		value := msg.GetValue(tf.conf.Object.SourceKey)
		if skipMessage(msg, value) {
			return []*message.Message{msg}, nil
		}

		msg.SetData(value.Bytes())
		return []*message.Message{msg}, nil
	}

	if tf.hasObjectSetKey {
		if len(msg.Data()) == 0 {
			return []*message.Message{msg}, nil
		}

		outMsg := message.New().SetMetadata(msg.Metadata())
		if err := outMsg.SetValue(tf.conf.Object.TargetKey, msg.Data()); err != nil {
			return nil, fmt.Errorf("transform %s: %v", tf.conf.ID, err)
		}

		return []*message.Message{outMsg}, nil
	}

	value := msg.GetValue(tf.conf.Object.SourceKey)
	if skipMessage(msg, value) {
		return []*message.Message{msg}, nil
	}

	if err := msg.SetValue(tf.conf.Object.TargetKey, value); err != nil {
		return nil, fmt.Errorf("transform %s: %v", tf.conf.ID, err)
	}

	return []*message.Message{msg}, nil
}

func (tf *objectCopy) String() string {
	b, _ := json.Marshal(tf.conf)
	return string(b)
}
