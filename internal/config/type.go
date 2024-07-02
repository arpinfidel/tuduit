package config

import "sync"

type Configuration struct {
	Config

	mut *sync.RWMutex
}

type Config struct {
	OpenAI OpenAIConfig `json:"open_ai"`
}

func (c *Configuration) Get() Config {
	if c.mut == nil {
		c.mut = &sync.RWMutex{}
	}
	c.mut.RLock()
	cp := c.Config
	c.mut.RUnlock()
	return cp
}

func (c *Configuration) Set(nc Config) {
	if c.mut == nil {
		c.mut = &sync.RWMutex{}
	}
	c.mut.Lock()
	c.Config = nc
	c.mut.Unlock()
}

func (c *Configuration) Lock() (unlock func()) {
	c.mut.Lock()
	return c.mut.Unlock
}

type OpenAIConfig struct {
	APIKey string `json:"api_key"`
}
