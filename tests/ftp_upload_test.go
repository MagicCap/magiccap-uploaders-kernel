package tests

import (
	"encoding/json"
	"io/ioutil"
	"magiccap-uploaders-kernel"
	"magiccap-uploaders-kernel/standards"
	"os"
	"testing"
)

func TestFTPUpload(t *testing.T) {
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
	Uploader := k.Uploaders["ftp"]
	if Uploader == nil {
		t.Error("FTP is nil.")
	}
	if os.Getenv("FTP_HOSTNAME") == "" {
		t.Skip("FTP_HOSTNAME is not set. Skipping test!")
		return
	}
	map_ := map[string]interface{}{
		"ftp_hostname": os.Getenv("FTP_HOSTNAME"),
		"ftp_port": 21,
		"ftp_username": "anonymous",
		"ftp_password": "anonymous",
		"ftp_directory": "/",
		"ftp_domain": "http://example.com/",
	}
	_, err = k.Uploaders["ftp"].Upload(map_, b, "magiccap.png")
	if err != nil {
		t.Error(err)
	}
	t.Log("Uploaded to the FTP server successfully.")
}
