package grawler

type IAnalyzer[T any] interface {
	GetUrls(source string) []string
	GetModel(source string) *T
}

type NewAnalyzer[T any] func(html *string) IAnalyzer[T]

type MockAnalyzer struct {
	GetUrls_  func() []string
	GetModel_ func() *ExapleModel
}

type ExapleModel struct {
	Title   string
	Content string
}

func (m *MockAnalyzer) GetUrls(source string) []string {
	return m.GetUrls_()
}

func (m *MockAnalyzer) GetModel(source string) *ExapleModel {
	return m.GetModel_()
}
