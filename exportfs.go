// Copyright 2015 Bowery, Inc.

package exportfs

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var (
	comReg      = regexp.MustCompile(`#.*$`)
	lineContReg = regexp.MustCompile(`\\[ \s]*$`)
	whiteReg    = regexp.MustCompile(`\s+`)
)

// Parse parses the exports content from the given reader, and returns a map
// with the path as the key, and the value being a list of machine exports
// for the path.
func Parse(input io.Reader) (map[string][]*Export, error) {
	scanner := bufio.NewScanner(input)
	exports := make(map[string][]*Export, 20)

	for scanner.Scan() {
		scannedLine := strings.TrimLeftFunc(scanner.Text(), unicode.IsSpace)
		line, lineExports, err := parseLine(scannedLine)
		if err != nil {
			return nil, err
		}

		// If it's a line continuation read all of the continued lines.
		if line != "" && lineExports == nil {
			for scanner.Scan() {
				newline := scanner.Text()

				line, lineExports, err = parseLine(line + newline)
				if err != nil {
					return nil, err
				}

				if lineExports != nil {
					break
				}
			}

			if line != "" && lineExports == nil {
				line, lineExports, err = parseLine(line)
				if err != nil {
					return nil, err
				}
			}
		}

		if lineExports != nil && len(lineExports) > 0 {
			path := lineExports[0].Path

			list, ok := exports[path]
			if !ok {
				exports[path] = lineExports
				continue
			}

			exports[path] = append(list, lineExports...)
		}
	}

	return exports, scanner.Err()
}

// parseLine parses the given line, if a line continuation is found, the
// line is returned with no export structure. Otherwise the line is parsed
// and the resulting exports are returned.
func parseLine(line string) (string, []*Export, error) {
	line = stripComments(line)
	if line == "" {
		return "", nil, nil
	}

	// If the line ends with \ then it's a line continuation.
	if lineContReg.MatchString(line) {
		return lineContReg.ReplaceAllString(line, ""), nil, nil
	}
	var exports []*Export
	var defaults []*Option
	line, path := parsePath(line)
	machines := strings.Fields(line)

	// Parse each machine and it's options.
	for _, machine := range machines {
		// A default option list, merge into existing list.
		if machine[0] == '-' {
			defaults = mergeOpts(defaults, parseOpts(machine[1:]))
			continue
		}

		exports = append(exports, parseMachine(machine, path, defaults))
	}

	return "", exports, nil
}

func parseMachine(line, path string, defaults []*Option) *Export {
	export := &Export{Path: path, Options: defaults}
	optIdx := strings.IndexByte(line, '(')
	if optIdx == -1 {
		export.Machine = line
		return export
	}

	export.Machine = line[:optIdx]
	export.Options = parseOpts(line[optIdx+1 : len(line)-1])
	return export
}

// parseOpts parses the comma separated option list given.
func parseOpts(line string) []*Option {
	vals := strings.Split(line, ",")
	opts := make([]*Option, len(vals))

	for idx, val := range vals {
		opt := new(Option)
		list := strings.SplitN(val, "=", 2)

		opt.Key = list[0]
		if len(list) > 1 {
			opt.Value = list[1]
		}

		opts[idx] = opt
	}

	return opts
}

// parsePath gets the path from a line and returns the line stripping the
// path from it.
func parsePath(line string) (newline string, path string) {
	var pathList []rune
	isQuoted := false
	idx := 0

	for i, r := range line {
		if (r == ' ' && !isQuoted) || (r == '"' && isQuoted) {
			if isQuoted {
				i++
			}
			idx = i
			break
		}

		if r == '"' && !isQuoted {
			isQuoted = true
			continue
		}

		pathList = append(pathList, r)
	}

	return line[idx:], string(pathList)
}

// stripComments removes comments from the given line.
func stripComments(line string) string {
	if comReg.MatchString(line) {
		return comReg.ReplaceAllString(line, "")
	}

	return line
}

// mergeOpts updates options replacing existing keys.
func mergeOpts(base, changes []*Option) []*Option {
	for _, opt := range changes {
		idx := -1
		for i, o := range base {
			if o.Key == opt.Key {
				idx = i
				break
			}
		}

		if idx > -1 {
			base[idx] = opt
		} else {
			base = append(base, opt)
		}
	}

	return base
}
