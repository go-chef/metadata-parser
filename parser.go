/* The MIT License (MIT)

Portions of this code Copyright (c) 2013-2015 Errplane Inc.
Copyright (c) 2015 Jesse Nelson

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package metadata

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-version"
)

// Metadata represents a coookbook metadata file
type Metadata struct {
	Depends []Dependency
	Name    string
	Version version.Version
}

type Dependency struct {
	Name       string
	Constraint version.Constraints // ... maybe look at hashi version.Constraint
}

// Parser represents a parser.
type Parser struct {
	s *bufScanner
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: newBufScanner(r)}
}

// Parse parses a metadata file
func (p *Parser) Parse() (md *Metadata, err error) {
	// return decl inits to nil, we need this to be something we can use
	md = &Metadata{}
	for {
		tok, _, _ := p.scanIgnoreWhitespace()

		switch tok {
		case EOF:
			return md, nil
		case NAME:
			md.Name, err = p.parseName()
			if err != nil {
				return md, err
			}
			p.unscan()
		case DEPENDS:
			d, err := p.parseDepends()
			if err != nil {
				return md, err
			}
			md.Depends = append(md.Depends, d)
			p.unscan()
		case VERSION:
			v, err := p.parseVersion()
			if err != nil {
				return nil, err
			}
			md.Version = v
			p.unscan()
		}
	}
}

func (p *Parser) parseVersion() (version.Version, error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != STRING {
		return version.Version{}, fmt.Errorf("found %s expected STRING at %d, %s", tok, pos, lit)
	}

	v, err := version.NewVersion(lit)
	if err != nil {
		return *v, fmt.Errorf("error parsing version at %d, %s : %s", pos, lit, err)
	}

	if len(v.Segments()) < 3 {
		return *v, fmt.Errorf("error parsing version,not enough segments at %d, %s : %s", pos, lit, err)
	}

	return *v, nil
}

func (p *Parser) parseName() (string, error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != STRING {
		return "", fmt.Errorf("found %s expected STRING at %d, %s", tok, pos, lit)
	}

	return lit, nil
}

func (p *Parser) parseDepends() (Dependency, error) {
	tok, pos, lit := p.scanIgnoreWhitespace()
	if tok != STRING {
		return Dependency{}, fmt.Errorf("found %s expected STRING at %d, %s", tok, pos, lit)
	}

	d := Dependency{}
	d.Name = lit
	// scan forward for the expected comma if there is a dep, if no comma no dep!
	tok, pos, lit = p.scanIgnoreWhitespace()
	if tok != COMMA {
		return d, nil
	}

	// scan forward for the dep
	tok, pos, lit = p.scanIgnoreWhitespace()
	if tok != STRING {
		return d, fmt.Errorf("found %s expected constraint STRING at %d, %s", tok, pos, lit)
	}

	// we have a dependency here
	c, err := version.NewConstraint(lit)
	if err != nil {
		return d, fmt.Errorf("error parsing constriant at %d, %s : %s", pos, lit, err)
	}
	d.Constraint = c

	return d, nil
}

// scan returns the next token from the underlying scanner.
func (p *Parser) scan() (tok Token, pos Pos, lit string) { return p.s.Scan() }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, pos Pos, lit string) {
	tok, pos, lit = p.scan()
	if tok == WS {
		tok, pos, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.s.Unscan() }
