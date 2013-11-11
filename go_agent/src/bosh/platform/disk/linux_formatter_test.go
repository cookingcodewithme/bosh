package disk

import (
	testsys "bosh/system/testhelpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinuxFormatWhenUsingSwapFs(t *testing.T) {
	fakeRunner := &testsys.FakeCmdRunner{}
	fakeFs := &testsys.FakeFileSystem{}
	fakeRunner.CommandResults = map[string][]string{
		"blkid -p /dev/xvda1": []string{`xxxxx TYPE="ext4" yyyy zzzz`, ""},
	}

	formatter := NewLinuxFormatter(fakeRunner, fakeFs)
	formatter.Format("/dev/xvda1", FileSystemSwap)

	assert.Equal(t, 2, len(fakeRunner.RunCommands))
	assert.Equal(t, []string{"mkswap", "/dev/xvda1"}, fakeRunner.RunCommands[1])
}

func TestLinuxFormatWhenUsingSwapFsAndPartitionIsSwap(t *testing.T) {
	fakeRunner := &testsys.FakeCmdRunner{}
	fakeFs := &testsys.FakeFileSystem{}
	fakeRunner.CommandResults = map[string][]string{
		"blkid -p /dev/xvda1": []string{`xxxxx TYPE="swap" yyyy zzzz`, ""},
	}

	formatter := NewLinuxFormatter(fakeRunner, fakeFs)
	formatter.Format("/dev/xvda1", FileSystemSwap)

	assert.Equal(t, 1, len(fakeRunner.RunCommands))
	assert.Equal(t, []string{"blkid", "-p", "/dev/xvda1"}, fakeRunner.RunCommands[0])
}

func TestLinuxFormatWhenUsingExt4FsWithLazyItableSupport(t *testing.T) {
	fakeRunner := &testsys.FakeCmdRunner{}
	fakeFs := &testsys.FakeFileSystem{}
	fakeFs.WriteToFile("/sys/fs/ext4/features/lazy_itable_init", "")
	fakeRunner.CommandResults = map[string][]string{
		"blkid -p /dev/xvda1": []string{`xxxxx TYPE="ext2" yyyy zzzz`, ""},
	}

	formatter := NewLinuxFormatter(fakeRunner, fakeFs)
	formatter.Format("/dev/xvda2", FileSystemExt4)

	assert.Equal(t, 2, len(fakeRunner.RunCommands))
	assert.Equal(t, []string{"mke2fs", "-t", "ext4", "-j", "-E", "lazy_itable_init=1", "/dev/xvda2"}, fakeRunner.RunCommands[1])
}

func TestLinuxFormatWhenUsingExt4FsWithoutLazyItableSupport(t *testing.T) {
	fakeRunner := &testsys.FakeCmdRunner{}
	fakeFs := &testsys.FakeFileSystem{}
	fakeRunner.CommandResults = map[string][]string{
		"blkid -p /dev/xvda1": []string{`xxxxx TYPE="ext2" yyyy zzzz`, ""},
	}

	formatter := NewLinuxFormatter(fakeRunner, fakeFs)
	formatter.Format("/dev/xvda2", FileSystemExt4)

	assert.Equal(t, 2, len(fakeRunner.RunCommands))
	assert.Equal(t, []string{"mke2fs", "-t", "ext4", "-j", "/dev/xvda2"}, fakeRunner.RunCommands[1])
}

func TestLinuxFormatWhenUsingExt4FsAndPartitionIsExt4(t *testing.T) {
	fakeRunner := &testsys.FakeCmdRunner{}
	fakeFs := &testsys.FakeFileSystem{}

	fakeRunner.CommandResults = map[string][]string{
		"blkid -p /dev/xvda1": []string{`xxxxx TYPE="ext4" yyyy zzzz`, ""},
	}

	formatter := NewLinuxFormatter(fakeRunner, fakeFs)
	formatter.Format("/dev/xvda1", FileSystemExt4)

	assert.Equal(t, 1, len(fakeRunner.RunCommands))
	assert.Equal(t, []string{"blkid", "-p", "/dev/xvda1"}, fakeRunner.RunCommands[0])
}
