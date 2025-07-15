package testdata

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/volodymyrzuyev/goCsInspect/pkg/inspect"
	"gopkg.in/yaml.v3"
)

var InspectParamsLocation = filepath.Join(GetTestDirectory(), "inspect")

func GetInspectParams() map[string]inspect.Params {
	fs, err := os.ReadDir(InspectParamsLocation)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]inspect.Params)

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		name := strings.ReplaceAll(f.Name(), ".yaml", "")

		toParse, err := os.ReadFile(filepath.Join(InspectParamsLocation, f.Name()))
		if err != nil {
			panic(err)
		}

		params := inspect.Params{}
		err = yaml.Unmarshal(toParse, &params)
		if err != nil {
			panic(err)
		}

		ret[name] = params
	}

	return ret
}
