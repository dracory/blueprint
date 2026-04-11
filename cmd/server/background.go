package main

import (
	"context"
	"sync"
)

type backgroundGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	doneCh chan struct{}
	once   sync.Once
}

func newBackgroundGroup(parent context.Context) *backgroundGroup {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &backgroundGroup{ctx: ctx, cancel: cancel, doneCh: make(chan struct{})}
}

func (g *backgroundGroup) Go(fn func(ctx context.Context)) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		fn(g.ctx)
	}()
}

func (g *backgroundGroup) stop() {
	g.once.Do(func() {
		g.cancel()
		g.wg.Wait()
		close(g.doneCh)
	})
}

func (g *backgroundGroup) Done() <-chan struct{} {
	return g.doneCh
}
