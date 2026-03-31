package integration

import (
	"Today-Todo/models"
	"Today-Todo/routers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) http.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB failed: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if err := db.AutoMigrate(&models.Todo{}, &models.WaterRecord{}, &models.StandRecord{}, &models.ShortVideoRecord{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	models.DB = db
	return routers.SetupRouter()
}

func performJSONRequest(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body failed: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	return resp
}

func TestTodoEndpointsIntegration(t *testing.T) {
	handler := setupTestRouter(t)

	createResp := performJSONRequest(t, handler, http.MethodPost, "/api/v1/todos", map[string]any{
		"title":       "write tests",
		"description": "minimum integration test",
		"priority":    "high",
	})
	if createResp.Code != http.StatusOK {
		t.Fatalf("create todo status=%d, body=%s", createResp.Code, createResp.Body.String())
	}

	listResp := performJSONRequest(t, handler, http.MethodGet, "/api/v1/todos", nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("list todos status=%d, body=%s", listResp.Code, listResp.Body.String())
	}

	var todos []map[string]any
	if err := json.Unmarshal(listResp.Body.Bytes(), &todos); err != nil {
		t.Fatalf("unmarshal todos failed: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("todos len=%d, want 1", len(todos))
	}
}

func TestInsightsEndpointIntegration(t *testing.T) {
	handler := setupTestRouter(t)

	_ = performJSONRequest(t, handler, http.MethodPost, "/api/v1/todos", map[string]any{
		"title":    "deep work",
		"priority": "high",
	})
	_ = performJSONRequest(t, handler, http.MethodPost, "/api/v1/health/water", map[string]any{
		"user_id": 1,
		"amount":  200,
	})
	_ = performJSONRequest(t, handler, http.MethodPost, "/api/v1/health/short-video", map[string]any{
		"user_id": 1,
		"count":   5,
	})

	insightResp := performJSONRequest(t, handler, http.MethodGet, "/api/v1/insights/today?user_id=1", nil)
	if insightResp.Code != http.StatusOK {
		t.Fatalf("insights status=%d, body=%s", insightResp.Code, insightResp.Body.String())
	}

	var payload map[string]map[string]any
	if err := json.Unmarshal(insightResp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal insights failed: %v", err)
	}

	data, ok := payload["data"]
	if !ok {
		t.Fatalf("missing data field: %s", insightResp.Body.String())
	}

	if risk, _ := data["risk_level"].(string); risk != "high" {
		t.Fatalf("risk_level=%v, want high", data["risk_level"])
	}

	if action, _ := data["suggested_action"].(string); action == "" {
		t.Fatalf("suggested_action should not be empty: %#v", data)
	}
}
