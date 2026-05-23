package service

import "fmt"

func BuildCreateCardMutation(name, email string, assetValue float64) string {
	return fmt.Sprintf(`mutation {
    createCard(input: {
        pipe_id: 301498182
        title: "%s"
        fields_attributes: [
            { field_id: "cliente_nome"     field_value: "%s" }
            { field_id: "cliente_email"    field_value: "%s" }
            { field_id: "valor_patrimonio" field_value: "%.2f" }
            { field_id: "status"           field_value: "Aguardando Análise" }
        ]
    }) {
        card {
            id
            title
            current_phase {
                name
            }
        }
    }
}`, name, name, email, assetValue)
}

func BuildUpdateCardMutation(cardID, status, priority string) string {
	return fmt.Sprintf(`mutation {
    updateFieldsValues(input: {
        nodeId: %s
        values: [
            { fieldId: "status"     value: "%s" }
            { fieldId: "prioridade" value: "%s" }
        ]
    }) {
        success
    }
}`, cardID, status, priority)
}
