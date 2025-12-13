package collaboration

import "table_collab/internal/domain"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ApplyTextUpdate(currentText string, update domain.Event) (string, int) {
	if payload, ok := update.Payload.(map[string]interface{}); ok {
		if text, ok := payload["text"].(string); ok {
			return text, update.Version + 1
		}
	}
	return currentText, 0
}

func (s *Service) ValidateEvent(event domain.Event) bool {
	switch event.Type {
	case domain.EventJoinRoom, domain.EventLeaveRoom,
		domain.EventCursorMove, domain.EventTextUpdate,
		domain.EventChatMessage:
		return true
	default:
		return false
	}
}
