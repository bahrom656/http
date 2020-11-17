package banners

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
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
	Image 	string
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
var ID int64 = 0

func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item.ID == 0 {
		ID++
		item.ID = ID

		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image
			err := uFile(file, "./web/banners/"+item.Image)
			if err != nil {
				return nil, err
			}
		}
		s.items = append(s.items, item)
		return item, nil
	}
	for index, value := range s.items {
		if value.ID == item.ID {
			if item.Image != "" {
				item.Image = fmt.Sprint(item.ID) + "." + item.Image
				err := uFile(file, "./web/banners/"+item.Image)
				if err != nil {
					return nil, err
				}
			} else {
				item.Image = s.items[index].Image
			}
			s.items[index] = item
			return item, nil
		}
	}
	return nil, errors.New("item not found")
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
	return nil, errors.New("item not found")
}

func uFile(file multipart.File, path string) error {
	var data, err = ioutil.ReadAll(file)
	if err != nil {
		return errors.New("not readble data")
	}

	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return errors.New("not saved from folder ")
	}

	return nil
}
