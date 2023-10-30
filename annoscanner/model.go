package annoscanner

type Model struct {
	ImportPackagePath string
	Package           string
	Name              string
	Annotation        string
	Attributes        map[string]ModelAttribute
}
