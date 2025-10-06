package store

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
	"github.com/valkey-io/valkey-go"
)

type StateStore interface {
	Store(ctx context.Context, state string, expiry time.Duration) error
	VerifyAndDelete(ctx context.Context, state string) (bool, error)
	Shutdown() error
}

type valkeyStateStore struct {
	logger ezutil.Logger
	client valkey.Client
}

func NewStateStore(logger ezutil.Logger, cfg config.Valkey) (StateStore, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		InitAddress: []string{cfg.Addr},
		Password:    cfg.Password,
		SelectDB:    cfg.Db,
	})
	if err != nil {
		return nil, eris.Wrap(err, "error initializing valkey client")
	}
	return &valkeyStateStore{
		logger,
		client,
	}, nil
}

func (vss *valkeyStateStore) Store(ctx context.Context, state string, expiry time.Duration) error {
	key := vss.constructKey(state)

	cmd := vss.client.B().Set().Key(key).Value("1").ExSeconds(int64(expiry.Seconds())).Build()
	if err := vss.client.Do(ctx, cmd).Error(); err != nil {
		return eris.Wrap(err, "error storing state in valkey")
	}

	return nil
}

func (vss *valkeyStateStore) VerifyAndDelete(ctx context.Context, state string) (bool, error) {
	key := vss.constructKey(state)

	getCmd := vss.client.B().Get().Key(key).Build()
	if err := vss.client.Do(ctx, getCmd).Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return false, nil
		}
		return false, eris.Wrap(err, "failed to get state in valkey")
	}

	delCmd := vss.client.B().Del().Key(key).Build()
	if err := vss.client.Do(ctx, delCmd).Error(); err != nil {
		err = eris.Wrap(err, "error deleting key from state store")
		vss.logger.Error(eris.ToString(err, true))
	}

	return true, nil
}

func (vss *valkeyStateStore) Shutdown() error {
	vss.client.Close()
	return nil
}

func (vss *valkeyStateStore) constructKey(state string) string {
	return fmt.Sprintf("state:%s", state)
}
