package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDBConnectorMock(t *testing.T) {
	connectorMock, errConnector := InitDBConnectorMock(t)
	if errConnector != nil {
		t.Fatal(errConnector)
	}
	assert.NoError(t, errConnector)
	assert.NotNil(t, connectorMock)
}
