package store

import (
	"context"
	"fmt"
	"github.com/767829413/advanced-go/open-platform/generates"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var (
	jsonMarshal   = jsoniter.Marshal
	jsonUnmarshal = jsoniter.Unmarshal
)

// NewRedisStoreWithCli create an instance of a redis store
func NewRedisStoreWithCli(cli *redis.Client, keyNamespace ...string) TokenStore {
	return NewRedisStoreWithInterface(cli, keyNamespace...)
}

// NewRedisClusterStore create an instance of a redis cluster store
func NewRedisClusterStore(opts *redis.ClusterOptions, keyNamespace ...string) TokenStore {
	if opts == nil {
		panic("options cannot be nil")
	}
	return NewRedisClusterStoreWithCli(redis.NewClusterClient(opts), keyNamespace...)
}

// NewRedisClusterStoreWithCli create an instance of a redis cluster store
func NewRedisClusterStoreWithCli(cli *redis.ClusterClient, keyNamespace ...string) TokenStore {
	return NewRedisStoreWithInterface(cli, keyNamespace...)
}

// NewRedisStoreWithInterface create an instance of a redis store
func NewRedisStoreWithInterface(cli clienter, keyNamespace ...string) TokenStore {
	store := &TokenStoreRedis{
		cli: cli,
	}

	if len(keyNamespace) > 0 {
		store.ns = keyNamespace[0]
	}
	return store
}

type clienter interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Exists(ctx context.Context, key ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

// TokenStore redis token store
type TokenStoreRedis struct {
	cli clienter
	ns  string
}

// Close close the store
func (s *TokenStoreRedis) Close() error {
	return s.cli.Close()
}

func (s *TokenStoreRedis) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.ns, key)
}

func (s *TokenStoreRedis) checkError(result redis.Cmder) (bool, error) {
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// remove
func (s *TokenStoreRedis) remove(ctx context.Context, key string) error {
	result := s.cli.Del(ctx, s.wrapperKey(key))
	_, err := s.checkError(result)
	return err
}

func (s *TokenStoreRedis) removeToken(
	ctx context.Context,
	tokenString string,
	isRefresh bool,
) error {
	basicID, err := s.getBasicID(ctx, tokenString)
	if err != nil {
		return err
	} else if basicID == "" {
		return nil
	}

	err = s.remove(ctx, tokenString)
	if err != nil {
		return err
	}

	token, err := s.getToken(ctx, basicID)
	if err != nil {
		return err
	} else if token == nil {
		return nil
	}

	checkToken := token.GetRefresh()
	if isRefresh {
		checkToken = token.GetAccess()
	}
	iresult := s.cli.Exists(ctx, s.wrapperKey(checkToken))
	if err := iresult.Err(); err != nil && err != redis.Nil {
		return err
	} else if iresult.Val() == 0 {
		return s.remove(ctx, basicID)
	}

	return nil
}

func (s *TokenStoreRedis) parseToken(result *redis.StringCmd) (generates.TokenInfo, error) {
	if ok, err := s.checkError(result); err != nil {
		return nil, err
	} else if ok {
		return nil, nil
	}

	buf, err := result.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var token generates.Token
	if err := jsonUnmarshal(buf, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *TokenStoreRedis) getToken(ctx context.Context, key string) (generates.TokenInfo, error) {
	result := s.cli.Get(ctx, s.wrapperKey(key))
	return s.parseToken(result)
}

func (s *TokenStoreRedis) parseBasicID(result *redis.StringCmd) (string, error) {
	if ok, err := s.checkError(result); err != nil {
		return "", err
	} else if ok {
		return "", nil
	}
	return result.Val(), nil
}

func (s *TokenStoreRedis) getBasicID(ctx context.Context, token string) (string, error) {
	result := s.cli.Get(ctx, s.wrapperKey(token))
	return s.parseBasicID(result)
}

// Create Create and store the new token information
func (s *TokenStoreRedis) Create(ctx context.Context, info generates.TokenInfo) error {
	ct := time.Now()
	jv, err := jsonMarshal(info)
	if err != nil {
		return err
	}

	pipe := s.cli.TxPipeline()
	if code := info.GetCode(); code != "" {
		pipe.Set(ctx, s.wrapperKey(code), jv, info.GetCodeExpiresIn())
	} else {
		basicID := uuid.Must(uuid.NewRandom()).String()
		aexp := info.GetAccessExpiresIn()
		rexp := aexp

		if refresh := info.GetRefresh(); refresh != "" {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
			if aexp.Seconds() > rexp.Seconds() {
				aexp = rexp
			}
			pipe.Set(ctx, s.wrapperKey(refresh), basicID, rexp)
		}

		pipe.Set(ctx, s.wrapperKey(info.GetAccess()), basicID, aexp)
		pipe.Set(ctx, s.wrapperKey(basicID), jv, rexp)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

// RemoveByCode Use the authorization code to delete the token information
func (s *TokenStoreRedis) RemoveByCode(ctx context.Context, code string) error {
	return s.remove(ctx, code)
}

// RemoveByAccess Use the access token to delete the token information
func (s *TokenStoreRedis) RemoveByAccess(ctx context.Context, access string) error {
	return s.removeToken(ctx, access, false)
}

// RemoveByRefresh Use the refresh token to delete the token information
func (s *TokenStoreRedis) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.removeToken(ctx, refresh, true)
}

// GetByCode Use the authorization code for token information data
func (s *TokenStoreRedis) GetByCode(ctx context.Context, code string) (generates.TokenInfo, error) {
	return s.getToken(ctx, code)
}

// GetByAccess Use the access token for token information data
func (s *TokenStoreRedis) GetByAccess(
	ctx context.Context,
	access string,
) (generates.TokenInfo, error) {
	basicID, err := s.getBasicID(ctx, access)
	if err != nil || basicID == "" {
		return nil, err
	}
	return s.getToken(ctx, basicID)
}

// GetByRefresh Use the refresh token for token information data
func (s *TokenStoreRedis) GetByRefresh(
	ctx context.Context,
	refresh string,
) (generates.TokenInfo, error) {
	basicID, err := s.getBasicID(ctx, refresh)
	if err != nil || basicID == "" {
		return nil, err
	}
	return s.getToken(ctx, basicID)
}