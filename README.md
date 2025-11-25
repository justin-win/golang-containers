# Simple Golang Containers
This was a project to learn how containers are implemented in Linux utilizing namespaces.

## To use
Create a directory that holds some basic functionality for a file system e.g.
```
$ mkdir ubuntu-fs
$ docker export $(docker create ubuntu:latest) | sudo tar -x -C ubuntu-fs
```

This snippet will grab the latest image of Ubuntu from docker and extract into the directory `ubuntu-fs/`.

Now you can run `go run main.go <run/child> <commands>` similar to Docker. If it cannot create a clone from `/proc/self/exe` then it you will need to compile it and run with escalated privileges
```
$ go build main.go
$ sudo ./main <run/child> <commands>
```
For example
```
$ sudo ./main run /bin/bash
root@container:/#
root@container:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
root@container:/#
```

## What I've learned
I've learned how Linux namespaces can be utilized to allow for isolation such as network, file systems, and process IDs.

There were defintely challenges in understanding how namespaces interacted with the broad systems, especially namespaces. One such problem was the host file system "/" being mounted into the container. This made it lose the isolation and figuring out how to prevent that using unshare flags as well as reading manpages on things such as pivot_root, mount, unmount, and clone.

Not only that, security concepts such as chroot escape.

## What I plan for later
I plan on implementing more namespaces such as Cgroups but also trying to add more functionality such as packaging software into the container such as python, go, gcc, and other dependencies that Docker already does as well as network control

## Credit
https://www.youtube.com/watch?v=8fi7uSYlOdc
