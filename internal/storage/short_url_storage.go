package storage

import (
	"bufio"
	"encoding/json"
	"os"

	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"go.uber.org/zap"
)

type ShortURLStorageReader struct {
	File   *os.File
	Reader *bufio.Scanner
	log    *zap.SugaredLogger
}

func NewShortURLStorageReader(storagePath string) (*ShortURLStorageReader, error) {
	storage, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURLStorageReader{
		File:   storage,
		Reader: bufio.NewScanner(storage),
		log:    _log.Log.Sugar(),
	}, nil
}

type ShortURLStorageWriter struct {
	File   *os.File
	Writer *bufio.Writer
	log    *zap.SugaredLogger
}

func NewShortURLStorageWriter(storagePath string) (*ShortURLStorageWriter, error) {
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortURLStorageWriter{
		File:   storage,
		Writer: bufio.NewWriter(storage),
		log:    _log.Log.Sugar(),
	}, nil
}

func (r *ShortURLStorageReader) Close() error {
	return r.File.Close()
}

func (r *ShortURLStorageReader) LoadData(cache map[string]*_model.ShortURI) error {
	clear(cache)
	for r.Reader.Scan() {
		data := r.Reader.Bytes()
		var shortURL _model.ShortURI
		if err := json.Unmarshal(data, &shortURL); err != nil {
			r.log.Warn("Failed to unmarshal short URI")
		}
		r.log.Infof("Add short URI: %s", shortURL.ID)
		cache[shortURL.Key] = &shortURL
	}
	if err := r.Reader.Err(); err != nil {
		return err
	}

	return nil
}

func (w *ShortURLStorageWriter) Close() error {
	return w.File.Close()
}

func (w *ShortURLStorageWriter) SaveData(cache map[string]*_model.ShortURI) error {
	for id, shortURL := range cache {
		w.log.Infof("Saving short URI %s to %s", id, w.File.Name())
		data, err := json.Marshal(shortURL)
		if err != nil {
			w.log.Warnf("Failed to marshal short URI id [%s]", id)
			return err
		}
		if err := w.writeLine(data, shortURL.ID); err != nil {
			w.log.Warnf("Failed to write short URI id [%s]", id)
		}
	}

	return nil
}

func (w *ShortURLStorageWriter) writeLine(data []byte, id string) error {
	if _, err := w.Writer.Write(data); err != nil {
		w.log.Warnf("Failed to write short URI id [%s]", id)

		return err
	}
	if err := w.Writer.WriteByte('\n'); err != nil {
		w.log.Warnf("Failed to write term symbol short URI id [%s]", id)

		return err
	}

	if err := w.Writer.Flush(); err != nil {
		w.log.Warnf("Failed to flush short URI id [%s]", id)

		return err
	}

	return nil
}
