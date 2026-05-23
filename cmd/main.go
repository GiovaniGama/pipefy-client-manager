package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/pipefy-client-manager/internal/api/handlers"
	"github.com/seu-usuario/pipefy-client-manager/internal/database"
	"github.com/seu-usuario/pipefy-client-manager/internal/repository"
	"github.com/seu-usuario/pipefy-client-manager/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("pipefy.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("erro ao conectar ao banco: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("erro ao migrar banco: %v", err)
	}

	clientRepo := repository.NewClientRepository(db)
	eventRepo := repository.NewEventRepository(db)
	clientSvc := service.NewClientService(clientRepo, eventRepo)

	clientHandler := handlers.NewClientHandler(clientSvc)
	webhookHandler := handlers.NewWebhookHandler(clientSvc)

	r := gin.Default()
	r.POST("/clientes", clientHandler.CreateClient)
	r.POST("/webhooks/pipefy/card-updated", webhookHandler.CardUpdated)

	log.Println("Servidor rodando em :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
