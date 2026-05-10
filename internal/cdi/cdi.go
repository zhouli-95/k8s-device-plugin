package cdi

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	cdiapi "tags.cncf.io/container-device-interface/pkg/cdi"
	cdiparser "tags.cncf.io/container-device-interface/pkg/parser"
	"tags.cncf.io/container-device-interface/specs-go"
)

const (
	cdiRoot = "/var/run/cdi"
)

type CDIHandler struct {
	Vendor string
	Class  string
	spec   *specs.Spec
}

type CDIInjector func(*specs.Spec) error

func NewCDIHandler(opts ...Option) *CDIHandler {
	cdi := &CDIHandler{}
	for _, opt := range opts {
		opt(cdi)
	}
	return cdi
}

func (c *CDIHandler) CreateSpec(handler CDIInjector) error {
	spec := specs.Spec{
		Version:        "0.6.0",
		Kind:           fmt.Sprintf("%s/%s", c.Vendor, c.Class),
		ContainerEdits: specs.ContainerEdits{},
	}

	if err := handler(&spec); err != nil {
		return err
	}
	c.spec = &spec
	if err := c.WriteSpec(); err != nil {
		return err
	}
	return nil
}

func (c *CDIHandler) WriteSpec() error {
	specName, err := cdiapi.GenerateNameForSpec(c.spec)
	if err != nil {
		return err
	}

	cdiPath := filepath.Join(cdiRoot, specName+".json")
	specDir := filepath.Dir(cdiPath)
	cache, _ := cdiapi.NewCache(
		cdiapi.WithAutoRefresh(false),
		cdiapi.WithSpecDirs(specDir),
	)
	if err := cache.WriteSpec(c.spec, filepath.Base(cdiPath)); err != nil {
		return fmt.Errorf("failed to write spec: %v", err)
	}

	return nil
}

func (c *CDIHandler) QualifiedName(class, id string) string {
	return cdiparser.QualifiedName(c.Vendor, class, id)
}

func (c *CDIHandler) GetCDIAnnotation(devices []string) (map[string]string, error) {
	return cdiapi.UpdateAnnotations(map[string]string{}, fmt.Sprintf("%s.%s", c.Vendor, c.Class), uuid.New().String(), devices)
}
