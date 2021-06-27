package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type Status struct {
	Index    int       `json:"index"`
	Total    int       `json:"total"`
	Current  int       `json:"current"`
	UpdateAt time.Time `json:"update_at"`
	Body     string    `json:"body"`
}

type ResumeClient struct {
	ctx       context.Context
	c         *redis.Client
	prefix    string
	batchSize int64

	rs          *redsync.Redsync
	tickerStops *sync.Map
	lockers     *sync.Map
}

func NewResumeClient(opts *redis.Options, keyPrefix string) *ResumeClient {
	c := redis.NewClient(opts)

	pool := goredis.NewPool(c)
	rs := redsync.New(pool)

	return &ResumeClient{
		ctx:       context.Background(),
		c:         c,
		prefix:    keyPrefix,
		batchSize: 100,

		rs: rs,

		tickerStops: &sync.Map{},
		lockers:     &sync.Map{},
	}
}

func (r *ResumeClient) CheckResume() ([]Status, error) {
	keys, err := r.scanKeys()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var ss []Status
	for _, k := range keys {
		s, err := r.getStatusByKey(k)
		if err != nil {
			continue
		}

		if now.Sub(s.UpdateAt) > 5*time.Second {
			ss = append(ss, *s)
		}
	}

	return ss, nil
}

func (r *ResumeClient) SetStatus(index int, total int, body string) error {
	s := &Status{
		Index:    index,
		Total:    total,
		UpdateAt: time.Now(),
		Body:     body,
	}
	return r.setStatus(index, s)
}

func (r *ResumeClient) UpdateCurrent(index int, current int) error {
	return r.updateCurrent(index, current)
}

func (r *ResumeClient) Done(index int) error {
	// close feed use ticker
	v, ok := r.tickerStops.LoadAndDelete(index)
	if ok {
		tc, ok := v.(chan bool)
		if ok {
			tc <- true
		}
	}

	// clean lock
	r.lockers.LoadAndDelete(index)

	// delete status key
	_, err := r.c.Del(r.ctx, r.constructKey(index)).Result()

	return err
}

// Feed index, tell check resume, this is still alive
func (r *ResumeClient) Feed(index int) {
	ticker := time.NewTicker(5 * time.Second)

	defer func() {
		ticker.Stop()
	}()

	stopCh := make(chan bool)
	r.tickerStops.Store(index, stopCh)

	for {
		select {
		case <-ticker.C:
			// check key status
			s, err := r.getStatus(index)
			if err != nil {
				return
			}

			if s.Total == s.Current {
				r.Done(index)
				return
			}

			r.updateTime(index)

		case <-stopCh:
			return

		}
	}
}

func (r *ResumeClient) setStatus(index int, s *Status) error {
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = r.c.Set(r.ctx, r.constructKey(index), string(body), 12*time.Hour).Result()
	if err != nil {
		return err
	}

	m := r.rs.NewMutex(r.constructLockKey(index))
	r.lockers.Store(index, m)

	return nil
}

func (r *ResumeClient) getStatus(index int) (*Status, error) {
	s, err := r.c.Get(r.ctx, r.constructKey(index)).Result()
	if err != nil {
		return nil, err
	}

	var status Status
	err = json.Unmarshal([]byte(s), &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (r *ResumeClient) getStatusByKey(key string) (*Status, error) {
	s, err := r.c.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var status Status
	err = json.Unmarshal([]byte(s), &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (r *ResumeClient) updateTime(index int) error {
	if err := r.lock(index); err != nil {
		return err
	}
	defer r.unlock(index)

	s, err := r.getStatus(index)
	if err != nil {
		return err
	}

	s.UpdateAt = time.Now()

	body, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = r.c.Set(r.ctx, r.constructKey(index), string(body), 12*time.Hour).Result()
	return err
}

func (r *ResumeClient) updateCurrent(index int, current int) error {
	if err := r.lock(index); err != nil {
		return err
	}
	defer r.unlock(index)

	s, err := r.getStatus(index)
	if err != nil {
		return err
	}

	s.Current = current
	s.UpdateAt = time.Now()

	body, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = r.c.Set(r.ctx, r.constructKey(index), string(body), 12*time.Hour).Result()
	return err
}

func (r *ResumeClient) constructKey(index int) string {
	return r.prefix + ":" + strconv.Itoa(index)
}

func (r *ResumeClient) constructLockKey(index int) string {
	return r.prefix + ":" + strconv.Itoa(index) + "-lock"
}

func (r *ResumeClient) lock(index int) error {
	l, ok := r.lockers.Load(index)
	if !ok {
		m := r.rs.NewMutex(r.constructLockKey(index))
		r.lockers.Store(index, m)
		l = m
	}

	m, ok := l.(*redsync.Mutex)
	if !ok {
		return fmt.Errorf("not a valid mutex for index: %d", index)
	}

	return m.LockContext(r.ctx)
}

func (r *ResumeClient) unlock(index int) error {
	l, ok := r.lockers.Load(index)
	if !ok {
		return fmt.Errorf("no such key: %d", index)
	}

	m, ok := l.(*redsync.Mutex)
	if !ok {
		return fmt.Errorf("not a valid mutex for index: %d", index)
	}

	if ok, err := m.UnlockContext(r.ctx); !ok || err != nil {
		return err
	}

	return nil
}

func (r *ResumeClient) scanKeys() (keys []string, err error) {
	var cursor uint64
	var batchKeys []string

	for {
		batchKeys, cursor, err = r.c.Scan(r.ctx, cursor, r.prefix+":*", r.batchSize).Result()
		if err != nil {
			return
		}

		keys = append(keys, batchKeys...)

		if cursor == 0 || len(keys) == 0 {
			return
		}
	}
}
