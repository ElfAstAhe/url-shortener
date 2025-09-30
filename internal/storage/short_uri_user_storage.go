package storage

import (
	"bufio"
	"encoding/json"
	"os"

	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"go.uber.org/zap"
)

type ShortURIUserStorageReader struct {
	File   *os.File
	Reader *bufio.Scanner
	log    *zap.SugaredLogger
}

func NewShortURIUserStorageReader(storagePath string) (*ShortURIUserStorageReader, error) {
	storage, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURIUserStorageReader{
		File:   storage,
		Reader: bufio.NewScanner(storage),
		log:    _log.Log.Sugar(),
	}, nil
}

func (r *ShortURIUserStorageReader) Close() error {
	return r.File.Close()
}

type ShortURIUserStorageWriter struct {
	File   *os.File
	Writer *bufio.Writer
	log    *zap.SugaredLogger
}

func NewShortURIUserStorageWriter(storagePath string) (*ShortURIUserStorageWriter, error) {
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURIUserStorageWriter{
		File:   storage,
		Writer: bufio.NewWriter(storage),
		log:    _log.Log.Sugar(),
	}, nil
}

func (w *ShortURIUserStorageWriter) Close() error {
	return w.File.Close()
}

func (r *ShortURIUserStorageReader) LoadData(cache map[string]*_model.ShortURIUser) error {
	clear(cache)
	for r.Reader.Scan() {
		data := r.Reader.Bytes()
		var entity _model.ShortURIUser
		if err := json.Unmarshal(data, &entity); err != nil {
			r.log.Warn("Failed to unmarshal short URI User", zap.Error(err))
		}
		r.log.Infof("Add short URI suer: %s", entity.ID)
		cache[entity.ID] = &entity
	}
	if err := r.Reader.Err(); err != nil {
		return err
	}

	return nil
}

func (w *ShortURIUserStorageWriter) SaveData(cache map[string]*_model.ShortURIUser) error {
	for id, entity := range cache {
		w.log.Infof("Saving short URI %s to %s", id, w.File.Name())
		data, err := json.Marshal(entity)
		if err != nil {
			w.log.Warnf("Failed to marshal short URI id [%s]", id)
			return err
		}
		if err := w.writeLine(data, entity.ID); err != nil {
			w.log.Warnf("Failed to write short URI id [%s]", id)
		}
	}

	return nil
}

func (w *ShortURIUserStorageWriter) writeLine(data []byte, id string) error {
	if _, err := w.Writer.Write(data); err != nil {
		w.log.Warnf("Failed to write short URI user with id [%s]", id)

		return err
	}
	if err := w.Writer.WriteByte('\n'); err != nil {
		w.log.Warnf("Failed to write term symbol short URI user with id [%s]", id)

		return err
	}

	if err := w.Writer.Flush(); err != nil {
		w.log.Warnf("Failed to flush short URI user with id [%s]", id)

		return err
	}

	return nil
}
