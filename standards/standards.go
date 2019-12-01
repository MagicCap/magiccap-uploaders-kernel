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
	Icon          string                                                                            `json:"icon"`
	Description   string                                                                            `json:"description"`
	Name          string                                                                            `json:"name"`
	ConfigOptions map[string]ConfigOption                                                           `json:"configOptions"`
	Upload        func(Config map[string]interface{}, Data []byte, Filename string) (string, error) `json:"-"`
}

// UploaderStructure defines the structure that an imported uploaders JSON will follow.
type UploaderStructure struct {
	Icon        string                  `json:"icon"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Config      map[string]ConfigOption `json:"config"`
	Spec        interface{}             `json:"spec"`
}

// GetStandardsMap gets the standards map.
func GetStandardsMap() map[string]func(Structure UploaderStructure) (*Uploader, error) {
	return map[string]func(Structure UploaderStructure) (*Uploader, error){
		"http": HTTPInit,
		"ftp":  FTPInit,
		"s3":   S3Init,
	}
}
