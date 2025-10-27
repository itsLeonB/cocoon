package provider

import (
	"errors"
	"net/http"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/store"
	"github.com/itsLeonB/ezutil/v2"
)

type Provider struct {
	Logger ezutil.Logger
	*DBs
	*Repositories
	*Services
	Store      store.StateStore
	httpClient *http.Client
}

func All(logger ezutil.Logger, configs config.Config) (*Provider, error) {
	dbs, err := ProvideDBs(configs.DB)
	if err != nil {
		return nil, err
	}

	repos := ProvideRepositories(dbs.GormDB)
	httpClient := configs.NewClient()

	store, err := store.NewStateStore(logger, configs.Valkey)
	if err != nil {
		if e := dbs.Shutdown(); e != nil {
			logger.Errorf("error cleaning up DB resources: %v", e)
		}
		return nil, err
	}

	services, err := ProvideServices(configs, repos, logger, store, httpClient)
	if err != nil {
		if e := dbs.Shutdown(); e != nil {
			logger.Errorf("error cleaning up DB resources: %v", e)
		}
		if e := store.Shutdown(); e != nil {
			logger.Errorf("error cleaning up store resources: %v", e)
		}
		return nil, err
	}

	return &Provider{
		Logger:       logger,
		DBs:          dbs,
		Repositories: repos,
		Services:     services,
		Store:        store,
		httpClient:   httpClient,
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
	if p.httpClient != nil {
		p.httpClient.CloseIdleConnections()
	}
	return err
}
