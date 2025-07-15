package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	mainhelpers "github.com/volodymyrzuyev/goCsInspect/internal/mainHelpers"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/inspect"
	"github.com/volodymyrzuyev/goCsInspect/pkg/logger"
	"github.com/volodymyrzuyev/goCsInspect/tests/testdata"
	"gopkg.in/yaml.v3"
)

var (
	cfgLocation string
	skip        bool
)

func init() {
	flag.StringVar(
		&cfgLocation, "config", config.DefaultConfigLocation, "configuration file used for the api",
	)

	flag.BoolVar(&skip, "skip", false, "skips already fetched data")
}

type resourceFetcher struct {
	TmLink      string
	InspectLink string
}

func main() {
	flag.Parse()

	cfg, err := config.ParseConfig(cfgLocation)
	if err != nil {
		fmt.Println("invalid configuration location, stopping")
		os.Exit(1)
	}

	l := slog.New(logger.NewHandler(cfg.GetLogLevel(), os.Stdout))
	slog.SetDefault(l)
	lt := l.WithGroup("Main")

	if err = os.MkdirAll(testdata.InspectParamsLocation, 0761); err != nil {
		lt.Error("unable to find/create inspect params directory, stopping", "error", err)
		os.Exit(1)
	}

	if err = os.MkdirAll(testdata.ResponseProtosLocation, 0761); err != nil {
		lt.Error("unable to find/create response proto directory, stopping", "error", err)
		os.Exit(1)
	}

	dataFetchingResources := make(map[string]resourceFetcher)
	f, err := os.ReadFile(testdata.DataToFetchProtosLocation)
	if err != nil {
		lt.Error("can't open data to fetch, stopping", "error", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(f, dataFetchingResources)
	if err != nil {
		lt.Error("can't unmarshal data to fetch, stopping", "error", err)
		os.Exit(1)
	}

	if skip {
		slog.Info("since \"--skip\" was passed, skipping files that already exist")
		fs, err := os.ReadDir(testdata.ResponseProtosLocation)
		if err != nil {
			lt.Error("unable to find/create response proto directory, stopping", "error", err)
			os.Exit(1)
		}
		for name, n := range fs {
			lt.Debug("skipping existing test", "name", name)
			delete(dataFetchingResources, strings.ReplaceAll(n.Name(), ".yaml", ""))
		}
	}

	estimatedTotalTime := time.Duration(len(dataFetchingResources)) * cfg.ClientCooldown

	if estimatedTotalTime > cfg.RequestTTl {
		lt.Info(
			"Total estimated request time exceeds Request TTL; adjusting to prevent timeouts.",
		)
		cfg.RequestTTl = estimatedTotalTime
	}

	cm := mainhelpers.InitClientManagerNoStorage(cfg, lt, l)

	for _, a := range cfg.Accounts {
		err := cm.AddClient(a)
		if err != nil {
			lt.Warn(
				fmt.Sprintf("client %v unable to login, won't be used", a.Username),
				"error",
				err,
			)
		}
	}

	var wg sync.WaitGroup
	for name, test := range dataFetchingResources {
		params, err := inspect.ParseInspectLink(test.InspectLink)
		if err != nil {
			lt.Warn("error parsing inspect link for test, skipping", "name", name)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			proto, err := cm.GetProto(params)
			if err != nil {
				lt.Warn("got error when fetching proto, skipping", "test_name", name, "error", err)
				return
			}
			lt.Debug("got response", "test_name", name, "proto", fmt.Sprintf("%+v", proto))
			storeProtos(name, test, proto, lt)
			storeInspectParams(name, test, proto, params, lt)
			lt.Info("finished fetching data", "name", name)
			return
		}()
	}
	wg.Wait()
}

func storeProtos(
	name string,
	r resourceFetcher,
	repProto *protobuf.CEconItemPreviewDataBlock,
	lt *slog.Logger,
) {
	storeLocation := filepath.Join(testdata.ResponseProtosLocation, name+".yaml")

	output, err := os.Create(storeLocation)
	if err != nil {
		lt.Error("error opening file, skipping", "test_name", name, "error", err)
		return
	}
	defer output.Close()

	fmt.Fprintf(output, "# inspectLink: %s\n", r.InspectLink)
	fmt.Fprintf(output, "# tmLink: %s\n", r.TmLink)
	fmt.Fprintf(output, "# proto: %+v\n", repProto)

	encoder := yaml.NewEncoder(output)
	encoder.SetIndent(4)
	err = encoder.Encode(repProto)
	if err != nil {
		lt.Error("error encoding test, skipping", "test_name", name, "error", err)
		return
	}
}

func storeInspectParams(
	name string,
	r resourceFetcher,
	repProto *protobuf.CEconItemPreviewDataBlock,
	params inspect.Parameters,
	lt *slog.Logger,
) {
	storeLocation := filepath.Join(testdata.InspectParamsLocation, name+".yaml")
	output, err := os.Create(storeLocation)
	if err != nil {
		lt.Error("error opening file, skipping", "test_name", name, "error", err)
		return
	}
	defer output.Close()

	fmt.Fprintf(output, "# inspectLink: %s\n", r.InspectLink)
	fmt.Fprintf(output, "# tmLink: %s\n", r.TmLink)
	fmt.Fprintf(output, "# proto: %+v\n", repProto)

	encoder := yaml.NewEncoder(output)
	encoder.SetIndent(4)
	err = encoder.Encode(params)
	if err != nil {
		lt.Error("error encoding test, skipping", "test_name", name, "error", err)
		return
	}
}
