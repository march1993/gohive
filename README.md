## What's gohive?
Gohive is a tool to help manage the micro-services driven by golang. Though nowadays dockers should be the best integration friend for golang, it's still useful to run golang applications like PHP(cPanel, Plesk) does where the golang applications are small and simple.


## Depends
1. OS support
Only Debian/stretch is supported so far.

2. nginx

3. git

4. MariaDB
Currently only sock connection is supported.

5. CertBot

6. Golang
`GOPATH` should be correctly configured.

## Install
```shell
go get github.com/march1993/gohive
cd $GOPATH/src/github.com/march1993/gohive
go build
./gohive install
```
