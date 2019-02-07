package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/status-im/status-go/sdk"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type Msg struct {
	sdk.Msg
}

func (m Msg) TimeString() string {
	fmt.Println("m.Timestamp =", m.Timestamp)
	t := time.Unix(m.Timestamp, 0)
	return humanize.RelTime(t, time.Now(), "earlier", "later")
}

// ByTimestamp implements sort.Interface for []Person based on
// the Timestamp field.
type ByTimestamp []sdk.Msg

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }

type messagesStore struct {
	db       *leveldb.DB
	messages []sdk.Msg
	maxCount int
}

func NewMessagesStore(maxCount int) *messagesStore {
	db, err := leveldb.OpenFile("/data/sg_bots/sg_spectator/messages_store", nil)
	if err != nil {
		log.Fatal("can't open levelDB file. ERR: %v", err)
	}

	// Load data
	messages := messagesFromIterator(db.NewIterator(nil, nil))
	sort.Sort(ByTimestamp(messages))
	log.Println("loaded messages", len(messages))

	return &messagesStore{db, messages, maxCount}
}

func (ms *messagesStore) Add(message sdk.Msg) error {
	if len(ms.messages) >= ms.maxCount {
		if message.Timestamp < ms.messages[0].Timestamp {
			log.Println("Message is too old, ignoring", message.ID, message.Timestamp)
			return nil
		}
	}
	_, err := ms.db.Get([]byte(key(message)), nil)
	if err == nil {
		log.Println("Message already exists, ignoring", message.ID, message.Timestamp)
	}

	toRemove := make([]sdk.Msg, 0)
	fromIdx := 0
	if len(ms.messages) >= ms.maxCount {
		fromIdx = len(ms.messages) - ms.maxCount + 1
		toRemove = make([]sdk.Msg, fromIdx)
		copy(toRemove, ms.messages[:fromIdx])
	}
	messages := append(ms.messages[fromIdx:], message)
	sort.Sort(ByTimestamp(messages))
	ms.messages = messages

	return ms.persist(toRemove, message)
}

func (ms *messagesStore) Messages() []Msg {
	messages := messagesFromIterator(ms.db.NewIterator(nil, nil))
	sort.Sort(ByTimestamp(messages))
	result := make([]Msg, len(messages))
	for idx, msg := range messages {
		result[idx] = Msg{msg}
	}
	return result
}

func (ms *messagesStore) Close() {
	ms.db.Close()
}

func (ms *messagesStore) persist(toRemove []sdk.Msg, toAdd sdk.Msg) error {
	batch := new(leveldb.Batch)

	data, err := json.Marshal(toAdd)
	if err != nil {
		return err
	}
	log.Println("Adding: ", toAdd.ID, toAdd.Timestamp, toAdd.ChannelName)
	log.Printf("Adding: %#v", toAdd)
	batch.Put([]byte(key(toAdd)), []byte(data))

	for _, messageToRemove := range toRemove {
		log.Println("Garbage-colleting: ", messageToRemove.ID, messageToRemove.Timestamp, messageToRemove.ChannelName)
		batch.Delete([]byte(key(messageToRemove)))
	}
	err = ms.db.Write(batch, nil)

	messages := messagesFromIterator(ms.db.NewIterator(nil, nil))
	log.Println("new len (db):", len(messages))

	return err
}

/* Helper functions */

func messagesFromIterator(iter iterator.Iterator) []sdk.Msg {
	messages := make([]sdk.Msg, 0)

	for iter.Next() {
		// Use key/value.
		var message sdk.Msg
		if err := json.Unmarshal(iter.Value(), &message); err != nil {
			log.Println("Cound not unmarshal JSON. ERR: %v", err)
		} else {
			messages = append(messages, message)
		}
	}
	iter.Release()

	return messages
}

func key(message sdk.Msg) string {
	return fmt.Sprintf("%s", message.ID())
}
