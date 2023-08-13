package public

import "go.mongodb.org/mongo-driver/mongo/options"

func (this *Pagination) ToOptions() *options.FindOptions {
	skip := int64(this.PageKey)
	limit := int64(this.PageSize)
	return &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
}
