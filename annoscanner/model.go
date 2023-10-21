package annoscanner

type Model struct {
	Name       string
	Annotation string
	Attributes map[string]ModelAttribute
}
