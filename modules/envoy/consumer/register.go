package consumer

import (
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
)

var stream = redistream.NewConsumer(storage.Redis(), &redistream.ConsumerOptions{
	GroupName: "Envoy",
})

func Serve() {
	emailCaptcha()
	entityOfflineAlert()
	entityOnlineAlert()
	stream.Serve()
}
