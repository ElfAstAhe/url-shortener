package repository

import "errors"

type DBConnCheckImMemRepo struct {
}

func newDBConnCheckImMemRepo() (*DBConnCheckImMemRepo, error) {
	return &DBConnCheckImMemRepo{}, nil
}

func (D DBConnCheckImMemRepo) CheckDBConn() error {
	return errors.New("IN_MEMORY DB, where is no connection")
}
