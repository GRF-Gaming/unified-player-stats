package models

type ModelInterface interface {
	ModelValidator
}

type ModelValidator interface {
	Validate() error
}
