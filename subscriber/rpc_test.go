package subscriber

import (
	"testing"
	"time"
)

type Parser struct{}

func (parser Parser) ParseResponse(data []byte) ([]Event, bool) {
	return []Event{data}, true
}

func TestRpcSubscriber_SubscribeToEvents(t *testing.T) {
	t.Run("subscribes to rpc endpoint", func(t *testing.T) {
		u := *rpcMockUrl
		u.Path = "/test/1"
		rpc := RpcSubscriber{Endpoint: u.String(), Parser: Parser{}, Interval: 1 * time.Second}

		events := make(chan Event)
		filter := TestsMockFilter{true}

		sub, err := rpc.SubscribeToEvents(events, filter)
		if err != nil {
			t.Errorf("SubscribeToEvents() error = %v", err)
			return
		}
		defer sub.Unsubscribe()

		event := <-events
		mockevent := string(event)
		if mockevent != "1" {
			t.Errorf("SubscribeToEvents() got unexpected first message = %v", mockevent)
			return
		}
		event = <-events
		mockevent = string(event)
		if mockevent != "2" {
			t.Errorf("SubscribeToEvents() got unexpected second message = %v", mockevent)
			return
		}
		return
	})
}

func TestSendGetRequest(t *testing.T) {
	t.Run("succeeds on normal response", func(t *testing.T) {
		u := *rpcMockUrl
		u.Path = "/test/2"

		_, err := sendGetRequest(u.String())
		if err != nil {
			t.Errorf("sendGetRequest() got unexpected error = %v", err)
			return
		}
	})

	t.Run("fails on bad status", func(t *testing.T) {
		u := *rpcMockUrl
		u.Path = "/fails"

		_, err := sendGetRequest(u.String())
		if err == nil {
			t.Error("sendGetRequest() expected error, but got nil")
			return
		}
	})
}

func TestRpcSubscriber_Test(t *testing.T) {
	type fields struct {
		Endpoint string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"succeeds connecting to valid endpoint",
			fields{Endpoint: rpcMockUrl.String()},
			false,
		},
		{
			"fails connecting to invalid endpoint",
			fields{Endpoint: "http://localhost:9999/invalid"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpc := RpcSubscriber{
				Endpoint: tt.fields.Endpoint,
			}
			if err := rpc.Test(); (err != nil) != tt.wantErr {
				t.Errorf("Test() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
