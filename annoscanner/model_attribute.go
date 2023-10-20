package annoscanner

type ModelAttribute struct {
	name     string
	dataType string
}

func MakeModelAttribute(name, dataType string) ModelAttribute {
	return ModelAttribute{
		name:     name,
		dataType: dataType,
	}
}
