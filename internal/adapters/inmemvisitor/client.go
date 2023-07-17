package inmemvisitor

import (
	"demoscraper/internal/core"
	"demoscraper/internal/core/entities"
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

func (r *Client) ToVisitMap() entities.VisitMap {
	cp := make(entities.VisitMap, len(r.visitMap))

	r.mutex.Lock()
	for k, v := range r.visitMap {
		cp[k] = v
	}
	r.mutex.Unlock()

	return cp
}

func (r *Client) Merge(visitMap entities.VisitMap) {
	r.mutex.Lock()
	for k, v := range visitMap {
		r.visitMap[k] = v
	}
	r.mutex.Unlock()
}
