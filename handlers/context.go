package handlers

import (
	"github.com/stekostas/fwd-dog/cache"
	"time"
)

type Context struct {
	Renderer     *Renderer
	CacheAdapter cache.Adapter
	TtlOptions   []time.Duration
}

func NewContext(renderer *Renderer, adapter cache.Adapter, ttlOptions []time.Duration) *Context {
	return &Context{Renderer: renderer, CacheAdapter: adapter, TtlOptions: ttlOptions}
}
