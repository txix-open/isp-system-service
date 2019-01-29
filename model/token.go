package model

import (
	"isp-system-service/entity"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type TokenRepository struct {
	DB orm.DB
}

var (
	emptyToken = &entity.Token{}
	TokenRep   TokenRepository
)

func (tr *TokenRepository) SaveToken(token entity.Token) (entity.Token, error) {
	_, err := tr.DB.Model(&token).Insert()
	return token, err
}

func (tr *TokenRepository) GetTokenById(id string) (*entity.Token, error) {
	token := &entity.Token{Token: id}
	err := tr.DB.Select(token)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return token, err
}

func (tr *TokenRepository) GetTokensByAppId(appId ...int32) ([]entity.Token, error) {
	var res = make([]entity.Token, 0)
	err := tr.DB.Model(&res).Where("app_id IN (?)", pg.In(appId)).Order("created_at DESC").Select()
	return res, err
}

func (tr *TokenRepository) DeleteTokens(tokens []string) (int, error) {
	res, err := tr.DB.Model(emptyToken).Where("token IN (?)", pg.In(tokens)).Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
