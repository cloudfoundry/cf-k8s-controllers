package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"

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
    Request body:
        %+v
    HTTP Status code: %d
    Response body:
        %s`,
		m.Request.Method, m.Request.URL,
		objectToPrettyJson(m.Request.Body),
		m.StatusCode(),
		formatAsPrettyJson(m.Body()),
	)
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
