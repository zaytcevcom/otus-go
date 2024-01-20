package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(10)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		c.Clear()

		wasInCache = c.Set("aaa", 100)
		require.False(t, wasInCache)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		for i := 1; i <= 4; i++ {
			c.Set(Key("k"+strconv.Itoa(i)), i)
		} // items: {"k2": 2, "k3": 3, "k4": 4} queue: [2, 3, 4]

		val, ok := c.Get("k1")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("k2")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("k4")
		require.True(t, ok)
		require.Equal(t, 4, val)
	})

	t.Run("purge oldest logic", func(t *testing.T) {
		c := NewCache(3)

		for i := 1; i <= 3; i++ {
			c.Set(Key("k"+strconv.Itoa(i)), i)
		} // items: {"k1": 1, "k2": 2, "k3": 3}  queue: [1, 2, 3]

		val, ok := c.Get("k2") // items: {"k1": 1, "k2": 2, "k3": 3}  queue: [2, 1, 3]
		require.True(t, ok)
		require.Equal(t, 2, val)

		wasInCache := c.Set("k3", 333) // items: {"k1": 1, "k2": 2, "k3": 333}  queue: [3, 2, 1]
		require.True(t, wasInCache)

		wasInCache = c.Set("k4", 4) // items: {"k4": 4, "k2": 2, "k3": 333}  queue: [4, 3, 2]
		require.False(t, wasInCache)

		val, ok = c.Get("k1") // items: {"k4": 4, "k2": 2, "k3": 333}  queue: [4, 3, 2]
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
