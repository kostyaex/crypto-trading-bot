package ui

type ResourceField struct {
	Name  string
	Title string
	//Value string
}

type Resource struct {
	Name        string
	Title       string
	FieldsOrder []string
	Fields      map[string]*ResourceField
}
