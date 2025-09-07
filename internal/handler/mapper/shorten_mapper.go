package mapper

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_dto "github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
)

func ResponseFromKey(key string) (*_dto.ShortenCreateResponse, error) {
	if key == "" {
		return nil, nil
	}

	return &_dto.ShortenCreateResponse{
		Result: _utl.BuildNewURI(_cfg.AppConfig.BaseURL, key),
	}, nil
}

func ResponseFromEntity(entity *_model.ShortURI) (*_dto.ShortenCreateResponse, error) {
	if entity == nil {
		return nil, nil
	}

	return ResponseFromKey(entity.ID)
}
