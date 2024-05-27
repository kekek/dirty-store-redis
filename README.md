# Golang Dirty Filter 

redis store 

> 支持底层存储按key前缀分组，实现不同的过滤实例

# Golang Dirty Filter

[![GoDoc](https://godoc.org/github.com/antlinker/go-dirtyfilter?status.svg)](https://godoc.org/github.com/antlinker/go-dirtyfilter)

> 基于DFA算法；
> 支持动态修改敏感词，同时支持特殊字符的筛选；
> 敏感词的存储支持内存存储及MongoDB存储。
> 敏感词redis 存储

## 获取

``` bash
$ go get -v github.com/antlinker/go-dirtyfilter
$ go get github.com/kekek/dirty-store-redis
```

## 使用

``` go
package main

import (
	"fmt"
	
	"github.com/antlinker/go-dirtyfilter"
	"github.com/kekek/dirty-store-redis"

	redisCli "github.com/redis/go-redis/v9"
)

var (
	filterText = `我是需要过滤的内容，内容为：**文*@@件**名，需要过滤。。。`
)

func main() {

	rdb := redisCli.NewClient(&redisCli.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

	redisStore := redis.NewRedisDirtyStore(rdb, "sample", []string{"文件", "内容"})

	filterManage := filter.NewDirtyManager(redisStore)
	resultCount, err := filterManage.Filter().FilterResult(filterText, '*', '@')
	if err != nil {
		panic(err)
	}
	fmt.Println(resultCount)
}
```

## 输出结果

```
map[内容:2 文件:1]
```