// package handler_test // Để _test để đảm bảo tính đóng gói độc lập

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/homework/lab/internal/api"
// 	"github.com/homework/lab/internal/config"
// 	"github.com/stretchr/testify/assert"
// )

// func TestHealthCheck_Integration(t *testing.T) {
// 	t.Parallel()

// 	testCases := []struct {
// 		name string

// 		// Style của bạn: Nhận vào router engine thực tế và trả về kết quả HTTPRecorder
// 		setupTestHTTP func(router *api.Engine) *httptest.ResponseRecorder

// 		expectedStatusCode      int
// 		expectedResponseContain string
// 	}{
// 		{
// 			name: "Health check successfully",
// 			setupTestHTTP: func(router *api.Engine) *httptest.ResponseRecorder {
// 				req := httptest.NewRequest(http.MethodGet, "/health_check", nil)
// 				respRec := httptest.NewRecorder()
// 				router.ServeHTTP(respRec, req)
// 				return respRec
// 			},
// 			expectedStatusCode:      http.StatusOK,
// 			expectedResponseContain: `{"message":"OK"`, // Kiểm tra cấu trúc JSON thô
// 		},
// 		{
// 			name: "Wrong health endpoint",
// 			setupTestHTTP: func(router *api.Engine) *httptest.ResponseRecorder {
// 				req := httptest.NewRequest(http.MethodGet, "/health_check_wrong", nil)
// 				respRec := httptest.NewRecorder()
// 				router.ServeHTTP(respRec, req)
// 				return respRec
// 			},
// 			expectedStatusCode: http.StatusNotFound,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc // Chốt giữ context vòng lặp cho Parallel Test

// 		t.Run(tc.name, func(testItem *testing.T) {
// 			testItem.Parallel()
// 			cfg, err := config.LoadConfig()
// 			if err != nil {
// 				panic(err)
// 			}
// 			apiEngine := api.NewEngine(cfg)

// 			rec := tc.setupTestHTTP(apiEngine)

// 			assert.Equal(t, tc.expectedStatusCode, rec.Code)
// 			if len(tc.expectedResponseContain) > 0 {
// 				assert.Equal(t, len(rec.Body.String()), len(tc.expectedResponseContain)+12+3)
// 			}
// 			assert.Contains(t, rec.Body.String(), tc.expectedResponseContain)
// 		})
// 	}
// }
