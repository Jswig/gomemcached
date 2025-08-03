package internal

import (
	"testing"
	"time"

	"github.com/Jswig/gomemcached/internal/util"
)

func TestCommand(t *testing.T) {
	t.Run("get an item from empty cache", func(t *testing.T) {
		cache := NewCache()
		cmd := Get{keys: []string{"my_key"}}
		got := string(cmd.Resolve(cache))
		want := retrievalEnd
		assertStringsEqual(t, got, want)
	})

	t.Run(
		"set one item in the cache then get it before it expires",
		func(t *testing.T) {
			cache := NewCache()
			value := []byte("hi jeff")
			key := "greeting"
			expiresAt := util.NowUTC().Add(time.Hour)
			cmdSet := Set{key, value, expiresAt}
			cmdSet.Resolve(cache)

			cmdGet := Get{keys: []string{"greeting"}}
			got := string(cmdGet.Resolve(cache))

			want := "VALUE greeting 7 \r\nhi jeff\r\nEND\r\n"
			assertStringsEqual(t, got, want)
		},
	)

	t.Run("set an item in the cache then get it after it expires", func(t *testing.T) {
		cache := NewCache()
		value := []byte("hi rob")
		key := "greeting"
		expiresIn := time.Millisecond
		expiresAt := util.NowUTC().Add(expiresIn)

		cmdSet := Set{key, value, expiresAt}
		cmdSet.Resolve(cache)

		cmdGet := Get{keys: []string{key}}
		// sleep to make sure the cache expires
		time.Sleep(10 * expiresIn)
		got := string(cmdGet.Resolve(cache))

		assertStringsEqual(t, got, retrievalEnd)
	})

	t.Run("set an item in the cache with no expiration then get it", func(t *testing.T) {
		cache := NewCache()
		value := []byte("hey leslie")
		key := "greeting"
		cmdSet := Set{key, value, NeverExpires()}
		cmdSet.Resolve(cache)

		cmdGet := Get{keys: []string{key}}
		got := string(cmdGet.Resolve(cache))

		want := "VALUE greeting 10 \r\nhey leslie\r\nEND\r\n"
		assertStringsEqual(t, got, want)
	})

	t.Run("set one item in cache, set it again then get it", func(t *testing.T) {
		cache := NewCache()
		value1 := []byte("hi clara")
		key := "greeting"
		value2 := []byte("hi bradley")

		cmdSet := Set{key, value1, NeverExpires()}
		cmdSet.Resolve(cache)
		cmdSet = Set{key, value2, NeverExpires()}
		cmdSet.Resolve(cache)

		cmdGet := Get{keys: []string{key}}
		got := string(cmdGet.Resolve(cache))

		want := "VALUE greeting 10 \r\nhi bradley\r\nEND\r\n"
		assertStringsEqual(t, got, want)
	})

	t.Run("add a key not yet in the cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value := []byte("hello, friend")

		cmdAdd := Add{key, value, NeverExpires()}
		got := string(cmdAdd.Resolve(cache))
		assertStringsEqual(t, got, storageStored)

		cmdGet := Get{keys: []string{key}}
		got = string(cmdGet.Resolve(cache))
		want := "VALUE greeting 13 \r\nhello, friend\r\nEND\r\n"
		assertStringsEqual(t, got, want)
	})

	t.Run("add a key aready in cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value1 := []byte("hello, fiend")
		value2 := []byte("hi, my enemy")

		cmdSet := Set{key, value1, NeverExpires()}
		cmdSet.Resolve(cache)
		cmdAdd := Add{key, value2, NeverExpires()}
		got := string(cmdAdd.Resolve(cache))
		assertStringsEqual(t, got, storageNotStored)

		cmdGet := Get{keys: []string{key}}
		got = string(cmdGet.Resolve(cache))
		want := "VALUE greeting 12 \r\nhello, fiend\r\nEND\r\n"
		assertStringsEqual(t, got, want)
	})

	t.Run("replace a key not yet in cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value := []byte("hello, brother")

		cmdReplace := Replace{key, value, NeverExpires()}
		got := string(cmdReplace.Resolve(cache))
		assertStringsEqual(t, got, storageNotStored)
	})

	t.Run("replace a key already in cache", func(t *testing.T) {
		cache := NewCache()
		key := "greeting"
		value1 := []byte("hello, brother")
		value2 := []byte("hellow, sister")

		cmdSet := Set{key, value1, NeverExpires()}
		cmdSet.Resolve(cache)
		cmdReplace := Replace{key, value2, NeverExpires()}
		got := string(cmdReplace.Resolve(cache))
		assertStringsEqual(t, got, storageStored)

		cmdGet := Get{keys: []string{key}}
		got = string(cmdGet.Resolve(cache))

		want := "VALUE greeting 14 \r\nhellow, sister\r\nEND\r\n"
		assertStringsEqual(t, got, want)
	})
}

func assertStringsEqual(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got:\n%q\nwanted:\n%q", got, want)
	}
}
