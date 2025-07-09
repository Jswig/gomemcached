package internal

import (
	"fmt"
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
			cmdSet := Set{key: key, value: value, expiresAt: expiresAt}
			cmdSet.Resolve(cache)

			cmdGet := Get{keys: []string{"greeting"}}
			got := string(cmdGet.Resolve(cache))

			want := (fmt.Sprintf(retrievalTextLine, key, len(value)) +
				string(value) + "\r\n" +
				retrievalEnd)

			assertStringsEqual(t, got, want)
		},
	)
	t.Run("set an item in the cache with no expiration then get it", func(t *testing.T) {
		cache := NewCache()
		value := []byte("hey leslie")
		key := "greeting"
		cmdSet := Set{key: key, value: value, expiresAt: NeverExpires()}
		cmdSet.Resolve(cache)

		cmdGet := Get{keys: []string{key}}
		got := string(cmdGet.Resolve(cache))

		want := (fmt.Sprintf(retrievalTextLine, key, len(value)) +
			string(value) + "\r\n" +
			retrievalEnd)
		assertStringsEqual(t, got, want)
	})
}

func assertStringsEqual(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got:\n%q\nwanted:\n%q", got, want)
	}
}
