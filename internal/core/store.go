package core

import (
	"context"
	"fmt"
	"os"

	"demoscraper/internal/core/entities"
)

type Store struct {
	marshaller Marshaller
}

func NewStore(marshaller Marshaller) *Store {
	return &Store{
		marshaller: marshaller,
	}
}

func (r *Store) Save(ctx context.Context, name string, entries <-chan entities.CrawlEntry) error {
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() { _ = file.Close() }()

	for {
		select {
		case <-ctx.Done():
			return nil
		case entry, ok := <-entries:
			if !ok {
				return nil
			}

			marshalledBytes, err := r.marshaller.Marshal(entry)
			if err != nil {
				return fmt.Errorf("marshal: %w", err)
			}

			if _, err = file.Write(marshalledBytes); err != nil {
				return fmt.Errorf("write: %w", err)
			}

			if _, err = file.WriteString("\n"); err != nil {
				return fmt.Errorf("write: %w", err)
			}
		}
	}
}
