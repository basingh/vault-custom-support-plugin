package main

import (
	"context"
	//"errors"
	//"reflect"
	"testing"

	//"github.com/hashicorp/vault/sdk/framework"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

// create a backend for testing
func getBackend(t *testing.T) (logical.Backend, logical.Storage) {
	b, storage, _ := getBackendWithEvents(t)
	return b, storage
}

type mockEventsSender struct {
	eventsProcessed []*logical.EventReceived
}

func (m *mockEventsSender) Send(ctx context.Context, eventType logical.EventType, event *logical.EventData) error {
	if m == nil {
		return nil
	}
	m.eventsProcessed = append(m.eventsProcessed, &logical.EventReceived{
		EventType: string(eventType),
		Event:     event,
	})

	return nil
}

func getBackendWithEvents(t *testing.T) (logical.Backend, logical.Storage, *mockEventsSender) {
	events := &mockEventsSender{}

	config := &logical.BackendConfig{
		Logger:       logging.NewVaultLogger(log.Trace),
		System:       &logical.StaticSystemView{},
		StorageView:  &logical.InmemStorage{},
		BackendUUID:  "test",
		EventsSender: events,
	}

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("unable to create backend: %v", err)
	}

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   config.StorageView,
	}
	req.ClientToken = "root"

	resp, err := b.HandleRequest(context.Background(), req)
	// if err != nil {
	// 	t.Fatalf("unable to read config: %s", err.Error())
	// 	return nil, nil, nil
	// }

	if resp == nil || resp.IsError() {
		//	t.Fatalf("Error during mount creation: %x", resp.Error().Error())
	}

	return b, config.StorageView, events

}

func TestHandleRead(t *testing.T) {
	b, storage := getBackend(t)

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"bar": "baz",
		},
	}

	requ := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "data/foo",
		Storage:   storage,
		Data:      data,
	}
	requ.ClientToken = "root"

	respo, err := b.HandleRequest(context.Background(), requ)
	if err != nil || (respo != nil && respo.IsError()) {
		t.Fatalf("data CreateOperation request failed, err: %s, resp %#v", err, respo)
	}

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "data/foo",
		Storage:   storage,
	}
	req.ClientToken = "root"
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("data ReadOperation request failed, err: %s, resp %#v", err, resp)
	}

}
