package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/ElfAstAhe/url-shortener/pkg/client/audit/dto"
)

type IncomeAuditStorageWriter struct {
	locker sync.Mutex
	File   *os.File
	Writer *bufio.Writer
}

func NewIncomeAuditStorageWriter(storagePath string) (*IncomeAuditStorageWriter, error) {
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &IncomeAuditStorageWriter{
		File:   storage,
		Writer: bufio.NewWriter(storage),
	}, nil
}

// Closer interface

func (w *IncomeAuditStorageWriter) Close() error {
	w.locker.Lock()
	defer w.locker.Unlock()

	err := w.Writer.Flush()
	if err != nil {
		return err
	}

	return w.File.Close()
}

// implementation

func (w *IncomeAuditStorageWriter) SaveData(ctx context.Context, dto *dto.IncomeAuditDto) error {
	if dto == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		{
			jsonBytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}

			err = w.writeLine(jsonBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *IncomeAuditStorageWriter) writeLine(jsonBytes []byte) error {
	w.locker.Lock()
	defer w.locker.Unlock()

	_, err := w.Writer.Write(jsonBytes)
	if err != nil {
		return err
	}
	_, err = w.Writer.Write([]byte("\n"))
	if err != nil {
		return err
	}

	err = w.Writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
