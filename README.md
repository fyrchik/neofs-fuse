# FUSE filesystem wrapper over NeoFS

## Usage
For usage with devenv:
`go run ./main.go -addr s03.neofs.devenv:8080 -key ../neofs-dev-env/wallets/wallet.key -target ~/mounthere`

The first level contains containers as directories.
The second level contains individual objects. Object attributes
are mapped as extended attributes (xattr).

### Example

#### Create container and push object
```
dzeta@wpc ~/r/neofs-node (master)> bin/neofs-cli --rpc-endpoint s03.neofs.devenv:8080 container create --await --policy 'REP 1 CBF 1 SELECT 2 FROM * AS X' --binary-key ../neofs-dev-env/wallets/wallet.key
container ID: 75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV
awaiting...
container has been persisted on sidechain
dzeta@wpc ~/r/neofs-node (master)> echo 123 >lamao
dzeta@wpc ~/r/neofs-node (master)> bin/neofs-cli --rpc-endpoint s01.neofs.devenv:8080 object put --cid 75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV --binary-key ../neofs-dev-env/wallets/wallet.key --file lamao --attributes someKey=someValue
[lamao] Object successfully stored
  ID: 8r2JvRVksgz5spRahh5dnVix4ikk82nFqXvYho4o1uLP
  CID: 75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV
```

#### View files in a mounted directory (`~/kek` in my case)
```
dzeta@wpc ~> cd ~/kek
dzeta@wpc ~/kek> ls
75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV

dzeta@wpc ~/kek> cd 75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV/
dzeta@wpc ~/k/75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV> ls
8r2JvRVksgz5spRahh5dnVix4ikk82nFqXvYho4o1uLP

dzeta@wpc ~/k/75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV> cat 8r2JvRVksgz5spRahh5dnVix4ikk82nFqXvYho4o1uLP
123
dzeta@wpc ~/k/75AZHN1PdYfA638DCDChgjZ67U4k2kvgC6CuMonNwpZV [2]> getfattr -d 8r2JvRVksgz5spRahh5dnVix4ikk82nFqXvYho4o1uLP
# file: 8r2JvRVksgz5spRahh5dnVix4ikk82nFqXvYho4o1uLP
user.FileName="lamao"
user.Timestamp="1626528167"
user.someKey="someValue"
```