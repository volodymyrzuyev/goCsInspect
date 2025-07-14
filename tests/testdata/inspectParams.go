package testdata

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/volodymyrzuyev/goCsInspect/pkg/inspectParams"
	"gopkg.in/yaml.v3"
)

func GetInspectParams() map[string]inspectParams.InspectParameters {
	protoPath := filepath.Join(GetTestDirectory(), "inspectParams")

	fs, err := os.ReadDir(protoPath)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]inspectParams.InspectParameters)

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		name := strings.ReplaceAll(f.Name(), ".yaml", "")

		toParse, err := os.ReadFile(filepath.Join(protoPath, f.Name()))
		if err != nil {
			panic(err)
		}

		params := inspectParams.InspectParameters{}
		err = yaml.Unmarshal(toParse, &params)
		if err != nil {
			panic(err)
		}

		ret[name] = params
	}

	return ret
}
