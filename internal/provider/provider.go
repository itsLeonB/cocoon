package provider

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/ezutil/v2"
)

type Provider struct {
	Logger ezutil.Logger
	*DBs
	*Repositories
	*Services
}

func All(logger ezutil.Logger, configs config.Config) *Provider {
	dbs := ProvideDBs(configs.DB)
	repos := ProvideRepositories(dbs.GormDB)

	return &Provider{
		Logger:       logger,
		DBs:          dbs,
		Repositories: repos,
		Services:     ProvideServices(configs, repos, logger),
	}
}
