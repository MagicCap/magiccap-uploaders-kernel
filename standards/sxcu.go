package MagicCapKernelStandards

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/dlclark/regexp2"
	"github.com/jakemakesstuff/structuredhttp"
	"github.com/yalp/jsonpath"
)

// SXCUSpec defines the SXCU spec for this uploader.
type SXCUSpec struct {
	DestinationType *string
	RequestMethod   *string // Newer version of below.
	RequestType     *string
	Body            *string
	Parameters      *map[string]string // Newer version of below.
	Arguments       *map[string]string
	FileFormName    *string // Needed for form request types.
	RequestURL      string
	URL             string
	Headers         *map[string]string
	RegexList       []string
}

// DollarSyntaxMatch is used to repersent a dollar syntax match.
type DollarSyntaxMatch struct {
	Start int
	End   int
	Match string
}

// GetDollarSyntax is used to get dollar syntax in the string.
func GetDollarSyntax(Origin string) *DollarSyntaxMatch {
	var MatchObj *DollarSyntaxMatch
	ToEscape := false
	for i, v := range Origin {
		if v == '\\' {
			if ToEscape {
				// Set the escaped char.
				ToEscape = false
				if MatchObj != nil {
					MatchObj.Match += string(v)
				}
				continue
			} else {
				// Set ToEscape to true and continue.
				ToEscape = true
				continue
			}
		} else if v == '$' {
			// Check for the start/end of the string.
			if ToEscape {
				// Check if the match object exists. If it's there, add an escaped dollar.
				if MatchObj != nil {
					MatchObj.Match += string(v)
				}
			} else if MatchObj == nil {
				// The start of the dollar syntax.
				MatchObj = &DollarSyntaxMatch{
					Match: "",
					Start: i,
					End:   0,
				}
			} else {
				// The end of the dollar syntax.
				MatchObj.End = i
				return MatchObj
			}
		} else {
			// Check if the match object exists. If it's there, add an character.
			if MatchObj != nil {
				MatchObj.Match += string(v)
			}
			ToEscape = false
		}
	}
	return MatchObj
}

// ShareXParamHandler is used to handle ShareX params.
func ShareXParamHandler(Data string, Filename string, Response string, ResponseURL string, Header *map[string]string, Regex []string) (string, error) {
	if Data == "response" {
		// Returns the response.
		return Response, nil
	} else if Data == "responseurl" || Data == "input" {
		// Returns the response URL. "input" is just here to keep compatibility.
		return ResponseURL, nil
	} else if Data == "filename" {
		// Returns the filename.
		return Filename, nil
	} else if strings.HasPrefix(Data, "header:") {
		// Returns the header.
		return (*Header)[strings.TrimPrefix(Data, "header:")], nil
	} else if strings.HasPrefix(Data, "json:") {
		// Returns the JSON path specified.
		s := strings.TrimPrefix(Data, "json:")
		var i interface{}
		err := json.Unmarshal([]byte(Response), &i)
		if err != nil {
			return "", err
		}
		res, err := jsonpath.Read(i, "$."+s)
		if err != nil {
			return "", err
		}
		s, ok := res.(string)
		if ok {
			return s, nil
		}
		b, err := json.Marshal(&i)
		if err != nil {
			return "", err
		}
		return string(b), nil
	} else if strings.HasPrefix(Data, "xml:") {
		// Returns the XML path specified.
		s := strings.TrimPrefix(Data, "xml:")
		doc, err := xmlquery.Parse(strings.NewReader(Response))
		if err != nil {
			return "", err
		}
		node := xmlquery.FindOne(doc, s)
		if node == nil {
			return "", nil
		}
		return node.Data, nil
	} else if strings.HasPrefix(Data, "base64:") {
		// Returns the string as base 64.
		s := strings.TrimPrefix(Data, "base64:")
		return base64.StdEncoding.EncodeToString([]byte(s)), nil
	} else if strings.HasPrefix(Data, "regex:") {
		// Returns the regex match.
		split := strings.Split(strings.TrimPrefix(Data, "regex:"), "|")
		if len(split) != 2 {
			return "", errors.New("Not enough arguments for regex.")
		}
		n, err := strconv.Atoi(split[0])
		if err != nil {
			return "", err
		}
		group := split[1]
		if n >= len(Regex) {
			return "", errors.New("Regex index does not exist.")
		}
		re, err := regexp2.Compile(Regex[n], regexp2.None)
		if err != nil {
			return "", err
		}
		m, err := re.FindStringMatch(Response)
		if err != nil {
			return "", err
		}
		capn, err := strconv.Atoi(group)
		if err == nil {
			g := m.GroupByNumber(capn)
			if g == nil {
				return "", errors.New("regex group does not exist")
			}
			return g.String(), nil
		} else {
			g := m.GroupByName(group)
			if g == nil {
				return "", errors.New("regex group does not exist")
			}
			return g.String(), nil
		}
	} else if strings.HasPrefix(Data, "random:") {
		// Handles random-ness.
		split := strings.Split(strings.TrimPrefix(Data, "regex:"), "|")
		return split[rand.Intn(len(split))], nil
	} else if strings.HasPrefix(Data, "prompt:") {
		// This just doesn't fail due to compatibility but adds a blank string instead.
		return "", nil
	} else if strings.HasPrefix(Data, "select:") {
		// This just doesn't fail due to compatibility but returns the first option instead.
		split := strings.Split(strings.TrimPrefix(Data, "regex:"), "|")
		return split[0], nil
	}

	// Unknown function.
	return "", errors.New("unknown function: " + Data)
}

// ShareXParamParse is used to parse a ShareX param.
func ShareXParamParse(Origin string, Filename string, Response string, ResponseURL string, Header *map[string]string, Regex []string) (string, error) {
	for {
		DollarSyntax := GetDollarSyntax(Origin)
		if DollarSyntax == nil {
			return Origin, nil
		} else {
			d, err := ShareXParamHandler(DollarSyntax.Match, Filename, Response, ResponseURL, Header, Regex)
			if err != nil {
				return "", err
			}
			Start := Origin[:DollarSyntax.Start]
			End := Origin[DollarSyntax.End+1:]
			Origin = Start + d + End
		}
	}
}

// SXCUInit defines the SXCU standard.
func SXCUInit(Structure UploaderStructure) (*Uploader, error) {
	SpecData, ok := Structure.Spec.(map[string]interface{})["sxcu_data"].(string)
	if !ok {
		return nil, errors.New("SXCU could not be found in the spec.")
	}
	return &Uploader{
		Description:   Structure.Description,
		Name:          Structure.Name,
		ConfigOptions: Structure.Config,
		Icon:          Structure.Icon,
		Upload: func(Config map[string]interface{}, Data []byte, Filename string) (string, error) {
			// Loads the spec.
			var spec SXCUSpec
			specr := strings.Replace(SpecData, "{sxcu_data}", Config["sxcu_data"].(string), 1)
			err := json.Unmarshal([]byte(specr), &spec)
			if err != nil {
				return "", err
			}

			// Deal with deprecated ShareX stuff.
			RequestMethodPtr := spec.RequestMethod
			if RequestMethodPtr == nil {
				RequestMethodPtr = spec.RequestType
			}
			RequestMethod := "POST"
			if RequestMethodPtr != nil {
				RequestMethod = *RequestMethodPtr
			}
			ParametersPtr := spec.Parameters
			if ParametersPtr == nil {
				ParametersPtr = spec.Arguments
			}
			Parameters := map[string]string{}
			if ParametersPtr != nil {
				Parameters = *ParametersPtr
			}

			// Check the destination type is compatible.
			DestinationType := "ImageUploader, TextUploader, FileUploader"
			if spec.DestinationType != nil {
				DestinationType = *spec.DestinationType
			}
			SupportedDestinationTypes, err := csv.NewReader(strings.NewReader(DestinationType)).Read()
			if err != nil {
				return "", err
			}
			Supported := false
			fnparts := strings.Split(Filename, ".")
			ext := strings.ToLower(fnparts[len(fnparts)-1])
			for _, v := range SupportedDestinationTypes {
				v = strings.Trim(v, " ")
				if v == "FileUploader" {
					// FileUploader implies the user can upload anything.
					Supported = true
					break
				} else if v == "ImageUploader" {
					// ImageUploader restricts the user to uploading only images if there is no other conditions.
					switch ext {
					case "png":
					case "jpeg":
					case "jpg":
					case "gif":
					case "tiff":
						Supported = true
						break
					default:
						// Do nothing.
					}
					if Supported {
						break
					}
				} else if v == "TextUploader" && (ext == "txt" || ext == "md") {
					// ImageUploader restricts the user to uploading only text if there is no other conditions.
					Supported = true
					break
				}
			}

			// Gets the request URL.
			RequestURL, err := ShareXParamParse(spec.RequestURL, Filename, "", "", spec.Headers, spec.RegexList)
			if err != nil {
				return "", err
			}

			// Creates the request handler.
			Handler := &structuredhttp.Request{
				URL:     RequestURL,
				Method:  RequestMethod,
				Headers: map[string]string{},
			}

			// Handles URL params.
			for k, v := range Parameters {
				Key, err := ShareXParamParse(k, Filename, "", "", spec.Headers, spec.RegexList)
				if err != nil {
					return "", err
				}
				Value, err := ShareXParamParse(v, Filename, "", "", spec.Headers, spec.RegexList)
				if err != nil {
					return "", err
				}
				Handler = Handler.Query(Key, Value)
			}

			// Handles headers.
			if spec.Headers != nil {
				for k, v := range *spec.Headers {
					Key, err := ShareXParamParse(k, Filename, "", "", spec.Headers, spec.RegexList)
					if err != nil {
						return "", err
					}
					Value, err := ShareXParamParse(v, Filename, "", "", spec.Headers, spec.RegexList)
					if err != nil {
						return "", err
					}
					Handler = Handler.Header(Key, Value)
				}
			}

			// Adds a 10 second timeout.
			Handler = Handler.Timeout(10 * time.Second)

			// Handle the body/content type.
			ContentType := "multipart/form-data"
			if spec.Body != nil {
				if *spec.Body == "FormURLEncoded" {
					ContentType = "application/x-www-form-urlencoded"
				} else if *spec.Body == "JSON" {
					ContentType = "application/json"
				} else if *spec.Body == "XML" {
					ContentType = "application/xml"
				} else if *spec.Body == "Binary" {
					ContentType = "application/octet-stream"
				}
			}
			if ContentType != "application/octet-stream" && spec.FileFormName == nil {
				return "", errors.New("Blank file form name.")
			}
			FileFormName := *spec.FileFormName
			IsMultipart := false
			if ContentType == "multipart/form-data" {
				// Handle normal encoded data.
				buffer := new(bytes.Buffer)
				writer := multipart.NewWriter(buffer)
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, FileFormName, Filename))
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
				Handler = Handler.MultipartForm(buffer, writer.FormDataContentType())
				IsMultipart = true
			} else if ContentType == "application/x-www-form-urlencoded" {
				// Handle adding URL encoded data.
				Handler = Handler.URLEncodedForm(url.Values{
					FileFormName: []string{
						url.QueryEscape(string(Data)),
					},
				})
			} else if ContentType == "application/json" {
				// Creates the JSON body.
				Handler = Handler.JSON(&map[string]string{
					FileFormName: base64.StdEncoding.EncodeToString(Data),
				})
			} else if ContentType == "application/xml" {
				// Creates the XML body.
				x, err := xml.Marshal(&map[string]string{
					FileFormName: base64.StdEncoding.EncodeToString(Data),
				})
				if err != nil {
					return "", err
				}
				Handler = Handler.Bytes(x)
			} else {
				// Handle a binary stream. Very simple.
				Handler = Handler.Bytes(Data)
			}
			if !IsMultipart {
				Handler = Handler.Header("Content-Type", ContentType)
			}

			// Run the request.
			r, err := Handler.Run()
			if err != nil {
				return "", err
			}

			// Handle the status.
			err = r.RaiseForStatus()
			s, terr := r.Text()
			if terr != nil {
				return "", err
			}
			if err != nil {
				return s, err
			}

			// Return the URL.
			loc := r.RawResponse.Request.URL.String()
			Headers := map[string]string{}
			for k, v := range r.RawResponse.Header {
				Headers[k] = v[0]
			}
			Result, err := ShareXParamParse(spec.URL, Filename, s, loc, &Headers, spec.RegexList)
			if err != nil {
				return "", err
			}
			return Result, nil
		},
	}, nil
}
