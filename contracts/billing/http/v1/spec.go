package v1

import _ "embed"

//go:embed openapi.yaml
var openAPISpec []byte

func OpenAPI() []byte { return openAPISpec }
