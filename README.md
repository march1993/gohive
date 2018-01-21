## What's gohive?
Gohive is a tool to help manage the micro-services driven by golang. Though nowadays dockers should be the best integration friend for golang, it's still useful to run golang applications like PHP(cPanel, Plesk) does where the golang applications are small and simple.


## Depends
1. OS support

Only Debian/stretch is supported so far. Directory `/gohive` would be used to store all the application and `/gohive.go` would be used to store all golang installtions.

2. nginx

The symbol link of `/etc/nginx/sites-enabled/default` would be canceled.

3. git

4. MariaDB

Currently only sock connection is supported. A database called `gohive` would be occupied.

5. CertBot

6. Golang

7. members

`members` is used to list users under a certain group.

`GOPATH` should be correctly configured. `1.9.4` is tested and recommended.

## Install
I recommend to install gohive on a fresh server.
```shell
go get github.com/march1993/gohive
cd $GOPATH/src/github.com/march1993/gohive
go build
./gohive -install
```
Use the generated token to login the web panel.


## Upgrade
```shell
cd $GOPATH/src/github.com/march1993/gohive
git pull
go get .
go build
sudo service gohive restart
```

## TODO List
1. Web UI
2. App management (git repository, pubkey importation, ulimit and etc.)
3. Golang support
4. Nginx settings (including SSL) support
5. Mysql support
