package MagicCapKernelStandards

// ConfigOption defines a config option.
type ConfigOption struct {
	Required *bool        `json:"required"`
	Type     string       `json:"type"`
	Value    string       `json:"value"`
	Default  *interface{} `json:"default"`
}

// Uploader defines the uploader in its final form.
type Uploader struct {
	Icon          string
	Description   string
	Name          string
	ConfigOptions map[string]ConfigOption
	Upload        func(Config map[string]interface{}, Data []byte, Filename string) (string, error)
}

// UploaderStructure defines the structure that an imported uploaders JSON will follow.
type UploaderStructure struct {
	Icon           string                  `json:"icon"`
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	Config         map[string]ConfigOption `json:"config"`
	Spec           interface{}             `json:"spec"`
}

// GetStandardsMap gets the standards map.
func GetStandardsMap() map[string]func(Structure UploaderStructure) (*Uploader, error) {
	return map[string]func(Structure UploaderStructure) (*Uploader, error){
		"http": HTTPInit,
	}
}
