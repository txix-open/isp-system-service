package docs

import _ "embed"

// Embedded swagger files.
var (
	//go:embed swagger.yaml
	SwaggerYaml []byte
)
