package jso

type CombinedMap[K comparable, V any] struct {
	Map[K, V]
	Multi MultiMap[K, V]
}

func (s *CombinedMap[K, V]) Put(key K, value V) {
	s.Map.Put(key, value)
	s.Multi.Put(key, value)
}
func (s *CombinedMap[K, V]) Delete(key K) {
	s.Map.Delete(key)
	s.Multi.Delete(key)
}
func (s *CombinedMap[K, V]) Clear() {
	s.Map.Clear()
	s.Multi.Clear()
}
func (s *CombinedMap[K, V]) Clone() CombinedMap[K, V] {
	return CombinedMap[K, V]{
		Map:   *s.Map.Clone(),
		Multi: *s.Multi.Clone(),
	}
}
