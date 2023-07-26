package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/toxyl/glog"
	"github.com/toxyl/spider/log"
)

func scpUpload(host, fileSrc, fileDst string) (string, error) {
	cmd := exec.Command(
		"/usr/bin/scp",
		"-i", inventory.Credentials.Key,
		"-C", fileSrc,
		fmt.Sprintf("%s@%s:%s", inventory.Credentials.User, host, fileDst),
	)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Error("Uploading %s to %s:%s failed: %s", glog.File(fileSrc), glog.Auto(host), glog.File(fileDst), glog.Error(err))
		return fileDst, err
	}
	return fileDst, nil
}

func scpServiceToHost(host, service string) error {
	dstPath := "/etc/systemd/system/" + service + ".service"
	ftmp := getTempFilePath(service + "-service")
	uploadedFile, err := scpUpload(host, ftmp, ftmp)
	if err != nil {
		return err
	}

	return sshMv(host, uploadedFile, dstPath)
}

func scpFileToHost(host, fileSrc, fileDst string) error {
	tmpDst := filepath.Base(fileSrc) + ".tmp"
	uploadedFile, err := scpUpload(host, fileSrc, tmpDst)
	if err != nil {
		return err
	}

	err = sshMv(host, uploadedFile, fileDst)
	if err != nil {
		log.Error("Moving %s to %s failed: %s", glog.File(uploadedFile), glog.File(fileDst), glog.Error(err))
		return err
	}
	return nil
}

func scpExecutableToHost(host, executable, destination string) error {
	dstPath := filepath.Join(destination, executable)
	if err := scpFileToHost(host, getTempFilePath(executable), dstPath); err != nil {
		return err
	}

	return sshChown(host, dstPath)
}
