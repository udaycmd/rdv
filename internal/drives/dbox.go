package drives

import (
	"errors"
	"io"
	"net/http"
)

type Dbox struct {
	client *http.Client
}

func NewDbox(client *http.Client) *Dbox {
	return &Dbox{client: client}
}

func (d *Dbox) View(id string) ([]Meta, error) {
	return nil, errors.New("dbox: View not implemented")
}

func (d *Dbox) Get(id string) (io.ReadCloser, error) {
	return nil, errors.New("dbox: Get not implemented")
}

func (d *Dbox) Put(r io.Reader, parentId string, name string) (*Meta, error) {
	return nil, errors.New("dbox: Put not implemented")
}

func (d *Dbox) Delete(id string) error {
	return errors.New("dbox: Delete not implemented")
}

func (d *Dbox) MkDir(parentId string, name string) (*Meta, error) {
	return nil, errors.New("dbox: MkDir not implemented")
}
