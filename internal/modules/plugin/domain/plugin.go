package domain

const TableName string = "plugin_configs"

type PluginConfig struct {
	ID       uint64 `gorm:"primaryKey;"`
	SlugName string `gorm:"column:slug_name;uniqueIndex;not null;size:128;check:LENGTH(slug_name) >= 4;"`
	Value    string `gorm:"column:value;type:text;"`
}

func (PluginConfig) TableName() string {
	return TableName
}
