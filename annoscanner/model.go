package annoscanner

type Model struct {
	Package    string
	Name       string
	Annotation string
	Attributes map[string]ModelAttribute
}
