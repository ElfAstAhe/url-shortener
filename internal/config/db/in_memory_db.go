package db

import (
	"bufio"
	"encoding/json"
	"os"

	_log "github.com/ElfAstAhe/url-shortener/internal/logger"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
)

type InMemoryDB struct {
	ShortURI map[string]*_model.ShortURI
}

var InMemoryDBInstance *InMemoryDB

func newInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		ShortURI: make(map[string]*_model.ShortURI),
	}
}

/*
func (db *InMemoryDB) GetShortURI(id string) (*_model.ShortURI, error) {
    // ..

    return nil, nil
}

func (db *InMemoryDB) PutShortURI(id string, shortURL *_model.ShortURI) error {
    // ..

    return nil
}

func (db *InMemoryDB) DeleteShortURI(id string) error {
    // ..

    return nil
}
*/

func (db *InMemoryDB) ShortURIExists(id string) (bool, error) {
	_, ok := db.ShortURI[id]

	return ok, nil
}

func (db *InMemoryDB) LoadData(storagePath string) error {
	log := _log.Log.Sugar()
	storage, err := os.OpenFile(storagePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer storage.Close()

	scanner := bufio.NewScanner(storage)
	for scanner.Scan() {
		data := scanner.Bytes()
		var shortURL _model.ShortURI
		if err := json.Unmarshal(data, &shortURL); err != nil {
			log.Warn("Failed to unmarshal short URI")
		}
		log.Infof("Add short URI: %s", shortURL.ID)
		InMemoryDBInstance.ShortURI[shortURL.Key] = &shortURL
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (db *InMemoryDB) SaveData(storagePath string) error {
	log := _log.Log.Sugar()
	storage, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer storage.Close()

	writer := bufio.NewWriter(storage)
	for id, shortURL := range db.ShortURI {
		log.Infof("Saving short URI %s to %s", id, storagePath)
		data, err := json.Marshal(shortURL)
		if err != nil {
			log.Warnf("Failed to marshal short URI id [%s]", id)
			continue
		}
		if _, err = writer.Write(data); err != nil {
			log.Warnf("Failed to write short URI id [%s]", id)
			continue
		}
		if err = writer.WriteByte('\n'); err != nil {
			log.Warnf("Failed to write term symbol short URI id [%s]", id)
			continue
		}

		if err = writer.Flush(); err != nil {
			log.Warnf("Failed to flush short URI id [%s]", id)
		}
	}

	return nil
}

func init() {
	InMemoryDBInstance = newInMemoryDB()
}
