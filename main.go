package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if argc := len(os.Args); argc < 2 {
		fmt.Println("./main <run/arg> <commands>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		fmt.Println("./main <run/arg> <commands>")
		os.Exit(1)
	}
}

func parent() {
	// Create a clone of our program
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// This clone should have a separate namespace for PID, MNT, and network for isolation
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER, // User NS to allow for rootless containers
		Unshareflags: syscall.CLONE_NEWNS, // CRITICAL: Unshare the mount namespace so it doesn't share with parent
												// This is so when we unmount our running FS, it won't impact the parent/host FS
		Credential: &syscall.Credential{
			Uid: 0,
			Gid: 0,
		},
		UidMappings: []syscall.SysProcIDMap{ // Map our child userID to namespace
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

func child() {
	root := "/home/justin/projects/containers-go/rootfs"
	syscall.Sethostname([]byte("container"))

	// Mount our desired location into a file system for our child
	must(syscall.Mount(root, root, "", syscall.MS_BIND, ""))

	// Create a directory that will hold our current running file system
	must(os.MkdirAll(root + "/oldrootfs", 0707))

	// Our current running FS (host "/") becomes root (arg1)
	// Our old FS (host "/") is put into <root>/oldrootfs
		// This is better than chroot since there is no escape
		// Whereas pivot_root allows us to set the root of our FS to our current namespace
	must(syscall.PivotRoot(root, root + "/oldrootfs"))

	// Change into our current running FS
	must(os.Chdir("/"))

	// Mount proc to be able to view processes
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	// Unmount our host FS (host "/") from our container
		// Because we've unshared(CLONE_NEWNS), it won't unmount from our host OS
	must(syscall.Unmount("oldrootfs", syscall.MNT_DETACH))
	must(os.Remove("oldrootfs"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

	// Cleanup
	must(syscall.Unmount("proc", syscall.MNT_DETACH))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
