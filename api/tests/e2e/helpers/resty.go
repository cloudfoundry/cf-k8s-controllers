package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-http-utils/headers"
	"github.com/go-resty/resty/v2"
)

type actualRestyResponse struct {
	*resty.Response
}

func NewActualRestyResponse(response *resty.Response) actualRestyResponse {
	return actualRestyResponse{Response: response}
}

func (m actualRestyResponse) GomegaString() string {
	return fmt.Sprintf(`
    Request: %s %s
    Request headers:
        %s
    Request body:
        %+v
    HTTP Status code: %d
    Response headers:
        %s
    Response body:
        %s`,
		m.Request.Method, m.Request.URL,
		headersString(m.Request.Header),
		objectToPrettyJson(m.Request.Body),
		m.StatusCode(),
		headersString(m.Header()),
		formatAsPrettyJson(m.Body()),
	)
}

func headersString(headers http.Header) string {
	var s []string
	for k := range headers {
		s = append(s, formatHeader(k, headers.Get(k)))
	}
	return strings.Join(s, "\n        ")
}

func formatHeader(key, value string) string {
	if key == headers.Authorization {
		return fmt.Sprintf("%s: %s", key, formatAuthorizationValue(value))
	}

	return fmt.Sprintf("%s: %s", key, value)
}

func formatAuthorizationValue(value string) string {
	substrings := strings.Split(value, " ")
	if len(substrings) != 2 {
		return value
	}

	return fmt.Sprintf("%s %s", substrings[0], "<Redacted>")
}

func objectToPrettyJson(obj interface{}) string {
	prettyJson, err := json.MarshalIndent(obj, "        ", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}

	return string(prettyJson)
}

func formatAsPrettyJson(b []byte) string {
	var prettyBuf bytes.Buffer
	if err := json.Indent(&prettyBuf, b, "        ", "  "); err != nil {
		return string(b)
	}

	return prettyBuf.String()
}
