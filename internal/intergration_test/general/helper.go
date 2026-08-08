package general_helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func MakeJSONRequest(method, target string, body interface{}) *http.Request {
	var bodyBytes []byte
	switch typedBody := body.(type) {
	case nil:
		bodyBytes = nil
	case string:
		bodyBytes = []byte(typedBody)
	default:
		bodyBytes, _ = json.Marshal(typedBody)
	}

	req := httptest.NewRequest(method, target, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}
