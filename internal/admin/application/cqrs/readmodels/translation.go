package readmodels

// Translation is the admin read model for a single localised string entry.
type Translation struct {
	Key    string `json:"key"`
	Locale string `json:"locale"`
	Value  string `json:"value"`
}
