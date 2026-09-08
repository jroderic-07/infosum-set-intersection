package datastore

type FrequencyStore struct {
	frequencies map[string]uint64
	totalCount  uint64
}

func (s *FrequencyStore) Add(key string) error {
	s.frequencies[key]++
	s.totalCount++

	return nil
}

func (s *FrequencyStore) Get(key string) (uint64, bool) {
	count, exists := s.frequencies[key]

	return count, exists
}

func (s *FrequencyStore) GetDistinctCount() uint64 {
	return uint64(len(s.frequencies))
}

func (s *FrequencyStore) Keys() []string {
	keys := make([]string, 0, len(s.frequencies))
	for key := range s.frequencies {
		keys = append(keys, key)
	}

	return keys
}

func (s *FrequencyStore) GetTotalCount() uint64 {
	return s.totalCount
}

func NewFrequencyStore() *FrequencyStore {
	return &FrequencyStore{
		frequencies: make(map[string]uint64),
	}
}
