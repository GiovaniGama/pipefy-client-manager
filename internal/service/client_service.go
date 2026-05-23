package service

import (
	"log"

	"github.com/GiovaniGama/pipefy-client-manager/internal/models"
	"github.com/GiovaniGama/pipefy-client-manager/internal/repository"
)

type ClientService struct {
	clientRepo *repository.ClientRepository
	eventRepo  *repository.EventRepository
}

func NewClientService(cr *repository.ClientRepository, er *repository.EventRepository) *ClientService {
	return &ClientService{clientRepo: cr, eventRepo: er}
}

type CreateClientInput struct {
	Name        string  `json:"cliente_nome" binding:"required"`
	Email       string  `json:"cliente_email" binding:"required,email"`
	RequestType string  `json:"tipo_solicitacao" binding:"required"`
	AssetValue  float64 `json:"valor_patrimonio" binding:"required"`
}

func (s *ClientService) CreateClient(input CreateClientInput) (*models.Client, error) {
	client := &models.Client{
		Name:        input.Name,
		Email:       input.Email,
		RequestType: input.RequestType,
		AssetValue:  input.AssetValue,
		Status:      "Aguardando Análise",
	}

	if err := s.clientRepo.Create(client); err != nil {
		return nil, err
	}

	mutation := BuildCreateCardMutation(client.Name, client.Email, client.AssetValue)
	log.Printf("[Pipefy] createCard mutation:\n%s\n", mutation)

	return client, nil
}

type WebhookInput struct {
	EventID     string `json:"event_id" binding:"required"`
	CardID      string `json:"card_id" binding:"required"`
	ClientEmail string `json:"cliente_email" binding:"required,email"`
	Timestamp   string `json:"timestamp" binding:"required"`
}

type WebhookResult struct {
	AlreadyProcessed bool
	Client           *models.Client
}

func (s *ClientService) ProcessWebhook(input WebhookInput) (*WebhookResult, error) {
	exists, err := s.eventRepo.Exists(input.EventID)
	if err != nil {
		return nil, err
	}
	if exists {
		return &WebhookResult{AlreadyProcessed: true}, nil
	}

	client, err := s.clientRepo.FindByEmail(input.ClientEmail)
	if err != nil {
		return nil, err
	}

	priority := calculatePriority(client.AssetValue)

	mutation := BuildUpdateCardMutation(input.CardID, "Processado", priority)
	log.Printf("[Pipefy] updateCardField mutation:\n%s\n", mutation)

	client.Status = "Processado"
	client.Priority = priority

	if err := s.clientRepo.Update(client); err != nil {
		return nil, err
	}

	if err := s.eventRepo.Save(input.EventID); err != nil {
		return nil, err
	}

	return &WebhookResult{AlreadyProcessed: false, Client: client}, nil
}

func calculatePriority(assetValue float64) string {
	if assetValue >= 200000 {
		return "prioridade_alta"
	}
	return "prioridade_normal"
}
