package xvalidator

type Validatable interface {
	GetErrors() Errors
}

// the map key should be in form of FieldName.ValidateTagName
//
// e.g.
//   - type Demo struct { Name string `validate:"required"`}
//   - the key should be Name.required
type Errors map[string]string
