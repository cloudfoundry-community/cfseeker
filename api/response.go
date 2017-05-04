package api

import (
	"encoding/json"
	"fmt"
)

//Response is the structured response of the CFSeeker API
type Response struct {
	Meta     *Metadata   `json:"meta,omitempty"`
	Contents interface{} `json:"contents,omitempty"`
}

//Metadata has additional information about an API Response, like error messages
type Metadata struct {
	Error   string `json:"error,omitempty"`
	Warning string `json:"warning,omitempty"`
}

//NewResponse returns a pointer to an empty Response struct
func NewResponse() *Response {
	return &Response{}
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
func (r *Response) AttachContents(c interface{}) *Response {
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
