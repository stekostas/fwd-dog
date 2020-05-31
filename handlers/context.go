package handlers

import (
	"github.com/stekostas/fwd-dog/cache"
	"time"
)

type Context struct {
	Renderer     *Renderer
	CacheAdapter cache.Adapter
	TtlOptions   map[time.Duration]string
}

func NewContext(renderer *Renderer, adapter cache.Adapter, ttlOptions map[time.Duration]string) *Context {
	return &Context{Renderer: renderer, CacheAdapter: adapter, TtlOptions: ttlOptions}
}
