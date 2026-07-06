package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codenomdev/viona/internal/modules/setting/constant"
	"github.com/codenomdev/viona/internal/modules/setting/dto"
	"github.com/codenomdev/viona/internal/modules/setting/repository"
	"github.com/codenomdev/viona/pkg/response"
	"gorm.io/gorm"
)

type (
	Service interface {
		GetValueByKey(ctx context.Context, key string) (any, error)
		GetSettingsPerGroup(ctx context.Context) (dto.SettingsPerGroupResponse, error)
		GetAllBySystemGroup(ctx context.Context) (dto.SettingPerSystemGroupResponse, error)
	}
	service struct {
		db          *gorm.DB
		settingRepo repository.Repository
	}
)

func NewService(
	db *gorm.DB,
	settingRepo repository.Repository,
) Service {
	return &service{
		db:          db,
		settingRepo: settingRepo,
	}
}

func (s *service) GetValueByKey(
	ctx context.Context,
	key string,
) (any, error) {
	tx := s.db.WithContext(ctx)

	result, err := s.settingRepo.GetByKey(ctx, tx, key)
	if err != nil {
		return nil, response.NewHttpNotFound(
			[]string{
				fmt.Sprintf("get setting key: %s not found", key),
			},
			nil,
		)
	}

	return decodeJSON(result.Values), nil
}

// Get all setting with group system
// Only status == 1
func (s *service) GetAllBySystemGroup(ctx context.Context) (dto.SettingPerSystemGroupResponse, error) {
	tx := s.db.WithContext(ctx)
	systems, err := s.settingRepo.GetByGroup(ctx, tx, string(constant.GroupSystem))

	if err != nil {
		return nil, response.NewHttpUnprocessedEntity([]string{"failed to get setting by system"})
	}

	result := make(dto.SettingPerSystemGroupResponse)

	for _, item := range systems {
		result[item.Key] = decodeJSON(item.Values)
	}

	return result, nil
}

func (s *service) GetSettingsPerGroup(
	ctx context.Context,
) (dto.SettingsPerGroupResponse, error) {
	tx := s.db.WithContext(ctx)

	settings, err := s.settingRepo.GetAll(ctx, tx)
	if err != nil {
		return nil, response.NewHttpUnprocessedEntity(
			[]string{"failed to get setting"},
		)
	}

	result := make(dto.SettingsPerGroupResponse)

	for _, item := range settings {
		if _, exists := result[item.Group]; !exists {
			result[item.Group] = make(map[string]any)
		}

		result[item.Group][item.Key] = decodeJSON(item.Values)
	}

	return result, nil
}

func decodeJSON(data []byte) any {
	if len(data) == 0 {
		return nil
	}

	var value any

	if err := json.Unmarshal(data, &value); err != nil {
		return nil
	}

	return value
}
