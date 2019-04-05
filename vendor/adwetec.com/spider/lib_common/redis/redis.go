package redis

import (
	"fmt"

	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

//改自revel框架的rediscache
//http://redis.readthedocs.org/en/2.6/
var (
	ErrCacheMiss = errors.New("rediscache: key not found")
	ErrNotStored = errors.New("rediscache: not stored")
	DEFAULT      = time.Duration(0)
	FOREVER      = time.Duration(-1)
)

//Cache Wraps the Redis client to meet the Cache interface.
type Cache struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

// Getter the content associated with the given key. decoding it into the given
// pointer.
//
// Returns:
//   - nil if the value was successfully retrieved and ptrValue set
//   - ErrCacheMiss if the value was not in the cache
//   - an implementation specific error otherwise
type Getter interface {
	Get(key string, ptrValue interface{}) error
}

// NewRedisCache until redigo supports sharding/clustering, only one host will be in hostList
func NewRedisCache(db, host, password string, defaultExpiration time.Duration) Cache {

	var pool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   0,
		IdleTimeout: time.Duration(240) * time.Second,
		Dial: func() (redis.Conn, error) {
			protocol := "tcp"
			toc := time.Millisecond * time.Duration(10000)
			tor := time.Millisecond * time.Duration(500000)
			tow := time.Millisecond * time.Duration(5000)
			c, err := redis.DialTimeout(protocol, host, toc, tor, tow)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					c.Close()
					return nil, err
				}
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}

			if _, err := c.Do("SELECT", db); err != nil {
				return err
			}

			return nil
		},
	}

	return Cache{pool, defaultExpiration}
}

//Set 设置字符串值，expires为过期时间
func (c Cache) Set(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	return c.invoke(conn.Do, key, value, expires)
}

//Add 添加 若key不存在则报错
func (c Cache) Add(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	existed, err := exists(conn, key)
	if err != nil {
		return err
	} else if existed {
		return ErrNotStored
	}
	return c.invoke(conn.Do, key, value, expires)
}

//Exists 判断key是否存在
func (c Cache) Exists(key string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("EXISTS", key))
}

//Replace 替换key值
func (c Cache) Replace(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	existed, err := exists(conn, key)
	if err != nil {
		return err
	} else if !existed {
		return ErrNotStored
	}
	err = c.invoke(conn.Do, key, value, expires)
	if value == nil {
		return ErrNotStored
	}
	return err
}

//Get 获取key值
func (c Cache) Get(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("GET", key)
	if err != nil {
		return err
	} else if raw == nil {
		return ErrCacheMiss
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}
	return Deserialize(item, ptrValue)
}

//GetStr 获取key值string类型
func (c Cache) GetStr(key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	} else if raw == nil {
		return "", ErrCacheMiss
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return "", err
	}
	return string(item), nil
}

func generalizeStringSlice(strs []string) []interface{} {
	ret := make([]interface{}, len(strs))
	for i, str := range strs {
		ret[i] = str
	}
	return ret
}

//GetMulti 获取多个key的值
func (c Cache) GetMulti(keys ...string) (Getter, error) {
	conn := c.pool.Get()
	defer conn.Close()

	items, err := redis.Values(conn.Do("MGET", generalizeStringSlice(keys)...))
	if err != nil {
		return nil, err
	} else if items == nil {
		return nil, ErrCacheMiss
	}

	m := make(map[string][]byte)
	for i, key := range keys {
		m[key] = nil
		if i < len(items) && items[i] != nil {
			s, ok := items[i].([]byte)
			if ok {
				m[key] = s
			}
		}
	}
	return ItemMapGetter(m), nil
}

func exists(conn redis.Conn, key string) (bool, error) {
	return redis.Bool(conn.Do("EXISTS", key))
}

//Delete 删除key值
func (c Cache) Delete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	existed, err := redis.Bool(conn.Do("DEL", key))
	if err == nil && !existed {
		err = ErrCacheMiss
	}
	return err
}

//Increment 自增
func (c Cache) Increment(key string, delta uint64) (uint64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	// Check for existance *before* increment as per the cache contract.
	// redis will auto create the key, and we don't want that. Since we need to do increment
	// ourselves instead of natively via INCRBY (redis doesn't support wrapping), we get the value
	// and do the exists check this way to minimize calls to Redis
	val, err := conn.Do("GET", key)
	if err != nil {
		return 0, err
	} else if val == nil {
		return 0, ErrCacheMiss
	}
	currentVal, err := redis.Int64(val, nil)
	if err != nil {
		return 0, err
	}
	sum := currentVal + int64(delta)
	_, err = conn.Do("SET", key, sum)
	if err != nil {
		return 0, err
	}
	return uint64(sum), nil
}

//Decrement 自减
func (c Cache) Decrement(key string, delta uint64) (newValue uint64, err error) {
	conn := c.pool.Get()
	defer conn.Close()
	// Check for existance *before* increment as per the cache contract.
	// redis will auto create the key, and we don't want that, hence the exists call
	existed, err := exists(conn, key)
	if err != nil {
		return 0, err
	} else if !existed {
		return 0, ErrCacheMiss
	}
	// Decrement contract says you can only go to 0
	// so we go fetch the value and if the delta is greater than the amount,
	// 0 out the value
	currentVal, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}
	if delta > uint64(currentVal) {
		tempint, err := redis.Int64(conn.Do("DECRBY", key, currentVal))
		return uint64(tempint), err
	}
	tempint, err := redis.Int64(conn.Do("DECRBY", key, delta))
	return uint64(tempint), err
}

//Flush 关闭
func (c Cache) Flush() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

func (c Cache) invoke(f func(string, ...interface{}) (interface{}, error), key string, value interface{}, expires time.Duration) error {

	switch expires {
	case DEFAULT:
		expires = c.defaultExpiration
	case FOREVER:
		expires = time.Duration(0)
	}

	b, err := Serialize(value)
	if err != nil {
		return err
	}
	conn := c.pool.Get()
	defer conn.Close()
	if expires > 0 {
		_, err := f("SETEX", key, int32(expires/time.Second), b)
		return err
	}
	_, err = f("SET", key, b)
	return err
}

//ItemMapGetter Implement a Getter on top of the returned item map.
type ItemMapGetter map[string][]byte

//Get 获取key值
func (g ItemMapGetter) Get(key string, ptrValue interface{}) error {
	item, ok := g[key]
	if !ok {
		return ErrCacheMiss
	}
	if len(item) == 0 {
		return ErrCacheMiss
	}
	return Deserialize(item, ptrValue)
}

//Lpush redis queue model source value
//BRPOPLPUSH source distination 300s-> value
//ExecEnd then LREM distination 1 value
func (c Cache) Lpush(key string, value interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(value)
	if err != nil {
		return err
	}
	_, err = conn.Do("LPUSH", key, b)
	return err
}

//BRpopLpush 移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止
func (c Cache) BRpopLpush(key, bkkey string, second int, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("BRPOPLPUSH", key, bkkey, second)
	if err != nil {
		return err
	} else if raw == nil {
		return nil
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}
	return Deserialize(item, ptrValue)
}

//Lpop 移除并返回列表的第一个元素
func (c Cache) Lpop(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("LPOP", key)
	// if isaffected == 0 {
	// 	// fmt.Println("LPOP:", key, isaffected, ",err:", err)
	// }
	return err
}

//Lrange 返回列表中指定区间内的元素，区间以偏移量begin和end指定
//其中 0 表示列表的第一个元素
//以 -1 表示列表的最后一个元素
func (c Cache) Lrange(key string, begin int, end int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("LRANGE", key, begin, end)
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

//Llen 返回列表的长度
//如果列表 key 不存在，则 key 被解释为一个空列表，返回 0
//如果 key 不是列表类型，返回一个错误
func (c Cache) Llen(key string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("LLEN", key)
	if err != nil {
		return 0, err
	} else if raw == nil {
		return 0, nil
	}
	return redis.Int(raw, err)
}

//Rpop 移除并返回列表的最后一个元素
func (c Cache) Rpop(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	raw, err := conn.Do("RPOP", key)
	if err != nil {
		return err
	} else if raw == nil {
		return nil
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}
	return Deserialize(item, ptrValue)
}

//Lrem 根据参数 COUNT 的值，移除列表中与参数 VALUE 相等的元素。
//count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
//count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
//count = 0 : 移除表中所有与 VALUE 相等的值。
func (c Cache) Lrem(key string, value interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(value)
	if err != nil {
		return err
	}
	//count 设置为1 从表头开始搜索
	isaffected, err := conn.Do("LREM", key, 1, b)
	if isaffected == 0 {
		// fmt.Println("LREM:", key, isaffected, ",err:", err, ",b:", b)
	}
	return err
}

//Sadd 向集合添加一个元素
func (c Cache) Sadd(key string, value interface{}) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Sadd", key, value)
	if err != nil {
		return 0, err
	}
	return redis.Int64(raws, err)
}

//Sadds 向集合添加多个元素
func (c Cache) Sadds(fieldValue ...interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Sadd", fieldValue...)
	if err != nil {
		return false, err
	}
	return redis.Bool(raws, err)
}

//Sismember 判断 member 元素是否是集合 key 的成员
func (c Cache) Sismember(key string, value interface{}) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Sismember", key, value)
	if err != nil {
		return "", err
	}
	return redis.String(raws, err)
}

//Smembers 返回集合中的所有的成员 string切片
//不存在的集合 key 被视为空集合
func (c Cache) Smembers(key string) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("SMEMBERS", key)
	if err != nil {
		return nil, err
	}
	return redis.Strings(raws, err)
}

//SmembersInterface 返回集合中的所有的成员
//不存在的集合 key 被视为空集合
func (c Cache) SmembersInterface(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("SMEMBERS", key)
	if err != nil {
		return err
	}

	item, err := redis.Bytes(raws, err)
	if err != nil {
		return err
	}

	return Deserialize(item, ptrValue)
}

//Srem 移除集合中一个元素
func (c Cache) Srem(key string, value interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("Srem", key, value))
}

//Scard 获取集合的成员数
func (c Cache) Scard(key string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Scard", key)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Spop 移除并返回集合中的一个随机元素
func (c Cache) Spop(key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Spop", key)
	if err != nil {
		return "", err
	}
	return redis.String(raws, err)
}

//Zadd 向有序集合添加一个或多个成员，或者更新已存在成员的分数
func (c Cache) Zadd(key string, value interface{}, score int64) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGEBYSCORE", key, score, score)
	if err != nil {
		return false, err
	}
	switch raws := raws.(type) {
	case []interface{}:
		if len(raws) > 0 {
			return false, nil
		}
	}
	b, err := Serialize(value)
	if err != nil {
		return false, err
	}
	reply, err := conn.Do("ZADD", key, score, b)
	return redis.Bool(reply, err)
}

func (c Cache) ZaddNoserialize(key string, value interface{}, score int64) error {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGEBYSCORE", key, score, score)
	if err != nil {
		return err
	}
	switch raws := raws.(type) {
	case []interface{}:
		if len(raws) > 0 {
			return nil
		}
	}
	_, err = conn.Do("ZADD", key, score, value)
	return err
}

//ZaddJSON 向有序集合添加一个或多个JSON成员
func (c Cache) ZaddJSON(key string, value []byte, score int64) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", key, score, value)
	if err != nil {
		return false, err
	}
	return true, err
}

//ZaddUpgrade 更新已存在成员的分数
func (c Cache) ZaddUpgrade(key string, value []byte, score int64) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGEBYSCORE", key, score, score)
	if err != nil {
		return false, err
	}
	switch raws := raws.(type) {
	case []interface{}:
		if len(raws) > 0 {
			return false, nil
		}
	}
	_, err = conn.Do("ZADD", key, score, value)
	if err != nil {
		return false, err
	}
	return true, err
}

//Zrangewithscores 返回有序集中，指定区间内的成员。
//其中成员的位置按分数值递增(从小到大)来排序。
//具有相同分数值的成员按字典序(lexicographical order )来排列。
func (c Cache) Zrangewithscores(key string, begin, end int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGE", key, begin, end, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

//ZaddMultiScore 添加多个元素
func (c Cache) ZaddMultiScore(key string, value interface{}, score int64) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(value)
	if err != nil {
		return false, err
	}
	_, err = conn.Do("ZADD", key, score, b)
	if err != nil {
		return false, err
	}
	return true, err
}

//ZaddInt64 添加int64
func (c Cache) ZaddInt64(key string, value int64, score int64) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", key, score, value)
	if err != nil {
		return false, err
	}
	return true, err
}

//ZaddMultiInt64 添加多个int64
func (c Cache) ZaddMultiInt64(key string, score_value ...int64) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", key, score_value)
	if err != nil {
		return false, err
	}
	return true, err
}

// ZaddMultiInt64  如："k", 1, 2, 3, 4  注：k 为key， 1 为score 2为值  to....
func (c Cache) ZaddMulti(val []interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", val...)
	return err
}

//Zincrby 有序集合中对指定成员的分数加上增量 increment
func (c Cache) Zincrby(key, member string, inc float64) (float64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZINCRBY", key, inc, member)
	if err != nil {
		return 0.0, err
	}
	return redis.Float64(raws, err)
}

//ZincrInt64by 有序集合中对指定成员的分数加上增量 increment
func (c Cache) ZincrInt64by(key string, member int64, inc int64) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZINCRBY", key, inc, member)
	if err != nil {
		return 0, err
	}
	return redis.Int64(raws, err)
}

//Zscore 返回有序集中，成员的分数值
func (c Cache) Zscore(key string, member interface{}) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(member)
	if err != nil {
		return 0, err
	}
	raws, err := conn.Do("ZSCORE", key, b)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Zscoreforstring 返回有序集中，成员的分数值 string型成员
func (c Cache) Zscoreforstring(key string, member string) (float64, error) {
	conn := c.pool.Get()
	defer conn.Close()

	raws, err := conn.Do("ZSCORE", key, member)
	if err != nil {
		return 0.0, err
	}
	return redis.Float64(raws, err)
}

//ZscoreJSON 返回有序集中，成员的分数值 JSON成员
func (c Cache) ZscoreJSON(key string, value []byte) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZSCORE", key, value)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//ZscoreInt 返回有序集中，成员的分数值 int成员
func (c Cache) ZscoreInt64(key string, member int64) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZSCORE", key, member)
	if err != nil {
		return 0, err
	}

	return redis.Int64(raws, err)
}

//Zrank 返回有序集合中指定成员的索引
func (c Cache) Zrank(key string, value interface{}) int {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(value)
	if err != nil {
		return -1
	}
	raws, err := conn.Do("Zrank", key, b)
	rankInt := 0
	if raws == nil {
		rankInt = -1
	} else {
		rankInt = int(raws.(int64))
	}

	return rankInt
}

//ZrankNoSerialize
func (c Cache) Zranknoserialize(key string, value interface{}) int64 {
	conn := c.pool.Get()
	defer conn.Close()

	raws, err := conn.Do("Zrank", key, value)
	if err != nil {
		fmt.Printf("Zranknoserialize err:%#+v:", err)
		return -1
	}
	var rankInt int64
	if raws == nil {
		rankInt = -1
	} else {
		rankInt = raws.(int64)
	}
	return rankInt
}

//Zrevrank 返回有序集合中指定成员的排名，有序集成员按分数值递减(从大到小)排序
func (c Cache) Zrevrank(key string, value interface{}) int {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(value)
	if err != nil {
		return -1
	}
	raws, err := conn.Do("Zrevrank", key, b)
	rankInt := 0
	if raws == nil {
		rankInt = -1
	} else {
		rankInt = int(raws.(int64))
	}

	return rankInt
}

func (c Cache) ZrangeByScoresLimit(key string, min, max interface{}, offset, count int, withScores ...bool) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	var raws interface{}
	var err error
	if len(withScores) > 0 && withScores[0] {
		raws, err = conn.Do("ZRANGEBYSCORE", key, min, max, "LIMIT", offset, count, "WITHSCORES")
	} else {
		raws, err = conn.Do("ZRANGEBYSCORE", key, min, max, "LIMIT", offset, count)
	}
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

//Zrevrankforstring 返回有序集合中指定成员的排名，有序集成员按分数值递减(从大到小)排序 string成员
func (c Cache) Zrevrankforstring(key string, value string) int {
	conn := c.pool.Get()
	defer conn.Close()

	raws, _ := conn.Do("Zrevrank", key, value)
	rankInt := 0
	if raws == nil {
		rankInt = -1
	} else {
		rankInt = int(raws.(int64))
	}

	return rankInt
}

//ZremJSON 移除有序集合中的一个 JSON成员
func (c Cache) ZremJSON(key string, value []byte) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	reply, zremerr := conn.Do("ZREM", key, value)
	str, zremerr := redis.String(reply, zremerr)
	return str == "1", zremerr
}

//Zrem 移除有序集合中的一个
func (c Cache) Zrem(key string, value interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(value)
	if err != nil {
		return false, err
	}
	reply, zremerr := conn.Do("ZREM", key, b)
	str, zremerr := redis.String(reply, zremerr)
	return str == "1", zremerr
}

//Zrembyscore 移除有序集合中给定的分数区间的所有成员
func (c Cache) Zrembyscore(key string, begin int64, end int64) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZREMRANGEBYSCORE", key, begin, end)
	return err
}

//Zrembyrank 移除有序集合中给定的排名区间的所有成员
func (c Cache) Zrembyrank(key string, begin int64, end int64) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZREMRANGEBYRANK", key, begin, end)
	return err
}

//Zrangebyscore 通过分数返回有序集合指定区间内的成员
func (c Cache) Zrangebyscore(key string, begin int64, end int64) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGEBYSCORE", key, begin, end)
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

//Zrangebyscore 通过分数返回有序集合指定区间内的成员和分数
func (c Cache) ZrangebyscoreWithscores(key string, begin int64, end int64) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGEBYSCORE", key, begin, end, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	return raws, nil
}

//Zrangebyscoreforstring 通过分数返回有序集合指定区间内的成员 string类型
func (c Cache) Zrangebyscoreforstring(key string, begin int64, end int64) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGEBYSCORE", key, begin, end)
	if err != nil {
		return nil, err
	}
	objAry := make([]string, 0)
	for _, raw := range raws.([]interface{}) {
		item, err := redis.Bytes(raw, err)
		if err != nil {
			continue
		}
		obj := ""
		Deserialize(item, &obj)
		if obj != "" {
			objAry = append(objAry, obj)
		}
	}
	return objAry, nil
}

//Zrangeforint 通过索引区间返回有序集合成指定区间内的成员 int类型
func (c Cache) Zrangeforint(key string, begin, end int) ([]int, error) {
	raws, err := c.Zrange(key, begin, end)
	if err != nil {
		return nil, err
	}

	objAry := make([]int, 0)
	for _, raw := range raws {
		item, err := redis.Bytes(raw, err)
		if err != nil {
			continue
		}
		obj := 0
		Deserialize(item, &obj)
		if obj != 0 {
			objAry = append(objAry, obj)
		}
	}
	return objAry, nil
}

//Zrangeforstring 通过索引区间返回有序集合成指定区间内的成员 string类型
func (c Cache) Zrangeforstring(key string, begin, end int) ([]string, error) {
	raws, err := c.Zrange(key, begin, end)
	if err != nil {
		return nil, err
	}

	objAry := make([]string, 0)
	for _, raw := range raws {
		item, err := redis.Bytes(raw, err)
		if err != nil {
			continue
		}
		obj := ""
		Deserialize(item, &obj)
		if obj != "" {
			objAry = append(objAry, obj)
		}
	}
	return objAry, nil
}

//Zrange 通过索引区间返回有序集合成指定区间内的成员
func (c Cache) Zrange(key string, begin int, end int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZRANGE", key, begin, end)
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
	// switch raws := raws.(type) {
	// case []interface{}:
	// 	return raws, nil
	// default:
	// 	return nil, nil
	// }
}

//Zrevrange 返回有序集中指定区间内的成员，通过索引，分数从高到底
func (c Cache) Zrevrange(key string, begin int, end int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZREVRANGE", key, begin, end)
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

//ZrevrangeInt64 通过分数返回有序集合指定区间内的成员
func (c Cache) ZrevrangeInt64(key string, begin int, end int) ([]int64, error) {
	raws, err := c.Zrevrange(key, begin, end)
	if err != nil {
		return nil, err
	}

	objAry := make([]int64, 0)
	for _, raw := range raws {
		item, err := redis.Bytes(raw, err)
		if err != nil {
			continue
		}
		obj := int64(0)
		Deserialize(item, &obj)
		if obj != 0 {
			objAry = append(objAry, obj)
		}
	}

	return objAry, nil
}

//ZrevrangeByScores 返回有序集中指定分数区间内的成员，分数从高到低排序
func (c Cache) ZrevrangeByScores(key string, begin, end int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZREVRANGEBYSCORE", key, begin, end)
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

func (c Cache) ZrevrangeByScoresLimit(key string, max, min interface{}, offset, count int, withScores bool) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	var raws interface{}
	var err error

	if withScores {
		raws, err = conn.Do("ZREVRANGEBYSCORE", key, max, min, "LIMIT", offset, count, "WITHSCORES")
	} else {
		raws, err = conn.Do("ZREVRANGEBYSCORE", key, max, min, "LIMIT", offset, count)
	}

	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

//Zrevrangewithscores 返回有序集中指定区间内的成员，通过索引，分数从高到底
func (c Cache) Zrevrangewithscores(key string, begin int, end int) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZREVRANGE", key, begin, end, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), nil
}

//Zcount 计算在有序集合中指定区间分数的成员数
func (c Cache) Zcount(key string, minScore, maxScore int64) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZCOUNT", key, minScore, maxScore)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Zcard 获取有序集合的成员数
func (c Cache) Zcard(key string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("ZCARD", key)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Hgetall 获取在哈希表中指定 key 的所有字段和值
func (c Cache) Hgetall(key string) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("HGETALL", key)
	if err != nil {
		return nil, err
	}
	return redis.Strings(raws, err)
}

//Hmget 获取所有给定字段的值
func (c Cache) Hmget(keyField ...string) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hmget", generalizeStringSlice(keyField)...)
	if err != nil {
		return nil, err
	}
	return redis.Strings(raws, err)
}

//HmgetMulti 获取所有给定多个字段的值
func (c Cache) HmgetMulti(key string, field ...interface{}) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hmget", key, field)
	if err != nil {
		return nil, err
	}
	return redis.Strings(raws, err)
}

//Hget 获取存储在哈希表中指定字段的值。
func (c Cache) Hget(key string, field interface{}) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hget", key, field)
	if err != nil {
		return "", err
	}

	return redis.String(raws, err)
}

//Hget 获取存储在哈希表中指定字段的值,返回[]byte
func (c Cache) HgetBytes(key string, field interface{}) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hget", key, field)
	if err != nil {
		return nil, err
	}

	return redis.Bytes(raws, err)
}

//Hdel 删除一个字段
func (c Cache) Hdel(key, field string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hdel", key, field)
	if err != nil {
		return 0, err
	}

	return redis.Int64(raws, err)
}

//HdelInt64 删除一个字段
func (c Cache) HdelInt64(key, field string) (int64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hdel", key, field)
	if err != nil {
		return 0, err
	}
	return redis.Int64(raws, err)
}

//Hincrby 为哈希表 key 中的指定字段的整数值加上增量 increment
func (c Cache) Hincrby(key, field string, count int64) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("HINCRBY", key, field, count)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//HsetMulti 同时将多个 fieldValue (域-值)对设置到哈希表 key 中
func (c Cache) HsetMulti(fieldValue ...interface{}) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	raws, err := conn.Do("Hmset", fieldValue...)
	if err != nil {
		return "", err
	}
	return redis.String(raws, err)
}

//HdelMulti 删除多个哈希表字段
func (c Cache) HdelMulti(keyFields ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("Hdel", keyFields...)
	return err
}

//Hmgetstrings 获取多个哈希key的值
func (c Cache) Hmgetstrings(args ...interface{}) (int, []*string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("Hmget", args...)
	if err != nil && err != redis.ErrNil {
		return 0, nil, err
	}
	if reply == nil {
		return 0, nil, ErrCacheMiss
	}
	switch raws := reply.(type) {
	case []interface{}:

		n := 0
		ret := make([]*string, 0, len(raws))
		for _, raw := range raws {
			if raw == nil {
				ret = append(ret, nil)
				continue
			}
			s, err := redis.String(raw, nil)
			if err != nil {
				return 0, nil, err
			}
			ret = append(ret, &s)
			n++
		}
		return n, ret, nil
	case nil:
		err = ErrCacheMiss
	case redis.Error:
		err = raws
	default:
		err = fmt.Errorf("unexpected type for result of HMGET: %T", reply)
	}
	return 0, nil, err
}

//Hmgetstrings 获取多个哈希key的值
func (c Cache) Hmgets(args ...interface{}) ([][]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("Hmget", args...)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if reply == nil {
		return nil, ErrCacheMiss
	}
	switch raws := reply.(type) {
	case []interface{}:

		var ret = make([][]byte, 0, len(raws))
		for _, raw := range raws {
			if raw == nil {
				ret = append(ret, nil)
				continue
			}

			s, err := redis.Bytes(raw, nil)
			if err != nil {
				return nil, err
			}
			ret = append(ret, s)
		}
		return ret, nil
	case nil:
		err = ErrCacheMiss
	case redis.Error:
		err = raws
	default:
		err = fmt.Errorf("unexpected type for result of HMGET: %T", reply)
	}
	return nil, err
}

//HgetallObj 获取所有给定字段的值
func (c Cache) HgetallObj(key string) ([]interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("HGETALL", key)
	if err != nil {
		return nil, err
	}
	return raws.([]interface{}), err
}

//HgetObj 获取一个字段值
func (c Cache) HgetObj(key, field interface{}) (interface{}, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hget", key, field)
	if err != nil {
		return "", err
	}
	return raws, err
}

//HgetObj 获取一个字段值
func (c Cache) HmgetObj(key string, fields ...interface{}) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := redis.Bytes(conn.Do("Hmget", key, fields))
	if err != nil {
		return []byte{}, err
	}

	return raws, err
}

func (c Cache) HgetToObj(key, field, ptr interface{}) error {
	dbBytes, err := redis.Bytes(c.HgetObj(key, field))
	if err != nil {
		return err
	}
	return Deserialize(dbBytes, ptr)
}

//HsetObj 设置一个字段值
func (c Cache) HsetObj(key, field string, value interface{}) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := Serialize(value)
	if err != nil {
		return 0, err
	}
	raws, err := conn.Do("Hset", key, field, b)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

func (c Cache) HmsetObj(args []interface{}) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()

	raws, err := conn.Do("Hmset", args...)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Hset 设置一个字段值
func (c Cache) Hset(key, field string, value int) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hset", key, field, value)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Hset 设置一个字段值
func (c Cache) HsetInt64(key, field string, value int64) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hset", key, field, value)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Hsetstring 设置一个字段值 string类型
func (c Cache) Hsetstring(key, field interface{}, value string) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hset", key, field, value)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Expire 为给定 key 设置过期时间
func (c Cache) Expire(key string, expire int) (int, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("EXPIRE", key, expire)
	if err != nil {
		return 0, err
	}
	return redis.Int(raws, err)
}

//Keys 查找所有符合给定模式( pattern)的 key
func (c Cache) Keys(key string) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Keys", key)
	if err != nil {
		return nil, err
	}
	return redis.Strings(raws, err)
}

/*Subscribe Publish and Subscribe

Use the Send, Flush and Receive methods to implement Pub/Sub subscribers.

c.Send("SUBSCRIBE", "example")
c.Flush()
for {
    reply, err := c.Receive()
    if err != nil {
        return err
    }
    // process pushed message
}

The PubSubConn type wraps a Conn with convenience methods for implementing subscribers. The Subscribe, PSubscribe, Unsubscribe and PUnsubscribe methods send and flush a subscription management command. The receive method converts a pushed message to convenient session for use in a type switch.

psc := redis.PubSubConn{c}
psc.Subscribe("example")
for {
    switch v := psc.Receive().(type) {
    case redis.Message:
        fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
    case redis.Subscription:
        fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
    case error:
        return v
    }
}

*/
func (c Cache) Subscribe(key string, exec func(string, error)) {
	conn := c.pool.Get()
	defer conn.Close()
	conn.Send("SUBSCRIBE", key)
	conn.Flush()
	for {
		reply, err := conn.Receive()
		if err != nil {
			fmt.Println("Subscribe conn.Receive() err", key, err)
			return
		}
		for _, raw := range reply.([]interface{}) {
			item, err := redis.Bytes(raw, err)
			if err != nil {
				fmt.Println("Subscribe redis.Bytes err", key, err)
				continue
			}
			obj := ""
			err = Deserialize(item, &obj)
			if err != nil {
				fmt.Println("Subscribe Deserialize err", key, err)
			}
			exec(obj, err)
		}
	}
}

//Publish 将信息发送到指定的频道
func (c Cache) Publish(key, value string) error {
	conn := c.pool.Get()
	defer conn.Close()

	b, err := Serialize(value)
	if err != nil {
		return err
	}
	err = conn.Send("PUBLISH", key, b)
	return err
}

// func (c Cache) Zrange(key string, begin int, end int, valuetype interface{}) ([]interface{}, error) {
// 	conn := c.pool.Get()
// 	defer conn.Close()
// 	raws, err := conn.Do("ZRANGE", key, begin, end)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// if v := reflect.ValueOf(ptrValue); v.Kind() == reflect.Ptr {
// 	// 	p := v.Elem()
// 	// 	fmt.Println("p.Kind():", p.Kind(), "reflect.TypeOf(ptrvalue):", reflect.TypeOf(ptrValue).Elem().Elem())
// 	// }
// 	switch raws := raws.(type) {
// 	case []interface{}:
// 		// switch ptype := reflect.TypeOf(ptrValue).Elem().Elem().(type) {
// 		// case interface{}:
// 		// 	fmt.Println("ptype:", ptype)
// 		// 	// ptrValue = make([]ptype, len(raws))
// 		// 	user := new(ptype)
// 		// 	ptrValue
// 		// 	fmt.Println("user:", reflect.TypeOf(user))
// 		// 	return nil, nil
// 		// default:
// 		// 	return nil, nil
// 		// }
// 		// fmt.Println("len(raws):", len(raws))

// 		res := make([]interface{}, len(raws))
// 		for i, raw := range raws {
// 			// fmt.Println("raw.type:", reflect.TypeOf(raw))
// 			item, err1 := redis.Bytes(raw, err)
// 			if err1 != nil {
// 				fmt.Println("redis.Bytes error:%v", err1)
// 			}
// 			Deserialize(item, valuetype)
// 			res[i] = valuetype
// 		}
// 		return res, nil
// 	default:
// 		return nil, nil
// 	}
// }

//INCRBYFLOAT  为 key 中的指定字段的浮点数值加上增量 increment
func (c Cache) INCRBYFLOAT(key, field string, count float64) (float64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("INCRBYFLOAT", key, field, count)
	if err != nil {
		return 0.0, err
	}
	//raw, err := redis.Int(raws, err)

	//return float64(raw), err
	return redis.Float64(raws, err)
}

//Hincrbyfloat 为哈希表 key 中的指定字段的浮点数值加上增量 increment
func (c Cache) Hincrbyfloat(key, field string, count float64) (float64, error) {
	conn := c.pool.Get()
	defer conn.Close()
	raws, err := conn.Do("Hincrbyfloat", key, field, count)
	if err != nil {
		return 0.0, err
	}
	//raw, err := redis.Int(raws, err)

	//return float64(raw), err
	return redis.Float64(raws, err)
}
