package consumer

import (
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
)

var stream = redistream.NewConsumer(storage.Redis(), &redistream.ConsumerOptions{
	GroupName: "Dispatcher",
})

func Serve() {
	stream.Serve()
}
