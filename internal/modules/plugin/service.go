package plugin

import (
	"context"

	"github.com/codenomdev/viona/internal/modules/plugin/constant"
	settingSvc "github.com/codenomdev/viona/internal/modules/setting/service"
	"github.com/codenomdev/viona/pkg/log"
	"github.com/codenomdev/viona/plugins"
	"go.uber.org/zap"
)

type (
	Service interface{}
	service struct {
		log            log.Logger
		settingService settingSvc.Service
		repo           Repository
	}
)

func NewService(
	log log.Logger,
	repo Repository,
) Service {
	s := &service{
		log:  log,
		repo: repo,
	}

	s.initPluginData()
	return s
}

func (s *service) initPluginData() {
	// init plugin status
	val, err := s.settingService.GetValueByKey(context.Background(), constant.PluginStatus)

	pluginStatus, ok := val.([]byte)
	if !ok {
		return
	}
	if err != nil {
		s.log.Error("init plugin status", zap.Error(err))
	} else {
		if err := plugins.StatusManager.UnmarshalJSON([]byte(pluginStatus)); err != nil {
			s.log.Error("init plugin status", zap.Error(err))
		}
	}

	// init plugin config
	pluginConfigs, err := s.repo.GetAll(context.Background(), nil)

	if err != nil {
		s.log.Error("init plugin config", zap.Error(err))
	} else {
		for _, pluginConfig := range pluginConfigs {
			err := plugins.CallConfig(func(fn plugins.Config) error {
				if fn.Info().SlugName == pluginConfig.SlugName {
					return fn.ConfigReceiver([]byte(pluginConfig.Value))
				}
				return nil
			})
			if err != nil {
				s.log.Error("parse plugin config failed: ", zap.String("pluginConfig.PluginSlugName", pluginConfig.SlugName), zap.Error(err))
			}
		}

		// _ = plugins.CallCache(func(cache plugins.Cache) error {
		// 	ps.data.Cache = cache
		// 	return nil
		// })
	}
}
