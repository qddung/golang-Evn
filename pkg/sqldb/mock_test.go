package sqldb

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitMockNewMiniPostgresEmbedded(t *testing.T) {
	for i := 0; i < 10; i++ {
		testName := "TestInitMockDB" + strconv.Itoa(i)
		t.Run(testName, func(testItem *testing.T) {
			testItem.Parallel()
			_, err := NewMiniPostgresEmbedded(testItem)
			assert.NoError(testItem, err)
		})

	}

}

func TestInitMockMemPostgres(t *testing.T) {
	_, err := NewMiniPostgres(t)
	assert.NoError(t, err)
}
