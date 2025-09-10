package mapper

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_dto "github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_utl "github.com/ElfAstAhe/url-shortener/internal/utils"
)

func ShortenCreateResponseFromKey(key string) (*_dto.ShortenCreateResponse, error) {
	if key == "" {
		return nil, nil
	}

	return &_dto.ShortenCreateResponse{
		Result: _utl.BuildNewURI(_cfg.AppConfig.BaseURL, key),
	}, nil
}

func ShortenCreateResponseFromEntity(entity *_model.ShortURI) (*_dto.ShortenCreateResponse, error) {
	if entity == nil {
		return nil, nil
	}

	return ShortenCreateResponseFromKey(entity.Key)
}

func ShortenBatchResponseFromKeys(source map[string]string) ([]*_dto.ShortenBatchResponseItem, error) {
	if source == nil {
		return make([]*_dto.ShortenBatchResponseItem, 0), nil
	}
	res := make([]*_dto.ShortenBatchResponseItem, len(source))
	for key, value := range source {
		res = append(res, &_dto.ShortenBatchResponseItem{
			CorrelationID: key,
			ShortURL:      value,
		})
	}

	return res, nil
}

func ShortenBatchResponseFromEntity(source map[string]*_model.ShortURI) ([]*_dto.ShortenBatchResponseItem, error) {
	if source == nil {
		return make([]*_dto.ShortenBatchResponseItem, 0), nil
	}
	res := make([]*_dto.ShortenBatchResponseItem, 0, len(source))
	for key, value := range source {
		res = append(res, &_dto.ShortenBatchResponseItem{
			CorrelationID: key,
			ShortURL:      _utl.BuildNewURI(_cfg.AppConfig.BaseURL, value.Key),
		})
	}

	return res, nil
}
