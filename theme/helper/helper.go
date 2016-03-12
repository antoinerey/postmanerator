package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"text/template"

	"regexp"

	"github.com/aubm/postmanerator/postman"
	"github.com/russross/blackfriday"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"findRequest":  findRequest,
		"findResponse": findResponse,
		"markdown":     markdown,
		"randomID":     randomID,
		"indentJSON":   indentJSON,
		"curlSnippet":  curlSnippet,
	}
}

func findRequest(requests []postman.Request, ID string) *postman.Request {
	for _, r := range requests {
		if r.ID == ID {
			return &r
		}
	}
	return nil
}

func findResponse(req postman.Request, name string) *postman.Response {
	for _, res := range req.Responses {
		if res.Name == name {
			return &res
		}
	}
	return nil
}

func markdown(input string) string {
	return string(blackfriday.MarkdownBasic([]byte(input)))
}

func randomID() int {
	return rand.Intn(999999999)
}

func indentJSON(input string) (string, error) {
	dest := new(bytes.Buffer)
	src := []byte(input)
	err := json.Indent(dest, src, "", "    ")
	return dest.String(), err
}

func curlSnippet(request postman.Request) string {
	var curlSnippet string
	payloadReady, _ := regexp.Compile("POST|PUT|PATCH|DELETE")
	curlSnippet += fmt.Sprintf("curl -X %v", request.Method)

	if payloadReady.MatchString(request.Method) {
		if request.DataMode == "urlencoded" {
			curlSnippet += ` -H "Content-Type: application/x-www-form-urlencoded"`
		} else if request.DataMode == "params" {
			curlSnippet += ` -H "Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW"`
		}
	}

	for _, header := range request.Headers() {
		curlSnippet += fmt.Sprintf(` -H "%v: %v"`, header.Name, header.Value)
	}

	if payloadReady.MatchString(request.Method) {
		if request.DataMode == "raw" && request.RawModeData != "" {
			curlSnippet += fmt.Sprintf(` -d '%v'`, request.RawModeData)
		} else if len(request.Data) > 0 {
			if request.DataMode == "urlencoded" {
				var dataList []string
				for _, data := range request.Data {
					dataList = append(dataList, fmt.Sprintf("%v=%v", data.Key, data.Value))
				}
				curlSnippet += fmt.Sprintf(` -d "%v"`, strings.Join(dataList, "&"))
			} else if request.DataMode == "params" {
				for _, data := range request.Data {
					curlSnippet += fmt.Sprintf(` -F "%v=%v"`, data.Key, data.Value)
				}
			}
		}
	}

	curlSnippet += fmt.Sprintf(` "%v"`, request.URL)
	return curlSnippet
}