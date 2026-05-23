package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GiovaniGama/pipefy-client-manager/internal/api/handlers"
	"github.com/GiovaniGama/pipefy-client-manager/internal/database"
	"github.com/GiovaniGama/pipefy-client-manager/internal/models"
	"github.com/GiovaniGama/pipefy-client-manager/internal/repository"
	"github.com/GiovaniGama/pipefy-client-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.Migrate(db))

	clientRepo := repository.NewClientRepository(db)
	eventRepo := repository.NewEventRepository(db)
	svc := service.NewClientService(clientRepo, eventRepo)

	r := gin.New()
	r.POST("/clientes", handlers.NewClientHandler(svc).CreateClient)
	r.POST("/webhooks/pipefy/card-updated", handlers.NewWebhookHandler(svc).CardUpdated)

	return r, db
}

func postJSON(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateValidClient(t *testing.T) {
	r, db := setupRouter(t)

	w := postJSON(r, "/clientes", map[string]any{
		"cliente_nome":     "João Silva",
		"cliente_email":    "joao.silva@example.com",
		"tipo_solicitacao": "Atualização cadastral",
		"valor_patrimonio": 250000,
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Client
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Aguardando Análise", response.Status)

	var clientInDB models.Client
	require.NoError(t, db.Where("email = ?", "joao.silva@example.com").First(&clientInDB).Error)
	assert.Equal(t, "Aguardando Análise", clientInDB.Status)
}
