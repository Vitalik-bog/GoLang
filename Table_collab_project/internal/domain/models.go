package domain

import "time"

type RoomType string

const (
	RoomTypeDocument   RoomType = "document"
	RoomTypeWhiteboard RoomType = "whiteboard"
	RoomTypeTable      RoomType = "table"
)

type Room struct {
	ID          string
	Name        string
	Type        RoomType
	OwnerID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsActive    bool
	MaxClients  int
	ClientCount int
	Content     string
	Version     int
	TableData   map[string]interface{}
}

type User struct {
	ID       string
	Username string
	Color    string
	Cursor   CursorPosition
	JoinedAt time.Time
	IsActive bool
}

type CursorPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
