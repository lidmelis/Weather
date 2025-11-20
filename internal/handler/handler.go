package handler

import (
	"context"

	"weatherProject/internal/storage"
	weatherpb "weatherProject/pkg/weather/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type WeatherHandler struct {
	weatherpb.UnimplementedWeatherAPIServer
	storage *storage.WeatherStorage
}

func NewWeatherHandler(storage *storage.WeatherStorage) *WeatherHandler {
	return &WeatherHandler{
		storage: storage,
	}
}

func (h *WeatherHandler) Create(ctx context.Context, req *weatherpb.CreateRequest) (*weatherpb.CreateResponse, error) {
	id := h.storage.Create(req.GetInfo())
	return &weatherpb.CreateResponse{Id: id}, nil
}

func (h *WeatherHandler) Get(ctx context.Context, req *weatherpb.GetRequest) (*weatherpb.GetResponse, error) {
	weather, err := h.storage.GetByID(req.GetId())
	if err != nil {
		return nil, err
	}
	return &weatherpb.GetResponse{Weather: weather}, nil
}

func (h *WeatherHandler) List(ctx context.Context, req *emptypb.Empty) (*weatherpb.ListResponse, error) {
	weathers := h.storage.GetAll()
	return &weatherpb.ListResponse{Weathers: weathers}, nil
}

func (h *WeatherHandler) Update(ctx context.Context, req *weatherpb.UpdateRequest) (*emptypb.Empty, error) {
	var city *string
	var temperature *float64

	if req.City != nil {
		city = &req.City.Value
	}
	if req.Tempereture != nil {
		temperature = &req.Tempereture.Value
	}

	err := h.storage.Update(req.GetId(), city, temperature)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *WeatherHandler) Delete(ctx context.Context, req *weatherpb.DeleteRequest) (*emptypb.Empty, error) {
	err := h.storage.Delete(req.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
