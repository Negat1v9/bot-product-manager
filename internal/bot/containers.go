package telegram

import (
	"sync"
	"time"
)

type TypeUserCommand struct {
	TypeCmd  int
	ListID   *int
	GroupID  *int
	CreateAt time.Time
}

type Container struct {
	RecoverList map[int]bool
	UsersCmd    map[int64]TypeUserCommand
	mu          sync.Mutex
}

func NewContainer() *Container {
	c := &Container{
		RecoverList: make(map[int]bool),
		mu:          sync.Mutex{},
		UsersCmd:    make(map[int64]TypeUserCommand),
	}
	go c.startCheckOldData()
	return c
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

func (c *Container) AddUserCmd(userID int64, typeCmd TypeUserCommand) {
	typeCmd.CreateAt = time.Now()
	c.mu.Lock()
	c.UsersCmd[userID] = typeCmd
	c.mu.Unlock()
}

func (c *Container) DeleteUserCmd(userID int64) {
	c.mu.Lock()
	delete(c.UsersCmd, userID)
	c.mu.Unlock()
}

func (c *Container) startCheckOldData() {
	for {
		time.Sleep(time.Minute)
		for userID, t := range c.UsersCmd {
			if time.Now().Sub(t.CreateAt) > time.Minute {
				c.mu.Lock()
				delete(c.UsersCmd, userID)
				c.mu.Unlock()
			}
		}
	}
}
