# go-ratticcli

CLI for RatticWeb

**WARNING: This product is quite unstable now**

# Usage

- list: list Cred
- show: show Cred

Typical usecase : use with [peco](https://github.com/peco/peco)

```
rattic list | peco | rattic show | pbcopy
```

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

# Build

```
dep ensure
go build -o rattic -ldflags "-w -s"
```

# Datastore

boltdb

- Bucket: Config
    - Key: `Token` , `LastUpdated`
- Bucket: Creds
    - Key: `Cred.id`
    - Value: `Cred.title`

