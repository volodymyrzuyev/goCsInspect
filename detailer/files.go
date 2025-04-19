package detailer

import (
	"os"
	"strconv"

	"github.com/andygrunwald/vdf"
	"github.com/volodymyrzuyev/goCsInspect/logger"
)

type files struct {
	csgoEnglish *os.File
	itemsGame   *os.File
}

func (d *detailer) populateAssets(f files) {

	p := parser{log: d.log}

	csgoEnglishParser := vdf.NewParser(f.csgoEnglish)
	csgoEnglish, err := csgoEnglishParser.Parse()
	if err != nil {
		panic("Could not parse csgo_english.txt")
	}
	englishTokensProto, ok := csgoEnglish["lang"].(map[string]interface{})["Tokens"].(map[string]interface{})
	if !ok {
		panic("Could not parse csgo_english.txt")
	}
	d.englishTokens = convertMapUnsafe(englishTokensProto)

	itemsGameParser := vdf.NewParser(f.itemsGame)
	itemsGame, err := itemsGameParser.Parse()
	if err != nil {
		panic("Could not parser items_game.txt")
	}

	raritieProtos, ok := itemsGame["items_game"].(map[string]interface{})["rarities"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error getting rarities.")
	}

	colorProtos, ok := itemsGame["items_game"].(map[string]interface{})["colors"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing colors.")
	}
	d.raritieDefenitions = p.createRaritieMap(raritieProtos, colorProtos)
	if len(d.raritieDefenitions) == 0 {
		panic("No rarities loaded, check items_game.txt")
	}

	paintKitProtos, ok := itemsGame["items_game"].(map[string]interface{})["paint_kits"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing paint kits.")
	}
	d.paintKits = p.parsePaintKits(paintKitProtos, d.englishTokens)
	if len(d.paintKits) == 0 {
		panic("No paint_kits loaded, check items_game.txt")
	}

	qualitieDefenitions, ok := itemsGame["items_game"].(map[string]interface{})["qualities"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing qualities.")
	}
	d.qualities = p.parseQualities(qualitieDefenitions, d.englishTokens)
	if len(d.qualities) == 0 {
		panic("No qualities loaded, check items_game.txt")
	}

	itemDefenitions, ok := itemsGame["items_game"].(map[string]interface{})["items"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing item defenitions.")
	}

	stickerDefenitions, ok := itemsGame["items_game"].(map[string]interface{})["sticker_kits"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing sticker defenitions.")
	}

	itemSets, ok := itemsGame["items_game"].(map[string]interface{})["item_sets"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing item sets.")
	}

	keyChains, ok := itemsGame["items_game"].(map[string]interface{})["keychain_definitions"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing Keychains.")
	}

	paintKitRarities, ok := itemsGame["items_game"].(map[string]interface{})["paint_kits_rarity"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Error parsing paint kit parities.")
	}

	prefabs, ok := itemsGame["items_game"].(map[string]interface{})["prefabs"].(map[string]interface{})
	if !ok {
		panic("Could not parse items_game.txt. Errors parsing prefabs.")
	}

	d.itemDefenitions = itemDefenitions
	d.stickerDefenitions = stickerDefenitions
	d.paintKitRarities = paintKitRarities
	d.keyChains = keyChains
	d.itemSets = itemSets
	d.prefabs = prefabs
}

type parser struct {
	log logger.Logger
}

type raritie struct {
	weapon    string
	character string
	regular   string
	color     color
}

type color struct {
	colorName string
	hexColor  string
}

func (p *parser) createRaritieMap(raritieProtos map[string]any, colorProtos map[string]any) map[uint32]raritie {
	raritieMap := make(map[uint32]raritie)

	for raritieNameID, raritieInterface := range raritieProtos {
		curRaritie := raritieInterface.(map[string]interface{})

		r := raritie{}

		idString, ok := curRaritie["value"]
		if !ok {
			p.log.Debug("Could not parse idString for raritie: %v. Skipping...", raritieNameID)
			continue
		}

		if r.weapon, ok = curRaritie["loc_key_weapon"].(string); !ok {
			p.log.Debug("Could not parse loc_key_weapon for raritie: %v. Skipping...", raritieNameID)
			continue
		}

		if r.character, ok = curRaritie["loc_key_character"].(string); !ok {
			p.log.Debug("Could not parse loc_key_character for raritie: %v. Skipping...", raritieNameID)
			continue
		}

		if r.regular, ok = curRaritie["loc_key"].(string); !ok {
			p.log.Debug("Could not parse loc_key for raritie: %v. Skipping...", raritieNameID)
			continue
		}

		colorID, ok := curRaritie["color"].(string)
		if !ok {
			p.log.Debug("Could not parse color: %v. Skipping...", raritieNameID)
			continue
		}

		colorProto, ok := colorProtos[colorID].(map[string]any)
		if !ok {
			p.log.Debug("Could not find color %v for raritie: %v. Skipping...", colorID, raritieNameID)
			continue
		}

		if r.color.colorName, ok = colorProto["color_name"].(string); !ok {
			p.log.Debug("Could not parse color_name for color: %v for raritie %v. Skipping...", colorID, raritieNameID)
		}

		if r.color.hexColor, ok = colorProto["hex_color"].(string); !ok {
			p.log.Debug("Could not parse hex_color for color: %v for raritie: %v. Skipping...", colorID, raritieNameID)
		}

		if id, err := strconv.ParseUint(idString.(string), 10, 32); err == nil {
			raritieMap[uint32(id)] = r
		} else {
			p.log.Debug("Could not conver raritieID to uint32 for raritie: %v. Skipping...", raritieNameID)
		}

	}

	p.log.Debug("Parsed and loaded rarities succesufully")
	return raritieMap
}

type paintKit struct {
	name              string
	descriptionString string
	descriptionTag    string
	wearRemapMin      string
	wearRemapMax      string
	itemName          string
}

func (p *parser) parsePaintKits(paintKitProto map[string]any, englishTokens map[string]string) map[uint32]paintKit {
	paintKitMap := make(map[uint32]paintKit)

	for kitID, paintKitInterface := range paintKitProto {

		currentKit, ok := paintKitInterface.(map[string]any)
		if !ok {
			continue
		}

		k := paintKit{}

		if k.name, ok = currentKit["name"].(string); !ok {
			p.log.Debug("Could not find name for paintKitID: %v. Skipping...", kitID)
			continue
		}

		if k.descriptionString, ok = currentKit["description_string"].(string); !ok {
			p.log.Debug("Could not find description_string for paintKitID: %v. Skipping...", kitID)
			continue
		}

		if k.descriptionTag, ok = currentKit["description_tag"].(string); !ok {
			p.log.Debug("Could not find description_tag for paintKitID: %v. Skipping...", kitID)
			continue
		}

		// if it does not exist, it's fine. Paint kits do not have define max and min float
		// will just be blank
		k.wearRemapMin, ok = currentKit["wear_remap_min"].(string)
		k.wearRemapMax, ok = currentKit["wear_remap_max"].(string)

		if k.itemName, ok = englishTokens[k.descriptionTag[1:]]; !ok {
			p.log.Debug("Could not find description tag for paintKitID: %v. Skipping...")
			continue
		}

		if kitIDUint, err := strconv.ParseUint(kitID, 10, 32); err == nil {
			paintKitMap[uint32(kitIDUint)] = k
		} else {
			p.log.Debug("Could not part kitID to uint for %v", kitID)
		}
	}

	return paintKitMap
}

type qualitie struct {
	name        string
	englishName string
}

func (p *parser) parseQualities(qualities map[string]any, englishTokens map[string]string) map[uint32]qualitie {
	qualitieMap := make(map[uint32]qualitie)

	for qName, q := range qualities {
		curQual := q.(map[string]any)

		qual := qualitie{name: qName}

		idString, ok := curQual["value"]
		if !ok {
			p.log.Debug("Could not parse idString for qualitie: %v. Skipping...", qName)
			continue
		}

		if qual.englishName, ok = englishTokens[qName]; !ok {
			p.log.Debug("Could not parse english name for qualitie: %v. Skipping...", qName)
			continue
		}

		if id, err := strconv.ParseUint(idString.(string), 10, 32); err == nil {
			qualitieMap[uint32(id)] = qual
		} else {
			p.log.Debug("Could not convert id to uint32 for qualitie: %v. Skipping...", qName)
			continue
		}

	}

	return qualitieMap
}

func convertMapUnsafe(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}
