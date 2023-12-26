package redistream

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type ConsumerOptions struct {
	ConsumerName    string        // default: os.Hostname
	GroupName       string        // default: redistream
	PendingTimeout  time.Duration // default: 4 seconds
	BlockTimeout    time.Duration // XRead timeout, default: 1 seconds
	ReclaimInterval time.Duration // default: 3 seconds
	// concurrency workers, default: 0
	// 0: temporary worker for every single message
	// >=1: fix workers that can limit consume speed
	Workers      int
	BufferSize   int // default: 16 (too big will cause pending timeout)
	MaxRetry     int // default: 10 times
	ErrorHandler func(err error)
}

type Consumer struct {
	r        redis.UniversalClient
	handlers map[string]MessageHandler
	streams  []string
	lock     sync.Mutex
	options  *ConsumerOptions
	queue    chan *Message
	rest     chan struct{}
	cancel   func()
	wg       sync.WaitGroup
}

func NewConsumer(r redis.UniversalClient, options ...*ConsumerOptions) *Consumer {
	if options == nil {
		options = append(options, &ConsumerOptions{})
	}
	c := &Consumer{
		r:       r,
		options: options[0],
	}
	if c.options.ConsumerName == "" {
		c.options.ConsumerName, _ = os.Hostname()
	}
	if c.options.GroupName == "" {
		c.options.GroupName = "redistream"
	}
	if c.options.PendingTimeout == 0 {
		c.options.PendingTimeout = 4 * time.Second
	}
	if c.options.BlockTimeout == 0 {
		c.options.BlockTimeout = 1 * time.Second
	}
	if c.options.ReclaimInterval == 0 {
		c.options.ReclaimInterval = 3 * time.Second
	}
	if c.options.Workers < 0 {
		panic("Redistream workers cant be negative number")
	}
	if c.options.MaxRetry <= 0 {
		c.options.MaxRetry = 10
	}
	if c.options.BufferSize < 0 {
		panic("Redistream buffer size cant be negative number")
	} else if c.options.BufferSize == 0 {
		c.options.BufferSize = 16
	}
	if c.options.ErrorHandler == nil {
		c.options.ErrorHandler = func(err error) {
			fmt.Println(err)
		}
	}
	c.handlers = make(map[string]MessageHandler)
	c.queue = make(chan *Message, c.options.BufferSize)
	c.rest = make(chan struct{}, c.options.BufferSize)
	return c
}

func (this *Consumer) Register(stream string, handler MessageHandler) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.handlers[stream] = handler
	this.streams = append(this.streams, stream)
	this.r.XGroupCreateMkStream(context.Background(), stream, this.options.GroupName, "0")
}

func (this *Consumer) Serve() {
	this.lock.Lock()
	defer this.lock.Unlock()
	var ctx context.Context
	ctx, this.cancel = context.WithCancel(context.Background())
	if this.options.Workers == 0 {
		this.wg.Add(1)
		go this.infiniteWork(ctx)
	} else {
		this.wg.Add(this.options.Workers)
		for i := 0; i < this.options.Workers; i++ {
			go this.work(ctx)
		}
	}

	for _, stream := range this.streams {
		go this.fetch(ctx, stream)
	}
	go this.reclaim(ctx)

	for i := 0; i < this.options.BufferSize; i++ {
		this.rest <- struct{}{}
	}

	this.wg.Wait()
}

func (this *Consumer) Stop() {
	this.cancel()
}

func (this *Consumer) fetch(ctx context.Context, stream string) {
	for {
		select {
		case <-ctx.Done():
			break
		case <-this.rest:
			cmd := this.r.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    this.options.GroupName,
				Consumer: this.options.ConsumerName,
				Streams: []string{
					stream,
					">",
				},
				Count: 1,
				Block: this.options.BlockTimeout,
				NoAck: false,
			})
			result, err := cmd.Result()
			if err != nil {
				time.Sleep(300 * time.Millisecond)
				this.rest <- struct{}{}
				if err, ok := err.(net.Error); ok && err.Timeout() {
					continue
				}
				if err == redis.Nil {
					continue
				}
				this.options.ErrorHandler(errors.Wrapf(err, "cant read group for %+v", stream))
				continue
			}
			if result == nil || len(result) == 0 || result[0].Messages == nil || len(result[0].Messages) == 0 {
				time.Sleep(300 * time.Millisecond)
				this.rest <- struct{}{}
				continue
			}
			this.queue <- &Message{
				ID:     result[0].Messages[0].ID,
				Stream: result[0].Stream,
				Values: result[0].Messages[0].Values,
			}
		}
	}
}

func (this *Consumer) reclaim(ctx context.Context) {
	ticker := time.NewTicker(this.options.ReclaimInterval)
	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			for stream := range this.handlers {
				start := "-"
				end := "+"

				for {
					res, err := this.r.XPendingExt(ctx, &redis.XPendingExtArgs{
						Stream: stream,
						Group:  this.options.GroupName,
						Start:  start,
						End:    end,
						Count:  int64(this.options.BufferSize - len(this.queue)),
					}).Result()
					if err != nil && err != redis.Nil {
						this.options.ErrorHandler(errors.Wrap(err, "error listing pending messages"))
						break
					}

					if len(res) == 0 {
						break
					}

					msgs := make([]string, 0)

					for _, r := range res {
						if r.RetryCount >= int64(this.options.MaxRetry) {
							err = this.r.XAck(ctx, stream, this.options.GroupName, r.ID).Err()
							if err != nil {
								this.options.ErrorHandler(errors.Wrapf(err, "error acknowledging after max retry for %q stream and %q message", stream, r.ID))
								continue
							}
						}
						if r.Idle >= this.options.PendingTimeout {
							claimres, err := this.r.XClaim(ctx, &redis.XClaimArgs{
								Stream:   stream,
								Group:    this.options.GroupName,
								Consumer: this.options.ConsumerName,
								MinIdle:  this.options.PendingTimeout,
								Messages: []string{r.ID},
							}).Result()
							if err != nil && err != redis.Nil {
								this.options.ErrorHandler(errors.Wrapf(err, "error claiming %d message(s)", len(msgs)))
								break
							}
							// If the Redis nil error is returned, it means that
							// the message no longer exists in the stream.
							// However, it is still in a pending state. This
							// could happen if a message was claimed by a
							// consumer, that consumer died, and the message
							// gets deleted (either through a XDEL call or
							// through MAXLEN). Since the message no longer
							// exists, the only way we can get it out of the
							// pending state is to acknowledge it.
							if err == redis.Nil {
								err = this.r.XAck(ctx, stream, this.options.GroupName, r.ID).Err()
								if err != nil {
									this.options.ErrorHandler(errors.Wrapf(err, "error acknowledging after failed claim for %q stream and %q message", stream, r.ID))
									continue
								}
							}
							for _, claimer := range claimres {
								<-this.rest
								this.queue <- &Message{
									ID:     claimer.ID,
									Stream: stream,
									Values: claimer.Values,
								}
							}
						}
					}

					newID, err := incrementMessageID(res[len(res)-1].ID)
					if err != nil {
						this.options.ErrorHandler(err)
						break
					}

					start = newID
				}
			}
		}
	}
}

func (this *Consumer) infiniteWork(ctx context.Context) {
	defer this.wg.Done()
	for {
		select {
		case <-ctx.Done():
			break
		case message := <-this.queue:
			go this.safeProcess(ctx, message)
		}
	}
}

func (this *Consumer) work(ctx context.Context) {
	defer this.wg.Done()
	for {
		select {
		case <-ctx.Done():
			break
		case message := <-this.queue:
			this.safeProcess(ctx, message)
		}
	}
}

func (this *Consumer) safeProcess(ctx context.Context, message *Message) {
	defer func() {
		this.rest <- struct{}{}
	}()
	defer func() {
		if r := recover(); r != nil {
			this.options.ErrorHandler(errors.Errorf("ConsumerFunc panic: %v", r))
		}
	}()
	err := this.handlers[message.Stream](ctx, message)
	if err != nil {
		this.options.ErrorHandler(errors.Wrapf(err, "ConsumerFunc error: %v", message))
		return
	}
	err = this.r.XAck(ctx, message.Stream, this.options.GroupName, message.ID).Err()
	if err != nil {
		this.options.ErrorHandler(errors.Wrap(err, "Ack failed"))
	}
}

// incrementMessageID takes in a message ID (e.g. 1564886140363-0) and
// increments the index section (e.g. 1564886140363-1). This is the next valid
// ID value, and it can be used for paging through messages.
func incrementMessageID(id string) (string, error) {
	parts := strings.Split(id, "-")
	index := parts[1]
	parsed, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		return "", errors.Wrapf(err, "error parsing message ID %q", id)
	}
	return fmt.Sprintf("%s-%d", parts[0], parsed+1), nil
}

var consumer *Consumer

func Instance() *Consumer {
	if consumer == nil {
		consumer = NewConsumer(storage.Redis(), &ConsumerOptions{
			Workers: 1,
			ErrorHandler: func(err error) {
				glgf.Error(err)
			},
		})
	}
	return consumer
}
