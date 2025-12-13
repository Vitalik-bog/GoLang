package memory

import (
	"sync"
	"time"

	"table_collab/internal/domain"
)

type RoomStore struct {
	rooms map[string]*domain.Room
	mu    sync.RWMutex
}

func NewRoomStore() *RoomStore {
	return &RoomStore{
		rooms: make(map[string]*domain.Room),
	}
}

func (s *RoomStore) Save(room *domain.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room.UpdatedAt = time.Now()
	s.rooms[room.ID] = room
	return nil
}

func (s *RoomStore) Get(id string) (*domain.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[id]
	if !exists {
		return nil, ErrNotFound
	}

	return room, nil
}

func (s *RoomStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.rooms, id)
	return nil
}

func (s *RoomStore) GetAll() ([]*domain.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]*domain.Room, 0, len(s.rooms))
	for _, room := range s.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

var ErrNotFound = struct {
	error
}{}
