package domain

import (
	"time"

	"github.com/codenomdev/viona/internal/modules/setting/constant"
	"gorm.io/datatypes"
)

const TableName string = "site_settings"

type SiteSetting struct {
	ID        uint64                `gorm:"primaryKey;"`
	Group     constant.SettingGroup `gorm:"column:group_name;size:64;not null;uniqueIndex:idx_group_key;check:length(group_name) >= 3"`
	Key       string                `gorm:"column:key;size:100;not null;uniqueIndex:idx_group_key;check:length(key) >= 4"`
	Values    datatypes.JSON        `gorm:"column:values;type:jsonb;default:null;"`
	Sort      int                   `gorm:"column:sort;not null;default:0;index;"`
	Status    int8                  `gorm:"column:status;not null;default:1;index;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SiteSetting) TableName() string {
	return TableName
}
