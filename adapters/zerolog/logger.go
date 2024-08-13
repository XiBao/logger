package zerolog

import (
	"path"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		dir, fileName := path.Split(file)
		_, lastPath := path.Split(strings.TrimSuffix(dir, "/"))
		filePath := fileName
		if lastPath != "" {
			filePath = path.Join(lastPath, fileName)
		}
		return filePath + ":" + strconv.Itoa(line)
	}
}
