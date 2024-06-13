# ssh-sentinel-server

[![CodeQL](https://github.com/ST2Projects/ssh-sentinel-server/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ST2Projects/ssh-sentinel-server/actions/workflows/codeql-analysis.yml) [![Go](https://github.com/ST2Projects/ssh-sentinel-server/actions/workflows/go.yml/badge.svg)](https://github.com/ST2Projects/ssh-sentinel-server/actions/workflows/go.yml)

A simple to use and deploy SSH CA server.

**This is a Work In Progress** - the project is in its early days. It is functional, and I'm using it for my own hosts, **but** you use it at your own risk

## Goals

There are a couple of SSH CA servers out there - I have found them all difficult to use and have specific platform
requirements. This projects aims to:

- Be simple to use and deploy
- Use sensible secure defaults

I'm also using this project to learn go, so if you come across it and notice something dumb please let me know by opening an issue!

## Installation

The release archive contains the binary `ssh-sentinel-server` and a samples directory containing a config template and a systemd service file.

You **will** need to edit the `config.json` to suit your needs. You may need to edit the service file depending on your OS.

To install, unpack the archive into the `/opt` directory then run the then install

```shell
mkdir /opt/sentinel
# Copy archive into directory
tar xvzf ssh-sentinel-server_$VERSION_$ARCH.tar.gz

cp samples/config.json .
cp samples/ssh-sentinel.service /etc/systemd/system/
systemctl daemon-reload
```

I'm working on an ansible role to do this and will update the readme when complete

### Configuration

Configuration is defined in the `config.json`. Properties are explained below. Full paths must be provided

- `CAPrivateKey` - Name of the CA private key. The key must be unencrypted - a future enhancement will allow encrypted keys
- `CAPublicKey` - Name of the CA public key.
- `maxValidTime` - Maximum lifespan of signed keys, in the normal [go duration format](https://pkg.go.dev/time#ParseDuration)
- `defaultExtensions` - A list of extensions to add to a key if the request does not contain any
- `db.dialect` - Must be `sqlite3`. A future release will add support for other DBs
- `db.username` - Username of the DB user
- `db.password` - Password of the DB user
- `db.connection` - Connection URL for the DB. For sqlite3 this is a file path
- `db.dbName` - Name of the DB
- `tls.local` - When set to `true` the server will generate a local TLS certificate. When `false` the server will generate a Let's Encrypt cert
- `tls.certDir` - Directory in which the generated certificate will be generated
- `tls.certDomains` - A list of domains to be included in the certificate.
- `tls.certEmail` - Needed when generating a certificate with let's encrypt
- `tls.dnsProvider` - Only `cloudflare` is supported at the moment. A future release will open up support for other providers
- `tls.dnsAPIToken` - The zone API token from cloudflare

### Generate a CA

You can generate a SSH CA with `ssh-keygen`. I suggest using ECDSA keys as they are smaller but this is not a requirement.

```shell
ssh-keygen -t ed25519 -f sentinel-ca -C sentinel-CA
```

**The key must not have a password** - this will be improved in a future release

### Adding users

Once you have the service installed you'll need to add some users. I hope to improve this process later but for now you can do it via the `admin` command

```shell
./ssh-sentinel-server admin -h
Create / delete users

Usage:
  ssh-sentinel-server admin [flags]

Flags:
  -c, --config string        Config file
  -C, --create               If set a new user will be created
  -h, --help                 help for admin
  -n, --name string          User's name
  -P, --principals strings   A list of principals for the user
  -U, --username string      Username
```

So to add a user

```shell
./ssh-sentinel-server admin -c config.json -C -n test -P test1 test2 -U test
```

Not that the username is the user associated with this service. The principals list the allowed usernames on the server you will ssh to.

## Usage

Here are some high level usage details

### Clients

The server stands up as a restful HTTP/S service. You can post requests via curl ( see [api docs](./api-docs.yaml) for the API ) or you can use the [CLI client](https://github.com/ST2Projects/ssh-sentinel-client)

### Servers

Servers require some configuration to use the CA. In short:

- Copy the CA **public key** to the server and save it in `/etc/ssh/ca.pub`
- Edit `/etc/ssh/sshd_config` and add `TrustedUserCAKeys /etc/ssh/ca.pub`
- Restart SSHD `service sshd restart`

The easiest way to do this across an estate is with ansible. There is an [ansible galaxy role available](https://galaxy.ansible.com/neo1908/install_ssh_ca)

Create a new playbook called `deploy-sentinel-ca.yml` with something like 

```yaml
- hosts: sentinel_managed
  vars:
    auth_url: https://auth.your.domain.com
    key_file: /etc/ssh/ca.pub

  roles:
    - { role: neo1908.install_ssh_ca }
```

Pull the role with `ansible-galaxy install neo1908.install_ssh_ca`

Then run your playbook `ansible-playbook deploy-sentinel-ca.yml`

## Releases

All releases are signed with [signify](https://github.com/aperezdc/signify) using this key:

```
untrusted comment: st2projects-code-signing public key
RWT+c7SH0ADLx3ndyVVDHySn8E+tsRqvVRmNxkeaU3wxG6pzRDtWMp4k
```

You can verify the checksums using

```shell
signify -V -p path/to/public/key.pub -m checksums.txt
```

## Development Dependencies

- lego
- mkcert
