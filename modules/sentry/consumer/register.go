package consumer

import (
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
)

var stream = redistream.NewConsumer(storage.Redis(), &redistream.ConsumerOptions{
	GroupName: "Sentry",
})

func Serve() {
	entityUpdate()
	entityClaim()
	entityDelete()
	stream.Serve()
}
