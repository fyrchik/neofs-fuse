package neofs

import (
	"context"
	"strings"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/nspcc-dev/neofs-api-go/pkg/client"
	"github.com/nspcc-dev/neofs-api-go/pkg/object"
)

// Object represents object inode.
type Object struct {
	fs.Inode

	Address *object.Address
}

type objectHandle struct {
	*object.Object
}

var (
	_ fs.InodeEmbedder   = (*Object)(nil)
	_ fs.NodeReader      = (*Object)(nil)
	_ fs.NodeOpener      = (*Object)(nil)
	_ fs.NodeGetattrer   = (*Object)(nil)
	_ fs.NodeGetxattrer  = (*Object)(nil)
	_ fs.NodeListxattrer = (*Object)(nil)
)

func (o *Object) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if f == nil {
		return 0
	}
	oh := f.(*objectHandle)
	out.Size = uint64(len(oh.Payload()))
	return 0
}

func (o *Object) Listxattr(ctx context.Context, dest []byte) (uint32, syscall.Errno) {
	r := o.Root().Operations().(*Root)

	ps := new(client.GetObjectParams).WithAddress(o.Address)
	obj, err := r.Client.GetObject(ctx, ps)
	if err != nil {
		return 0, syscall.ENOENT
	}

	var buf []byte
	buf = append(buf, "system.object_id"...)
	buf = append(buf, 0)

	for _, a := range obj.Attributes() {
		buf = append(buf, "user."...)
		buf = append(buf, a.Key()...)
		buf = append(buf, 0)
	}

	return copyBuffer(dest, buf)
}

func (o *Object) Getxattr(ctx context.Context, attr string, dest []byte) (uint32, syscall.Errno) {
	r := o.Root().Operations().(*Root)

	ps := new(client.GetObjectParams).WithAddress(o.Address)
	obj, err := r.Client.GetObject(ctx, ps)
	if err != nil {
		return 0, syscall.ENOENT
	}

	switch {
	case strings.HasPrefix(attr, "system."):
		attr = strings.TrimPrefix(attr, "system.")
		switch attr {
		case "object_id":
			return copyBuffer(dest, append([]byte(obj.ID().String()), 0))
		}
	case strings.HasPrefix(attr, "user."):
		attr = strings.TrimPrefix(attr, "user.")
		for _, a := range obj.Attributes() {
			if a.Key() == attr {
				return copyBuffer(dest, append([]byte(a.Value()), 0))
			}
		}
	}
	return 0, syscall.ENODATA // ENOATTR is synonym to ENODATA
}

func (o *Object) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	r := o.Root().Operations().(*Root)

	ps := new(client.GetObjectParams).WithAddress(o.Address)
	obj, err := r.Client.GetObject(ctx, ps)
	if err != nil {
		return nil, 0, syscall.ENOENT
	}

	return &objectHandle{obj}, 0, 0
}

func (o *Object) Read(ctx context.Context, f fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	data := f.(*objectHandle).Payload()
	dest = append(dest[:0], data...)
	return fuse.ReadResultData(dest), 0
}
