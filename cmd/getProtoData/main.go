package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/joho/godotenv"
	"github.com/volodymyrzuyev/goCsInspect/internal/client"
	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common"
	"github.com/volodymyrzuyev/goCsInspect/pkg/config"
	"github.com/volodymyrzuyev/goCsInspect/pkg/creds"
	"github.com/volodymyrzuyev/goCsInspect/pkg/inspect"
	"github.com/volodymyrzuyev/goCsInspect/pkg/logger"
	"gopkg.in/yaml.v3"
)

type resourceFetcher struct {
	TmLink      string
	InspectLink string
}

var dataLocation = common.GetAbsolutePath("tests")

func main() {
	log := slog.New(logger.NewHandler(slog.LevelInfo, os.Stdout))
	slog.SetDefault(log)

	err := godotenv.Load()
	if err != nil {
		slog.Error("Unable to load .evn")
		panic(err)
	}

	dance := creds.Credentials{
		SharedSecret: os.Getenv("GenDetailerTestDataSharedSecret"),
		Username:     os.Getenv("GenDetailerTestDataUserName"),
		Password:     os.Getenv("GenDetailerTestDataPassword"),
	}

	gcHandler := gcHandler.NewGcHandler(log)

	client, err := client.NewInspectClient(
		dance,
		config.DefaultConfig.ClientCooldown,
		gcHandler,
		log,
	)
	if err != nil {
		slog.Error("Can't create client")
		panic(err)
	}

	err = client.LogIn()
	if err != nil {
		slog.Error("Can't login as client")
		panic(err)
	}

	resources := make(map[string]resourceFetcher)
	f, err := os.ReadFile(filepath.Join(dataLocation, "protoFetchData.yaml"))
	if err != nil {
		slog.Error("Can't open data to fetch")
		panic(err)
	}

	err = yaml.Unmarshal(f, resources)
	if err != nil {
		slog.Error("Can't Unmarshal fetch date")
		panic(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "-skip" {
		slog.Info("Due to the (-skip) flag, skinping fetching data for protos that exist")
		fs, err := os.ReadDir(filepath.Join(dataLocation, "protos"))
		if err != nil {
			slog.Error("Can't get proto directory")
			panic(err)
		}
		for _, n := range fs {
			delete(resources, strings.ReplaceAll(n.Name(), ".yaml", ""))
		}
	}

	i := 1
	for name, r := range resources {
		time.Sleep(config.DefaultConfig.ClientCooldown + 2*time.Second)

		params, err := inspect.ParseInspectLink(r.InspectLink)
		if err != nil {
			slog.Error("Err parsing inspect link", "name", name)
		}
		requestProto, _ := params.GenerateGcRequestProto()

		ctx, cancel := context.WithTimeout(context.TODO(), config.DefaultConfig.RequestTTl)

		repProto, err := client.InspectItem(ctx, requestProto)
		if err != nil {
			slog.Error("Err getting skin", "error", name, "InspectLink", r.InspectLink)
			panic(err)
		}

		slog.Info(fmt.Sprintf("Status: Got response for (%s) %+v", name, repProto))
		storeProtos(name, r, repProto)
		storeInspectParams(name, r, repProto, params)
		slog.Info(
			fmt.Sprintf(
				"Status: Finished (%v), %3.2f%% done!",
				name,
				float64(i)/float64(len(resources))*100,
			),
		)
		i++

		cancel()
	}
}

func storeProtos(name string, r resourceFetcher, repProto *protobuf.CEconItemPreviewDataBlock) {
	storeLocation := path.Join(
		filepath.Join(dataLocation, filepath.Join("responseProtos", name+".yaml")),
	)
	output, err := os.OpenFile(storeLocation, os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		slog.Error(fmt.Sprintf("Error opening store file: %s", storeLocation))
		panic(err)
	}
	defer output.Close()

	fmt.Fprintf(output, "# inspectLink: %s\n", r.InspectLink)
	fmt.Fprintf(output, "# tmLink: %s\n", r.TmLink)
	fmt.Fprintf(output, "# proto: %+v\n", repProto)

	encoder := yaml.NewEncoder(output)
	encoder.SetIndent(4)
	err = encoder.Encode(repProto)
	if err != nil {
		slog.Error(fmt.Sprintf("Error encoding proto %s", storeLocation))
		panic(err)
	}
}

func storeInspectParams(
	name string,
	r resourceFetcher,
	repProto *protobuf.CEconItemPreviewDataBlock,
	params inspect.Parameters,
) {
	storeLocation := path.Join(
		filepath.Join(dataLocation, filepath.Join("inspectParams", name+".yaml")),
	)
	output, err := os.OpenFile(storeLocation, os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		slog.Error(fmt.Sprintf("Error opening store file: %s", storeLocation))
		panic(err)
	}
	defer output.Close()

	fmt.Fprintf(output, "# inspectLink: %s\n", r.InspectLink)
	fmt.Fprintf(output, "# tmLink: %s\n", r.TmLink)
	fmt.Fprintf(output, "# proto: %+v\n", repProto)

	encoder := yaml.NewEncoder(output)
	encoder.SetIndent(4)
	err = encoder.Encode(params)
	if err != nil {
		slog.Error(fmt.Sprintf("Error encoding proto %s", storeLocation))
		panic(err)
	}
}
