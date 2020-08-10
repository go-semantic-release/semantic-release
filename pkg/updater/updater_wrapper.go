package updater

import (
	"context"
	"errors"
)

const FilesUpdaterPluginName = "files_updater"

type FilesUpdaterServer struct {
	Impl FilesUpdater
	UnimplementedFilesUpdaterPluginServer
}

func (f *FilesUpdaterServer) ForFiles(ctx context.Context, request *FilesUpdaterForFiles_Request) (*FilesUpdaterForFiles_Response, error) {
	return &FilesUpdaterForFiles_Response{Files: f.Impl.ForFiles()}, nil
}

func (f *FilesUpdaterServer) Apply(ctx context.Context, request *FilesUpdaterApply_Request) (*FilesUpdaterApply_Response, error) {
	err := f.Impl.Apply(request.File, request.NewVersion)
	res := &FilesUpdaterApply_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

type FilesUpdaterClient struct {
	Impl FilesUpdaterPluginClient
}

func (f *FilesUpdaterClient) ForFiles() string {
	res, err := f.Impl.ForFiles(context.Background(), &FilesUpdaterForFiles_Request{})
	if err != nil {
		panic(err)
	}
	return res.Files
}

func (f *FilesUpdaterClient) Apply(file, newVersion string) error {
	res, err := f.Impl.Apply(context.Background(), &FilesUpdaterApply_Request{
		File:       file,
		NewVersion: newVersion,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}
