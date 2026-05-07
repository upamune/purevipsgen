//go:build !cgo

package introspection

func NewIntrospection(isDebug bool) *Introspection {
	panic("purevipsgen code generation requires cgo for GObject introspection; generated bindings can be used with CGO_ENABLED=0")
}

func (v *Introspection) GetVipsVersion() string {
	return "0.0.0"
}

func (v *Introspection) DiscoverImageTypes() []ImageTypeInfo {
	return nil
}

func (v *Introspection) DiscoverOperations() []Operation {
	return nil
}

func (v *Introspection) DiscoverEnumTypes() []EnumTypeInfo {
	return nil
}
