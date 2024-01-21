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

package jattach_test

import (
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/jattach/jattach"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/libpak"
)

func testJAttach(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Layers.Path = t.TempDir()
		Expect(err).NotTo(HaveOccurred())
	})

	it("contributes JAttach", func() {
		dep := libpak.BuildpackDependency{
			URI:    "https://localhost/jattach.tgz",
			SHA256: "61f13e4c60ad63295453a5e1c92e04ad34bab09871b622b5aac0e72236431000",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		j, _ := jattach.NewJAttach(dep, dc)
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = j.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.LayerTypes.Build).To(BeFalse())
		Expect(layer.LayerTypes.Cache).To(BeFalse())
		Expect(layer.LayerTypes.Launch).To(BeTrue())

		Expect(filepath.Join(layer.Path, "bin", "jattach")).To(BeARegularFile())
	})

}
