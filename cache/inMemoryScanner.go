package cache

type inMemoryScanner struct {
	pair
	pairCh  chan *pair
	closeCh chan struct{}
}

func (s *inMemoryScanner) Scan() bool {
	p, ok := <-s.pairCh
	if ok {
		s.key, s.value = p.key, p.value
	}
	return ok
}

func (s *inMemoryScanner) Key() string {
	return s.key
}

func (s *inMemoryScanner) Value() []byte {
	return s.value
}

func (s *inMemoryScanner) Close() {
	close(s.closeCh)
}
