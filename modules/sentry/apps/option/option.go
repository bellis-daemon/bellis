package option

import (
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

func ToOption[T any](source bson.M) (target T) {
	err := mapstructure.Decode(source, &target)
	if err != nil {
		panic(err)
	}
	return
}
