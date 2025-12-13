package domain

type EventType string

const (
	EventJoinRoom    EventType = "join_room"
	EventLeaveRoom   EventType = "leave_room"
	EventCursorMove  EventType = "cursor_move"
	EventTextUpdate  EventType = "text_update"
	EventElementAdd  EventType = "element_add"
	EventChatMessage EventType = "chat_message"
	EventError       EventType = "error"
	EventSync        EventType = "sync"
)

type Event struct {
	Type      EventType   `json:"type"`
	RoomID    string      `json:"room_id,omitempty"`
	UserID    string      `json:"user_id,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp int64       `json:"timestamp"`
	Version   int         `json:"version,omitempty"`
}

type JoinRoomPayload struct {
	Username string `json:"username"`
	Color    string `json:"color,omitempty"`
}

type TextUpdatePayload struct {
	Text    string `json:"text"`
	Version int    `json:"version"`
}
