package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bruce-Sakura/UploadMyself/backend/handler"
	"github.com/Bruce-Sakura/UploadMyself/backend/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.Skill{}, &model.Voice{}, &model.Avatar{}, &model.Task{})
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := handler.New(db)

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	v1 := r.Group("/api/v1")
	{
		s := v1.Group("/skills")
		s.POST("", h.CreateSkill)
		s.GET("", h.ListSkills)
		s.GET("/:id", h.GetSkill)
	}
	return r
}

func TestHealth(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateAndGetSkill(t *testing.T) {
	db := setupTestDB(t)
	r := setupRouter(db)

	// Create
	body, _ := json.Marshal(map[string]string{"name": "test", "corpus": "hello world"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("create: expected 201, got %d", w.Code)
	}

	var created model.Skill
	json.Unmarshal(w.Body.Bytes(), &created)
	if created.Name != "test" {
		t.Fatalf("expected name=test, got %s", created.Name)
	}

	// Get
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/skills/"+created.ID, nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Fatalf("get: expected 200, got %d", w2.Code)
	}
}
