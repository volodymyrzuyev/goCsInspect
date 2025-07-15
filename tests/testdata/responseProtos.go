package testdata

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"gopkg.in/yaml.v3"
)

var ResponseProtosLocation = filepath.Join(GetTestDirectory(), "responseProtos")

func GetResponseProtos() map[string]*protobuf.CEconItemPreviewDataBlock {
	fs, err := os.ReadDir(ResponseProtosLocation)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]*protobuf.CEconItemPreviewDataBlock)

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		name := strings.ReplaceAll(f.Name(), ".yaml", "")

		toParse, err := os.ReadFile(filepath.Join(ResponseProtosLocation, f.Name()))
		if err != nil {
			panic(err)
		}

		newProto := protobuf.CEconItemPreviewDataBlock{}
		err = yaml.Unmarshal(toParse, &newProto)
		if err != nil {
			panic(err)
		}

		ret[name] = &newProto
	}

	return ret
}
