package kcache

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 创建一个KCache, 默认本地缓存过期时间 5s
func TestKCacheNew(t *testing.T) {
	kc := New()
	d := kc.Get("myKey", GetData())
	if d.err != nil {
		t.Error("get key err:", d.err)
		return
	}
	data := d.d.(map[string]string)
	t.Log("key found", data)
	fmt.Println("finish")
}

// 创建一个KCache, 自定义本地缓存过期时间
func TestKCacheNewWithExp(t *testing.T) {
	kc := NewWithExp(2 * time.Second)
	d := kc.Get("myKey", GetData())
	if d.err != nil {
		t.Error("get key err:", d.err)
		return
	}
	data := d.d.(map[string]string)
	t.Log("key found", data)
	fmt.Println("finish")
}

// SingleGet 获取缓存，函数不带参数
func TestKCacheSingleGet(t *testing.T) {
	kc := New()
	d := kc.Get("myKey", GetData())
	if d.err != nil {
		t.Error("get key err:", d.err)
		return
	}
	data := d.d.(map[string]string)
	t.Log("key found", data)
	fmt.Println("finish")
}

// Get 获取缓存，函数带参数
func TestKCacheGet(t *testing.T) {
	kc := New()
	params := map[string]string{
		"k1": "value1",
		"k2": "value2",
	}
	d := kc.Get("myKey", GetDataV2("myKey", params))
	if d.err != nil {
		t.Error("get key err:", d.err)
		return
	}
	t.Log("key found", d)
	fmt.Println("finish")
}

// GetWithExp 获取缓存，自定义本地缓存时间
func TestKCacheGetWithExp(t *testing.T) {
	kc := New()
	exp := 2 * time.Second
	params := map[string]string{
		"k1": "value1",
		"k2": "value2",
	}
	d := kc.GetWithExp("myKey", exp, GetDataV2("myKey", params))
	if d.err != nil {
		t.Error("get key err:", d.err)
		return
	}
	t.Log("key found", d)
	fmt.Println("finish")
}

// Set 设置缓存
func TestKCacheSet(t *testing.T) {
	kc := New()
	params := map[string]string{
		"k1": "value1",
		"k2": "value2",
	}
	kc.Set("myKey", params)
	fmt.Println("finish")
}

// SetWithExp 设置缓存，自定义本地缓存时间。
func TestKCacheSetWithExp(t *testing.T) {
	kc := New()
	exp := 2 * time.Second
	params := map[string]string{
		"k1": "value1",
		"k2": "value2",
	}
	kc.SetWithExp("myKey", exp, params)
	fmt.Println("finish")
}

// Delete 删除缓存
func TestKCacheDelete(t *testing.T) {
	kc := New()
	params := map[string]string{
		"k1": "value1",
		"k2": "value2",
	}
	kc.Set("myKey", params)
	kc.Delete("myKey")
	fmt.Println("finish")
}

// LocalCacke 单纯使用本地缓存，不需要自维护缓存数据
func TestKCacheLocalCacke(t *testing.T) {
	kc := New()
	// SET
	kc.lc.Set("myKey", "myValue", 2*time.Second)
	// GET
	d, f := kc.lc.Get("myKey")
	if !f {
		t.Error("get key not found:")
		return
	}
	// ...
	// kc.lc => *cache.Cache /patrickmn/go-cache下的所有方法都可以使用。
	data := d.(string)
	t.Log("key found", data)
	fmt.Println("finish")
}

// 测试并发性能
func TestKCacheConcurrency(t *testing.T) {
	startTime := time.Now()
	kc := New()
	var sw sync.WaitGroup
	for i := 0; i < 100000; i++ {
		for j := 0; j < 100; j++ {
			sw.Add(1)
			go func(j int) {
				defer sw.Done()
				key := "myKey" + strconv.Itoa(j)
				params := map[string]string{
					"k1": "value1",
					"k2": "value2",
				}
				d := kc.Get(key, GetDataV2(key, params))
				if d.err != nil {
					t.Error("get key err:", d.err)
					return
				}
			}(j)
		}
	}
	sw.Wait()
	fmt.Println(time.Now().Sub(startTime).Milliseconds(), "ms")
	fmt.Println("finish")
}

// 测试持续并发性能
func TestKCacheContinuousConcurrency(t *testing.T) {
	startTime := time.Now()
	kc := New()
	var sw sync.WaitGroup

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-ticker.C:
			go func() {
				sw.Add(1)
				defer sw.Done()
				for i := 0; i < 10000; i++ {
					for j := 0; j < 100; j++ {
						sw.Add(1)
						go func(j int) {
							defer sw.Done()
							key := "myKey" + strconv.Itoa(j)
							params := map[string]string{
								"k1": "value1",
								"k2": "value2",
							}
							d := kc.Get(key, GetDataV2(key, params))
							if d.err != nil {
								t.Error("get key err:", d.err)
								return
							}
						}(j)
					}
				}
			}()
		}

		fmt.Println("Tick at", time.Now(), "count", count)
		count++
		if count >= 10 {
			break
		}
	}
	sw.Wait()
	fmt.Println(time.Now().Sub(startTime).Milliseconds(), "ms")
	fmt.Println("finish")
}

// 获取缓存数据函数
func GetData() GetKcDatafunc {
	return func() KcData {
		// sleep 模拟从 Redis、DB 中获取数据
		time.Sleep(20 * time.Millisecond)
		d := map[string]string{
			"k1": "value1",
			"k2": "value2",
		}
		return KcData{d: d, err: nil}
	}
}

// 获取缓存数据函数, 带参数
func GetDataV2(key string, params map[string]string) GetKcDatafunc {
	return func() KcData {
		// sleep 模拟从 Redis、DB 中获取数据，也可以先从 redis 获取数据, 如果获取不到，再从 DB 中获取。
		time.Sleep(20 * time.Millisecond)
		data := make(map[string]string)
		for k, v := range params {
			data[k+key] = v
		}
		return KcData{d: data, err: nil}
	}
}
