package page_loader

type IPageLoader interface {
	LoadPage(url string) (string, error)
}

type MockPageLoader struct {
	LoadPage_ func(url string) (string, error)
}

func (m *MockPageLoader) LoadPage(url string) (string, error) {
	return m.LoadPage_(url)
}
