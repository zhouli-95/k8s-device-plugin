/*
Copyright 2026 zhouli-95

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
