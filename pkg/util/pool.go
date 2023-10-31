package util

import (
	"sync"
)

func NewPool() *Pool {
	return &Pool{}
}

type Pool struct {
	pool sync.Map
}

func (p *Pool) Set(key string, val interface{}) {
	p.pool.Store(key, val)
}

func (p *Pool) Get(key string) (interface{}, bool) {
	return p.pool.Load(key)
}
