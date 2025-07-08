package internal

import (
	"slices"
	"testing"
	"time"

	"github.com/Jswig/gomemcached/internal/util"
)

func TestCache(t *testing.T) {
	t.Run("get item from empty cache", func(t *testing.T) {
		cache := NewCache()
		_, isValidItem := cache.Get("my_key")

		assertItemInvalid(t, isValidItem)
	})

	t.Run(
		"set one item in the cache then get it before it expires",
		func(t *testing.T) {
			cache := NewCache()
			value := []byte("hi jeff")
			key := "greeting"
			expiresAt := util.NowUTC().Add(time.Hour)
			cache.Set(key, value, expiresAt)
			gotValue, isValidItem := cache.Get(key)

			assertItemValid(t, isValidItem)
			assertItemValue(t, gotValue, value)
		})

	t.Run("set an item in the cache with no expiration then get it", func(t *testing.T) {
		cache := NewCache()
		value := []byte("hey leslie")
		key := "greeting"
		cache.Set(key, value, util.ZeroTime())

		gotValue, isValidItem := cache.Get(key)

		assertItemValid(t, isValidItem)
		assertItemValue(t, gotValue, value)
	})

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

		assertItemInvalid(t, isValidItem)
	})

	t.Run("set one item in cache, set it again then get it", func(t *testing.T) {
		cache := NewCache()
		value1 := []byte("hi clara")
		key := "greeting"
		value2 := []byte("hi bradley")

		expiresAt := util.NowUTC().Add(time.Hour)

		cache.Set(key, value1, expiresAt)
		cache.Set(key, value2, expiresAt)
		gotValue, isValidItem := cache.Get(key)

		assertItemValid(t, isValidItem)
		assertItemValue(t, gotValue, value2)
	})
}

func assertItemInvalid(t *testing.T, isValid bool) {
	t.Helper()
	if isValid {
		t.Fatal("item is valid, wanted invalid")
	}
}

func assertItemValid(t *testing.T, isValid bool) {
	t.Helper()
	if !isValid {
		t.Fatal("item is invalid, wanted valid")
	}
}

func assertItemValue(t *testing.T, got []byte, want []byte) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("got %s want %s", got, want)
	}
}
