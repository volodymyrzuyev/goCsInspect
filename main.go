package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/volodymyrzuyev/goCsInspect/accounts"
	"github.com/volodymyrzuyev/goCsInspect/detailer"
	"github.com/volodymyrzuyev/goCsInspect/logger"
	"github.com/volodymyrzuyev/goCsInspect/types"
)

// itemid:42829028921 defindex:4751 paintindex:0 rarity:6 quality:4 inventory:296 origin:23
func main() {
	log := logger.NewLogger(os.Stdout)

	detailer, err := detailer.NewDetailer(log)
	if err != nil {
		panic(err.Error())
	}

	err = godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dance := types.Credentials{SharedSecret: os.Getenv("sharedSecret"), Username: os.Getenv("userName"), Password: os.Getenv("password")}

	clientManager := accounts.NewClientManager(log)

	err = clientManager.AddClient(dance)
	params, _ := types.ParseInspectLink("steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M640180078099397070A8839911627D9837170024513331701")

	skin, err := clientManager.InspectSkin(params)
	if err != nil {
		log.Error("%v", err)
		return
	}

	log.Debug("%+v", skin)
	log.Debug("%+v", detailer.GetDetails(skin, params))
}
