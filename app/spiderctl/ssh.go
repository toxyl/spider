package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/toxyl/glog"
	"github.com/toxyl/spider/log"
)

func sshExec(host string, args ...string) error {
	args = append([]string{
		"-o", "ConnectTimeout=5",
		"-q",
		"-i", inventory.Credentials.Key,
		fmt.Sprintf("%s@%s", inventory.Credentials.User, host),
		"-tt"}, args...)
	cmd := exec.Command("/usr/bin/ssh", args...)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func sshExecInteractive(host string, args ...string) error {
	args = append([]string{
		"-o", "ConnectTimeout=5",
		"-q",
		"-i", inventory.Credentials.Key,
		fmt.Sprintf("%s@%s", inventory.Credentials.User, host),
		"-tt"}, args...)
	cmd := exec.Command("/usr/bin/ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func sshCommandExists(host, command string) bool {
	return nil == sshExec(host, "command", "-v", command)
}

func sshChown(host, file string) error {
	return sshExec(host, "sudo /usr/bin/chown", fmt.Sprintf("%s:%s", inventory.Credentials.User, inventory.Credentials.Group), file)
}

func sshMv(host, src, dst string) error {
	return sshExec(host, "sudo /usr/bin/mv", src, dst)
}

func sshMkdir(host, dir string) error {
	return sshExec(host, "sudo /usr/bin/mkdir", "-p", dir)
}

func sshRm(host, file string) error {
	log.Info("Removing %s on %s...", glog.File(file), glog.Auto(host))
	err := sshExec(host, "sudo /usr/bin/rm", "-rf", file)
	if err != nil {
		log.Error("Could not remove %s on %s...", glog.File(file), glog.Auto(host))
	}
	return nil
}

func sshServiceEnable(host, service string) error {
	return sshExec(host, "sudo /usr/bin/systemctl", "enable", "/etc/systemd/system/"+service+".service")
}

func sshServiceStart(host, service string) error {
	return sshExec(host, "sudo /usr/sbin/service", service, "start")
}

func sshServiceStop(host, service string) error {
	return sshExec(host, "sudo /usr/sbin/service", service, "stop")
}

func sshOpenPortUFW(host string, port int) error {
	return sshExec(host, "sudo /usr/sbin/ufw allow", fmt.Sprint(port), ">/dev/null")
}

func sshOpenPortIPTables(host string, port int) error {
	return sshExec(host, "sudo /usr/sbin/iptables -I INPUT -p tcp --dport", fmt.Sprint(port), "-j ACCEPT >/dev/null")
}
