package sqldb

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitMockDB(t *testing.T) {
	for i := 0; i < 10; i++ {
		testName := "TestInitMockDB" + strconv.Itoa(i)
		t.Run(testName, func(testItem *testing.T) {
			testItem.Parallel()
			_, err := NewMiniPostgres(testItem)
			assert.NoError(testItem, err)
		})

	}

}
