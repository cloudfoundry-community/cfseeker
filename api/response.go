package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/cfseeker/seeker"
)

//Response is the structured response of the CFSeeker API
type Response struct {
	Meta     *Metadata     `json:"meta,omitempty"`
	Contents seeker.Output `json:"contents,omitempty"`
	writer   http.ResponseWriter
	code     int
}

//Metadata has additional information about an API Response, like error messages
type Metadata struct {
	Error   string `json:"error,omitempty"`
	Warning string `json:"warning,omitempty"`
	Message string `json:"message,omitempty"`
}

//NewResponse returns a pointer to an empty Response struct
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{writer: w}
}

//Code sets the HTTP response code, which will be used if write is called.
func (r *Response) Code(code int) *Response {
	r.code = code
	return r
}

// Message takes the receiver Response, attaches the given message as a
// informational message, and then returns the Response object.
func (r *Response) Message(m string) *Response {
	if r.Meta == nil {
		r.Meta = &Metadata{}
	}
	r.Meta.Message = m
	return r
}

//Err takes the receiver Response, attaches the given message as an error,
// and then returns the Response object.
func (r *Response) Err(e string) *Response {
	if r.Meta == nil {
		r.Meta = &Metadata{}
	}
	r.Meta.Error = e
	return r
}

//Warn takes the receiver Response, attaches the given message as a warning,
// and then returns the Response object.
func (r *Response) Warn(w string) *Response {
	if r.Meta == nil {
		r.Meta = &Metadata{}
	}
	r.Meta.Warning = w
	return r
}

//AttachContents takes the given interface and assigns it as the response contents
func (r *Response) AttachContents(c seeker.Output) *Response {
	r.Contents = c
	return r
}

//Bytes marshals the Response to JSON and returns the result as a byte array
func (r Response) Bytes() []byte {
	respBytes, err := json.Marshal(&r)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal response from object %#v", r))
	}
	return respBytes
}

//String marshals the Response to JSON and returns the result as a string
func (r Response) String() string {
	respBytes, err := json.Marshal(&r)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal response from object %#v", r))
	}
	return string(respBytes)
}

func (r *Response) Write() {
	r.writer.Header().Set("Content-Type", "application/json")
	if r.code != 0 {
		r.writer.WriteHeader(r.code)
	}
	r.writer.Write(r.Bytes())
}
