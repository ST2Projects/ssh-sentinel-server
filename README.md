# ssh-sentinel-server

A simple to use and deploy SSH CA server.

**This is a Work In Progress** - the project is in its early days. It is functional and I'm using it for my own hosts **but** you use it at your own risk

## Installation

The release archive contains the binary `ssh-sentinel-server` and a samples directory containing a config template and a systemd service file.

You **will** need to edit the `config.json` to suit your needs. You may need to edit the service file depending on your OS.

To install, unpack the archive into the `/opt` directory then run the then install

```shell
mkdir /opt/sentinel
# Copy archive into directory
tar xvzf ssh-sentinel-server_$VERSION_$ARCH.tar.gz

make install
```

## Configuration

Configuration is defined in the `config.json`. Properties are explained below. All paths are relative to the `resources` directory

- `CAPrivateKey` - Name of the CA private key. The key must be unencrypted - a future enhancement will allow encrypted keys
- `CAPublicKey` - Name of the CA public key.
- `MaxValidTime` - Maximum lifespan of signed keys, in the normal [go duration format](https://pkg.go.dev/time#ParseDuration)
- `db.dialect` - Must be `sqlite3`. A future release will add support for other DBs
- `db.username` - Username of the DB user
- `db.password` - Password of the DB user
- `db.connection` - Connection URL for the DB. For sqlite3 this is a file path
- `db.dbName` - Name of the DB

## Goals

There are a couple of SSH CA servers out there - I have found them all difficult to use and have specific platform
requirements. This projects aims to:

- Be simple to use and deploy
- Use sensible secure defaults

## Releases

All releases are signed with [signify](https://github.com/aperezdc/signify) using this key:

```
untrusted comment: signify public key
RWTZI8XJdh5lJ5/cWa8yZry/28x5frzb/PZqp4PL3IfFQ154BaX8ja4Q
```

You can verify the checksums using

```shell
signify -V -p path/to/public/key.pub -m checksums.txt
```

## Development Dependencies

- lego
- mkcert
