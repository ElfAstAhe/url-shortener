package repository

import (
	"database/sql"
	"errors"

	_db "github.com/ElfAstAhe/url-shortener/internal/config/db"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
	"github.com/google/uuid"
)

type shortURIUserPgRepo struct {
	db _db.DB
}

const (
	getShortURIUserSQL               string = `select id, short_uri_id, user_id from short_uri_users where id = $1`
	getShortURIUserByUniqueSQL       string = `select id, short_uri_id, user_id from short_uri_users where user_id = $1 and short_uri_id = $2`
	listShortURIUserAllByUserSQL     string = `select id, short_uri_id, user_id from short_uri_users where user_id = $1`
	listShortURIUserAllByShortURISQL string = `select id, short_uri_id, user_id from short_uri_users where short_uri_id = $1`
	createShortURIUserSQL            string = `insert into short_uri_users(id, short_uri_id, user_id) values($1, $2, $3)`
)

func newShortURIUserPgRepo(db _db.DB) (*shortURIUserPgRepo, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &shortURIUserPgRepo{db: db}, nil
}

func (pgsu *shortURIUserPgRepo) Get(ID string) (*_model.ShortURIUser, error) {
	row := pgsu.db.GetDB().QueryRow(getShortURIUserSQL, ID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURIUser{}
	// id, short_uri_id, user_id
	err := row.Scan(&result.ID, &result.ShortURIID, &result.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgsu *shortURIUserPgRepo) GetByUnique(userID string, shortURIID string) (*_model.ShortURIUser, error) {
	row := pgsu.db.GetDB().QueryRow(getShortURIUserByUniqueSQL, userID, shortURIID)
	if row.Err() != nil && !errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}

	var result = _model.ShortURIUser{}
	// id, short_uri_id, user_id
	err := row.Scan(&result.ID, &result.ShortURIID, &result.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pgsu *shortURIUserPgRepo) ListAllByUser(userID string) ([]*_model.ShortURIUser, error) {
	res := make([]*_model.ShortURIUser, 0)
	rows, err := pgsu.db.GetDB().Query(listShortURIUserAllByUserSQL, userID)
	if err != nil {
		return res, err
	}
	defer _utl.CloseOnly(rows)

	for rows.Next() {
		var result = _model.ShortURIUser{}
		// id, short_uri_id, user_id
		err := rows.Scan(&result.ID, &result.ShortURIID, &result.UserID)
		if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgsu *shortURIUserPgRepo) ListAllByShortURI(shortURIID string) ([]*_model.ShortURIUser, error) {
	res := make([]*_model.ShortURIUser, 0)
	rows, err := pgsu.db.GetDB().Query(listShortURIUserAllByShortURISQL, shortURIID)
	if err != nil {
		return res, err
	}
	defer _utl.CloseOnly(rows)

	for rows.Next() {
		var result = _model.ShortURIUser{}
		// id, short_uri_id, user_id
		err := rows.Scan(&result.ID, &result.ShortURIID, &result.UserID)
		if err != nil {
			return nil, err
		}

		res = append(res, &result)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (pgsu *shortURIUserPgRepo) Create(entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	find, err := pgsu.GetByUnique(entity.UserID, entity.ShortURIID)
	if err != nil {
		return nil, err
	}
	if find != nil {
		return find, errors.New("entity already exists")
	}

	stmt, err := pgsu.db.GetDB().Prepare(createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmt)

	res, err := internalCreateShortURIUser(stmt, entity)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (pgsu *shortURIUserPgRepo) CreateTran(tx *sql.Tx, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	if err := _model.ValidateShortURIUser(entity); err != nil {
		return nil, err
	}

	find, err := pgsu.GetByUnique(entity.UserID, entity.ShortURIID)
	if err != nil {
		return nil, err
	}
	if find != nil {
		return find, errors.New("entity already exists")
	}

	stmt, err := tx.Prepare(createShortURIUserSQL)
	if err != nil {
		return nil, err
	}
	defer _utl.CloseOnly(stmt)

	res, err := internalCreateShortURIUser(stmt, entity)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func internalCreateShortURIUser(stmt *sql.Stmt, entity *_model.ShortURIUser) (*_model.ShortURIUser, error) {
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	entity.ID = newID.String()

	_, err = stmt.Exec(entity.ID, entity.ShortURIID, entity.UserID)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
