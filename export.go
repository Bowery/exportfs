// Copyright 2015 Bowery, Inc.

package exportfs

import (
	"fmt"
	"strings"
)

// Option is a single options key and optional value
type Option struct {
	Key   string
	Value string
}

// String implements Stringer returning the literal format for an option.
func (opt *Option) String() string {
	if opt.Value == "" {
		return opt.Key
	}

	return fmt.Sprintf("%s=%s", opt.Key, opt.Value)
}

// Export describes an export for a path and the machine able to access it with
// a number of options.
type Export struct {
	Path    string
	Machine string
	Options []*Option
}

// String implements Stringer providing the output as the literal format the
// export is parsed as.
func (ex *Export) String() string {
	// Format: path machine(key,key=val) machine(key,key=value)
	return fmt.Sprintf("%s %s%s", ex.formatPath(), ex.Machine, ex.formatOpts())
}

// formatPath formats the path correctly, adding double quotes if a space
// is found.
func (ex *Export) formatPath() string {
	if !strings.Contains(ex.Path, " ") {
		return ex.Path
	}

	return fmt.Sprintf(`"%s"`, ex.Path)
}

// formatOpts formats the literal format for the options list.
func (ex *Export) formatOpts() string {
	if ex.Options == nil || len(ex.Options) == 0 {
		return ""
	}

	list := make([]string, 0, len(ex.Options))
	for _, opt := range ex.Options {
		list = append(list, opt.String())
	}

	return "(" + strings.Join(list, ",") + ")"
}

// ExportsLiteral returns the literal format for a list of exports that belong
// to the same path.
func ExportsLiteral(exports []*Export) string {
	if len(exports) == 0 {
		return ""
	}
	path := exports[0].formatPath()
	machineList := make([]string, len(exports))

	for i, export := range exports {
		machineList[i] = fmt.Sprintf("%s%s", export.Machine, export.formatOpts())
	}

	return fmt.Sprintf("%s %s", path, strings.Join(machineList, " "))
}
