package metadata_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chef/metadata-parser"
	"github.com/hashicorp/go-version"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	var test_metadata = `
	name "awesomesauce"
	description 'some description'
	long_description File.Read('README.md')
	version "1.0.1"

	depends "yum", "> 3.0.0"

	depends "mysql", "3.0.0"
	depends "redis"
`

	c, _ := version.NewConstraint("> 3.0.0")
	yumDep := metadata.Dependency{
		Name:       "yum",
		Constraint: c,
	}
	c, _ = version.NewConstraint("3.0.0")
	mysqlDep := metadata.Dependency{
		Name:       "mysql",
		Constraint: c,
	}
	redisDep := metadata.Dependency{
		Name:       "redis",
		Constraint: version.Constraints{},
	}

	v, err := version.NewVersion("1.0.1")
	expect := metadata.Metadata{
		Depends: []metadata.Dependency{yumDep, mysqlDep, redisDep},
		Name:    "awesomesauce",
		Version: *v,
	}

	meta, err := metadata.NewParser(strings.NewReader(test_metadata)).Parse()
	if err != nil {
		fmt.Println(" Error parsing: ", err)
		t.Fail()
	}

	//	if !reflect.DeepEqual(meta, expect) {
	if meta.Name != expect.Name {
		fmt.Printf("Name Mismatch, expected '%s' got '%s'", expect.Name, meta.Name)
		t.Fail()
	}

	if meta.Version.String() != expect.Version.String() {
		fmt.Printf("Name Mismatch, expected '%s' got '%s'", expect.Version, meta.Version)
		t.Fail()
	}

	for i, _ := range []string{"yum", "mysql", "redis"} {
		if meta.Depends[i].Name != expect.Depends[i].Name {
			fmt.Printf("Dep Mismatch,\nexpected: %s\n got: %s", spew.Sdump(expect.Depends[i]), spew.Sdump(meta.Depends[i]))
			t.Fail()
		}
		if meta.Depends[i].Constraint.String() != expect.Depends[i].Constraint.String() {
			fmt.Printf("Dep Mismatch,\nexpected: %s\n got: %s", spew.Sdump(expect.Depends[i]), spew.Sdump(meta.Depends[i]))
			t.Fail()
		}
	}
}
