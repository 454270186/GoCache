package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu      sync.Mutex
	callMap map[string]*call // <key, that call for processing this key>
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.callMap == nil {
		g.callMap = make(map[string]*call)
	}
	if c, ok := g.callMap[key]; ok {
		// if there is alreay a process for the key, wait for the result
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.callMap[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.callMap, key)
	g.mu.Unlock()

	return c.val, c.err
}