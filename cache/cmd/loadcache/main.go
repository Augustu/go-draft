package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/Augustu/go-draft/cache/cmd/server/model"
	"github.com/Augustu/go-draft/cache/pkg/cache"
	"github.com/Augustu/go-draft/cache/pkg/config"
	"github.com/Augustu/go-draft/cache/pkg/store"
)

const (
	dsn string = "root:FuCkU@!@#$%^@tcp(127.0.0.1:30306)/cache?charset=utf8mb4&parseTime=True&loc=Local"
	// dsn string = "root:FuCkU@!@#$%^@tcp(118.31.14.196:30306)/cache?charset=utf8mb4&parseTime=True&loc=Local"

	host string = "127.0.0.1:32379"
	// host   string = "118.31.14.196:32379"
	passwd string = ""

	updateStats = `local key = KEYS[1]
	local min = ARGV[1]
	local max = ARGV[2]
	local num = ARGV[3]
	
	local res = redis.call('zrangebyscore', key, min, max)
	local data = res[1]
	if ( data == nil ) then
			local newmember = cjson.encode({s = num})
			redis.call('zadd', key, min, newmember)
			return 1
	end
	
	local json = cjson.decode(data)
	json.s = json.s + num
	
	local updated = cjson.encode(json)
	redis.call('zadd', key, min, updated)
	
	redis.call('zrem', key, data)
	return 0`
)

func newConfig() config.Config {
	c := config.Config{}
	c.URI = dsn
	c.Host = host
	c.Passwd = passwd
	return c
}

type cacheLoader struct {
	statsChan chan model.Stats

	rs   *redsync.Redsync
	lock *redsync.Mutex

	ctx context.Context
	st  *store.Store
	ca  *cache.Cache

	evalsha string
}

func newCacheLoader(c config.Config) *cacheLoader {
	cl := cacheLoader{}
	cl.statsChan = make(chan model.Stats, 10240)

	cl.ctx = context.Background()

	cl.st = store.New(&c.Store)
	cl.ca = cache.New(&c.Cache)

	pool := goredis.NewPool(cl.ca.Client)
	cl.rs = redsync.New(pool)

	mutexname := "global-mutex"
	cl.lock = cl.rs.NewMutex(mutexname,
		redsync.WithRetryDelay(100*time.Millisecond),
		redsync.WithTries(3))

	return &cl
}

func (cl *cacheLoader) init() {
	sha, err := cl.ca.Client.ScriptLoad(cl.ctx, updateStats).Result()
	if err != nil {
		log.Fatalf("load lua script failed: %s", err)
		return
	}

	cl.evalsha = sha
	log.Printf("eval sha: %s", cl.evalsha)
}

func (cl *cacheLoader) loadStats() {
	defer close(cl.statsChan)

	batchsize := 1000
	offset := 0

	cl.st.DB = cl.st.DB.Debug()

	for {
		log.Printf("find stats")

		var sts []model.Stats
		// err := cl.st.DB.Order("id asc").Offset(offset).Limit(batchsize).Find(&sts).Error
		err := cl.st.DB.Raw("select * from stats join ( select id from stats limit ? offset ?) as a on a.id = stats.id", batchsize, offset).Find(&sts).Error
		if err != nil {
			log.Printf("find batch failed: %s", err)
			continue
		}

		log.Printf("find stats: %d, len statschan: %d", len(sts), len(cl.statsChan))

		for _, s := range sts {
			cl.statsChan <- s
		}

		log.Printf("cached stats")

		if len(sts) < batchsize {
			log.Printf("load stats done, total: %d", offset+len(sts))
			return
		}

		log.Printf("load stats offset:  %d", offset)
		offset += batchsize
	}
}

type scoreKey struct {
	Score int `json:"s"`
}

func (cl *cacheLoader) cacheKeyWithLock(pkey string, key, value int) {
	// t1 := time.Now()

	// if err := cl.lock.Lock(); err != nil {
	// 	log.Printf("fetch global lock failed: %s", err)
	// 	return
	// }

	// defer func() {
	// 	if ok, err := cl.lock.Unlock(); !ok || err != nil {
	// 		log.Printf("unlock failed")
	// 	}
	// }()

	mutexKey := pkey + "-" + strconv.Itoa(value)
	mutex := cl.rs.NewMutex(mutexKey, redsync.WithRetryDelay(100*time.Millisecond), redsync.WithTries(3))
	if err := mutex.Lock(); err != nil {
		log.Printf("fetch global lock failed: %s", err)
		return
	}

	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			log.Printf("unlock failed")
		}
	}()

	// cl.ca.Client.Ping(cl.ctx).Result()

	vs := strconv.Itoa(value)

	res, err := cl.ca.Client.ZRangeByScore(cl.ctx, pkey, &redis.ZRangeBy{
		Min:    vs,
		Max:    vs,
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		log.Printf("find key: %s by score: %s failed: %s", pkey, vs, err)
		return
	}

	if len(res) == 0 {
		body, err := json.Marshal(scoreKey{
			Score: key,
		})
		if err != nil {
			log.Printf("marshal scoreKey: %d failed: %s", key, err)
			return
		}

		_, err = cl.ca.Client.ZAdd(cl.ctx, pkey, &redis.Z{
			Score:  float64(value),
			Member: body,
		}).Result()
		if err != nil {
			// TODO catch this case
			log.Printf("add key: %d value: %d failed: %s", key, value, err)
			return
		}
	}

	if len(res) == 2 {
		var k scoreKey
		err := json.Unmarshal([]byte(res[0]), &k)
		if err != nil {
			log.Printf("unmarshal scoreKey: %s failed: %s", res[0], err)
			return
		}

		k.Score += value

		body, err := json.Marshal(k)
		if err != nil {
			log.Printf("marshal scoreKey: %#v failed: %s", k, err)
			return
		}

		cl.ca.Client.ZAdd(cl.ctx, pkey, &redis.Z{
			Score:  float64(value),
			Member: body,
		})
	}

	// t2 := time.Now()
	// log.Printf("deal one in: %s", t2.Sub(t1).String())
}

func (cl *cacheLoader) cacheKey(pkey string, key, value int) {
	vkey := strconv.Itoa(key)
	vs := strconv.Itoa(value)

	res, err := cl.ca.Client.EvalSha(cl.ctx, cl.evalsha, []string{pkey}, vs, vs, vkey).Result()
	if err != nil {
		log.Printf("set pkey: %s key: %d value: %d failed: %s, res: %s", pkey, key, value, err, res)
	}
}

func (cl *cacheLoader) migrateCache() {
	// t1 := time.Now()
	// count := 0

	for s := range cl.statsChan {
		// if count >= 1000 {
		// 	t2 := time.Now()
		// 	sub := t2.Sub(t1).Milliseconds()
		// 	log.Printf("handle 1000 in %d ms", sub)

		// 	log.Printf("sleep %d ms", time.Duration(1000-sub))
		// 	time.Sleep(time.Duration(1000-sub) * time.Millisecond)

		// 	t1 = time.Now()
		// 	count = 0
		// }

		// count++
		pkey := "stats:a:" + strconv.Itoa(s.A)
		value := int(s.OccurredAt.Unix())

		cl.cacheKey(pkey+":b", s.B, value)
		// cl.cacheKey(pkey+":c", s.C, value)
		// cl.cacheKey(pkey+":d", s.D, value)
	}
}

func main() {
	c := newConfig()
	cl := newCacheLoader(c)

	cl.init()

	log.Printf("start")
	go cl.loadStats()

	// go cl.migrateCache()
	// go cl.migrateCache()
	// go cl.migrateCache()
	// go cl.migrateCache()
	// go cl.migrateCache()
	// go cl.migrateCache()
	// go cl.migrateCache()
	cl.migrateCache()
}
