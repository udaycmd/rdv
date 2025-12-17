package drives

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/internal/oauth/providers"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Meta struct {
	Id           string
	Name         string
	Size         int64
	MimeType     string
	LastModified time.Time
	IsDir        bool
}

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

var SupportedDriveProviders = []providers.OauthProvider{
	providers.NewGdriveAuthProvider(),
	providers.NewDboxAuthProvider(),
}

func GetDriveOauthProvider(name string) providers.OauthProvider {
	for _, p := range SupportedDriveProviders {
		if p.Name() == name {
			return p
		}
	}

	return nil
}

func NewDriveFromProvider(ctx context.Context, provider string) (Drive, error) {
	p := GetDriveOauthProvider(provider)
	if p == nil {
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	config := p.GetConfig()

	t, err := oauth.GetToken(config.ClientID)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, t)

	switch p.Name() {
	case "gdrive":
		srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return nil, err
		}

		return NewGdrive(srv), nil
	case "dbox":
		return NewDbox(client), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", p.Name())
	}
}
