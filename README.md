# ssh-sentinel-server

A simple to use and deploy SSH CA server.

**This is a Work In Progress** - the project is in its early days and is not ready to use anywhere.

Once ready, I will update the README and provide some more info in terms of usage and deployment

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
