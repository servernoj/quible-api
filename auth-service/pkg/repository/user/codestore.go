package user

import (
	"errors"
)

// CodeStore is an interface for storing and retrieving verification codes.
type CodeStore interface {
	SaveCode(email string, code int) error
	GetCode(email string) (int, error)
}

// InMemoryCodeStore is an implementation of CodeStore that stores codes in memory.
type InMemoryCodeStore struct {
	codes map[string]int
}

func NewInMemoryCodeStore() *InMemoryCodeStore {
	return &InMemoryCodeStore{
		codes: make(map[string]int),
	}
}

func (s *InMemoryCodeStore) SaveCode(email string, code int) error {
	s.codes[email] = code
	return nil
}

func (s *InMemoryCodeStore) GetCode(email string) (int, error) {
	code, ok := s.codes[email]
	if !ok {
		return 0, errors.New("code not found")
	}
	return code, nil
}

// VerifyCode checks whether the provided code matches the code associated with the specified email address.
func VerifyCode(email string, code int, store CodeStore) (bool, error) {
	savedCode, err := store.GetCode(email)
	if err != nil {
		return false, err
	}
	return savedCode == code, nil
}
