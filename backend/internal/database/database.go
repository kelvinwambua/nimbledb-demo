package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kelvinwambua/nimbledb/network"
)

type Service interface {
	Health() map[string]string
	Close() error
	Query(query string) ([]string, [][]interface{}, error)
	Execute(query string) error
	QueryRow(query string) ([]interface{}, error)
}

type service struct {
	client *network.Client
	addr   string
}

var dbInstance *service

func New(addr string) Service {
	if dbInstance != nil {
		return dbInstance
	}

	client := network.NewClient(addr)
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to NimbleDB at %s: %v", addr, err)
	}

	dbInstance = &service{
		client: client,
		addr:   addr,
	}

	log.Printf("Connected to NimbleDB at %s", addr)
	return dbInstance
}

func (s *service) Health() map[string]string {
	stats := make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- s.client.Ping()
	}()

	select {
	case err := <-done:
		if err != nil {
			stats["status"] = "down"
			stats["error"] = fmt.Sprintf("db down: %v", err)
			return stats
		}
		stats["status"] = "up"
		stats["message"] = "NimbleDB is healthy"
	case <-ctx.Done():
		stats["status"] = "down"
		stats["error"] = "health check timeout"
	}

	return stats
}

func (s *service) Close() error {
	log.Println("Disconnecting from NimbleDB")
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *service) Query(query string) ([]string, [][]interface{}, error) {
	return s.client.Query(query)
}

func (s *service) Execute(query string) error {
	_, _, err := s.client.Query(query)
	return err
}

func (s *service) QueryRow(query string) ([]interface{}, error) {
	_, rows, err := s.client.Query(query)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found")
	}
	return rows[0], nil
}
