package v3client

import (
	"context"
	"errors"
	"fmt"
)

type Client interface {
	GetApplication(ctx context.Context, appId string) (*MasheryApplication, error)
	GetApplicationPackageKeys(ctx context.Context, appId string) ([]MasheryPackageKey, error)
}

type PluggableClient struct {
	schema    ClientMethodSchema
	transport *HttpTransport
}

type ClientMethodSchema struct {
	GetApplicationContext     func(ctx context.Context, appId string, transport *HttpTransport) (*MasheryApplication, error)
	GetApplicationPackageKeys func(ctx context.Context, appId string, transport *HttpTransport) ([]MasheryPackageKey, error)
}

func (c *PluggableClient) GetApplication(ctx context.Context, appId string) (*MasheryApplication, error) {
	if c.schema.GetApplicationContext != nil {
		return c.schema.GetApplicationContext(ctx, appId, c.transport)
	} else {
		return nil, c.notImplemented("GetApplication")
	}
}

func (c *PluggableClient) GetApplicationPackageKeys(ctx context.Context, appId string) ([]MasheryPackageKey, error) {
	if c.schema.GetApplicationContext != nil {
		return c.schema.GetApplicationPackageKeys(ctx, appId, c.transport)
	} else {
		return nil, c.notImplemented("GetApplicationPackageKeys")
	}
}

func (c *PluggableClient) notImplemented(meth string) error {
	return errors.New(fmt.Sprintf("No implementation method was supplied for method %s", meth))
}
