# KCache
go 本地缓存解决方案，支持本地缓存过期、缓存过期自维护机制。

使用请参考 [使用实例](kcache_test.go)

### 创建KCache
- 创建一个KCache, 默认本地缓存过期时间 5s

  ```
  kc := New()
  ```


- 创建一个KCache, 自定义本地缓存过期时间

  ```
  kc := NewWithExp(2 * time.Second)
  ```

### 获取缓存
- GET 获取缓存，函数不带参数，本地缓存过期时间为创建 KCache 时设置的全局过期时间。

  ```
    kc := New()
    d := kc.Get("myKey", GetData())
  ```
    GET 方法包含两个参数，第一个参数为缓存的key，第二个参数为获取缓存数据的函数。当缓存不存在时，会调用函数获取数据，并将数据缓存起来。
    函数需符合 GetKcDatafunc 类型、返回值需符合 KcData 类型。
  ```
  type GetKcDatafunc func() KcData
  
  type KcData struct { 
    interface{} 
    error
  }
  ```
  
  示例：
  ```
  // 获取缓存数据
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
  ```

- Get 获取缓存，函数带参数

  ```
  kc := New()
  params := map[string]string{
    "k1": "value1",
    "k2": "value2",
  }
  d := kc.Get("myKey", GetDataV2("myKey", params))
  ```
  
  示例：
  ```
  // 获取缓存数据
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
  ```
    
- GetWithExp 获取缓存，自定义本地缓存时间

  ```
  kc := New() 
  exp := 2 * time.Second
  params := map[string]string{
    "k1": "value1",
    "k2": "value2",
  }
  d := kc.GetWithExp("myKey", exp, GetDataV2("myKey", params))
  ```

### GetKcDatafunc 实现
- Kcache 中间函数（强烈推荐）
  
  通过 Kcache 中间函数调用原有的获取数据函数，该函数内部不含任何业务代码，减少业务代码与缓存代码的耦合。
  ```
  kc := New()
  exp := 2 * time.Second
  params := map[string]string{
    "k1": "value1",
    "k2": "value2",
  }
  d := kc.GetWithExp("myKey", exp, GetDataKcache("myKey", params))
  ```

  ```
  // 获取缓存数据, Kcache 中间函数
  func GetDataKcache(key string, params map[string]string) GetKcDatafunc {
    return func() KcData {
      data, err := GetDataV2(key, params)
      return KcData{Data: data, Err: err}
    }
  }
  
  // 获取数据
  func GetDataV2(key string, params map[string]string) (map[string]string, error) {
    // sleep 模拟从 Redis、DB 中获取数据，也可以先从 redis 获取数据, 如果获取不到，再从 DB 中获取。
    time.Sleep(20 * time.Millisecond)
    data := make(map[string]string)
      for k, v := range params {
      data[k+key] = v
    }
    return data, nil
  }
  ```

- 闭包函数（推荐）

  简单获取数据的业务逻辑可以使用闭包函数。
  ```
  kc := New()
  params := map[string]string{
    "k1": "value1",
    "k2": "value2",
  }
  key := "myKey"
  fc := func() KcData {
    // sleep 模拟从 Redis、DB 中获取数据，也可以先从 redis 获取数据, 如果获取不到，再从 DB 中获取。
    time.Sleep(20 * time.Millisecond)
    data := make(map[string]string)
    for k, v := range params {
        data[k+key] = v
    }
    return KcData{Data: data, Err: nil}
  }
  d := kc.Get(key, fc)
  ```

- 业务混合
  ```
  kc := New()
  d := kc.Get("myKey", GetData())
  ```
  ```
  // 获取缓存数据
  func GetData() GetKcDatafunc {
    return func() KcData {
      // sleep 模拟从 Redis、DB 中获取数据
      time.Sleep(20 * time.Millisecond)
      d := map[string]string{
        "k1": "value1",
        "k2": "value2",
      }
      return KcData{Data: d, Err: nil}
    }
  }
  ```
### 设置缓存
- Set 设置缓存，本地缓存过期时间为创建 KCache 时设置的全局过期时间。

  正常情况下无需使用 Set 方法，因为 Get 方法会自动设置缓存。

  ```
  kc := New()
  params := map[string]string{
    "k1": "value1",
    "k2": "value2",
  }
  d := kc.Set("myKey", params)
  ```

- SetWithExp 设置缓存，自定义本地缓存时间。

  正常情况下无需使用 SetWithExp 方法，因为 Get 方法会自动设置缓存。

  ```
  kc := New()
  exp := 2 * time.Second
  params := map[string]string{
    "k1": "value1",
    "k2": "value2",
  }
  d := kc.SetWithExp("myKey", params, exp)
  ```
  
### 删除缓存
- Delete 删除本地缓存

  正常情况下无需使用 Delete 方法，因为有自动删除缓存机制。

  ```
  kc := New()
  params := map[string]string{
    "k1": "value1",
    "k2": "value2",
  }
  d := kc.Delete("myKey")
  ```
  
### 单纯使用本地缓存

- 不需要自维护缓存数据
- 底层使用的 [go-cache](https://github.com/patrickmn/go-cache)，go-cache下的所有方法都可以使用。
  
  ```
  kc := New()
  // SET
  kc.lc.Set("myKey", "myValue", 2*time.Second)
  // GET
  d, f := kc.lc.Get("myKey")
  // other
  ...
  
  ```

更多使用案列请参考 [使用实例](kcache_test.go)
