package index

import (
	"context"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/mongo"
)

var indexes = make(map[*mongo.Collection][]mongo.IndexModel)

func RegistrerIndex(col *mongo.Collection, idxs []mongo.IndexModel) {
	indexes[col] = idxs
}

func InitIndexes() {
	ctx := context.Background()
	for col, idxs := range indexes {
		_, err := col.Indexes().CreateMany(ctx, idxs)
		if err != nil {
			glgf.Error("create mongo indexes error: ", err)
		}
	}
}
