package neofs

import (
	"context"
	"crypto/ecdsa"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neofs-api-go/pkg/client"
	"github.com/nspcc-dev/neofs-api-go/pkg/owner"
)

// Root represents root of the NeoFS filesystem which contains
// all containers owned by a specific owner.
type Root struct {
	fs.Inode

	Client     client.Client
	PrivateKey *keys.PrivateKey
	Address    string
}

// OnAdd implements fs.NodeOnAdder.
func (r *Root) OnAdd(ctx context.Context) {
	cl, err := client.New(
		client.WithDefaultPrivateKey(&r.PrivateKey.PrivateKey),
		client.WithAddress(r.Address))
	if err != nil {
		panic(err)
	}
	r.Client = cl

	w, _ := owner.NEO3WalletFromPublicKey((*ecdsa.PublicKey)(r.PrivateKey.PublicKey()))
	cids, err := r.Client.ListContainers(ctx, owner.NewIDFromNeo3Wallet(w))
	if err != nil {
		panic(err)
	}

	for i := range cids {
		c := &Container{ID: cids[i]}
		cDir := r.NewPersistentInode(ctx, c, fs.StableAttr{Mode: syscall.S_IFDIR})

		r.AddChild(c.ID.String(), cDir, false)
	}
}
