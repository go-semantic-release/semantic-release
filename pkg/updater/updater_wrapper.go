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

func (f *FilesUpdaterServer) Init(_ context.Context, request *FilesUpdaterInit_Request) (*FilesUpdaterInit_Response, error) {
	err := f.Impl.Init(request.Config)
	res := &FilesUpdaterInit_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

func (f *FilesUpdaterServer) Name(_ context.Context, _ *FilesUpdaterName_Request) (*FilesUpdaterName_Response, error) {
	return &FilesUpdaterName_Response{Name: f.Impl.Name()}, nil
}

func (f *FilesUpdaterServer) Version(_ context.Context, _ *FilesUpdaterVersion_Request) (*FilesUpdaterVersion_Response, error) {
	return &FilesUpdaterVersion_Response{Version: f.Impl.Version()}, nil
}

func (f *FilesUpdaterServer) ForFiles(_ context.Context, _ *FilesUpdaterForFiles_Request) (*FilesUpdaterForFiles_Response, error) {
	return &FilesUpdaterForFiles_Response{Files: f.Impl.ForFiles()}, nil
}

func (f *FilesUpdaterServer) Apply(_ context.Context, request *FilesUpdaterApply_Request) (*FilesUpdaterApply_Response, error) {
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

func (f *FilesUpdaterClient) Init(m map[string]string) error {
	res, err := f.Impl.Init(context.Background(), &FilesUpdaterInit_Request{
		Config: m,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (f *FilesUpdaterClient) Name() string {
	res, err := f.Impl.Name(context.Background(), &FilesUpdaterName_Request{})
	if err != nil {
		panic(err)
	}
	return res.Name
}

func (f *FilesUpdaterClient) Version() string {
	res, err := f.Impl.Version(context.Background(), &FilesUpdaterVersion_Request{})
	if err != nil {
		panic(err)
	}
	return res.Version
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
