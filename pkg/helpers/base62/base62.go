package base62_helper

import (
	"github.com/ivanrad/base62"
)

//go:generate mockery --name=Base62Helper --filename=base62.go --outpkg=mocks
type Base62Helper interface {
	Encode(data string) string
}

type base62Helper struct {
}

func (b *base62Helper) Encode(data string) string {
	res := base62.EncodeToString([]byte(data))
	return res
}

func New() Base62Helper {
	return &base62Helper{}
}
