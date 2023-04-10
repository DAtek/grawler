package cache

type MockCache struct {
	Get_    func(key string) (string, error)
	Set_    func(key string, val string) error
	Has_    func(key string) bool
	Delete_ func(key string) error
}

func (m *MockCache) Get(key string) (string, error) {
	return m.Get_(key)
}

func (m *MockCache) Set(key string, val string) error {
	return m.Set_(key, val)
}

func (m *MockCache) Has(key string) bool {
	return m.Has_(key)
}

func (m *MockCache) Delete(key string) error {
	return m.Delete_(key)
}
