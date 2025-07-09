package internal

import (
	"bytes"
	"fmt"
	"time"
)

// --- response text for storage commands ---

// reply for when an item was stored successfully
const storageStored = "STORED\r\n"

// reply for when an item was not stored (e.g. replace or add when conditions are not met)
const storageNotStored = "NOT_STORED\r\n"

// --- response text for deletion commands ---

// reply for when an item was successfully deleted
const deletionDeleted = "DELETED\r\n"

// reply for when an item requested for deletion was not found
const deletionNotFound = "NOT_FOUND\r\n"

// --- response text for retrieval commands ---

// signals the end of a retrieval reply
const retrievalEnd = "END\r\n"

// template for the text line in a retrieval reply
const retrievalTextLine = "VALUE %s %d \r\n"

type Command interface {
	// gets the bytes reponse of the command in the format dictated by the
	// memcached protocol.
	// See https://raw.githubusercontent.com/memcached/memcached/refs/heads/master/doc/protocol.txt
	// for a reference on the memcached protocol.
	Resolve(*Cache) []byte
}

// memcached 'add' command
type Add struct {
	key       string
	value     []byte
	expiresAt time.Time
}

func (cmd *Add) Resolve(cache *Cache) []byte {
	wasAdded := cache.Add(cmd.key, cmd.value, cmd.expiresAt)
	if wasAdded {
		return []byte(storageStored)
	} else {
		return []byte(storageNotStored)
	}
}

// memcached 'delete' command
type Delete struct {
	key string
}

func (cmd *Delete) Resolve(cache *Cache) []byte {
	wasDeleted := cache.Delete(cmd.key)
	if wasDeleted {
		return []byte(deletionDeleted)
	} else {
		return []byte(deletionNotFound)
	}
}

// memcached 'get' command
type Get struct {
	keys []string
}

func (cmd *Get) Resolve(cache *Cache) []byte {
	// TODO: there might be a more efficient way of doing this in both cases by 
	// fetching the items, determining the total response size, then pre-allocating
	// a byte slice of the right length.
	result := &bytes.Buffer{}
	for _, key := range cmd.keys {
		value, hasValidItem := cache.Get(key)
		if hasValidItem {
			result.WriteString(fmt.Sprintf(retrievalTextLine, key, len(value)))
			result.Write(value)
			result.WriteString("\r\n")
		}
	}
	result.WriteString(retrievalEnd)
	return result.Bytes()
}

// memcached 'replace' command
type Replace struct {
	key       string
	value     []byte
	expiresAt time.Time
}

func (cmd *Replace) Resolve(cache *Cache) []byte {
	wasReplaced := cache.Replace(cmd.key, cmd.value, cmd.expiresAt)
	if wasReplaced {
		return []byte(storageStored)
	} else {
		return []byte(storageNotStored)
	}
}

// memcached 'set' command
type Set struct {
	key       string
	value     []byte
	expiresAt time.Time
}

func (cmd *Set) Resolve(cache *Cache) []byte {
	cache.Set(cmd.key, cmd.value, cmd.expiresAt)
	return []byte(storageStored)
}
