package get_gid

import "github.com/pefish/go-hack/pkg/g"

func GetGid() int64 {
	return g.GetG().Id
}
