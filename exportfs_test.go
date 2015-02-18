// Copyright 2015 Bowery, Inc.

package exportfs

import (
	"strings"
	"testing"
)

const testExports = `# sample /etc/exports file

/               master(rw) trusty(rw,no_root_squash)
     /projects       proj*.local.domain(rw) # Example comment
/usr            *.local.domain(ro) @trusted(rw)
/home/drake       pc001(rw,all_squash,anonuid=150,anongid=100)
"/home/dumb user" pc002(rw)
/pub            master(rw)                #Master should be able to write.
/pub            *(ro,insecure,all_squash)
                  /srv/www        -sync,rw server @trusted @external(ro)
/foo            2001:db8:9:e54::/64(rw) 192.0.2.0/24(rw)
"/var/log"        master(rw) \
                trusty(ro) \
                @trusted(ro,insecure)
/build          buildhost[0-9].local.domain(rw)
`

// Keep ordering the same as the string above.
var expected = map[string][]*Export{
	"/": {
		{"/", "master", []*Option{{"rw", ""}}},
		{"/", "trusty", []*Option{{"rw", ""}, {"no_root_squash", ""}}},
	},
	"/projects": {
		{"/", "proj*.local.domain", []*Option{{"rw", ""}}},
	},
	"/usr": {
		{"/usr", "*.local.domain", []*Option{{"ro", ""}}},
		{"/usr", "@trusted", []*Option{{"rw", ""}}},
	},
	"/home/drake": {
		{"/home/drake", "pc001", []*Option{
			{"rw", ""}, {"all_squash", ""}, {"anonuid", "150"}, {"anongid", "100"},
		}},
	},
	"/home/dumb user": {
		{"/home/dumb user", "pc002", []*Option{{"rw", ""}}},
	},
	"/pub": {
		{"/pub", "master", []*Option{{"rw", ""}}},
		{"/pub", "*", []*Option{{"ro", ""}, {"insecure", ""}, {"all_squash", ""}}},
	},
	"/srv/www": {
		{"/srv/www", "server", []*Option{{"sync", ""}, {"rw", ""}}},
		{"/srv/www", "@trusted", []*Option{{"sync", ""}, {"rw", ""}}},
		{"/srv/www", "@external", []*Option{{"ro", ""}}},
	},
	"/foo": {
		{"/foo", "2001:db8:9:e54::/64", []*Option{{"rw", ""}}},
		{"/foo", "192.0.2.0/24", []*Option{{"rw", ""}}},
	},
	"/var/log": {
		{"/var/log", "master", []*Option{{"rw", ""}}},
		{"/var/log", "trusty", []*Option{{"ro", ""}}},
		{"/var/log", "@trusted", []*Option{{"ro", ""}, {"insecure", ""}}},
	},
	"/build": {
		{"/build", "buildhost[0-9].local.domain", []*Option{{"rw", ""}}},
	},
}

func TestParse(t *testing.T) {
	exports, err := Parse(strings.NewReader(testExports))
	if err != nil {
		t.Fatal(err)
	}

	for path, exportsList := range exports {
		expectedExports, ok := expected[path]
		if !ok {
			t.Error("Path not found in expected list")
			continue
		}

		for idx, export := range exportsList {
			expectedExport := expectedExports[idx]
			if export.Path != path {
				t.Error("Exports path isn't the same as its parent")
				continue
			}

			if export.Machine != expectedExport.Machine {
				t.Error("Exports machine isn't the expected value")
				continue
			}

			for i, opt := range export.Options {
				if opt.Key != expectedExport.Options[i].Key {
					t.Error("Exports option key isn't the expected value")
				}

				if opt.Value != expectedExport.Options[i].Value {
					t.Error("Exports option value isn't the expected value")
				}
			}
		}
	}
}
