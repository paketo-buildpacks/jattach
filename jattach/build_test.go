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
	"testing"

	"github.com/paketo-buildpacks/jattach/jattach"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		build jattach.Build
		ctx   libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Application.Path = t.TempDir()
		Expect(err).NotTo(HaveOccurred())

		t.Setenv("BP_ARCH", "amd64")

		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "jattach"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "jattach",
					"version": "1.0.0",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []string{"cpe:2.3:a:jattach:jattach:1.0.0:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/jattach@1.0.0?arch=amd64",
				},
			},
		}
		ctx.StackID = "test-stack-id"
	})

	it("contributes JAttach", func() {
		result, err := build.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("jattach"))

		Expect(result.BOM.Entries).To(HaveLen(1))
	})
}
