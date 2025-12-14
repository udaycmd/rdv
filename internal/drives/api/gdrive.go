package api

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/udaycmd/rdv/internal/drives"
	"google.golang.org/api/drive/v3"
)

type Gdrive struct {
	service *drive.Service
}

func NewGdrive(srv *drive.Service) *Gdrive {
	return &Gdrive{service: srv}
}

func (g *Gdrive) View(ctx context.Context, id string) ([]drives.Meta, error) {
	if id == "" {
		id = "root"
	}

	q := fmt.Sprintf("'%s' in parents and trashed = false", id)
	call := g.service.Files.List().Q(q).Fields("files(id, name, size, mimeType, modifiedTime)").Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("gdrive ls: %w", err)
	}

	var meta []drives.Meta
	for _, f := range resp.Files {
		isDir := f.MimeType == "application/vnd.google-apps.folder"

		lm, _ := time.Parse(time.RFC3339, f.ModifiedTime)

		meta = append(meta, drives.Meta{
			Id:           f.Id,
			Name:         f.Name,
			Size:         f.Size,
			MimeType:     f.MimeType,
			LastModified: lm,
			IsDir:        isDir,
		})
	}

	return meta, nil
}

func (g *Gdrive) Get(ctx context.Context, id string) (io.ReadCloser, error) {
	resp, err := g.service.Files.Get(id).Context(ctx).Download()
	if err != nil {
		return nil, fmt.Errorf("gdrive get: %w", err)
	}

	return resp.Body, nil
}
