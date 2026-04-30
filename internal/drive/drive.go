package drive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/internal/oauth/providers"
)

type Meta struct {
	Id           string
	Name         string
	Size         int64
	MimeType     string
	LastModified time.Time
	IsDir        bool
}

type DriveFactory func(ctx context.Context, client *http.Client) (Drive, error)

type Drive interface {
	// Returns the contents of a directory.
	// If id is empty, the root directory of the drive is assumed.
	View(id string) ([]Meta, error)

	// Returns a file's content as a stream.
	// The caller is responsible for closing the returned stream.
	Get(id string) (io.ReadCloser, error)

	// Uploads a new file or updates an existing one.
	// 'r' is the data stream, 'parentId' is the target folder (optional),
	// and 'name' is the filename.
	// Returns the metadata of the created file.
	Put(r io.Reader, parentId string, name string) (*Meta, error)

	// Removes an object by its Id.
	Delete(id string) error

	// Creates a new directory.
	MkDir(parentId string, name string) (*Meta, error)
}

var driveRegistry = make(map[string]DriveFactory)

func Register(name string, factory DriveFactory) {
	driveRegistry[name] = factory
}

func NewDriveFromProvider(ctx context.Context, provider string) (Drive, error) {
	p := providers.Get(provider)
	if p == nil {
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	config := p.GetConfig()

	t, err := oauth.GetToken(config.ClientID)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, t)

	factory, ok := driveRegistry[p.Name()]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", p.Name())
	}

	return factory(ctx, client)
}
