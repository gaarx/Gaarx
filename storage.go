package gaarx

import (
	"errors"
	"sync"
)

type (
	storage struct {
		innerMap map[string]*sync.Map
	}
)

func (s *storage) Get(scope string, key interface{}) (interface{}, error) {
	if scopeMap, ok := s.innerMap[scope]; ok {
		if value, exist := scopeMap.Load(key); exist {
			return value, nil
		}
		return nil, errors.New("not found element in scope")
	}
	return nil, errors.New("invalid scope")
}

func (s *storage) Set(scope string, key interface{}, value interface{}) error {
	if scopeMap, ok := s.innerMap[scope]; ok {
		scopeMap.Store(key, value)
	}
	return errors.New("invalid scope")
}

// LoadOrSet returns value, loaded and error
func (s *storage) LoadOrSet(
	scope string,
	key interface{},
	value interface{},
) (actual interface{}, loaded bool, err error) {
	if scopeMap, ok := s.innerMap[scope]; ok {
		actual, loaded = scopeMap.LoadOrStore(key, value)
		return actual, loaded, nil
	}
	return nil, false, errors.New("invalid scope")
}

func (s *storage) Delete(scope string, key interface{}) error {
	if scopeMap, ok := s.innerMap[scope]; ok {
		scopeMap.Delete(key)
		return nil
	}
	return errors.New("invalid scope")
}

func (s *storage) ClearScope(scope string) error {
	if scopeMap, ok := s.innerMap[scope]; ok {
		scopeMap.Range(func(k, v interface{}) bool {
			scopeMap.Delete(k)
			return true
		})
		return nil
	}
	return errors.New("invalid scope")
}

func (s *storage) Range(scope string, rangeFunc func(k, v interface{}) bool) error {
	if scopeMap, ok := s.innerMap[scope]; ok {
		scopeMap.Range(rangeFunc)
		return nil
	}
	return errors.New("invalid scope")
}

func (s *storage) GetAll(scope string) (map[string]interface{}, error) {
	if scopeMap, ok := s.innerMap[scope]; ok {
		all := make(map[string]interface{})
		scopeMap.Range(func(k, v interface{}) bool {
			all[k.(string)] = v
			return true
		})
		return all, nil
	}
	return nil, errors.New("invalid scope")
}
