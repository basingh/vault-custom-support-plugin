package vaultcustomsupportplugin

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// factory function to setup a backend in vault
// this is mandatory function and should include call to setup backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// create a backend object to secret engine

type mordorBackend struct {
	*framework.Backend
	// check what RWMutex do?
	// calling https://pkg.go.dev/sync#RWMutex.
	lock   sync.RWMutex
	client *mordorClient
}

// backend function which returns mordorBackend object

func backend() *mordorBackend {
	var b = mordorBackend{}

	// backend is an implementation of logical.Backend in Vault SDK
	// this provides framework to handle routing and validation
	// https://github.com/hashicorp/vault/blob/main/sdk/framework/backend.go
	b.Backend = &framework.Backend{

		Help: strings.TrimSpace(backendHelp),
		// pathAppend() is the helper function in Vault SDK to handle list of paths into same list
		Paths: framework.pathAppend(),
		// this is setup a skeleton structure of secret
		// https://github.com/hashicorp/vault/blob/main/sdk/framework/secret.go#L11-L40
		Secrets: []*framework.Secret{},
		// this is just setting backend type as secret
		// https://github.com/hashicorp/vault/blob/main/sdk/logical/logical.go#L21-L30
		BackendType: logical.TypeLogical,
		Invalidate:  b.invalidate,
	}

	return &b

}

// reset method to lock the backend while target API client object is reset
func (b *mordorBackend) reset() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.client = nil
}

// invalidate method which call reset to reset configuration
func (b *mordorBackend) invalidate(ctx context.Context, key string) {
	if key == "config" {
		b.reset()
	}
}

const backendHelp = `
		The Mordor secrets backend dynamically generates user tokens.
		After mounting this backend, credentials to manage Mordor user tokens
		must be configured with the "config/" endpoints.
		`
