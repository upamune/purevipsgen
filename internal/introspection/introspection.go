package introspection

// #cgo pkg-config: vips
// #include "introspection.h"
import "C"
import (
	"log"
)

// NewIntrospection creates a new Introspection instance for analyzing libvips
// operations, initializing the libvips library in the process.
func NewIntrospection(isDebug bool) *Introspection {
	// Initialize libvips
	if C.vips_init(C.CString("purevipsgen")) != 0 {
		log.Fatal("Failed to initialize libvips")
	}
	defer C.vips_shutdown()

	return &Introspection{
		discoveredEnumTypes:  make(map[string]string),
		discoveredImageTypes: map[string]ImageTypeInfo{},
		isDebug:              isDebug,
	}
}

// GetVipsVersion returns the libvips version string
func (v *Introspection) GetVipsVersion() string {
	return C.GoString(C.vips_version_string())
}
