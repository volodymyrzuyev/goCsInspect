package detailer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert/yaml"
	"github.com/volodymyrzuyev/goCsInspect/pkg/item"
	"github.com/volodymyrzuyev/goCsInspect/tests/testdata"
)

type expected struct {
	Item *item.Item
	Err  error
}

func getExpected() map[string]*expected {
	protoPath := filepath.Join(testdata.GetTestDirectory(), "detailerGoldenOutput")

	fs, err := os.ReadDir(protoPath)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]*expected)

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		name := strings.ReplaceAll(f.Name(), ".yaml", "")

		toParse, err := os.ReadFile(filepath.Join(protoPath, f.Name()))
		if err != nil {
			panic(err)
		}

		newExpected := expected{}
		err = yaml.Unmarshal(toParse, &newExpected)
		if err != nil {
			panic(err)
		}

		ret[name] = &newExpected
	}

	return ret
}

func TestDetailSkin(t *testing.T) {
	slog.SetDefault(slog.New(slog.DiscardHandler))
	detailer, err := NewDetailerGameFiles(
		filepath.Join(
			filepath.Join(filepath.Dir(testdata.GetTestDirectory()), "game_files"),
			"csgo_english.txt",
		),
		filepath.Join(
			filepath.Join(filepath.Dir(testdata.GetTestDirectory()), "game_files"),
			"items_game.txt",
		),
	)
	if err != nil {
		t.Fatalf("could not create detailer %v", err)
	}

	tests := testdata.GetResponseProtos()
	expected := getExpected()

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {

			actual, actErr := detailer.DetailProto(input)

			cur, ok := expected[name]
			if !ok {
				t.Fatal("A test case does not have a golden output")
			}

			assert.Equal(t, actErr, cur.Err, "errors should be the same. Resp error %v", actErr)

			if actErr != nil {
				return
			}

			// proto comparison
			assert.Equal(t, cur.Item.Accountid, actual.Accountid, "Proto should be the same")
			assert.Equal(t, cur.Item.Itemid, actual.Itemid, "Proto should be the same")
			assert.Equal(t, cur.Item.Defindex, actual.Defindex, "Proto should be the same")
			assert.Equal(t, cur.Item.Paintindex, actual.Paintindex, "Proto should be the same")
			assert.Equal(t, cur.Item.Rarity, actual.Rarity, "Proto should be the same")
			assert.Equal(t, cur.Item.Quality, actual.Quality, "Proto should be the same")
			assert.Equal(t, cur.Item.Paintwear, actual.Paintwear, "Proto should be the same")
			assert.Equal(t, cur.Item.Paintseed, actual.Paintseed, "Proto should be the same")
			assert.Equal(
				t,
				cur.Item.Killeaterscoretype,
				actual.Killeaterscoretype,
				"Proto should be the same",
			)
			assert.Equal(
				t,
				cur.Item.Killeatervalue,
				actual.Killeatervalue,
				"Proto should be the same",
			)
			assert.Equal(t, cur.Item.Customname, actual.Customname, "Proto should be the same")
			assert.Equal(t, cur.Item.Inventory, actual.Inventory, "Proto should be the same")
			assert.Equal(t, cur.Item.Origin, actual.Origin, "Proto should be the same")
			assert.Equal(t, cur.Item.Questid, actual.Questid, "Proto should be the same")
			assert.Equal(t, cur.Item.Dropreason, actual.Dropreason, "Proto should be the same")
			assert.Equal(t, cur.Item.Musicindex, actual.Musicindex, "Proto should be the same")
			assert.Equal(t, cur.Item.Entindex, actual.Entindex, "Proto should be the same")
			assert.Equal(t, cur.Item.Petindex, actual.Petindex, "Proto should be the same")

			assert.Equal(
				t,
				fmt.Sprintf("%.15f", cur.Item.FloatValue),
				fmt.Sprintf("%.15f", actual.FloatValue),
				"Float values should be same",
			)
			assert.Equal(t, cur.Item.MinFloat, actual.MinFloat, "MinFloat should be the same")
			assert.Equal(t, cur.Item.MaxFloat, actual.MaxFloat, "MaxFloat should be the same")
			assert.Equal(t, cur.Item.ItemName, actual.ItemName, "ItemName should be the same")
			assert.Equal(
				t,
				cur.Item.QualityName,
				actual.QualityName,
				"QualityName should be the same",
			)
			assert.Equal(t, cur.Item.WeaponType, actual.WeaponType, "WeaponType should be the same")
			assert.Equal(t, cur.Item.RarityName, actual.RarityName, "RarityName should be the same")
			assert.Equal(t, cur.Item.WearName, actual.WearName, "WearName should be the same")
			assert.Equal(
				t,
				cur.Item.FullItemName,
				actual.FullItemName,
				"FullItemName should be the same",
			)
			assert.Equal(t, cur.Item.Stickers, actual.Stickers, "Stickers should be the same")
			assert.Equal(t, cur.Item.Keychains, actual.Keychains, "Keychains should be the same")
		})
	}

}
