package entgqlplus

import "encoding/json"

type schmeaAnnotationOption = uint
type authFieldGroup = int8

type annotation struct {
	name          string
	SchemaOptions []schmeaAnnotationOption
}

const (
	Subscription schmeaAnnotationOption = iota << 0
	AuthSchema
)

const (
	schemaAnnotationKey = "entgqlSchema"
)

func (a *annotation) Name() string {
	return a.name
}

func (a *annotation) decode(v interface{}) error {
	buffer, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(buffer, a)
}

func SchemaAnnotation(options ...schmeaAnnotationOption) *annotation {
	return &annotation{
		name:          schemaAnnotationKey,
		SchemaOptions: options,
	}
}
