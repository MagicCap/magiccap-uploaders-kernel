package MagicCapKernelStandards

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/ajg/form.v1"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

// POSTAs defines what the type of the POST request is.
type POSTAs struct {
	Type string `json:"type"`
	Key string `json:"key"`
}

// HTTPSpec defines the HTTP spec for this uploader.
type HTTPSpec struct {
	Method string `json:"method"`
	URL string `json:"url"`
	POSTAs *POSTAs `json:"post_as"`
	Headers *map[string]string `json:"headers"`
	ResponseKey *string `json:"response_key"`
}

// HTTPInit defines the HTTP standard.
func HTTPInit(Structure UploaderStructure) (*Uploader, error) {
	b, err := json.Marshal(Structure.Spec)
	if err != nil {
		return nil, err
	}
	var spec HTTPSpec
	err = json.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}
	e := base64.Encoding{}
	return &Uploader{
		Description:   Structure.Description,
		Name:          Structure.Name,
		ConfigOptions: Structure.Config,
		Icon:          Structure.Icon,
		Upload: func(Config map[string]interface{}, Data []byte, Filename string) (string, error) {
			var URL string
			var POSTData *bytes.Buffer
			if spec.POSTAs.Type == "b64" {
				URL = spec.URL + "?" + spec.POSTAs.Key + "=" + e.EncodeToString(Data)
			} else if spec.POSTAs.Type == "raw" {
				POSTData = bytes.NewBuffer(Data)
			} else if spec.POSTAs.Type == "multipart" {
				buffer := new(bytes.Buffer)
				writer := multipart.NewWriter(buffer)
				if err != nil {
					return "", err
				}
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, spec.POSTAs.Key, Filename))
				h.Set("Content-Type", http.DetectContentType(Data))
				part, err := writer.CreatePart(h)
				if err != nil {
					return "", err
				}
				_, err = part.Write(Data)
				if err != nil {
					return "", err
				}
				err = writer.Close()
				if err != nil {
					return "", err
				}
				POSTData = buffer
				URL = spec.URL
			} else if spec.POSTAs.Type == "urlencoded" {
				u, err := form.EncodeToString(map[string]interface{}{
					spec.POSTAs.Key: Data,
				})
				if err != nil {
					return "", err
				}
				POSTData = bytes.NewBufferString(u)
				URL = spec.URL
			} else {
				return "", errors.New("POST type not defined.")
			}
			r, err := http.NewRequest(spec.Method, URL, POSTData)
			if err != nil {
				return "", err
			}
			if spec.Headers != nil {
				for k, v := range *spec.Headers {
					r.Header.Set(k, v)
				}
			}
			client := http.Client{
				Timeout: 30 * time.Second,
			}
			resp, err := client.Do(r)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", nil
			}
			ResponseType := math.Floor(float64(resp.StatusCode) / 100)
			if ResponseType == 4 || ResponseType == 5 {
				return "", errors.New("Uploader returned the status " + strconv.Itoa(resp.StatusCode) + ".")
			}
			if spec.ResponseKey == nil {
				return string(b), nil
			} else {
				var JSONMap map[string]interface{}
				err := json.Unmarshal(b, JSONMap)
				if err != nil {
					return "", err
				}
				Key := strings.Split(*spec.ResponseKey, ".")
				MapContext := JSONMap
				Last, Key := Key[len(Key)-1], Key[:len(Key)-1]
				var ok bool
				for _, v := range Key {
					MapContext, ok = MapContext[v].(map[string]interface{})
					if !ok {
						return "", errors.New("A value in the uploader is not a string map.")
					}
				}
				s, ok := MapContext[Last].(string)
				if !ok {
					return "", errors.New("The final value in the uploader is not a string.")
				}
				return s, nil
			}
		},
	}, nil
}
