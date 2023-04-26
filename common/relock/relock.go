package relock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Mutex struct {
	m        sync.Mutex
	client   *redis.Client
	cancel   func()
	name     string
	expire   time.Duration
	interval time.Duration
}

func (this *Mutex) Lock() error {
	return this.LockContext(context.Background())
}

func (this *Mutex) LockContext(ctx context.Context) error {
	ctx2, cancel := context.WithCancel(ctx)
	for {
		result, err := this.client.SetNX(ctx2, this.name, "", this.expire).Result()
		if err != nil {
			cancel()
			return err
		}
		if result {
			break
		}
	}
	this.cancel = cancel
	go func() {
		t := time.NewTicker(this.interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				result, err := this.client.Expire(ctx2, this.name, this.expire).Result()
				if err != nil || !result {
					this.cancel()
					return
				}
			case <-ctx2.Done():
				return
			}
		}
	}()
	return nil
}

func (this *Mutex) Unlock() {
	this.client.Del(context.Background(), this.name)
	this.cancel()
}

func NewMutex(client *redis.Client, name string) *Mutex {
	m := &Mutex{
		name:     name,
		expire:   5 * time.Second,
		interval: 2 * time.Second,
		client:   client,
	}
	return m
}
