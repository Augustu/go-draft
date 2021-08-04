package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Augustu/go-draft/cache/cmd/server/model"
	"github.com/Augustu/go-draft/cache/pkg/cache"
	"github.com/Augustu/go-draft/cache/pkg/config"
	"github.com/Augustu/go-draft/cache/pkg/store"
	"github.com/go-redis/redis/v8"
)

const (
	dsn string = "root:FuCkU@!@#$%^@tcp(10.10.15.11:32306)/cache?charset=utf8mb4&parseTime=True&loc=Local"

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

	return &cl
}

func (cl *cacheLoader) loadStats() {
	defer close(cl.statsChan)

	batchsize := 1000
	offset := 0

	for {
		var sts []model.Stats
		err := cl.st.DB.Find(&sts).Order("id asc").Offset(offset).Limit(batchsize).Error
		if err != nil {
			log.Printf("find batch failed: %s", err)
			continue
		}

		for _, s := range sts {
			cl.statsChan <- s
		}

		if len(sts) < batchsize {
			log.Printf("load stats done, total: %d", offset+len(sts))
			return
		}

		offset += batchsize
	}
}

type scoreKey struct {
	Score int `json:"s"`
}

func (cl *cacheLoader) cacheKey(pkey string, key, value int) {
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

	go cl.loadStats()
	cl.migrateCache()
}
