package internal

import (
	"slices"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("get item from empty cache", func(t *testing.T) {
		cache := NewCache()
		_, isValidItem := cache.Get("my_key")

		assertItemInvalid(t, isValidItem)
	})

	t.Run("put one item in the cache then retrieve it before it expires", func(t *testing.T) {
		cache := NewCache()
		value := []byte("hi jeff")
		key := "greeting"

		cache.Set(key, value, time.Hour)
		gotValue, isValidItem := cache.Get(key)

		assertItemValid(t, isValidItem)
		assertItemValue(t, gotValue, value)
	})

	t.Run("put one item in the cache then retrieve it after it expires", func(t *testing.T) {
		cache := NewCache()
		value := []byte("hi rob")
		key := "greeting"
		expiresIn := time.Millisecond

		cache.Set(key, value, expiresIn)
		// sleep to make sure the cache expires
		time.Sleep(10 * expiresIn)
		_, isValidItem := cache.Get(key)

		assertItemInvalid(t, isValidItem)
	})

	t.Run("put one item in cache, replace the same key then retrieve it", func(t *testing.T) {
		cache := NewCache()
		value1 := []byte("hi clara")
		key := "greeting"
		value2 := []byte("hi bradley")
		
		cache.Set(key, value1, time.Hour)
		cache.Set(key, value2, time.Hour)
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
