package consumer

import (
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
)

var consumer *redistream.Consumer

func instance() *redistream.Consumer {
	if consumer == nil {
		consumer = redistream.NewConsumer(storage.Redis(), &redistream.ConsumerOptions{
			Workers: 8,
			ErrorHandler: func(err error) {
				glgf.Error(err)
			},
		})
	}
	return consumer
}

func Serve() {
	consumer.Serve()
}
