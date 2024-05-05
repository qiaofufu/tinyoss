package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func NewGetStream(server, objectName string) (*GetStream, error) {
	if server == "" || objectName == "" {
		return nil, fmt.Errorf("server or objectName is empty")
	}
	resp, err := http.Get("http://" + server + "/objects/" + objectName)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("data server return http code %d", resp.StatusCode)
	}

	return &GetStream{reader: resp.Body}, nil
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
