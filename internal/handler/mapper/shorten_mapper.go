package mapper

import (
	_cfg "github.com/ElfAstAhe/url-shortener/internal/config"
	_dto "github.com/ElfAstAhe/url-shortener/internal/handler/dto"
	_model "github.com/ElfAstAhe/url-shortener/internal/model"
	_srv "github.com/ElfAstAhe/url-shortener/internal/service"
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
	if len(source) == 0 {
		return make([]*_dto.ShortenBatchResponseItem, 0), nil
	}
	res := make([]*_dto.ShortenBatchResponseItem, 0)
	for key, value := range source {
		res = append(res, &_dto.ShortenBatchResponseItem{
			CorrelationID: key,
			ShortURL:      value,
		})
	}

	return res, nil
}

func ShortenBatchResponseFromEntity(source map[string]*_model.ShortURI) ([]*_dto.ShortenBatchResponseItem, error) {
	if len(source) == 0 {
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

func ShortenBatchFromDto(source []*_dto.ShortenBatchCreateItem) (map[string]string, error) {
	res := make(map[string]string)
	if len(source) == 0 {
		return res, nil
	}

	for _, item := range source {
		res[item.CorrelationID] = item.OriginalURL
	}

	return res, nil
}

func UserShortensFromModel(source map[string]string) ([]*_dto.UserShorten, error) {
	res := make([]*_dto.UserShorten, 0)
	if len(source) == 0 {
		return nil, nil
	}
	for key, value := range source {
		res = append(res, _dto.NewUserShorten(value, key))
	}

	return res, nil
}

func UserBatchDeletesFromDto(source _dto.ShortenBatchDeleteRequest) (_srv.UserBatchDeletes, error) {
	res := make(_srv.UserBatchDeletes, 0)
	if len(source) == 0 {
		return res, nil
	}

	for _, item := range source {
		res = append(res, item)
	}

	return res, nil
}
