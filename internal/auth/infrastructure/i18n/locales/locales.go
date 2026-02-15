package locales

import "embed"

// Files contains the embedded translation JSON files for Auth service.
//
//go:embed *.json
var Files embed.FS
