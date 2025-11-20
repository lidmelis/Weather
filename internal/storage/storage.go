package storage

import (
	"errors"
	"sync"
	"time"

	weatherpb "weatherProject/pkg/weather/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrWeatherNotFound = errors.New("weather not found")
)

type WeatherStorage struct {
	mu      sync.RWMutex
	weather map[int64]*weatherpb.Weather
	lastID  int64
}

func NewWeatherStorage() *WeatherStorage {
	return &WeatherStorage{
		weather: make(map[int64]*weatherpb.Weather),
		lastID:  0,
	}
}

func (s *WeatherStorage) Create(info *weatherpb.WeatherInfo) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	now := timestamppb.New(time.Now())

	s.weather[s.lastID] = &weatherpb.Weather{
		Id:        s.lastID,
		Info:      info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.lastID
}

func (s *WeatherStorage) GetByID(id int64) (*weatherpb.Weather, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	weather, ok := s.weather[id]
	if !ok {
		return nil, ErrWeatherNotFound
	}
	return weather, nil
}

func (s *WeatherStorage) GetAll() []*weatherpb.Weather {
	s.mu.RLock()
	defer s.mu.RUnlock()

	weathers := make([]*weatherpb.Weather, 0, len(s.weather))
	for _, weather := range s.weather {
		weathers = append(weathers, weather)
	}
	return weathers
}

func (s *WeatherStorage) Update(id int64, city *string, temperature *float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	weather, ok := s.weather[id]
	if !ok {
		return ErrWeatherNotFound
	}

	if city != nil {
		weather.Info.City = *city
	}
	if temperature != nil {
		weather.Info.Tempereture = *temperature
	}

	weather.UpdatedAt = timestamppb.New(time.Now())
	s.weather[id] = weather

	return nil
}

func (s *WeatherStorage) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.weather[id]
	if !ok {
		return ErrWeatherNotFound
	}
	delete(s.weather, id)
	return nil
}
