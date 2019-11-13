package crud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The main point of these tests is to catch any case where the interface
// changes. But we also provide a mock for testing.
var _ Store = &MockStore{}

func TestMockStore(t *testing.T) {
	s := NewMockStore()
	is := assert.New(t)
	is.NoError(s.Save("test", []byte("data")))
	list, err := s.List()
	is.NoError(err)
	is.Len(list, 1)
	data, err := s.Read("test")
	is.NoError(err)
	is.Equal(data, []byte("data"))

}

type MockStore struct {
	data map[string][]byte
}

func NewMockStore() *MockStore {
	return &MockStore{data: map[string][]byte{}}
}

func (s *MockStore) Connect() error {
	countB, ok := s.data["connect-count"]
	if !ok {
		countB = []byte("0")
	}

	count, err := strconv.Atoi(string(countB))
	if err != nil {
		return fmt.Errorf("could not convert connect-count %s to int: %v", string(countB), err)
	}

	s.data["connect-count"] = []byte(strconv.Itoa(count + 1))

	return nil
}

func (s *MockStore) Close() error {
	countB, ok := s.data["close-count"]
	if !ok {
		countB = []byte("0")
	}

	count, err := strconv.Atoi(string(countB))
	if err != nil {
		return fmt.Errorf("could not convert close-count %s to int: %v", string(countB), err)
	}

	s.data["close-count"] = []byte(strconv.Itoa(count + 1))

	return nil
}

func (s *MockStore) List() ([]string, error) {
	buf := make([]string, len(s.data))
	i := 0
	for k := range s.data {
		buf[i] = k
		i++
	}
	return buf, nil
}
func (s *MockStore) Save(name string, data []byte) error { s.data[name] = data; return nil }
func (s *MockStore) Read(name string) ([]byte, error)    { return s.data[name], nil }
func (s *MockStore) Delete(name string) error            { delete(s.data, name); return nil }
