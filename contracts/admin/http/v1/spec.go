package v1

import _ "embed"

//go:embed openapi.yaml
var openAPI []byte

func OpenAPI() []byte {
	return openAPI
}
