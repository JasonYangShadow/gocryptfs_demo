package demo

import (
	"os"
	"testing"
)

func TestEncryption(t *testing.T) {
	d := &Demo{}
	err := d.init()
	if err != nil {
		t.Fatal(err)
	}

	err = d.initGocryptfs()
	if err != nil {
		t.Fatal(err)
	}

	err = d.encrypt()
	if err != nil {
		t.Fatal(err)
	}

	err = d.fuserunmount(d.plainDir)
	if err != nil {
		t.Fatal(err)
	}

	os.RemoveAll(d.tempDir)
}

func TestSquashMount(t *testing.T) {
	d := &Demo{}
	err := d.init()
	if err != nil {
		t.Fatal(err)
	}

	err = d.createSquashfsArchive()
	if err != nil {
		t.Fatal(err)
	}

	err = d.squashfuseMount()
	if err != nil {
		t.Fatal(err)
	}

	err = d.fuserunmount(d.squashMountPoint)
	if err != nil {
		t.Fatal(err)
	}

	os.RemoveAll(d.tempDir)
}
