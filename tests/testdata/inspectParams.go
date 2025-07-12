package testdata

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/volodymyrzuyev/goCsInspect/pkg/common"
	"github.com/volodymyrzuyev/goCsInspect/pkg/types"
	"gopkg.in/yaml.v3"
)

func GetInspectParams() map[string]types.InspectParameters {
	protoPath := common.GetAbsolutePath(filepath.Join("tests", "inspectParams"))

	fs, err := os.ReadDir(protoPath)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]types.InspectParameters)

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		name := strings.ReplaceAll(f.Name(), ".yaml", "")

		toParse, err := os.ReadFile(filepath.Join(protoPath, f.Name()))
		if err != nil {
			panic(err)
		}

		params := types.InspectParameters{}
		err = yaml.Unmarshal(toParse, &params)
		if err != nil {
			panic(err)
		}

		ret[name] = params
	}

	return ret
}
