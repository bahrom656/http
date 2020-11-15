package banners

import (
	"context"
	"errors"
	"sync"
)

type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
}

func (s *Service) GetAll(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}

	return nil, errors.New("items not found")
}
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.items = append(s.items, item)
	return item, nil
}
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for index, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:index], s.items[index+1:]...)
			return banner, nil
		}
	}

	return nil, errors.New("items not found")
}
