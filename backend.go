package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	//"github.com/aerospike/aerospike-client-go/v5/logger"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// factory function to setup a backend in vaultw
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
	// calling https://pkg.go.dev/sync#RWMutex.
	lock sync.RWMutex
	//client *mordorClient
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
		Paths: framework.PathAppend(
			b.paths(),
		),
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
	//b.client = nil
}

// invalidate method which call reset to reset configuration
func (b *mordorBackend) invalidate(ctx context.Context, key string) {
	if key == "config" {
		b.reset()
	}
}

// Based on https://github.com/hashicorp/vault/blob/main/sdk/framework/path.go we will not setup operations on this path
// read, write and delete

func (b *mordorBackend) paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: framework.MatchAllRegex("path"),

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Defining path of secret",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRead,
					Summary:  "Read secrets",
				},
				// error due to https://github.com/hashicorp/vault/blob/main/sdk/framework/backend.go#L111 is looking for FieldData
				// data coming from request, need to know how to pass these hardcoded data to this func
				// there are some ideas here to push it directly physical storage https://github.com/hashicorp/vault/blob/main/sdk/logical/storage_inmem.go#L37-L45
				// using https://github.com/hashicorp/vault/blob/main/sdk/physical/entry.go
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Update secret on path",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Write secret on path",
				},
				//	logical.DeleteOperation: &framework.PathOperation{
				//		Callback: b.handleDelete,
				//		Summary:  "Delete secret from path",
				//	},
			},
			//ExistenceCheck: b.HandleExistenceCheck(),
		},
	}
}

// handleRead operation to read data from path

func (b *mordorBackend) handleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("Client token is empty")
	}

	path := data.Get("path").(string)
	var rawData map[string]interface{}
	entry, err := req.Storage.Get(ctx, req.ClientToken+"/"+path)

	if err != nil {
		return nil, err
	} else if entry == nil {
		err := logical.ErrorResponse("No value found")
		return err, nil
	}

	fetchedData := entry.Value
	if fetchedData == nil {
		resp := logical.ErrorResponse("No value found at %v%v", req.MountPoint, path)
		return resp, nil
	}

	if err := jsonutil.DecodeJSON(fetchedData, &rawData); err != nil {
		return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	}

	resp := &logical.Response{
		Data: rawData,
	}

	return resp, nil
}

// handleWrite operation to write on the path

func (b *mordorBackend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// store kv pair in required path

	if len(req.Data) == 0 {
		return nil, fmt.Errorf("Data cannot be nil")
	}
	path := data.Get("path").(string)

	buf, err := json.Marshal(req.Data)
	if err != nil {
		return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	}

	entry := &logical.StorageEntry{
		Key:      req.ClientToken + "/" + path,
		Value:    buf,
		SealWrap: false,
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if response, err := req.Storage.Get(ctx, entry.Key); err != nil {
		return nil, err
	} else {
		logger := hclog.New(&hclog.LoggerOptions{})
		logger.Info(string(response.Value))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			entry.Key: entry,
		},
	}
	return resp, nil
}

const backendHelp = `
		The Mordor secrets backend dynamically generates user tokens.
		After mounting this backend, credentials to manage Mordor user tokens
		must be configured with the "config/" endpoints.
		`
