package index

import (
	"context"

	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/mongo"
)

var indexes = make(map[**mongo.Collection][]mongo.IndexModel)

func RegistrerIndex(col **mongo.Collection, idxs []mongo.IndexModel) {
	indexes[col] = idxs
}

func InitIndexes() {
	ctx := context.Background()
	for col, idxs := range indexes {
		glgf.Debug(*col, idxs)
		c := col
		i := idxs
		go func() {
			_, err := (*c).Indexes().CreateMany(ctx, i)
			if err != nil {
				glgf.Error("create mongo indexes error: ", err)
			}
		}()
	}
}
