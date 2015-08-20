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

func TestParser_CommentMeta(t *testing.T) {
	var test_metadata = `
	#
	# Author:: Jesse Nelson <spheromak@gmail.com>
	#
	# Licensed under the Apache License, Version 2.0 (the "License");
	# you may not use this file except in compliance with the License.
	# You may obtain a copy of the License at
	#
	# http://www.apache.org/licenses/LICENSE-2.0
	#
	# Unless required by applicable law or agreed to in writing, software
	# distributed under the License is distributed on an "AS IS" BASIS,
	# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	# See the License for the specific language governing permissions and
	# limitations under the License.
	#

	name 'user'
	version '1.0.24'

	maintainer 'Jesse Nelson'
	maintainer_email 'spheromak@gmail.com'
	license 'Apache 2.0'
	description 'Create system user'
	long_description IO.read(File.join(File.dirname(__FILE__), 'README.md'))

	depends 'poise', '~> 1.0'
	depends 'sudo', '~> 2.5'
	depends 'user', '~> 0.4'
	depends 'ssh', '~> 0.10.4'
	depends 'dnf'
	`

	meta, err := metadata.NewParser(strings.NewReader(test_metadata)).Parse()
	if err != nil {
		fmt.Println(" Error parsing: ", err)
		t.Fail()
	}

	spew.Dump(meta)

}
