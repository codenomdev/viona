package service

import (
	"context"
	"encoding/json"

	"github.com/codenomdev/viona/internal/modules/setting/dto"
	"github.com/codenomdev/viona/internal/modules/setting/repository"
	"github.com/codenomdev/viona/pkg/response"
	"gorm.io/gorm"
)

type (
	Service interface {
		GetSettingsPerGroup(ctx context.Context) (*dto.SettingsPerGroupResponse, error)
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

func (s *service) GetSettingsPerGroup(
	ctx context.Context,
) (*dto.SettingsPerGroupResponse, error) {
	tx := s.db.WithContext(ctx)

	settings, err := s.settingRepo.GetAll(ctx, tx)
	if err != nil {
		return nil, response.NewHttpUnprocessedEntity(
			[]string{"unprocessed entity"},
		)
	}

	result := make(dto.SettingsPerGroupResponse)

	for _, item := range *settings {
		if _, exists := result[item.Group]; !exists {
			result[item.Group] = make(map[string]any)
		}

		result[item.Group][item.Key] = decodeJSON(item.Values)
	}

	return &result, nil
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
