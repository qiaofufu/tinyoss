package objectstream

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type TempObjectStream struct {
	uuid   string
	size   int64
	server string
	hash   string
}

func NewTempObjectStream(server, hash string, size int64) *TempObjectStream {
	p := &TempObjectStream{size: size, server: server, hash: hash}
	_ = p.createTemp()
	return p
}

func (w *TempObjectStream) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("http://%s/temp/%s", w.server, w.uuid), bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("data server return http code %d", resp.StatusCode)
	}
	return len(p), nil
}

func (w *TempObjectStream) createTemp() error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/temp/%s", w.server, w.hash), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Size", strconv.FormatInt(w.size, 10))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err

	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("data server return http code %d", resp.StatusCode)

	}
	uid, err := io.ReadAll(resp.Body)
	if err != nil {
		return err

	}
	w.uuid = string(uid)
	return nil
}

func (w *TempObjectStream) send(method string) error {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s/temp/%s", w.server, w.uuid), nil)
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send abort commit")
	}
	return nil
}

func (w *TempObjectStream) Abort() error {
	return w.send(http.MethodDelete)
}

func (w *TempObjectStream) Commit() error {
	return w.send(http.MethodPut)
}
