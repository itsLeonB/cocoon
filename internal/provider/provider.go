package provider

import (
	"errors"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/store"
	"github.com/itsLeonB/ezutil/v2"
)

type Provider struct {
	Logger ezutil.Logger
	*DBs
	*Repositories
	*Services
	Store store.StateStore
}

func All(logger ezutil.Logger, configs config.Config) (*Provider, error) {
	dbs := ProvideDBs(configs.DB)
	repos := ProvideRepositories(dbs.GormDB)

	store, err := store.NewStateStore(logger, configs.Valkey)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Logger:       logger,
		DBs:          dbs,
		Repositories: repos,
		Services:     ProvideServices(configs, repos, logger, store),
		Store:        store,
	}, nil
}

func (p *Provider) Shutdown() error {
	var err error
	if p.DBs != nil {
		if e := p.DBs.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
	}
	if p.Store != nil {
		if e := p.Store.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}
