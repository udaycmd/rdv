package api

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/udaycmd/rdv/internal/drives"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
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

		modTime, err := time.Parse(time.RFC3339, f.ModifiedTime)
		if err != nil {
			return nil, err
		}

		meta = append(meta, drives.Meta{
			Id:           f.Id,
			Name:         f.Name,
			Size:         f.Size,
			MimeType:     f.MimeType,
			LastModified: modTime,
			IsDir:        f.MimeType == "application/vnd.google-apps.folder",
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

func (g *Gdrive) Put(ctx context.Context, r io.Reader, parentId string, name string) (*drives.Meta, error) {
	file := &drive.File{
		Name:    name,
		Parents: []string{parentId},
	}

	call := g.service.Files.Create(file).
		Media(r, googleapi.ChunkSize(16*1024*1024)).
		Context(ctx).Fields("id, name, size, mimeType, modifiedTime").
		ProgressUpdater(func(current, total int64) {})

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	modTime, err := time.Parse(time.RFC3339, resp.ModifiedTime)
	if err != nil {
		return nil, err
	}

	return &drives.Meta{
		Id:           resp.Id,
		Name:         resp.Name,
		Size:         resp.Size,
		MimeType:     resp.MimeType,
		LastModified: modTime,
		IsDir:        resp.MimeType == "application/vnd.google-apps.folder",
	}, nil
}
