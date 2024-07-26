package kcache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

const (
	// DefaultLcExp 默认本地缓存过期时间 5s
	DefaultLcExp time.Duration = 5 * time.Second
)

type KCache struct {
	lc     *cache.Cache  // 本地缓存
	mu     sync.RWMutex  // 读写锁
	lcExp  time.Duration // 本地缓存过期时间
	entrys map[string]*entry
}

type entry struct {
	fc GetKcDatafunc // 本地缓存不存在时，获取缓存数据函数
	mu sync.RWMutex  // 读写锁
}

// GetKcDatafunc 获取缓存数据函数，本地缓存不存在时
type GetKcDatafunc func() KcData

type KcData struct {
	Data interface{}
	Err  error
}

// New 创建一个KCache, 默认本地缓存过期时间 5s
func New() *KCache {
	return &KCache{
		lc:     cache.New(DefaultLcExp, time.Second*10),
		mu:     sync.RWMutex{},
		lcExp:  DefaultLcExp,
		entrys: map[string]*entry{},
	}
}

// NewWithExp 创建一个KCache, 自定义本地缓存过期时间
// lcExp 本地缓存过期时间
func NewWithExp(lcExp time.Duration) *KCache {
	return &KCache{
		lc:     cache.New(lcExp, lcExp+time.Second*10),
		mu:     sync.RWMutex{},
		lcExp:  lcExp,
		entrys: map[string]*entry{},
	}
}

// Get 获取缓存
// k 缓存KEY
// fc 本地缓存不存在时，获取缓存数据函数
func (kc *KCache) Get(k string, fc GetKcDatafunc) KcData {
	kc.mu.RLock()
	e := kc.entrys[k]
	kc.mu.RUnlock()
	if e == nil {
		kc.mu.Lock()
		e = kc.entrys[k]
		if e == nil {
			e = &entry{fc: fc, mu: sync.RWMutex{}}
			kc.entrys[k] = e
		} else {
			e = kc.entrys[k]
		}
		kc.mu.Unlock()
	}

	d, f := kc.lc.Get(k)
	if !f {
		e.mu.Lock()
		defer e.mu.Unlock()
		d, f = kc.lc.Get(k)
		if !f {
			d = e.fc()
			kc.lc.Set(k, d, kc.lcExp)
		}
	}

	return d.(KcData)
}

// GetWithCtx 获取缓存, 支持上下文
// k 缓存KEY
// fc 本地缓存不存在时，获取缓存数据函数
func (kc *KCache) GetWithCtx(k string, fc GetKcDatafunc) KcData {
	d, f := kc.lc.Get(k)
	if !f {
		kc.mu.Lock()
		defer kc.mu.Unlock()
		d, f = kc.lc.Get(k)
		if !f {
			d = fc()
			kc.lc.Set(k, d, kc.lcExp)
		}
	}

	return d.(KcData)
}

// GetWithExp 获取缓存，自定义本地缓存时间
// k 缓存KEY
// t 自定义本地缓存时间
// fc 本地缓存不存在时，获取缓存数据函数
func (kc *KCache) GetWithExp(k string, t time.Duration, fc GetKcDatafunc) KcData {
	kc.mu.RLock()
	e := kc.entrys[k]
	kc.mu.RUnlock()
	if e == nil {
		kc.mu.Lock()
		e = kc.entrys[k]
		if e == nil {
			e = &entry{fc: fc, mu: sync.RWMutex{}}
			kc.entrys[k] = e
		} else {
			e = kc.entrys[k]
		}
		kc.mu.Unlock()
	}

	d, f := kc.lc.Get(k)
	if !f {
		e.mu.Lock()
		defer e.mu.Unlock()
		d, f = kc.lc.Get(k)
		if !f {
			d = e.fc()
			kc.lc.Set(k, d, t)
		}
	}

	return d.(KcData)
}

// GetWithCtx 获取缓存, 支持上下文
// k 缓存KEY
// t 自定义本地缓存时间
// fc 本地缓存不存在时，获取缓存数据函数
func (kc *KCache) GetWithExpCtx(k string, t time.Duration, fc GetKcDatafunc) KcData {
	d, f := kc.lc.Get(k)
	if !f {
		kc.mu.Lock()
		defer kc.mu.Unlock()
		d, f = kc.lc.Get(k)
		if !f {
			d = fc()
			kc.lc.Set(k, d, t)
		}
	}

	return d.(KcData)
}

// Set 设置本地缓存
// k 缓存KEY
// Data 缓存数据
func (kc *KCache) Set(k string, d interface{}) {
	kcd := KcData{Data: d, Err: nil}
	kc.lc.Set(k, kcd, kc.lcExp)
}

// SetWithExp 设置本地缓存,自定义本地缓存时间
// k 缓存KEY
// Data 缓存数据
// t 本地缓存时间
func (kc *KCache) SetWithExp(k string, t time.Duration, d interface{}) {
	kcd := KcData{Data: d, Err: nil}
	kc.lc.Set(k, kcd, t)
}

// Delete 删除本地缓存
// k 缓存KEY
func (kc *KCache) Delete(k string) {
	kc.lc.Delete(k)
}
