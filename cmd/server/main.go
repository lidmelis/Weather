package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"weatherProject/internal/handler"
	"weatherProject/internal/storage"
	weatherpb "weatherProject/pkg/weather/v1"
)

type server struct {
	weatherpb.UnimplementedWeatherAPIServer
}

func (s *server) Get(ctx context.Context, req *weatherpb.GetRequest) (*weatherpb.GetResponse, error) {
	log.Print(req.GetId())
	return &weatherpb.GetResponse{Weather: &weatherpb.Weather{Info: &weatherpb.WeatherInfo{City: "Moscow", Tempereture: -1.37}}}, nil
}

func main() {
	weatherStorage := storage.NewWeatherStorage()
	weatherHandler := handler.NewWeatherHandler(weatherStorage)

	weatherStorage.Create(&weatherpb.WeatherInfo{
		City:        "Moscow",
		Tempereture: -1.37,
	})
	weatherStorage.Create(&weatherpb.WeatherInfo{
		City:        "London",
		Tempereture: 8.5,
	})

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	reflection.Register(s)

	weatherpb.RegisterWeatherAPIServer(s, weatherHandler)

	log.Printf("weather gRPC server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
