package annoscanner

type ModelAttribute struct {
	Name     string
	DataType string
}

func MakeModelAttribute(name, dataType string) ModelAttribute {
	return ModelAttribute{
		Name:     name,
		DataType: dataType,
	}
}
