package main

import (
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/cobra"
)

func main() {
	var (
		retrievalPaths map[string]string
		outputPath     string
		pkg            string
	)

	cmd := cobra.Command{
		Use:          "jsg [flags] FILE ...",
		Short:        "json schema code generator",
		Long:         `json schema code generator`,
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var loader Loader
			for retrievalURI, retrievalPath := range retrievalPaths {
				loader.AddRetrievalPath(retrievalURI, retrievalPath)
			}

			compiler := jsonschema.NewCompiler()
			compiler.UseLoader(&loader)

			var schemas []*jsonschema.Schema
			for _, path := range args {
				schema, err := compiler.Compile(path)
				if err != nil {
					return err
				}
				schemas = append(schemas, schema)
			}

			var mapper Mapper
			file := File{Pkg: pkg}

			for _, schema := range schemas {
				typ, err := mapper.Map(schema)
				if err != nil {
					return err
				}
				file.Types = append(file.Types, typ)
			}

			var generator Generator
			return generator.Generate(os.Stdout, &file)
		},
	}

	cmd.Flags().StringToStringVarP(&retrievalPaths, "retrieval-path", "r", nil, "mappings of retrieval uris to local file paths for schema loading")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "path to write generated files to instead of stdout")
	cmd.Flags().StringVarP(&pkg, "package", "p", "main", "package name for generated types")
	cmd.Execute()
}
