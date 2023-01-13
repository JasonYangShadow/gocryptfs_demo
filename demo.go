package demo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Demo struct {
	tempDir          string // temp folder contains everything
	mksquashfsPath   string // mksquashfs path
	gocryptfsPath    string // gocryptfs path
	dataDir          string // data folder contains plain text
	plainDir         string // plain folder mount by gocryptfs(decrypted content folder)
	file             string // test file plain text
	squashfile       string // using mksquashfs created squash archive
	cipherDir        string // gocryptfs cipher folder
	fusermount       string // fusermount path
	squashfuse       string // squashfuse path
	squashMountPoint string // squashfuse mount point a folder
}

func (d *Demo) init() error {
	temp, err := os.MkdirTemp(os.TempDir(), "demo-")
	if err != nil {
		return err
	}
	d.tempDir = temp

	datapath := fmt.Sprintf("%s/data", temp)
	err = os.Mkdir(datapath, 0o755)
	if err != nil {
		return err
	}
	d.dataDir = datapath

	plainpath := fmt.Sprintf("%s/plain", temp)
	err = os.Mkdir(plainpath, 0o755)
	if err != nil {
		return err
	}
	d.plainDir = plainpath

	d.file = fmt.Sprintf("%s/hello", datapath)
	d.squashfile = fmt.Sprintf("%s/squashfs", plainpath)

	squashfs, err := exec.LookPath("mksquashfs")
	if err != nil {
		return err
	}
	d.mksquashfsPath = squashfs

	gocryptfs, err := exec.LookPath("gocryptfs")
	if err != nil {
		return err
	}
	d.gocryptfsPath = gocryptfs

	d.cipherDir = fmt.Sprintf("%s/cipher", temp)
	err = os.Mkdir(d.cipherDir, 0o755)
	if err != nil {
		return err
	}

	fusermount, err := exec.LookPath("fusermount")
	if err != nil {
		return err
	}
	d.fusermount = fusermount

	squashfuse, err := exec.LookPath("squashfuse")
	if err != nil {
		return err
	}
	d.squashfuse = squashfuse

	d.squashMountPoint = fmt.Sprintf("%s/squashmount", temp)
	err = os.Mkdir(d.squashMountPoint, 0o755)
	if err != nil {
		return err
	}

	return nil
}

func (d *Demo) createSquashfsArchive() error {
	file, err := os.Create(d.file)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("hello world")
	if err != nil {
		return err
	}

	cmd := exec.Command(d.mksquashfsPath, d.dataDir, d.squashfile, "-comp", "xz")
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (d *Demo) initGocryptfs() error {
	cmd := exec.Command(d.gocryptfsPath, "-init", "-deterministic-names", "-plaintextnames", d.cipherDir)
	cmd.Stdin = strings.NewReader("12345\n12345\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println(stdout.String())
	fmt.Println(stderr.String())

	return nil
}

func (d *Demo) encrypt() error {
	cmd := exec.Command(d.gocryptfsPath, d.cipherDir, d.plainDir)
	cmd.Stdin = strings.NewReader("12345\n")
	if err := cmd.Run(); err != nil {
		return err
	}

	return d.createSquashfsArchive()
}

func (d *Demo) decrypt() error {
	cmd := exec.Command(d.gocryptfsPath, d.cipherDir, d.plainDir)
	cmd.Stdin = strings.NewReader("12345\n")
	return cmd.Run()
}

func (d *Demo) squashfuseMount() error {
	cmd := exec.Command(d.squashfuse, d.squashfile, d.squashMountPoint)
	return cmd.Run()
}

func (d *Demo) fuserunmount(mountpoint string) error {
	cmd := exec.Command(d.fusermount, "-u", mountpoint)
	return cmd.Run()
}
