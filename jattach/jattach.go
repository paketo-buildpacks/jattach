/*
 * Copyright 2018-2024 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jattach

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/libpak/crush"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type JAttach struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJAttach(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (JAttach, libcnb.BOMEntry) {
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return JAttach{LayerContributor: contributor}, entry
}

func (j JAttach) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		binDir := filepath.Join(layer.Path, "bin")

		if err := os.MkdirAll(binDir, 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to mkdir\n%w", err)
		}
		j.Logger.Bodyf("Copying to %s", binDir)

		if err := crush.Extract(artifact, binDir, 0); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to expand jattach\n%w", err)
		}

		file := filepath.Join(binDir, "jattach")
		if err := os.Chmod(file, 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to chmod %s\n%w", file, err)
		}

		return layer, nil
	})
}

func (j JAttach) Name() string {
	return j.LayerContributor.LayerName()
}
