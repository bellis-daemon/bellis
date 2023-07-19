package public

import "go.mongodb.org/mongo-driver/mongo/options"

func (this *Pagination) ToOptions() *options.FindOptions {
	skip := int64(this.PageKey)
	if this.PageKey == 0 {
		this.PageKey = 1
	}
	offset := int64((this.PageKey - 1) * this.PageSize)
	if offset < 0 {
		offset = 0
	}
	return &options.FindOptions{
		Skip:  &skip,
		Limit: &offset,
	}
}
