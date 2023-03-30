package webmethods

import (
	coreapi "github.com/Axway/agent-sdk/pkg/api"
)

// MockClient is the mock client
type MockClient struct {
	SendFunc func(request coreapi.Request) (*coreapi.Response, error)
}

// Do is the mock client's `Do` func
func (m *MockClient) Send(req coreapi.Request) (*coreapi.Response, error) {
	return m.SendFunc(req)
}
