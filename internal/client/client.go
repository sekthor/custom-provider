package client

import "context"

type Client interface {
	GetResource(context.Context, string) (Resource, error)
	Create(context.Context, Resource) (Resource, error)
	Delete(context.Context, string) error
	Update(context.Context, Resource) (Resource, error)
}
