package locales

import "embed"

// Files contains the embedded translation JSON files for the Admin service.
//
//go:embed *.json
var Files embed.FS
