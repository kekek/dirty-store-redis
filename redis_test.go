package redis

import (
	"github.com/redis/go-redis/v9"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// 测试前，清空 redis
var _ = Describe("敏感词redis存储测试", func() {
	var (
		memStore *RedisDirtyStore
	)
	BeforeEach(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		})

		s := NewRedisDirtyStore(rdb, "test-", []string{"共产党"})
		memStore = s
	})
	It("Read Test", func() {
		for v := range memStore.Read() {
			Expect(v).To(Equal("共产党"))
		}
	})
	It("ReadAll Test", func() {
		result, err := memStore.ReadAll()
		if err != nil {
			Fail(err.Error())
			return
		}
		Expect(result).To(Equal([]string{"共产党"}))
	})
	It("Remove Test", func() {
		err := memStore.Remove("共产党")
		if err != nil {
			Fail(err.Error())
			return
		}
		result, err := memStore.ReadAll()
		if err != nil {
			Fail(err.Error())
			return
		}
		Expect(len(result)).To(Equal(0))
	})
	It("Write Test", func() {
		err := memStore.Write("党")
		if err != nil {
			Fail(err.Error())
			return
		}
		result, err := memStore.ReadAll()
		if err != nil {
			Fail(err.Error())
			return
		}
		Expect(len(result)).To(Equal(2))
	})
	It("Version Test", func() {
		Expect(memStore.Version()).To(Equal(uint64(7)))
		err := memStore.Write("党")
		if err != nil {
			Fail(err.Error())
			return
		}
		Expect(memStore.Version()).To(Equal(uint64(8)))
	})
})
