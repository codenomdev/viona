package dto

import "github.com/codenomdev/viona/internal/modules/setting/constant"

type SettingsPerGroupResponse map[constant.SettingGroup]map[string]any
