## What's gohive?
Gohive is a tool to help manage the micro-services driven by golang. Though nowadays dockers should be the best integration friend for golang, it's still useful to run golang applications like PHP(cPanel, Plesk) does where the golang applications are small and simple.


## Depends
### 1. OS support

Only Debian/stretch is supported so far. Directory `/gohive` would be used to store all the application and `/gohive.go` would be used to store all golang installtions.

### 2. nginx

The symbol link `/etc/nginx/sites-enabled/default` would be cancelled.

### 3. git

### 4. MariaDB

Currently only sock connection is supported. A database called `gohive` would be occupied.

### 5. CertBot

### 6. Golang
`GOPATH` should be correctly configured. `1.9.2` is tested and recommended.

### 7. members

`members` is used to list users under a certain group.

### 8. sudo
`sudo` is used to trigger restarting the corresponding service after the git repository is updated.

## Install
### 1. Deploy gohive
I recommend to install gohive on a fresh server.
```shell
go get github.com/march1993/gohive
cd $GOPATH/src/github.com/march1993/gohive
go build
./gohive -install
```
Use the generated token to login the web panel.

### 2. Set up server name
Navigate to `http://$SERVER_IP/` and change `server name`.

### 3. Set up DNS record
Add a wild-cards record. For example, `gohive.example.com` is your server, then add a CNAME record of `*.gohive.example.com` to `gohive.example.com`. Thus app `some_app` can be accessed through `some_app.gohive.example.com`.


## Upgrade
```shell
cd $GOPATH/src/github.com/march1993/gohive
git pull
go get .
go build
sudo service gohive restart
```

## TODO List
1. Better Web UI
2. Nginx SSL support
3. Mysql support
4. CPU/Mem/Disk Usage quota


## Customize
You can put `.gohive.bashrc` into the root of your repository to change environments like GOPATH.
