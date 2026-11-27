package client

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

var _ Client = &filesystemClient{}

type filesystemClient struct {
	dir string
}

func NewFilesystemClient(directory string) (*filesystemClient, error) {
	if _, err := os.Stat(directory); err != nil {
		return nil, err
	}
	return &filesystemClient{dir: directory}, nil
}

func (f *filesystemClient) GetResource(ctx context.Context, file string) (Resource, error) {
	var resource Resource
	path := filepath.Join(f.dir, file)
	data, err := os.ReadFile(path)
	if err != nil {
		return resource, err
	}
	err = yaml.Unmarshal(data, &resource)
	return resource, err
}

func (f *filesystemClient) Create(ctx context.Context, resource Resource) (Resource, error) {
	data, err := yaml.Marshal(&resource)
	if err != nil {
		return resource, err
	}

	path := filepath.Join(f.dir, resource.Id)
	err = os.WriteFile(path, data, 0644)
	return resource, err
}

func (f *filesystemClient) Delete(ctx context.Context, id string) error {
	path := filepath.Join(f.dir, id)
	return os.Remove(path)
}

func (f *filesystemClient) Update(ctx context.Context, resource Resource) (Resource, error) {

	if resource.Id == "" {
		return resource, errors.New("can not update resource without id")
	}

	path := filepath.Join(f.dir, resource.Id)

	if _, err := os.Stat(path); err != nil {
		return resource, errors.Join(
			errors.New("resource does not exist"),
			err,
		)
	}

	data, err := yaml.Marshal(&resource)
	if err != nil {
		return resource, err
	}

	err = os.WriteFile(path, data, 0644)
	return resource, err
}
