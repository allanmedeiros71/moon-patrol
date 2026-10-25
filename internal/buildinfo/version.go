// Package buildinfo exposes the version stamped into the binary at build time.
package buildinfo

// Version is set by the linker (-X) at build time.
var Version = "dev"

// String returns the program name and version.
func String() string {
	return "moonpatrol " + Version
}
