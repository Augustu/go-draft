package main

import (
	"context"
	"log"

	"github.com/Augustu/go-draft/utils"
	"github.com/go-redis/redis/v8"
)

type client struct {
	*redis.Client
	context context.Context
}

func (c *client) genZSetData(key string, num int) error {
	for i := 0; i < num; i++ {
		res, err := c.ZAdd(c.context, key, &redis.Z{
			Score:  utils.RandomFloat64(),
			Member: utils.RandString(12),
		}).Result()
		if err != nil {
			log.Printf("add into zset failed: %s code: %d", err.Error(), res)
			return err
		}
	}
	log.Printf("gen %d for key: %s", num, key)
	return nil
}

func (c *client) genZSetDataDup(key string, num int) error {
	for i := 0; i < num; i++ {
		res, err := c.ZAdd(c.context, key, &redis.Z{
			Score:  float64(utils.RandomInt(0, 100)),
			Member: utils.RandString(12),
		}).Result()
		if err != nil {
			log.Printf("add into zset failed: %s code: %d", err.Error(), res)
			return err
		}
	}
	log.Printf("gen %d for key: %s", num, key)
	return nil
}

// func (c *client) scanZSet(key string, count int64) error {
// 	keys, cursor, err := c.ZScan(c.context, key, 0, "", count).Result()
// 	if err != nil {
// 		log.Printf("scan zset by key: %s cursor: %d failed: %s", key, cursor, err.Error())
// 		return err
// 	}

// 	for i := 0; i < len(keys); i += 2 {
// 		log.Printf("%d\tkey: %s   score: %s", i/2, keys[i], keys[i+1])
// 	}

// 	log.Printf("scaned %d zset by key: %s\n", len(keys)/2, key)
// 	return nil
// }

func (c *client) scanZSet(key string) error {
	var batchsize int64 = 100
	var start int64 = 0

	count := 0

	for {
		stop := start + batchsize - 1
		zset, err := c.ZRevRangeWithScores(c.context, key, start, stop).Result()
		if err != nil {
			log.Printf("revrange key: %s at %d for %d failed: %s", key, start, batchsize, err.Error())
			break
		}

		if len(zset) == 0 {
			break
		}

		start = stop + 1
		count += len(zset)
	}

	log.Printf("scan zset: %s get %d", key, count)
	return nil
}

func (c *client) rankZSet(key string, count int64) error {
	zset, err := c.ZRevRangeWithScores(c.context, key, 0, count).Result()
	if err != nil {
		log.Printf("rank zset key: %s failed: %s", key, err.Error())
		return err
	}

	for i, z := range zset {
		log.Printf("\t%d\trank: \tkey: %s\tscore: %f\n", i, z.Member, z.Score)
	}

	return nil
}

func (c *client) genHash(key string, num int) error {
	batchsize := 100

	rest := num % batchsize
	nb := num / batchsize

	if rest > 0 {
		tmp := make(map[string]interface{})

		for i := 0; i < rest; i++ {
			tmp[utils.RandString(5)] = utils.RandString(12)
		}

		res, err := c.HSet(c.context, key, tmp).Result()
		if err != nil {
			log.Printf("create hash for key: %s failed: %s res: %d", key, err.Error(), res)
			return err
		}
	}

	for i := 0; i < nb; i++ {
		tmp := make(map[string]interface{})
		for i := 0; i < batchsize; i++ {
			tmp[utils.RandString(5)] = utils.RandString(12)
		}

		res, err := c.HSet(c.context, key, tmp).Result()
		if err != nil {
			log.Printf("create hash for key: %s failed: %s res: %d", key, err.Error(), res)
			return err
		}
	}

	return nil
}

func (c *client) scanHash(key string) error {
	var cursor uint64 = 0
	var count int = 0

	var h []string
	var cur uint64 = 0
	var err error

	for {
		h, cur, err = c.HScan(c.context, key, uint64(cursor), "", 100).Result()
		if err != nil {
			log.Printf("iterate key: %s failed: %s", key, err.Error())
			break
		}
		cursor = cur
		count += len(h) / 2

		if cursor == 0 {
			break
		}
	}

	log.Printf("iterate key: %s count: %d", key, count)
	return nil
}

func test(c *client) {
	var err error
	key := "test"

	// err = c.genZSetData(key, 100)
	// if err != nil {
	// 	log.Panic(err)
	// }

	err = c.scanZSet(key)
	if err != nil {
		log.Panic(err)
	}

	// err = c.rankZSet(key, 20)
	// if err != nil {
	// 	log.Panic(err)
	// }
}

func test1(c *client) {
	var err error
	key := "test-1"

	// err = c.genZSetDataDup(key, 100)
	// if err != nil {
	// 	log.Panic(err)
	// }

	err = c.scanZSet(key)
	if err != nil {
		log.Panic(err)
	}

	// err = c.rankZSet(key, 20)
	// if err != nil {
	// 	log.Panic(err)
	// }
}

func test2(c *client) {
	var err error
	key := "hashes"

	// err = c.genHash(key, 1000)
	// if err != nil {
	// 	log.Printf("get zhash failed: %s", err.Error())
	// }

	err = c.scanHash(key)
	if err != nil {
		log.Printf("scan zhash failed: %s", err.Error())
	}
}

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MaxRetries: 3,
	})

	c := &client{
		context: context.Background(),
		Client:  rc,
	}

	// test(c)
	test1(c)
	// test2(c)

}
