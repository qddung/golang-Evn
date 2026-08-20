package connection

import (
	"testing"

	"github.com/homework/lab/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
)

func TestInitDBConnectorMock(t *testing.T) {
	fix := fixture.NewUserTestCase(t)
	connectorMock, errConnector := InitDBConnectorMock(t, fix)
	if errConnector != nil {
		t.Fatal(errConnector)
	}
	assert.NoError(t, errConnector)
	assert.NotNil(t, connectorMock)
}
