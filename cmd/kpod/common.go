package main

import (
	is "github.com/containers/image/storage"
	"github.com/containers/storage"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func getStore(c *cli.Context) (storage.Store, error) {
	options := storage.DefaultStoreOptions
	if c.GlobalIsSet("storage-driver") {
		options.GraphDriverName = c.GlobalString("storage-driver")
	}
	if c.GlobalIsSet("storage-opt") {
		opts := c.GlobalStringSlice("storage-opt")
		if len(opts) > 0 {
			options.GraphDriverOptions = opts
		}
	}
	store, err := storage.GetStore(options)
	if err != nil {
		return nil, err
	}
	is.Transport.SetStore(store)
	return store, nil
}

func findImage(store storage.Store, image string) (*storage.Image, error) {
	var img *storage.Image
	ref, err := is.Transport.ParseStoreReference(store, image)
	if err == nil {
		img, err = is.Transport.GetStoreImage(store, ref)
	}
	if err != nil {
		img2, err2 := store.Image(image)
		if err2 != nil {
			if ref == nil {
				return nil, errors.Wrapf(err, "error parsing reference to image %q", image)
			}
			return nil, errors.Wrapf(err, "unable to locate image %q", image)
		}
		img = img2
	}
	return img, nil

}
