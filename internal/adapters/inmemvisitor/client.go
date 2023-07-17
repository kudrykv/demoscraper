package inmemvisitor

import (
	"demoscraper/internal/core"
	"sync"
)

type Client struct {
	mutex    *sync.Mutex
	visitMap map[string]struct{}
}

func New() core.Visitor {
	return &Client{
		mutex:    &sync.Mutex{},
		visitMap: make(map[string]struct{}),
	}
}

func (r *Client) Visit(some string) {
	r.mutex.Lock()
	r.visitMap[some] = struct{}{}
	r.mutex.Unlock()
}

func (r *Client) IsVisited(some string) bool {
	r.mutex.Lock()
	_, ok := r.visitMap[some]
	r.mutex.Unlock()

	return ok
}
