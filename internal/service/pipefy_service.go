package service

import "fmt"

func BuildCreateCardMutation(name, email string, assetValue float64) string {
	return fmt.Sprintf(`mutation {
  createCard(input: {
    pipe_id: "301498182"
    fields_attributes: [
      { field_id: "cliente_nome", field_value: "%s" }
      { field_id: "cliente_email", field_value: "%s" }
      { field_id: "valor_patrimonio", field_value: "%.2f" }
    ]
  }) {
    card {
      id
      title
    }
  }
}`, name, email, assetValue)
}

func BuildUpdateCardMutation(cardID, status, priority string) string {
	return fmt.Sprintf(`mutation {
  updateCardField(input: {
    card_id: "%s"
    field_id: "status"
    new_value: "%s"
  }) {
    card {
      id
    }
  }
}

mutation {
  updateCardField(input: {
    card_id: "%s"
    field_id: "prioridade"
    new_value: "%s"
  }) {
    card {
      id
    }
  }
}`, cardID, status, cardID, priority)
}
