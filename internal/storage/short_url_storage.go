package storage

import (
	"bufio"
	"encoding/json"
	"os"

	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	"go.uber.org/zap"
)

type ShortUrlStorageReader struct {
	File   *os.File
	Reader *bufio.Scanner
	log    *zap.SugaredLogger
}

func NewShortUrlStorageReader(storagePath string) (*ShortUrlStorageReader, error) {
	storage, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortUrlStorageReader{
		File:   storage,
		Reader: bufio.NewScanner(storage),
		log:    _log.Log.Sugar(),
	}, nil
}

type ShortUrlStorageWriter struct {
	File   *os.File
	Writer *bufio.Writer
	log    *zap.SugaredLogger
}

func NewShortUrlStorageWriter(storagePath string) (*ShortUrlStorageWriter, error) {
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &ShortUrlStorageWriter{
		File:   storage,
		Writer: bufio.NewWriter(storage),
		log:    _log.Log.Sugar(),
	}, nil
}

func (r *ShortUrlStorageReader) Close() error {
	return r.File.Close()
}

func (r *ShortUrlStorageReader) LoadData(cache map[string]*_model.ShortURI) error {
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

func (w *ShortUrlStorageWriter) Close() error {
	return w.File.Close()
}

func (w *ShortUrlStorageWriter) SaveData(cache map[string]*_model.ShortURI) error {
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

func (w *ShortUrlStorageWriter) writeLine(data []byte, id string) error {
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
