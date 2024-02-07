package store

import (
	"context"
	"errors"
	"sync"

	"github.com/767829413/advanced-go/open-platform/models"
)

// NewClientStore create client store
func NewClientStore() *ClientStoreIns {
	return &ClientStoreIns{
		data: make(map[string]models.ClientInfo),
	}
}

// ClientStore client information store
type ClientStoreIns struct {
	sync.RWMutex
	data map[string]models.ClientInfo
}

// GetByID according to the ID for the client information
func (cs *ClientStoreIns) GetByID(ctx context.Context, id string) (models.ClientInfo, error) {
	cs.RLock()
	defer cs.RUnlock()

	if c, ok := cs.data[id]; ok {
		return c, nil
	}
	return nil, errors.New("not found")
}

// Set set client information
func (cs *ClientStoreIns) Set(id string, cli models.ClientInfo) (err error) {
	cs.Lock()
	defer cs.Unlock()

	cs.data[id] = cli
	return
}
