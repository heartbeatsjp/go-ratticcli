# go-ratticcli

golang RatticCLI for RatticWeb

# Usage

- search: search Cred
- show: show Cred
- reload: reload token and local cache

Typical usecase : use with [peco](https://github.com/peco/peco)

```
rattic search | peco --exec "rattic show --fields=password" | pbcopy
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

- env `ENDPOINT` / option `--endpoint` (default: `https://localhost` )
- env `USER` / option `--user` (default: local username)

# Build

```
glide install
go build -o rattic
```

# Datastore

boltdb

- Bucket: Config
    - Key: `Token` , `LastUpdated`
- Bucket: Creds
    - Key: `Cred.id`
    - Value: `Cred.title`

