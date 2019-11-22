package MagicCapKernelStandards

import (
	"bytes"
	"encoding/json"
	"github.com/jlaffaye/ftp"
	"github.com/magiccap/magiccap-uploaders-kernel/utils"
	"time"
)

// FTPSpec defines the FTP spec for this uploader.
type FTPSpec struct {
	Hostname string `json:"hostname"`
	Port string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Directory string `json:"directory"`
	BaseURL string `json:"base_url"`
}

// FTPInit defines the FTP standard.
func FTPInit(Structure UploaderStructure) (*Uploader, error) {
	b, err := json.Marshal(Structure.Spec)
	if err != nil {
		return nil, err
	}
	var spec FTPSpec
	err = json.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}
	return &Uploader{
		Description:   Structure.Description,
		Name:          Structure.Name,
		ConfigOptions: Structure.Config,
		Icon:          Structure.Icon,
		Upload: func(Config map[string]interface{}, Data []byte, Filename string) (string, error) {
			PortStr, err := utils.SubString(spec.Port, Config, Filename)
			if err != nil {
				return "", err
			}
			Hostname, err := utils.SubString(spec.Hostname, Config, Filename)
			if err != nil {
				return "", err
			}
			Username, err := utils.SubString(spec.Username, Config, Filename)
			if err != nil {
				return "", err
			}
			Password, err := utils.SubString(spec.Password, Config, Filename)
			if err != nil {
				return "", err
			}
			Directory, err := utils.SubString(spec.Directory, Config, Filename)
			if err != nil {
				return "", err
			}
			BaseURL, err := utils.SubString(spec.BaseURL, Config, Filename)
			if err != nil {
				return "", err
			}

			conn, err := ftp.Dial(Hostname + ":" + PortStr, ftp.DialWithTimeout(5 * time.Second))
			if err != nil {
				return "", err
			}
			err = conn.Login(Username, Password)
			if err != nil {
				return "", err
			}
			err = conn.Stor(Directory + Filename, bytes.NewBuffer(Data))
			if err != nil {
				return "", err
			}
			err = conn.Quit()
			if err != nil {
				return "", err
			}
			return BaseURL + Filename, nil
		},
	}, nil
}
