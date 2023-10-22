package main

import (
	"log"
	"sync"
)

func NewManager() *Manager {
	return &Manager{}
}

type Manager struct {
	consumers sync.Map // key:value = { Owner }:{ []Consumer }
}

func (m *Manager) AddConsumerAndRun(owner string, consumers ...Consumer) {
	if len(consumers) == 0 {
		return
	}

	m.consumers.Store(owner, consumers)
	for _, consumer := range consumers {
		consumer := consumer
		go func() {
			consumer.RunConsume()
		}()
	}
}

func (m *Manager) StopConsumerByOwner(owner string) {
	value, exist := m.consumers.LoadAndDelete(owner)
	if !exist {
		return
	}

	for _, consumer := range value.([]Consumer) {
		err := consumer.Shutdown()
		if err != nil {
			log.Printf("consumer=%v: Shutdown: %v", consumer.Name, err)
		}
	}
	return
}

func (m *Manager) StopAll() {
	wg := sync.WaitGroup{}
	m.consumers.Range(
		func(key, value any) bool {
			for _, consumer := range value.([]Consumer) {
				consumer := consumer
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := consumer.Shutdown()
					if err != nil {
						log.Printf("consumer=%v: Shutdown: %v", consumer.Name, err)
					}
				}()
			}
			return true
		},
	)
	wg.Wait()
}
