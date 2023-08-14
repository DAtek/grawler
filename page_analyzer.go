package grawler

type IAnalyzer[T any] interface {
	GetUrls() []string
	GetModel() *T
}

type NewAnalyzer[T any] func(html, source *string) (IAnalyzer[T], error)

type MockAnalyzer struct {
	GetUrls_  func() []string
	GetModel_ func() *ExapleModel
}

type ExapleModel struct {
	Title   string
	Content string
}

func (m *MockAnalyzer) GetUrls() []string {
	return m.GetUrls_()
}

func (m *MockAnalyzer) GetModel() *ExapleModel {
	return m.GetModel_()
}
