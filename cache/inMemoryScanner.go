package cache

type inMemoryScanner struct {
	pair
	pairCh  chan *pair
	closeCh chan struct{}
}

func (s *inMemoryScanner) Scan() bool {
	p, ok := <-s.pairCh
	if ok {
		s.k, s.v = p.k, p.v
	}
	return ok
}

func (s *inMemoryScanner) Key() string {
	return s.k
}

func (s *inMemoryScanner) Value() []byte {
	return s.v
}

func (s *inMemoryScanner) Close() {
	close(s.closeCh)
}
