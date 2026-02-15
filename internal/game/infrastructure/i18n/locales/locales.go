package locales

import "embed"

// Files contains the embedded systemic translation JSON files.
//
//go:embed *.json
var Files embed.FS
