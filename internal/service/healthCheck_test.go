package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_PingAction(t *testing.T) {
	t.Parallel()
	healthService := NewHealthCheck()

	testCases := []struct {
		name string

		expectedMessage string
		expectedError   error
	}{
		{
			name: "normal case - length 12",

			expectedMessage: "OK",
			expectedError:   nil,
		},
		{
			name: "normal case - length 16",

			expectedMessage: "OK",
			expectedError:   nil,
		},
		{
			name: "normal case - length 1000",

			expectedMessage: "OK",
			expectedError:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			message, err := healthService.Ping()
			assert.Equal(testItem, tc.expectedError, err)
			assert.Equal(testItem, tc.expectedMessage, message)
		})
	}

}
