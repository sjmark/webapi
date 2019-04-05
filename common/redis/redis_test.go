package redis

import (
	"net"
	"testing"
	"time"
)

// These tests require redis webapi running on localhost:6379 (the default)
const redisTestServer = "127.0.0.1:6379"

var newRedisCache = func(t *testing.T, defaultExpiration time.Duration) Cache {

	c, err := net.Dial("tcp", redisTestServer)
	if err == nil {
		c.Write([]byte("flush_all\r\n"))
		c.Close()
		redisCache := NewRedisCache("1", redisTestServer, "", defaultExpiration)
		redisCache.Flush()
		return redisCache
	}
	t.Errorf("couldn't connect to redis on %s", redisTestServer)
	t.FailNow()
	panic("")
}

func TestRedisCache_TypicalGetSet(t *testing.T) {
	typicalGetSet(t, newRedisCache)
}

func TestRedisCache_IncrDecr(t *testing.T) {
	incrDecr(t, newRedisCache)
}

func TestRedisCache_Expiration(t *testing.T) {
	expiration(t, newRedisCache)
}

func TestRedisCache_EmptyCache(t *testing.T) {
	emptyCache(t, newRedisCache)
}

func TestRedisCache_Replace(t *testing.T) {
	testReplace(t, newRedisCache)
}

func TestRedisCache_Add(t *testing.T) {
	testAdd(t, newRedisCache)
}

func TestRedisCache_GetMulti(t *testing.T) {
	testGetMulti(t, newRedisCache)
}
