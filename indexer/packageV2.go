package indexer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hectorj2f/pkgind/transport"
)

// I store the dependencies as the packages that use this package to facilitate
// the operations Index, Query and Remove
type PackageV2 struct {
	Name         string
	Dependencies []*PackageV2
	Consumers    []*PackageV2
}

type PackageManagerV2 struct {
	sync.Mutex // guards the fields below
	packages   map[string]*PackageV2
}

func NewPackageManagerV2() *PackageManagerV2 {
	pkgMgmt := &PackageManagerV2{
		packages: make(map[string]*PackageV2, 0),
	}
	return pkgMgmt
}

func (m *PackageManagerV2) Size() int {
	return len(m.packages)
}

// Index index a package and adds its parents to the list
// Return a response code and an error
func (m *PackageManagerV2) Index(msg *transport.MessageRequest) (transport.MessageResponseCode, error) {
	m.Lock()
	defer m.Unlock()
	// package already exists
	if _, ok := m.packages[msg.Package]; ok {
		return transport.OK, nil
	}

	for _, dependency := range msg.Dependencies {
		if _, ok := m.packages[dependency]; !ok {
			return transport.FAIL, errors.New(fmt.Sprintf("dependency '%s' is required to be installed '%s'", dependency, msg.Package))
		}
	}

	m.packages[msg.Package] = &PackageV2{
		Name:         msg.Package,
		Dependencies: make([]*PackageV2, 0),
		Consumers:    make([]*PackageV2, 0),
	}
	for _, dependency := range msg.Dependencies {
		m.packages[msg.Package].Dependencies = append(m.packages[msg.Package].Dependencies, m.packages[dependency])
		m.packages[dependency].Consumers = append(m.packages[dependency].Consumers, m.packages[msg.Package])
	}

	return transport.OK, nil
}

// Remove deletes an indexed package from the list
// Return a response code and an error
func (m *PackageManagerV2) Remove(msg *transport.MessageRequest) (transport.MessageResponseCode, error) {
	m.Lock()
	defer m.Unlock()

	if pkg, ok := m.packages[msg.Package]; ok {
		if len(pkg.Consumers) > 0 {
			for _, p := range pkg.Consumers {
				if _, ok := m.packages[p.Name]; ok {
					return transport.FAIL, errors.New(fmt.Sprintf("package '%s' is required by '%v'", msg.Package, p.Name))
				}
			}
		}
		delete(m.packages, msg.Package)
	}

	return transport.OK, nil
}

// Query searches for a package in the indexed list
// Return a response code and an error
func (m *PackageManagerV2) Query(msg *transport.MessageRequest) (transport.MessageResponseCode, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.packages[msg.Package]; ok {
		return transport.OK, nil
	}

	return transport.FAIL, errors.New(fmt.Sprintf("package '%s' is not indexed", msg.Package))
}
