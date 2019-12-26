package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	MagicCapKernel "github.com/magiccap/magiccap-uploaders-kernel"
	MagicCapKernelStandards "github.com/magiccap/magiccap-uploaders-kernel/standards"
)

func TestSXCURLMEUpload(t *testing.T) {
	k := MagicCapKernel.Kernel{
		Uploaders: map[string]*MagicCapKernelStandards.Uploader{},
	}
	f, err := ioutil.ReadFile("../uploaders/v1.json")
	if err != nil {
		t.Error("Failed to open V1 uploaders. Does it exist?")
	}
	var j map[string]interface{}
	err = json.Unmarshal(f, &j)
	if err != nil {
		t.Error(err)
	}
	err = k.Load(j)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadFile("magiccap.png")
	if err != nil {
		t.Error(err)
	}
	Uploader := k.Uploaders["sharex"]
	if Uploader == nil {
		t.Error("SXCU uploader is nil.")
	}
	if os.Getenv("SXCU_DATA") == "" {
		t.Skip("SXCU_DATA is not set. Skipping test!")
		return
	}
	map_ := map[string]interface{}{
		"sxcu_data": os.Getenv("SXCU_DATA"),
	}
	url, err := k.Uploaders["sharex"].Upload(map_, b, "magiccap.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
}
