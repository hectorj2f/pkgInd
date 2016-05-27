package indexer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hectorj2f/pkgind/log"
	"github.com/hectorj2f/pkgind/transport"
)

// I store the dependencies as the packages that use this package to facilitate
// the operations Index, Query and Remove
type Dependencies []string

type PackageManager struct {
	sync.Mutex // guards the fields below
	packages   map[string]Dependencies
}

func NewPackageManager() *PackageManager {
	pkgMgmt := &PackageManager{
		packages: make(map[string]Dependencies, 0),
	}
	return pkgMgmt
}

func (m *PackageManager) Size() int {
	return len(m.packages)
}

// Index index a package and adds its parents to the list
// Return a response code and an error
func (m *PackageManager) Index(msg *transport.MessageRequest) (transport.MessageResponseCode, error) {
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

	m.packages[msg.Package] = make([]string, 0)
	for _, dependency := range msg.Dependencies {
		m.packages[dependency] = append(m.packages[dependency], msg.Package)
	}

	return transport.OK, nil
}

// Remove deletes an indexed package from the list
// Return a response code and an error
func (m *PackageManager) Remove(msg *transport.MessageRequest) (transport.MessageResponseCode, error) {
	m.Lock()
	defer m.Unlock()

	if pkg, ok := m.packages[msg.Package]; ok {
		if len(pkg) > 0 {
			for _, name := range pkg {
				if _, ok := m.packages[name]; ok {
					log.Logger().Infof("CANNOT the package exists  %s", name)
					return transport.FAIL, errors.New(fmt.Sprintf("package '%s' is required by '%v'", msg.Package, name))
				}
			}
		}
		delete(m.packages, msg.Package)
	}

	return transport.OK, nil
}

// Query searches for a package in the indexed list
// Return a response code and an error
func (m *PackageManager) Query(msg *transport.MessageRequest) (transport.MessageResponseCode, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.packages[msg.Package]; ok {
		return transport.OK, nil
	}

	return transport.FAIL, errors.New(fmt.Sprintf("package '%s' is not indexed", msg.Package))
}
