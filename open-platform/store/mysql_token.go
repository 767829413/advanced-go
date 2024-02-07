package store

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/767829413/advanced-go/open-platform/models"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

// StoreItem data item
type StoreItem struct {
	ID        int64  `gorm:"primaryKey;column:id" json:"id"`
	ExpiredAt int64  `gorm:"column:expired_at"    json:"expiredAt"`
	Code      string `gorm:"column:code"          json:"code"` // 255
	Access    string `gorm:"column:access"        json:"access"`
	Refresh   string `gorm:"column:refresh"       json:"refresh"`
	Data      string `gorm:"column:data"          json:"data"` // 2048
}

// TableName get sql table name.获取数据库表名
func (m *StoreItem) TableName() string {
	return "oauth2_token"
}

// Store mysql token store
type TokenStoreMysql struct {
	tableName string
	db        *gorm.DB
	stdout    io.Writer
	ticker    *time.Ticker
}

// SetStdout set error output
func (s *TokenStoreMysql) SetStdout(stdout io.Writer) *TokenStoreMysql {
	s.stdout = stdout
	return s
}

// Close close the store
func (s *TokenStoreMysql) Close() {
	s.ticker.Stop()
}

func (s *TokenStoreMysql) gc() {
	for range s.ticker.C {
		s.clean()
	}
}

func (s *TokenStoreMysql) clean() {
	var n int64
	err := s.db.Table(s.tableName).
		Table(s.tableName).
		Where("expired_at<= ? OR (code='' AND access='' AND refresh='')", time.Now().Unix()).
		Count(&n).Error
	if err != nil || n == 0 {
		if err != nil {
			s.errorf(err.Error())
		}
		return
	}

	err = s.db.Exec(
		"DELETE FROM ? WHERE expired_at<= ? OR (code='' AND access='' AND refresh='')",
		s.tableName,
		time.Now().Unix(),
	).Error
	if err != nil {
		s.errorf(err.Error())
	}
}

func (s *TokenStoreMysql) errorf(format string, args ...interface{}) {
	if s.stdout != nil {
		buf := fmt.Sprintf("[OAUTH2-MYSQL-ERROR]: "+format, args...)
		_, _ = s.stdout.Write([]byte(buf))
	}
}

// Create create and store the new token information
func (s *TokenStoreMysql) Create(ctx context.Context, info models.TokenInfo) error {
	buf, _ := jsoniter.Marshal(info)
	item := &StoreItem{
		Data: string(buf),
	}

	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Unix()
	} else {
		item.Access = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()

		if refresh := info.GetRefresh(); refresh != "" {
			item.Refresh = info.GetRefresh()
			item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Unix()
		}
	}

	return s.db.Create(&item).Error
}

// RemoveByCode delete the authorization code
func (s *TokenStoreMysql) RemoveByCode(ctx context.Context, code string) error {
	err := s.db.Exec("UPDATE ? SET code='' WHERE code=? LIMIT 1", s.tableName, code).Error
	return err
}

// RemoveByAccess use the access token to delete the token information
func (s *TokenStoreMysql) RemoveByAccess(ctx context.Context, access string) error {
	err := s.db.Exec("UPDATE ? SET access='' WHERE access=? LIMIT 1", s.tableName, access).Error
	return err
}

// RemoveByRefresh use the refresh token to delete the token information
func (s *TokenStoreMysql) RemoveByRefresh(ctx context.Context, refresh string) error {
	err := s.db.Exec("UPDATE ? SET refresh='' WHERE refresh=? LIMIT 1", s.tableName, refresh).Error
	return err
}

func (s *TokenStoreMysql) toTokenInfo(data string) models.TokenInfo {
	var tm models.Token
	_ = jsoniter.Unmarshal([]byte(data), &tm)
	return &tm
}

// GetByCode use the authorization code for token information data
func (s *TokenStoreMysql) GetByCode(ctx context.Context, code string) (models.TokenInfo, error) {
	if code == "" {
		return nil, nil
	}

	var item StoreItem
	err := s.db.Table(s.tableName).
		Select("*").
		Where("code=?", code).
		Limit(1).
		First(&item).Error
	if err != nil {
		return nil, nil
	}
	return s.toTokenInfo(item.Data), nil
}

// GetByAccess use the access token for token information data
func (s *TokenStoreMysql) GetByAccess(
	ctx context.Context,
	access string,
) (models.TokenInfo, error) {
	if access == "" {
		return nil, nil
	}

	var item StoreItem
	err := s.db.Table(s.tableName).
		Select("*").
		Where("access=?", access).
		Limit(1).
		First(&item).Error
	if err != nil {
		return nil, nil
	}
	return s.toTokenInfo(item.Data), nil
}

// GetByRefresh use the refresh token for token information data
func (s *TokenStoreMysql) GetByRefresh(
	ctx context.Context,
	refresh string,
) (models.TokenInfo, error) {
	if refresh == "" {
		return nil, nil
	}

	var item StoreItem
	err := s.db.Table(s.tableName).
		Select("*").
		Where("refresh=?", refresh).
		Limit(1).
		First(&item).Error
	if err != nil {
		return nil, nil
	}
	return s.toTokenInfo(item.Data), nil
}
