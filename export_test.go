// Copyright 2015 Bowery, Inc.

package exportfs

import (
	"fmt"
	"testing"
)

func TestExportSingle(t *testing.T) {
	export := &Export{
		Path:    "/pub",
		Machine: "*.test.com",
		Options: []*Option{
			{"ro", ""},
			{"no_wdelay", ""},
			{"mp", "/www/pub"},
		},
	}
	exportLit := "/pub *.test.com(ro,no_wdelay,mp=/www/pub)"

	if export.String() != exportLit {
		t.Error("Export literal is not the expected value")
	}
}

func TestExportList(t *testing.T) {
	exports := []*Export{
		{
			Path:    "/pub",
			Machine: "master",
			Options: []*Option{{"rw", ""}},
		},
		{
			Path:    "/pub",
			Machine: "*",
			Options: []*Option{{"ro", ""}, {"insecure", ""}, {"all_squash", ""}},
		},
	}
	exportsLit := "/pub master(rw) *(ro,insecure,all_squash)"

	if ExportsLiteral(exports) != exportsLit {
		fmt.Println(ExportsLiteral(exports))
		t.Error("Export list literal is not expected value")
	}
}
