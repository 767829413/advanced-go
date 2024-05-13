package cache_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/767829413/advanced-go/cache"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Role []int  `json:"role"`
	Name string `json:"name"`
}

func init() {
	options := &redis.Options{
		Addr:     "",
		Username: "",
		Password: "",
		DB:       1,
	}
	rdb := redis.NewClient(options)
	db, err := gorm.Open(
		mysql.Open(
			"user:passd@tcp(localhost:3306)/cache?charset=utf8mb4&parseTime=true&loc=Local",
		),
		&gorm.Config{
			Logger: logger.Default,
		},
	)
	if err != nil {
		log.Panic("gorm.Open error: ", err)
	}
	cache.InitByRedisMysql(rdb, db)
}

// 测试基本能力(缓存失效和删除缓存)
func TestCacheOperations(t *testing.T) {
	cacheManagerIns := cache.GetCacheManager()

	testData := map[string]string{"key1": "value1"}
	expectedByte, _ := json.Marshal(testData)
	expected := string(expectedByte)

	err := cacheManagerIns.Set("key1", testData, 2*time.Second)
	if err != nil {
		log.Panic("cache Set error: ", err)
	}
	time.Sleep(1 * time.Second)
	v, isExist := cacheManagerIns.Get("key1")
	log.Println("sleep 1 second ,get key1 value", v, isExist)
	assert.Equal(t, expected, v)
	assert.Equal(t, true, isExist)
	time.Sleep(2 * time.Second)

	v, isExist = cacheManagerIns.Get("key1")
	log.Println("sleep 3 second ,get key1 value", v, isExist)
	assert.Equal(t, "", v)
	assert.Equal(t, false, isExist)

	testData2 := []string{"value1"}
	expectedByte2, _ := json.Marshal(testData2)
	expected2 := string(expectedByte2)

	err = cacheManagerIns.Set("key2", []string{"value1"}, 5*time.Second)
	assert.Nil(t, err)

	v, isExist = cacheManagerIns.Get("key2")
	assert.Equal(t, expected2, v)
	assert.Equal(t, true, isExist)

	log.Println("get key2 value", v, isExist)

	err = cacheManagerIns.Del("key2")
	assert.Nil(t, err)
	v, isExist = cacheManagerIns.Get("key2")
	assert.Equal(t, "", v)
	assert.Equal(t, false, isExist)

	log.Println("delete key2 ,try get key2 value", v, isExist)
}

// Test 1: Test Setting Unsupported Data Type
func TestCacheSetUnsupportedDataType(t *testing.T) {
	ch := make(chan int)
	err := cache.GetCacheManager().Set("chanKey", ch, 1*time.Minute)
	assert.NotNil(t, err, "Expected an error when setting unsupported data type")
}

// Test 2: Test Getting Value with Incorrect Type Assertion
func TestCacheGetWithIncorrectTypeAssertion(t *testing.T) {
	testData := "stringValue"
	err := cache.GetCacheManager().Set("stringKey", testData, 1*time.Minute)
	assert.Nil(t, err)

	_, isExist := cache.GetCacheManager().Get("stringKey")
	assert.True(t, isExist, "Expected the key to exist")

	// Attempt to get the value with incorrect type assertion
	var intValue int
	v, isExist := cache.GetCacheManager().Get("stringKey")
	assert.True(t, isExist, "Expected the key to exist")
	err = json.Unmarshal([]byte(v), &intValue)
	assert.NotNil(t, err, "Expected an error due to incorrect type assertion")
}

// Test 3: Test Setting and Getting Struct without JSON Tags
func TestCacheStructWithoutJSONTags(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	user := User{Name: "John", Age: 30}
	err := cache.GetCacheManager().Set("userKey", user, 1*time.Minute)
	assert.Nil(t, err)

	v, isExist := cache.GetCacheManager().Get("userKey")
	assert.True(t, isExist, "Expected the key to exist")

	var retrievedUser User
	err = json.Unmarshal([]byte(v), &retrievedUser)
	assert.Nil(t, err, "Expected no error during unmarshal")
	assert.Equal(t, user, retrievedUser, "Expected retrieved struct to match the original")
}

// Test 4: Test Setting Complex Nested Structures
func TestCacheComplexNestedStructures(t *testing.T) {
	type Nested struct {
		Field1 string
		Field2 int
	}
	type Complex struct {
		Nested Nested
		Flag   bool
	}
	complexData := Complex{Nested: Nested{Field1: "data", Field2: 42}, Flag: true}
	err := cache.GetCacheManager().Set("complexKey", complexData, 1*time.Minute)
	assert.Nil(t, err)

	v, isExist := cache.GetCacheManager().Get("complexKey")
	assert.True(t, isExist, "Expected the key to exist")

	var retrievedComplex Complex
	err = json.Unmarshal([]byte(v), &retrievedComplex)
	assert.Nil(t, err, "Expected no error during unmarshal")
	assert.Equal(
		t,
		complexData,
		retrievedComplex,
		"Expected retrieved complex structure to match the original",
	)
}

// Test 5: Test Setting Nil Value
func TestCacheSetNilValue(t *testing.T) {
	cache.GetCacheManager().Set("nilKey", nil, 1*time.Minute)
	v, isExist := cache.GetCacheManager().Get("nilKey")
	assert.Equal(
		t,
		"",
		v,
		"Expected retrieved empty string",
	)
	assert.True(t, isExist, "Expected the key to exist")
}
