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

	redisStore := redis.NewRedisDirtyStore(rdb, "sample", []string{"文件", "内容", "Redis Incr 命令"})

	filterManage := filter.NewDirtyManager(redisStore)
	resultCount, err := filterManage.Filter().FilterResult(filterText, '*', '@')
	if err != nil {
		panic(err)
	}

	fmt.Println(resultCount)

	ch := filterManage.Store().Read()
	for w := range ch {
		fmt.Println(w)
	}

}
