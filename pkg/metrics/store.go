package metrics

import "sync"

type Store struct {
	store sync.Map
}

func NewStore() *Store {
	return &Store{
		store: sync.Map{},
	}
}
func (s *Store) Add(report *Report) {
	val, ok := s.store.Load(report.Key)
	if ok {
		loaded := val.(*Report)
		loaded.ErrorsCount += report.ErrorsCount
		loaded.ResponseVolume += report.ResponseVolume
		loaded.ResponseCount += report.ResponseCount
		loaded.RequestVolume += report.RequestVolume
		loaded.RequestCount += report.RequestCount
	} else {
		s.store.Store(report.Key, report.Clone())
	}
}

func (s *Store) Get(key string) *Report {
	val, ok := s.store.Load(key)
	if ok {
		return val.(*Report)
	}
	return nil
}

func (s *Store) List() []*Report {
	var list []*Report
	s.store.Range(func(key, value interface{}) bool {
		list = append(list, value.(*Report))
		return true
	})
	return list
}
