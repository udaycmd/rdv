package drives

import (
	"context"
	"io"
	"time"

	"github.com/udaycmd/rdv/internal/oauth"
	"github.com/udaycmd/rdv/internal/oauth/providers"
)

var SupportedDriveProviders = []oauth.OauthProvider{
	providers.NewGdriveAuthProvider(),
	providers.NewDboxAuthProvider(),
}

func GetDriveOauthProvider(name string) oauth.OauthProvider {
	for _, p := range SupportedDriveProviders {
		if p.GetConfig().Name == name {
			return p
		}
	}

	return nil
}

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
	View(ctx context.Context, id string) ([]Meta, error)

	// Returns a file's content as a stream.
	// The caller is responsible for closing the returned stream.
	Get(ctx context.Context, id string) (io.ReadCloser, error)

	// Uploads a new file or updates an existing one.
	// 'r' is the data stream, 'parentId' is the target folder (optional),
	// and 'name' is the filename.
	// Returns the metadata of the created file.
	Put(ctx context.Context, r io.Reader, parentId string, name string) (*Meta, error)

	// Removes an object by its Id.
	Delete(ctx context.Context, id string) error

	// Creates a new directory.
	MkDir(ctx context.Context, parentId string, name string) (*Meta, error)
}
