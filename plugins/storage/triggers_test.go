package storage

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
	"testing"

	. "github.com/onsi/gomega"
)

// TestRepairRegistryOwnership covers the install-time repair that exists
// to fix #8557 on systems that already ran the buggy 0.38.0 install. The
// repair must rewrite ownership without rewriting the file mode, and
// must tolerate a missing registry tree (clean install).
func TestRepairRegistryOwnership(t *testing.T) {
	RegisterTestingT(t)
	withTempLibRoot(t)

	Expect(EnsureEntriesDirectory()).To(Succeed())
	path := entryPath("legacy-deadbeef")
	Expect(os.WriteFile(path, []byte(`{"name":"legacy-deadbeef"}`), 0640)).To(Succeed())

	Expect(repairRegistryOwnership()).To(Succeed())

	info, err := os.Stat(path)
	Expect(err).NotTo(HaveOccurred())
	Expect(info.Mode().Perm()).To(Equal(os.FileMode(0640)))

	stat, ok := info.Sys().(*syscall.Stat_t)
	Expect(ok).To(BeTrue())

	current, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	Expect(strconv.Itoa(int(stat.Uid))).To(Equal(current.Uid))

	Expect(os.RemoveAll(RegistryDirectory())).To(Succeed())
	Expect(repairRegistryOwnership()).To(Succeed())
}
