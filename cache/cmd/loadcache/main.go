package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"

	"github.com/Augustu/go-draft/cache/cmd/server/model"
	"github.com/Augustu/go-draft/cache/pkg/cache"
	"github.com/Augustu/go-draft/cache/pkg/config"
	"github.com/Augustu/go-draft/cache/pkg/store"
)

const (
	dsn string = "root:FuCkU@!@#$%^@tcp(127.0.0.1:30306)/cache?charset=utf8mb4&parseTime=True&loc=Local"

	host   string = "10.10.15.11:32379"
	passwd string = ""
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
}

func newCacheLoader(c config.Config) *cacheLoader {
	cl := cacheLoader{}
	cl.statsChan = make(chan model.Stats, 10240)

	cl.ctx = context.Background()

	cl.st = store.New(&c.Store)
	cl.ca = cache.New(&c.Cache)

	// pool := goredis.NewPool(cl.ca.Client)
	// cl.rs = redsync.New(pool)

	// mutexname := "global-mutex"
	// cl.lock = cl.rs.NewMutex(mutexname)

	return &cl
}

func (cl *cacheLoader) loadStats() {
	defer close(cl.statsChan)

	batchsize := 10000
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

		log.Printf("find stats: %d", len(sts))

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

func (cl *cacheLoader) cacheKey(pkey string, key, value int) {
	// t1 := time.Now()

	// if err := cl.lock.Lock(); err != nil {
	// 	log.Printf("fetch global lock failed: %s", err)
	// 	return
	// }

	// mutexKey := pkey + "-" + strconv.Itoa(value)
	// mutex := cl.rs.NewMutex(mutexKey, redsync.WithRetryDelay(time.Second), redsync.WithExpiry(3*time.Second))
	// if err := mutex.Lock(); err != nil {
	// 	log.Printf("fetch global lock failed: %s", err)
	// 	return
	// }

	// defer func() {
	// 	// if ok, err := cl.lock.Unlock(); !ok || err != nil {
	// 	// 	log.Printf("unlock failed")
	// 	// }
	// 	if ok, err := mutex.Unlock(); !ok || err != nil {
	// 		log.Printf("unlock failed")
	// 	}
	// }()

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

func (cl *cacheLoader) migrateCache() {
	for s := range cl.statsChan {
		pkey := "stats:a:" + strconv.Itoa(s.A)
		value := int(s.OccurredAt.Unix())

		cl.cacheKey(pkey+":b", s.B, value)
		cl.cacheKey(pkey+":c", s.C, value)
		cl.cacheKey(pkey+":d", s.D, value)
	}
}

func main() {
	c := newConfig()
	cl := newCacheLoader(c)

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
