package MagicCapKernel

import (
	"encoding/json"

	MagicCapKernelStandards "github.com/magiccap/magiccap-uploaders-kernel/standards"
)

// Kernel defines the kernel structure.
type Kernel struct {
	Uploaders map[string]*MagicCapKernelStandards.Uploader
}

// Load loads a map of uploaders from V1.
func (k Kernel) Load(V1File map[string]interface{}) error {
	StandardsMap := MagicCapKernelStandards.GetStandardsMap()
	for Implementation, UploaderMap := range V1File {
		Loader := StandardsMap[Implementation]
		if Loader == nil {
			continue
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
