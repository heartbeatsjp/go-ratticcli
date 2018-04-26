# go-ratticcli

CLI for RatticWeb

**WARNING: This product is quite unstable now**

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

```
go get github.com/heartbeatsjp/go-ratticcli
```

or

```
curl -L <release_url>  # TODO
```

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

