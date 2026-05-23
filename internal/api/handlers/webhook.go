package handlers

import (
	"errors"
	"net/http"

	"github.com/GiovaniGama/pipefy-client-manager/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebhookHandler struct {
	svc *service.ClientService
}

func NewWebhookHandler(svc *service.ClientService) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

func (h *WebhookHandler) CardUpdated(c *gin.Context) {
	var input service.WebhookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.ProcessWebhook(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "cliente não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.AlreadyProcessed {
		c.JSON(http.StatusOK, gin.H{"message": "evento já processado"})
		return
	}

	c.JSON(http.StatusOK, result.Client)
}
