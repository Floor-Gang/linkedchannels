# Linked Channels
Hide text-channels until members join their corresponding voice channel.

## Setup
Download [go](https://golang.org)
```
$ go mod download
$ cd ./cmd/linkedchannels
$ go build
$ ./linkedchannels
# ... Setup the new config.yml ...
$ ./linkedchannels
```

## Bot Usage
Add a new link:
 * Copy the ID of the voice channel you want linked
 * `.link add <voice ID> #channel`

Remove a link:
 * `.link remove <voice ID or #channel>`
 
List linked channels:
 * `.link list`
