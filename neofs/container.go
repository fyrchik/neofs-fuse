package neofs

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/nspcc-dev/neofs-api-go/pkg/client"
	cid "github.com/nspcc-dev/neofs-api-go/pkg/container/id"
	"github.com/nspcc-dev/neofs-api-go/pkg/object"
)

// Container represents container inode.
type Container struct {
	fs.Inode

	ID *cid.ID
}

var (
	_ fs.InodeEmbedder = (*Container)(nil)
	_ fs.NodeReaddirer = (*Container)(nil)
)

func (c *Container) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	r := c.Root().Operations().(*Root)
	ps := new(client.SearchObjectParams).WithContainerID(c.ID)
	objectIDs, err := r.Client.SearchObject(ctx, ps)
	if err != nil {
		return nil, syscall.EREMOTEIO
	}

	for i := range objectIDs {
		addr := object.NewAddress()
		addr.SetContainerID(c.ID)
		addr.SetObjectID(objectIDs[i])

		obj := &Object{Address: addr}
		objInode := c.NewPersistentInode(ctx, obj, fs.StableAttr{Mode: syscall.S_IFREG})
		c.AddChild(objectIDs[i].String(), objInode, false)
	}

	dirEntries := make([]fuse.DirEntry, len(objectIDs))
	for i := range objectIDs {
		dirEntries[i].Name = objectIDs[i].String()
		dirEntries[i].Ino = uint64(i + 10)
		dirEntries[i].Mode = syscall.S_IFREG
	}

	return fs.NewListDirStream(dirEntries), 0
}
