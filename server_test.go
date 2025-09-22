package portfolio

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jojojojodr/portfolio/internal/routers"
	"github.com/gin-gonic/gin"
)

func TestServer_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svr := gin.Default()
	routers.FrontendRouter(svr)
	routers.V1Router(svr)
	routers.HandleRouter(svr)

	req := httptest.NewRequest("GET", "/v1/health", nil)
	w := httptest.NewRecorder()

	svr.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("unexpected health response: %v", resp)
	}
}
