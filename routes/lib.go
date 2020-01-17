package routes

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

func readRequestBody(r *http.Request) ([]byte, error) {
	var body []byte

	if r.Body != nil {
		defer r.Body.Close()
		if data, err := ioutil.ReadAll(r.Body); err != nil {
			io.Copy(ioutil.Discard, r.Body)
		} else {
			body = data
		}
	}

	if len(body) == 0 {
		return nil, errors.New("body of the request is empty")
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("received Content-Type=%s, Expected Content-Type is 'application/json'", contentType)
	}

	glog.Infof("Request received: \n %s \n", body)
	return body, nil
}

func writeError(writer http.ResponseWriter, message string, status int) {
	glog.Errorf(message)
	http.Error(writer, message, status)
}
