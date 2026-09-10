package base62_helper

import (
	"testing"

	"github.com/ivanrad/base62"
	_ "github.com/ivanrad/base62"
	"github.com/stretchr/testify/assert"
)

func TestBase62(t *testing.T) {
	testCase := []struct {
		name string
		data string
	}{
		{
			name: "normal case",
			data: "1234",
		},
	}

	for _, tc := range testCase {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			base62Helper := New()
			res := base62Helper.Encode(tc.data)
			assert.Equal(t, base62.EncodeToString([]byte(tc.data)), res)

		})

	}
}
