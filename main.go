package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neofs-fuse/neofs"
)

var (
	key  = flag.String("key", "", "private key to use")
	addr = flag.String("addr", "", "address of neofs node")
	mp   = flag.String("target", "", "mount point")
)

func main() {
	flag.Parse()

	opts := new(fs.Options)
	opts.Debug = true
	opts.Logger = log.Default()

	priv, err := loadPrivateKey(*key)
	if err != nil {
		opts.Logger.Fatalf("can't load private key: %v", err)
	}

	s, err := fs.Mount(*mp, &neofs.Root{
		PrivateKey: priv,
		Address:    *addr,
	}, opts)
	if err != nil {
		opts.Logger.Fatalf("can't mount file system: %v", err)
	}

	defer func() {
		if err := s.Unmount(); err != nil {
			opts.Logger.Fatalf("can't unmount file system: %v", err)
		}
	}()
	s.Wait()
}

func loadPrivateKey(s string) (*keys.PrivateKey, error) {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}
	return keys.NewPrivateKeyFromBytes(data)
}
