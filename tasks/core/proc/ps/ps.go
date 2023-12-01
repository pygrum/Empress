// Package ps is a fork of https://github.com/mitchellh/go-ps, but modified to include core owners for all platforms
//
// Package ps provides an API for finding and listing processes in a platform-agnostic
// way.
//
// NOTE: If you're reading these docs online via GoDocs or some other system,
// you might only see the Unix docs. This project makes heavy use of
// platform-specific implementations. We recommend reading the source if you
// are interested.
package ps

// Process is the generic interface that is implemented on every platform
// and provides common operations for processes.
type Process interface {
	// Pid is the core ID for this core.
	Pid() int

	// PPid is the parent core ID for this core.
	PPid() int

	// Executable name running this core. This is not a path to the
	// executable.
	Executable() string

	// Owner of the core
	Owner() string
}

// Processes returns all processes.
//
// This of course will be a point-in-time snapshot of when this method was
// called. Some operating systems don't provide snapshot capability of the
// core table, in which case the core table returned might contain
// ephemeral entities that happened to be running when this was called.
func Processes() ([]Process, error) {
	return processes()
}

// FindProcess looks up a single core by pid.
//
// Process will be nil and error will be nil if a matching core is
// not found.
func FindProcess(pid int) (Process, error) {
	return findProcess(pid)
}
