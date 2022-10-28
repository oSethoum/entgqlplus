package entgqlplus

import "encoding/json"

type schmeaAnnotationOption = uint

type annotation struct {
	name          string
	SchemaOptions []schmeaAnnotationOption
}

const (
	Subscription schmeaAnnotationOption = iota << 1
)

const (
	entgqlSchemaKey = "entgqlSchema"
)

func (a *annotation) Name() string {
	return entgqlSchemaKey
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
		name:          entgqlSchemaKey,
		SchemaOptions: options,
	}
}
