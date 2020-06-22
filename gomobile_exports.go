// +build ios android

package MagicCapKernel

import "C"
import (
	"encoding/json"

	MagicCapKernelStandards "github.com/magiccap/magiccap-uploaders-kernel/standards"
)

var cKernel *Kernel

// InitKernel is used to initialise a local variable with a kernel with the bytes specified.
//export InitKernel
func InitKernel(Data []byte) error {
	var x map[string]interface{}
	err := json.Unmarshal(Data, &x)
	if err != nil {
		return err
	}
	k := Kernel{Uploaders: map[string]*MagicCapKernelStandards.Uploader{}}
	err = k.Load(x)
	if err != nil {
		return err
	}
	cKernel = &k
	return nil
}

// GetUploaderIDs is used to get a byte array containing a JSON array of the ID's of uploaders.
//export GetUploaderIDs
func GetUploaderIDs() []byte {
	x := make([]string, len(cKernel.Uploaders))
	i := 0
	for k := range cKernel.Uploaders {
		x[i] = k
		i++
	}
	b, _ := json.Marshal(x)
	return b
}

// GetUploader is used to get bytes which is a JSON object of the uploader by ID.
//export GetUploader
func GetUploader(ID string) []byte {
	b, _ := json.Marshal(cKernel.Uploaders[ID])
	return b
}

// UploadFile is used to upload a file.
//export UploadFile
func UploadFile(UploaderID, Filename string, Data, Config []byte) (string, error) {
	var x map[string]interface{}
	err := json.Unmarshal(Config, &x)
	if err != nil {
		return "", err
	}
	return cKernel.Uploaders[UploaderID].Upload(x, Data, Filename)
}
