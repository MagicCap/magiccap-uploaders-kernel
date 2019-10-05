package MagicCapKernel

import (
	"encoding/json"
	"errors"
	"magiccap-uploaders-kernel/standards"
)

// Defines the kernel structure.
type Kernel struct {
	Uploaders map[string]*MagicCapKernelStandards.Uploader
}

// Loads a map of uploaders from V1.
func (k Kernel) Load(V1File map[string]interface{}) error {
	StandardsMap := MagicCapKernelStandards.GetStandardsMap()
	for Implementation, UploaderMap := range V1File {
		Loader := StandardsMap[Implementation]
		if Loader == nil {
			return errors.New("Loader not found.")
		}
		for key, v := range UploaderMap.(map[string]interface{}) {
			b, err := json.Marshal(v)
			if err != nil {
				return err
			}
			var std MagicCapKernelStandards.UploaderStructure
			err = json.Unmarshal(b, &std)
			if err != nil {
				return err
			}
			uploader, err := Loader(std)
			if err != nil {
				return err
			}
			k.Uploaders[key] = uploader
		}
	}
	return nil
}
