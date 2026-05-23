package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/GiovaniGama/pipefy-client-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestClient(t *testing.T, r *gin.Engine, email string, assetValue float64) {
	t.Helper()
	w := postJSON(r, "/clientes", map[string]any{
		"cliente_nome":     "Cliente Teste",
		"cliente_email":    email,
		"tipo_solicitacao": "Abertura de conta",
		"valor_patrimonio": assetValue,
	})
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestWebhookHighPriority(t *testing.T) {
	r, _ := setupRouter(t)
	email := "alta@example.com"
	createTestClient(t, r, email, 250000)

	w := postJSON(r, "/webhooks/pipefy/card-updated", map[string]any{
		"event_id":      "evt_alta_001",
		"card_id":       "card_001",
		"cliente_email": email,
		"timestamp":     "2026-05-18T12:00:00Z",
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var client models.Client
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &client))
	assert.Equal(t, "prioridade_alta", client.Priority)
	assert.Equal(t, "Processado", client.Status)
}

func TestWebhookNormalPriority(t *testing.T) {
	r, _ := setupRouter(t)
	email := "normal@example.com"
	createTestClient(t, r, email, 150000)

	w := postJSON(r, "/webhooks/pipefy/card-updated", map[string]any{
		"event_id":      "evt_normal_001",
		"card_id":       "card_002",
		"cliente_email": email,
		"timestamp":     "2026-05-18T12:00:00Z",
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var client models.Client
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &client))
	assert.Equal(t, "prioridade_normal", client.Priority)
	assert.Equal(t, "Processado", client.Status)
}

func TestWebhookIdempotency(t *testing.T) {
	r, _ := setupRouter(t)
	email := "idem@example.com"
	createTestClient(t, r, email, 300000)

	payload := map[string]any{
		"event_id":      "evt_idem_001",
		"card_id":       "card_003",
		"cliente_email": email,
		"timestamp":     "2026-05-18T12:00:00Z",
	}

	w1 := postJSON(r, "/webhooks/pipefy/card-updated", payload)
	assert.Equal(t, http.StatusOK, w1.Code)

	var first models.Client
	require.NoError(t, json.Unmarshal(w1.Body.Bytes(), &first))
	assert.Equal(t, "Processado", first.Status)

	w2 := postJSON(r, "/webhooks/pipefy/card-updated", payload)
	assert.Equal(t, http.StatusOK, w2.Code)

	var second map[string]string
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &second))
	assert.Equal(t, "evento já processado", second["message"])
}
