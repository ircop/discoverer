package logger

import (
	"testing"
	"os"
	"bytes"
)

func TestWrongDir(t *testing.T) {
	// proc is not writeable even for root
	err := InitLogger(false, "/proc/wrongdir")
	if err == nil {
		t.Fatal("worng dir: should return error")
	}
}

func TestNoPerms(t *testing.T) {
	// logger will try to create files in /proc
	err := InitLogger(false, "/proc")
	if err == nil {
		t.Fatal("log in /proc: should return error")
	}
}

func TestCreateDir(t *testing.T) {
	dir := "testdata/newdir"

	// cleanup
	if _, err := os.Stat("testdata/newdir"); err == nil {
		os.RemoveAll("testdata/newdir")
	}

	err := InitLogger(false, dir)
	if err != nil {
		t.Fatalf("creating testdata/newdir: should not return error (%s)", err.Error())
	}

	// cleanup
	if _, err := os.Stat(dir); err == nil {
		os.RemoveAll(dir)
	}
}

func TestWrite(t *testing.T) {
	var b bytes.Buffer
	if err := mlog.write("test", &b); err != nil {
		t.Fatal("write to buffer should not return error")
	}

	got := b.String()
	if got != "test" {
		t.Fatal("got is not 'test' string")
	}
}
