package openapi2go_test

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	openapi2go "github.com/unstoppablemango/openapi2go/pkg"
	"github.com/unstoppablemango/openapi2go/pkg/config"
	"github.com/unstoppablemango/openapi2go/pkg/openapi"
)

// generate parses spec and returns formatted Go source, or "" when no schemas.
func generate(fset *token.FileSet, spec string, conf config.Config) (string, error) {
	doc, err := openapi.ParseDocument([]byte(spec))
	if err != nil {
		return "", err
	}

	files, err := openapi2go.Generate(fset, doc, conf)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, files[0]); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func fieldSpec(title, fieldName, fieldType string) string {
	return fmt.Sprintf(`openapi: "3.1.0"
info:
  title: %s
  version: "1.0.0"
components:
  schemas:
    TestType:
      type: object
      properties:
        %s:
          type: %s
`, title, fieldName, fieldType)
}

var _ = Describe("Generate", func() {
	var fset *token.FileSet

	BeforeEach(func() {
		fset = token.NewFileSet()
	})

	Context("when spec has no components", func() {
		It("returns nil", func() {
			spec := `openapi: "3.1.0"
info:
  title: testpkg
  version: "1.0.0"
`
			src, err := generate(fset, spec, config.Config{
				PackageName:    "testpkg",
				FileNameSuffix: ".go",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(src).To(BeEmpty())
		})
	})

	Context("when spec has empty schemas", func() {
		It("returns a file with only the package declaration", func() {
			spec := `openapi: "3.1.0"
info:
  title: testpkg
  version: "1.0.0"
components:
  schemas: {}
`
			src, err := generate(fset, spec, config.Config{
				PackageName:    "testpkg",
				FileNameSuffix: ".go",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(src).To(HavePrefix("package testpkg"))
		})
	})

	Describe("field types", func() {
		DescribeTable("maps OpenAPI primitive types to Go identifiers",
			func(openAPIType, expectedGoType string) {
				src, err := generate(fset, fieldSpec("testpkg", "value", openAPIType), config.Config{
					PackageName:    "testpkg",
					FileNameSuffix: ".go",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(src).To(ContainSubstring(fmt.Sprintf("Value %s", expectedGoType)))
			},
			Entry("boolean maps to bool", "boolean", "bool"),
			Entry("integer maps to int", "integer", "int"),
			Entry("string maps to string", "string", "string"),
			Entry("object maps to any", "object", "any"),
		)

		Context("array field", func() {
			It("maps array of string items to []string", func() {
				spec := `openapi: "3.1.0"
info:
  title: testpkg
  version: "1.0.0"
components:
  schemas:
    TestType:
      type: object
      properties:
        items:
          type: array
          items:
            type: string
`
				src, err := generate(fset, spec, config.Config{
					PackageName:    "testpkg",
					FileNameSuffix: ".go",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(src).To(ContainSubstring("Items []string"))
			})

			It("maps array of integer items to []int", func() {
				spec := `openapi: "3.1.0"
info:
  title: testpkg
  version: "1.0.0"
components:
  schemas:
    TestType:
      type: object
      properties:
        counts:
          type: array
          items:
            type: integer
`
				src, err := generate(fset, spec, config.Config{
					PackageName:    "testpkg",
					FileNameSuffix: ".go",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(src).To(ContainSubstring("Counts []int"))
			})
		})
	})

	Describe("struct generation", func() {
		It("generates a type declaration for each schema", func() {
			spec := `openapi: "3.1.0"
info:
  title: testpkg
  version: "1.0.0"
components:
  schemas:
    Foo:
      type: object
      properties:
        name:
          type: string
    Bar:
      type: object
      properties:
        count:
          type: integer
`
			src, err := generate(fset, spec, config.Config{
				PackageName:    "testpkg",
				FileNameSuffix: ".go",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(src).To(ContainSubstring("type Foo struct"))
			Expect(src).To(ContainSubstring("type Bar struct"))
		})

		It("generates all fields in a schema", func() {
			spec := `openapi: "3.1.0"
info:
  title: testpkg
  version: "1.0.0"
components:
  schemas:
    Person:
      type: object
      properties:
        name:
          type: string
        age:
          type: integer
        active:
          type: boolean
`
			src, err := generate(fset, spec, config.Config{
				PackageName:    "testpkg",
				FileNameSuffix: ".go",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(src).To(MatchRegexp(`Name\s+string`))
			Expect(src).To(MatchRegexp(`Age\s+int`))
			Expect(src).To(MatchRegexp(`Active\s+bool`))
		})

		It("title-cases field names", func() {
			src, err := generate(fset, fieldSpec("testpkg", "firstname", "string"), config.Config{
				PackageName:    "testpkg",
				FileNameSuffix: ".go",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(src).To(ContainSubstring("Firstname string"))
		})
	})

	Describe("package name", func() {
		Context("when doc title contains no spaces", func() {
			It("uses the title as the package name", func() {
				src, err := generate(fset, fieldSpec("mypkg", "name", "string"), config.Config{
					PackageName:    "fallback",
					FileNameSuffix: ".go",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(src).To(HavePrefix("package mypkg"))
			})
		})

		Context("when doc title contains spaces", func() {
			It("falls back to config PackageName", func() {
				src, err := generate(fset, fieldSpec("My API", "name", "string"), config.Config{
					PackageName:    "myapi",
					FileNameSuffix: ".go",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(src).To(HavePrefix("package myapi"))
			})
		})

		Context("when doc title contains a tab", func() {
			It("falls back to config PackageName", func() {
				src, err := generate(fset, fieldSpec("My\tAPI", "name", "string"), config.Config{
					PackageName:    "tabfallback",
					FileNameSuffix: ".go",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(src).To(HavePrefix("package tabfallback"))
			})
		})
	})

	Describe("config", func() {
		Context("with a custom PackageName", func() {
			It("reflects PackageName in the generated file", func() {
				src, err := generate(fset, fieldSpec("My Spec", "name", "string"), config.Config{
					PackageName:    "customname",
					FileNameSuffix: ".go",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(src).To(HavePrefix("package customname"))
			})
		})

		Context("with per-type and per-field entries present", func() {
			It("generates without error", func() {
				spec := fieldSpec("testpkg", "name", "string")
				conf := config.Config{
					PackageName:    "testpkg",
					FileNameSuffix: ".go",
					Types: map[string]config.Type{
						"TestType": {
							Fields: map[string]config.Field{
								"name": {Type: "int64"},
							},
						},
					},
				}

				_, err := generate(fset, spec, conf)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
