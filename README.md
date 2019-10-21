# exlru

## summary

ExLru is implement of Lru thar add expire feature.
Lru is base on code https://github.com/golang/groupcache/tree/master/lru


### Example
```go

import (
	"context"
	"github.com/DCRcoder/exlru"
)

DefaultMemCache := exlru.NewMemCache()
callFunc = func(ctx contenx.Context) (interface{}, error){
    result, err := somefunction()
    return result, err
}

key := exlru.GenKey(param1, param2, param3)
name := "cache_name"
result, err := DefaultMemCache.Execute(hd.R.Context(), name, key, cacheFunc, 100, nil)

// if need expire
expire := 3 * time.Second
result, err := DefaultMemCache.Execute(hd.R.Context(), name, key, cacheFunc, 100, &expire)
```