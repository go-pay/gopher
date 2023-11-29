package redislock

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-pay/util/js"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL    = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
	luaObtain  = redis.NewScript(`if redis.call("set", KEYS[1], ARGV[1], "NX", "PX", ARGV[2]) then return redis.status_reply("OK") end`)
	luaGet     = redis.NewScript(`return redis.call("get", KEYS[1])`)
	luaSet     = redis.NewScript(`return redis.call("set", KEYS[1], ARGV[1], "PX", ARGV[2])`)
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("redislock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("redislock: lock not held")
)

// Client wraps a redis client.
type Client struct {
	client redis.Scripter
	//tmp    []byte
	//tmpMu  sync.Mutex
}

// New creates a new Client instance with a custom namespace.
func New(client redis.Scripter) *Client {
	return &Client{client: client}
}

type Data struct {
	Err  string
	Data string
}

func (c *Client) GetData(ctx context.Context, key string) (*Data, error) {
	result, err := luaGet.Run(ctx, c.client, []string{key}).Result()
	if err != nil {
		return nil, err
	}
	d := new(Data)
	if err = js.UnmarshalString(result.(string), d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *Client) SetData(ctx context.Context, key string, data any, err error, ttl time.Duration) error {
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	mm := Data{
		Data: js.MarshalString(data),
	}
	if err != nil {
		mm.Err = err.Error()
	}
	jsonStr := js.MarshalString(mm)
	return luaSet.Run(ctx, c.client, []string{key}, jsonStr, ttlVal).Err()
}

// Obtain tries to obtain a new lock using a key with the given TTL.
// May return ErrNotObtained if not successful.
func (c *Client) Obtain(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	value := strings.ReplaceAll(uuid.New().String(), "-", "")
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	_, err := c.obtain(ctx, key, value, ttlVal)
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotObtained
		}
		return nil, err
	}
	return &Lock{Client: c, key: key, value: value}, nil
}

func (c *Client) obtain(ctx context.Context, key, value string, ttlVal string) (bool, error) {
	_, err := luaObtain.Run(ctx, c.client, []string{key}, value, ttlVal).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

// --------------------------------------------------------------------

// Lock represents an obtained, distributed lock.
type Lock struct {
	*Client
	key   string
	value string
}

// Obtain is a short-cut for New(...).Obtain(...).
func Obtain(ctx context.Context, client redis.Scripter, key string, ttl time.Duration) (*Lock, error) {
	return New(client).Obtain(ctx, key, ttl)
}

// Key returns the redis key used by the lock.
func (l *Lock) Key() string {
	return l.key
}

// Token returns the token value set by the lock.
func (l *Lock) Token() string {
	return l.value
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (l *Lock) TTL(ctx context.Context) (time.Duration, error) {
	res, err := luaPTTL.Run(ctx, l.client, []string{l.key}, l.value).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	if num, ok := res.(int64); ok && num > 0 {
		return time.Duration(num) * time.Millisecond, nil
	}
	return 0, nil
}

// Refresh extends the lock with a new TTL.
// May return ErrNotObtained if refresh is unsuccessful.
func (l *Lock) Refresh(ctx context.Context, ttl time.Duration) error {
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	status, err := luaRefresh.Run(ctx, l.client, []string{l.key}, l.value, ttlVal).Result()
	if err != nil {
		return err
	}
	if status == int64(1) {
		return nil
	}
	return ErrNotObtained
}

// Release manually releases the lock.
// May return ErrLockNotHeld.
func (l *Lock) Release(ctx context.Context) error {
	if l == nil {
		return ErrLockNotHeld
	}
	res, err := luaRelease.Run(ctx, l.client, []string{l.key}, l.value).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrLockNotHeld
		}
		return err
	}
	if i, ok := res.(int64); !ok || i != 1 {
		return ErrLockNotHeld
	}
	return nil
}
