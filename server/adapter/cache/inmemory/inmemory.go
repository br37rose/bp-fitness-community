package inmemory

type InMemoryStorage struct {
	s map[uint64]string
}

func NewCache() *InMemoryStorage {
	return &InMemoryStorage{
		s: map[uint64]string{},
	}
}

func (ims *InMemoryStorage) Pull(id uint64) string {
	v, ok := ims.s[id]
	if !ok {
		return ""
	}
	delete(ims.s, id)
	return v
}

func (ims *InMemoryStorage) Push(id uint64, v string) {
	ims.s[id] = v
}

func (ims *InMemoryStorage) Look(id uint64) string {
	v, ok := ims.s[id]
	if !ok {
		return ""
	}
	return v
}

func (ims *InMemoryStorage) LookOrPush(id uint64, nv string) string {
	ov, ok := ims.s[id]
	if !ok {
		ims.Push(id, nv)
		return nv
	}
	return ov
}

func (ims *InMemoryStorage) IDs() []uint64 {
	keys := make([]uint64, 0, len(ims.s))
	for k := range ims.s {
		keys = append(keys, k)
	}
	return keys
}
