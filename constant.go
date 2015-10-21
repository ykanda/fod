package main

import "fmt"

const (
	FS_TYPE_FILE    = "f"
	FS_TYPE_DIR     = "d"
	FS_TYPE_SYMLINK = "s"
)

const (
	VERSION_MAJOUR = 0
	VERSION_MINOR  = 1
	VERSION_PATCH  = 0
)

func Version() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		VERSION_MAJOUR,
		VERSION_MINOR,
		VERSION_PATCH,
	)
}
