package plan

type Kind string

type Plan interface {
	Kind() Kind
}
