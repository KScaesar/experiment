package main

import (
	"sync"
)

type Resource interface {
	IsInUse() bool
	SetInUse(inUse bool)
	GetSubPool() ResourcePool
	UsageRate() int
}

type ResourcePool interface {
	GetMinUsageRateResource() Resource
	ReleaseResource(Resource)
}

func NewBasicResource(id int, subPool ResourcePool) *BasicResource {
	return &BasicResource{
		id:        id,
		subPool:   subPool,
		usageRate: 0,
	}
}

type BasicResource struct {
	id        int
	subPool   ResourcePool
	usageRate int
}

func (c *BasicResource) IsInUse() bool {
	return c.usageRate == 0
}

func (c *BasicResource) SetInUse(inUse bool) {
	if inUse {
		c.usageRate++
	} else {
		c.usageRate--
	}
}

func (c *BasicResource) GetSubPool() ResourcePool {
	return c.subPool
}

func (c *BasicResource) UsageRate() int {
	return c.usageRate
}

type BasicResourcePool struct {
	resources []Resource
	mu        sync.Mutex
}

func (pool *BasicResourcePool) GetMinUsageRateResource() Resource {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	var minUsageRateResource Resource

	for _, resource := range pool.resources {
		if minUsageRateResource == nil || resource.UsageRate() < minUsageRateResource.UsageRate() {
			minUsageRateResource = resource
		}
	}

	if minUsageRateResource != nil {
		minUsageRateResource.SetInUse(true)
	}

	return minUsageRateResource
}

func (pool *BasicResourcePool) ReleaseResource(connection Resource) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if connection != nil {
		connection.SetInUse(false)
	}
}
