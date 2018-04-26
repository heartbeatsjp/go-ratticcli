# go-ratticcli

[![wercker status](https://app.wercker.com/status/079626343f674daa3859a61ccdb3c22b/s/ "wercker status")](https://app.wercker.com/project/byKey/079626343f674daa3859a61ccdb3c22b)

CLI for RatticWeb

RatticWeb is Password Management for Humans.

Strongly recommend to use fork [netmarkjp/RatticWeb](https://github.com/netmarkjp/RatticWeb).

Original RatticWeb is not maintained, but netmarkjp fork is still maintaned.

# Usage

- `list`: list Creds
    - If cache expired, reload runs in background. While fetching new cred(for cache), command will not exit.
    - If you never want to use stole cache, run `reload` just before `list`
- `show`: show Cred
- `reload`: reload Creds

Typical usecase : use with [peco](https://github.com/peco/peco)

```
rattic list | peco | rattic show | pbcopy
```

CLI to clipbard tools

- Windows: `clip`
- MacOSX: `pbcopy`
- Linux: `xsel --clipbard --input`

# Install

## local build

```
go get github.com/heartbeatsjp/go-ratticcli
```

or

## Use pre-built binary

Download from [releases page](https://github.com/heartbeatsjp/go-ratticcli/releases)

- Windows: windows_amd64.tar.gz
- MacOSX: darwin_amd64.tar.gz
- Linux: linux_amd64.tar.gz

# Configuration

- env `RATTIC_ENDPOINT` / option `--endpoint` (default: `https://localhost` )
- env `RATTIC_USER` / option `--user` (default: local username)
- env `RATTIC_TOKEN` / option `--token`
- env `RATTIC_CACHE_TTL` / option `--cache-ttl` (default: 86400 (sec))

# Build

Recommend: use `wercker` cli.
( binaries are put on `.wercker/latest/output/` )

```
wercker build --artifacts
```

If build by hand localy

```
dep ensure
go build -o rattic -ldflags "-w -s"
```

# Datastore

boltdb

- Bucket: Config
    - Key: `LastUpdated`
- Bucket: Creds
    - Key: `Cred.id`
    - Value: `Cred.title`

