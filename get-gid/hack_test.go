package get_gid

import "testing"
import "github.com/pefish/go-test-assert"

func TestGetGid(t *testing.T) {
	gid := GetGid()
	test.Equal(t, true, gid > 0 && gid < 100)
}