package lfscommand

import (
	"github.com/iikira/lfsutil/lfs"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseOutName(input, full, out string) string {
	if input == full { // 都是文件
		return filepath.Join(out, filepath.Base(full))
	}

	return filepath.Join(out, strings.TrimPrefix(full, input))
}

func ParseInputObject(ogs ...string) []*lfs.Object {
	objects := make([]*lfs.Object, 0, len(ogs))
	for _, og := range ogs {
		li := strings.LastIndex(og, ":")
		if li == -1 {
			objects = append(objects, &lfs.Object{
				OID: og,
			})
			continue
		}
		size, err := strconv.ParseInt(og[li+1:], 10, 64)
		if err != nil {
			objects = append(objects, &lfs.Object{
				OID: og[:li],
			})
			continue
		}
		objects = append(objects, &lfs.Object{
			OID:  og[:li],
			Size: size,
		})
	}

	return objects
}
