package internal

import (
	"encoding/binary"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os"
	"sync"
	"time"
)

const (
	BucketConversations = "conversations"
	BucketNameIndex     = "names"
)

type ConversationNotFoundError struct {
	Key interface{}
}

func (e *ConversationNotFoundError) Error() string {
	return fmt.Sprintf("conversation '%v' is not found", e.Key)
}

type ConversationNameDuplicatedError struct {
	Name string
	Id   uint64
}

func (e *ConversationNameDuplicatedError) Error() string {
	return fmt.Sprintf("conversation name '%s' is already used by id %d", e.Name, e.Id)
}

type ConversationInvalidNameError struct {
	Name string
}

func (e *ConversationInvalidNameError) Error() string {
	return fmt.Sprintf("conversation name '%s' is invalid (using only numbers for a name is not allowed)", e.Name)
}

type StoreManager struct {
	DBPath string
	store  *Store
	lock   sync.RWMutex
}

func (m *StoreManager) Open() (*Store, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.store != nil {
		// already opened
		return m.store, nil
	}

	opt := *bolt.DefaultOptions
	opt.Timeout = 1 * time.Second
	db, err := bolt.Open(m.DBPath, os.FileMode(0600), &opt)
	if err != nil {
		return nil, err
	}

	m.store = &Store{
		m:  m,
		db: db,
	}

	return m.store, nil
}

func (m *StoreManager) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.store == nil {
		// already closed
		return nil
	}
	err := m.store.Close()
	if err != nil {
		return err
	}
	m.store = nil
	return nil
}

type Store struct {
	m  *StoreManager
	db *bolt.DB
}

func (s *Store) Close() error {
	if s.db == nil {
		// already closed
		return nil
	}
	err := s.db.Close()
	if err != nil {
		return err
	}
	s.db = nil
	s.m.store = nil
	return nil
}

func (s *Store) Init() error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(BucketConversations)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(BucketNameIndex)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateConversation(co *Conversation) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bc := tx.Bucket([]byte(BucketConversations))
		id, _ := bc.NextSequence()
		co.Id = id
		co.CreatedAt = time.Now().UTC()
		buf, err := serialize(co)
		if err != nil {
			return err
		}

		// update name index
		if co.Name != "" {
			if err := checkValidConversationName(co.Name); err != nil {
				return err
			}
			bni := tx.Bucket([]byte(BucketNameIndex))
			// check if name is already used
			if v := bni.Get([]byte(co.Name)); v != nil {
				return &ConversationNameDuplicatedError{Name: co.Name, Id: btouint64(v)}
			}
			if err := bni.Put([]byte(co.Name), uint64tob(co.Id)); err != nil {
				return err
			}
		}

		return bc.Put(uint64tob(co.Id), buf)
	})
}

func (s *Store) UpdateConversation(co *Conversation) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bc := tx.Bucket([]byte(BucketConversations))
		buf := bc.Get(uint64tob(co.Id))
		if buf == nil {
			return &ConversationNotFoundError{Key: co.Id}
		}

		old := NewConversation()
		if err := deserialize(buf, old); err != nil {
			return err
		}

		// update name index
		if old.Name != co.Name {
			bni := tx.Bucket([]byte(BucketNameIndex))
			if old.Name != "" {
				if err := bni.Delete([]byte(old.Name)); err != nil {
					return err
				}
			}
			if co.Name != "" {
				if err := checkValidConversationName(co.Name); err != nil {
					return err
				}
				// check if name is already used
				if v := bni.Get([]byte(co.Name)); v != nil {
					return &ConversationNameDuplicatedError{Name: co.Name, Id: btouint64(v)}
				}
				if err := bni.Put([]byte(co.Name), uint64tob(co.Id)); err != nil {
					return err
				}
			}
		}

		buf, err := serialize(co)
		if err != nil {
			return err
		}
		return bc.Put(uint64tob(co.Id), buf)
	})
}

func (s *Store) DeleteConversationById(id uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		co := NewConversation()
		buf := tx.Bucket([]byte(BucketConversations)).Get(uint64tob(id))
		if buf == nil {
			return &ConversationNotFoundError{Key: id}
		}
		if err := deserialize(buf, co); err != nil {
			return err
		}
		if co.Name != "" {
			if err := tx.Bucket([]byte(BucketNameIndex)).Delete([]byte(co.Name)); err != nil {
				return err
			}
		}

		return tx.Bucket([]byte(BucketConversations)).Delete(uint64tob(id))
	})
}

func (s *Store) GetConversationById(id uint64) (*Conversation, error) {
	co := NewConversation()
	err := s.db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket([]byte(BucketConversations)).Get(uint64tob(id))
		if buf == nil {
			return &ConversationNotFoundError{Key: id}
		}
		return deserialize(buf, co)
	})
	if err != nil {
		return nil, err
	}
	return co, nil
}

func (s *Store) GetConversationByName(name string) (*Conversation, error) {
	co := NewConversation()
	err := s.db.View(func(tx *bolt.Tx) error {
		id := tx.Bucket([]byte(BucketNameIndex)).Get([]byte(name))
		if id == nil {
			return &ConversationNotFoundError{Key: name}
		}
		buf := tx.Bucket([]byte(BucketConversations)).Get(id)
		if buf == nil {
			return &ConversationNotFoundError{Key: name}
		}
		return deserialize(buf, co)
	})
	if err != nil {
		return nil, err
	}
	return co, nil
}

func (s *Store) GetConversationByKey(key *ConversationKey) (*Conversation, error) {
	if key.IsId {
		return s.GetConversationById(key.IdValue)
	} else {
		return s.GetConversationByName(key.Value)
	}
}

func (s *Store) RenameConversationByKey(key *ConversationKey, name string) error {
	co, err := s.GetConversationByKey(key)
	if err != nil {
		return err
	}

	co.Name = name
	return s.UpdateConversation(co)
}

func (s *Store) ListConversations(query *ListConversationsQuery) (*ConversationList, error) {
	l := &ConversationList{
		query:         query,
		Conversations: []*Conversation{},
		HasNext:       false,
	}

	err := s.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(BucketConversations)).Cursor()
		if query.Reverse {
			if query.Begin != nil {
				begin := uint64tob(*query.Begin)
				var k, v []byte
				if cursor.Bucket().Get(begin) != nil {
					k, v = cursor.Seek(begin)
				} else {
					// If the seeking key does not exist then the next key is used.
					k, _ = cursor.Seek(begin)
					if k == nil {
						k, v = cursor.Last()
					} else {
						k, v = cursor.Prev()
					}
				}

				for ; k != nil; k, v = cursor.Prev() {
					c := NewConversation()
					if err := deserialize(v, c); err != nil {
						return err
					}

					l.TryAppendConversation(c)
					if l.IsLimitReached() {
						break
					}
				}

			} else {
				for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
					c := NewConversation()
					if err := deserialize(v, c); err != nil {
						return err
					}

					l.TryAppendConversation(c)
					if l.IsLimitReached() {
						break
					}
				}
			}

			if k, _ := cursor.Prev(); k != nil {
				l.HasNext = true
				n := binary.BigEndian.Uint64(k)
				l.Next = &n
			} else {
				l.HasNext = false
			}
		} else {
			if query.Begin != nil {
				begin := uint64tob(*query.Begin)
				for k, v := cursor.Seek(begin); k != nil; k, v = cursor.Next() {
					c := NewConversation()
					if err := deserialize(v, c); err != nil {
						return err
					}

					l.TryAppendConversation(c)
					if l.IsLimitReached() {
						break
					}
				}
			} else {
				for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
					c := NewConversation()
					if err := deserialize(v, c); err != nil {
						return err
					}

					l.TryAppendConversation(c)
					if l.IsLimitReached() {
						break
					}
				}
			}

			if k, _ := cursor.Next(); k != nil {
				l.HasNext = true
				n := binary.BigEndian.Uint64(k)
				l.Next = &n
			} else {
				l.HasNext = false
			}
		}

		l.Count = len(l.Conversations)
		return nil
	})

	return l, err
}
