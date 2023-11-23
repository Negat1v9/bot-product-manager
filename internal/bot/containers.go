package telegram

import "sync"

type Container struct {
	RecoverList map[int]bool
	mu          sync.Mutex
}

func NewContainer() *Container {
	return &Container{
		RecoverList: make(map[int]bool),
		mu:          sync.Mutex{},
	}
}

func (c *Container) SetRecoverList(listID int) {
	c.mu.Lock()
	c.RecoverList[listID] = true
	c.mu.Unlock()
}
func (c *Container) isInContainerList(listID int) bool {
	if _, ok := c.RecoverList[listID]; ok {
		return true
	}
	return false
}
func (c *Container) DeleteRecoverList(listID int) {
	c.mu.Lock()
	delete(c.RecoverList, listID)
	c.mu.Unlock()
}
