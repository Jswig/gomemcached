package internal

import (
	"slices"
	"testing"
	"time"

	"github.com/Jswig/gomemcached/internal/util"
)

func TestCache(t *testing.T) {
	t.Run("set an item in the cache then get it after it expires", func(t *testing.T) {
		cache := NewCache()
		value := []byte("hi rob")
		key := "greeting"
		expiresIn := time.Millisecond
		expiresAt := util.NowUTC().Add(expiresIn)
		cache.Set(key, value, expiresAt)
		// sleep to make sure the cache expires
		time.Sleep(10 * expiresIn)
		_, isValidItem := cache.Get(key)

		assertEqual(t, isValidItem, false)
	})

	t.Run("set one item in cache, set it again then get it", func(t *testing.T) {
		cache := NewCache()
		value1 := []byte("hi clara")
		key := "greeting"
		value2 := []byte("hi bradley")

		cache.Set(key, value1, NeverExpires())
		cache.Set(key, value2, NeverExpires())
		gotValue, isValidItem := cache.Get(key)

		assertEqual(t, isValidItem, true)
		assertItemValueEqual(t, gotValue, value2)
	})

	t.Run("add a key not yet in the cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value := []byte("hello, friend")

		wasAdded := cache.Add(key, value, NeverExpires())
		assertEqual(t, wasAdded, true)
		gotValue, isValid := cache.Get(key)
		assertEqual(t, isValid, true)
		assertItemValueEqual(t, gotValue, value)
	})

	t.Run("add a key aready in cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value1 := []byte("hello, fiend")
		value2 := []byte("hi, my enemy")

		cache.Set(key, value1, NeverExpires())
		wasAdded := cache.Add(key, value2, NeverExpires())
		assertEqual(t, wasAdded, false)

		gotValue, _ := cache.Get(key)
		assertItemValueEqual(t, gotValue, value1)
	})

	t.Run("replace a key not yet in cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value := []byte("hello, brother")

		wasReplaced := cache.Replace(key, value, NeverExpires())
		assertEqual(t, wasReplaced, false)
	})

	t.Run("replace a key already in cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value1 := []byte("hello, brother")
		value2 := []byte("hellow, sister")

		cache.Set(key, value1, NeverExpires())
		wasReplaced := cache.Replace(key, value2, NeverExpires())
		assertEqual(t, wasReplaced, true)

		gotValue, _ := cache.Get(key)
		assertItemValueEqual(t, gotValue, value2)
	})
}

func assertEqual[T comparable](t *testing.T, got T, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func assertItemValueEqual(t *testing.T, got []byte, want []byte) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("got %s wanted %s", got, want)
	}
}
