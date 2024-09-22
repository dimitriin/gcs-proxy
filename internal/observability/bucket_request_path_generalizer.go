package observability

import (
	"net/http"

	"github.com/gorilla/mux"
)

type BucketRequestPathGeneralizer struct{}

func NewBucketRequestPathGeneralizer() *BucketRequestPathGeneralizer {
	return &BucketRequestPathGeneralizer{}
}

func (b *BucketRequestPathGeneralizer) GetRequestGeneralizedPath(r *http.Request) string {
	return mux.Vars(r)["bucket"]
}
