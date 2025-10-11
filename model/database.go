package model

import (
	"github.com/pocketbase/pocketbase/core"
)

var (
	_ core.RecordProxy = (*Config)(nil)
)

const (
	DbNameConfigs     = "configs"
	ConfigsFieldKey   = "key"
	ConfigsFieldValue = "value"
)

type Config struct {
	core.BaseRecordProxy
}

func NewConfig(record *core.Record) *Config {
	config := new(Config)
	config.SetProxyRecord(record)
	return config
}

func NewConfigFromCollection(collection *core.Collection) *Config {
	record := core.NewRecord(collection)
	return NewConfig(record)
}

func (config *Config) Key() ConfigKey {
	return MustParseConfigKey(config.GetString("key"))
}

func (config *Config) SetKey(value ConfigKey) {
	config.Set("key", value)
}

func (config *Config) Value() string {
	return config.GetString("value")
}

func (config *Config) SetValue(value string) {
	config.Set("value", value)
}
