package introspection

// Introspection provides discovery and analysis of libvips operations
// through GObject Introspection, extracting operation metadata, argument
// details, and supported enum types.
type Introspection struct {
	discoveredEnumTypes  map[string]string
	enumTypeNames        []enumTypeName
	discoveredImageTypes map[string]ImageTypeInfo
	isDebug              bool
}

// Operation represents a libvips operation.
type Operation struct {
	Name               string
	GoName             string
	Description        string
	Arguments          []Argument
	RequiredInputs     []Argument
	OptionalInputs     []Argument
	RequiredOutputs    []Argument
	OptionalOutputs    []Argument
	HasThisImageInput  bool
	HasImageOutput     bool
	HasOneImageOutput  bool
	HasBufferInput     bool
	HasBufferOutput    bool
	HasArrayImageInput bool
	ImageTypeString    string
}

// Argument represents an argument to a libvips operation.
type Argument struct {
	Name         string
	GoName       string
	Type         string
	GoType       string
	CType        string
	Description  string
	IsRequired   bool
	IsInput      bool
	IsInputN     bool
	IsOutput     bool
	IsOutputN    bool
	IsSource     bool
	IsTarget     bool
	IsImage      bool
	IsBuffer     bool
	IsArray      bool
	Flags        int
	IsEnum       bool
	EnumType     string
	NInputFrom   string
	DefaultValue interface{}
}

// EnumTypeInfo holds information about a vips enum type.
type EnumTypeInfo struct {
	CName       string
	GoName      string
	Description string
	Values      []EnumValueInfo
}

// EnumValueInfo holds information about an enum value.
type EnumValueInfo struct {
	CName       string
	GoName      string
	Value       int
	Description string
	GoValue     string
}

type enumTypeName struct {
	CName  string
	GoName string
}

// ImageTypeInfo represents information about an image type.
type ImageTypeInfo struct {
	TypeName  string
	EnumName  string
	EnumValue string
	MimeType  string
	Order     int
	HasLoader bool
	HasSaver  bool
}
