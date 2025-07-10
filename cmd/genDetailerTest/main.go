package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/joho/godotenv"
	"github.com/volodymyrzuyev/goCsInspect/config"
	"github.com/volodymyrzuyev/goCsInspect/internal/client"
	"github.com/volodymyrzuyev/goCsInspect/internal/gcHandler"
	"github.com/volodymyrzuyev/goCsInspect/pkg/types"
)

const fileLocation = "./pkg/detailer/detailerDataNew_test.go"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dance := types.Credentials{
		SharedSecret: os.Getenv("GenDetailerTestDataSharedSecret"),
		Username:     os.Getenv("GenDetailerTestDataUserName"),
		Password:     os.Getenv("GenDetailerTestDataPassword"),
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	gcHandler := gcHandler.NewGcHandler(config.TimeOutDuration)

	client, err := client.NewInspectClient(config.DefaultClientConfig, gcHandler, dance)
	if err != nil {
		panic(err)
	}

	err = client.LogIn()
	if err != nil {
		panic(err)
	}

	links := make(map[string]string)
	tmLinks := make(map[string]string)

	tmLinks["Skin"] = "https://steamcommunity.com/market/listings/730/AWP%20%7C%20Printstream%20%28Factory%20New%29"
	links["Skin"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M645819192553877985A44825267139D16736572647846215177"

	tmLinks["Skin StatTrak"] = "https://steamcommunity.com/market/listings/730/StatTrak%E2%84%A2%20AWP%20%7C%20Printstream%20%28Factory%20New%29"
	links["Skin StatTrak"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M649195897111665513A44158728090D11837336287321839237"

	tmLinks["Skin StatTrak Stickers"] = "https://steamcommunity.com/market/listings/730/StatTrak%E2%84%A2%20AWP%20%7C%20Printstream%20%28Factory%20New%29"
	links["Skin StatTrak Stickers"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M659329991369725526A44816452463D640747018096259789"

	tmLinks["Skin Stickers Keychain"] = "https://steamcommunity.com/market/listings/730/AK-47%20%7C%20Vulcan%20%28Factory%20New%29"
	links["Skin Stickers Keychain"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M650322517130496214A44417931234D11961637775872526773"

	tmLinks["Knife Vanila"] = "https://steamcommunity.com/market/listings/730/%E2%98%85%20Bayonet"
	links["Knife Vanila"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M666085390842291119A45045640461D5362143296623681792"

	tmLinks["Knife Vanila StatTrak"] = "https://steamcommunity.com/market/listings/730/%E2%98%85%20StatTrak%E2%84%A2%20Bayonet"
	links["Knife Vanila StatTrak"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M660453062733259347A44410422468D776985712673648068"

	tmLinks["Knife StatTrack"] = "https://steamcommunity.com/market/listings/730/%E2%98%85%20StatTrak%E2%84%A2%20Bayonet%20%7C%20Gamma%20Doppler%20%28Factory%20New%29"
	links["Knife StatTrack"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M657078191632766063A24514899689D17025473724194038752"

	tmLinks["Knife"] = "https://steamcommunity.com/market/listings/730/%E2%98%85%20Bayonet%20%7C%20Gamma%20Doppler%20%28Factory%20New%29"
	links["Knife"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M653700491952159799A44851210636D9387414943281127756"

	tmLinks["Gloves"] = "https://steamcommunity.com/market/listings/730/%E2%98%85%20Bloodhound%20Gloves%20%7C%20Bronzed%20%28Factory%20New%29"
	links["Gloves"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M6664853361691228811A39592993407D14448091259802652413"

	tmLinks["Sticker Gold"] = "https://steamcommunity.com/market/listings/730/Sticker%20%7C%20somebody%20%28Gold%29%20%7C%20Shanghai%202024"
	links["Sticker Gold"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M648070992403718122A45045428824D14025976628708375625"

	tmLinks["Sticker Holo"] = "https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Lit%20%28Holo%29"
	links["Sticker Holo"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M650322792241404634A44809037288D922326464949352065"

	tmLinks["Sticker Glitter"] = "https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Lotus%20%28Glitter%29"
	links["Sticker Glitter"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M643567392820418477A42466192195D9975096757687244901"

	tmLinks["Sticker Foil"] = "https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Boom%20%28Foil%29"
	links["Sticker Foil"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M661581516166896654A44745305965D11845634205072979149"

	tmLinks["Sticker Paper"] = "https://steamcommunity.com/market/listings/730/Sticker%20%7C%20BIG%20%7C%20Krakow%202017"
	links["Sticker Paper"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M636811354732916948A44663924163D2360593803736594085"

	tmLinks["Sticker Lenticular"] = "https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Freeze%20%28Lenticular%29"
	links["Sticker Lenticular"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M658204091488801075A45029291629D11962765263242622181"

	tmLinks["Agent"] = "https://steamcommunity.com/market/listings/730/Sir%20Bloody%20Miami%20Darryl%20%7C%20The%20Professionals"
	links["Agent"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M657077916547623203A44924191790D2903522233577573364"

	tmLinks["Agent Patch"] = "https://steamcommunity.com/market/listings/730/Sir%20Bloody%20Miami%20Darryl%20%7C%20The%20Professionals"
	links["Agent Patch"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M653697663271621789A44360148598D433416971386187216"

	tmLinks["Keychain"] = "https://steamcommunity.com/market/listings/730/Charm%20%7C%20Baby%20Karat%20T"
	links["Keychain"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M658204091494906705A45048018784D10296998617952951029"

	tmLinks["Patch"] = "https://steamcommunity.com/market/listings/730/Patch%20%7C%20Elder%20God"
	links["Patch"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M644693292653772584A23430518132D3314979815666773509"

	tmLinks["Pin"] = "https://steamcommunity.com/market/listings/730/Alyx%20Pin"
	links["Pin"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M666085390839497849A43786144675D2747665804181709680"

	tmLinks["Graffiti"] = "https://steamcommunity.com/market/listings/730/Sealed%20Graffiti%20%7C%20Drug%20War%20Veteran"
	links["Graffiti"] = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M643567392817240427A44225554008D11566448795219559360"

	output, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(output, "package detailer\n")
	fmt.Fprintf(output, "import (\n")
	fmt.Fprintf(output, "	\"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf\"\n")
	fmt.Fprintf(output, "	\"github.com/volodymyrzuyev/goCsInspect/pkg/item\"")
	fmt.Fprintf(output, ")\n\n")
	fmt.Fprintf(output, "func getTestCasesNew() map[string]protoTestCase {\n")
	fmt.Fprintf(output, "	tests := make(map[string]protoTestCase)\n")
	fmt.Fprintf(output, "	var input *protobuf.CEconItemPreviewDataBlock\n")
	fmt.Fprintf(output, "	var expectedItem *item.Item\n")

	output.Close()

	i := 1
	for name, link := range links {
		params, _ := types.ParseInspectLink(link)
		requestProto, _ := params.GenerateGcRequestProto()
		repProto, err := client.InspectItem(requestProto)
		if err != nil {
			fmt.Printf("link: (%s), name: (%s)\n", link, name)
			panic(err)
		}

		fmt.Printf("%+v\n", repProto)
		generateProtoTestCaseReflectV2(name, link, tmLinks[name], repProto)
		fmt.Printf("Finished (%v), %3.2f%% done!\n", name, float64(i)/float64(len(links))*100)
		i++
		time.Sleep(config.RequestCooldown + 2*time.Second)
	}
	output, err = os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(output, "	return tests\n")
	fmt.Fprintf(output, "	}\n")
	output.Close()
	fmt.Println("Done!")
}

func generateProtoTestCaseReflectV2(name, link, tmLink string, repProto *protobuf.CEconItemPreviewDataBlock) {
	output, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	fmt.Fprintf(output, "// %s\n", tmLink)
	fmt.Fprintf(output, "// %s\n", link)
	fmt.Fprintf(output, "// %+v\n", repProto)
	fmt.Fprintf(output, "	input = &protobuf.CEconItemPreviewDataBlock{\n")

	if repProto == nil {
		fmt.Fprintf(output, "	// repProto was nil, no fields generated\n")
		fmt.Fprintf(output, "	},\n")
		fmt.Fprintf(output, "	// ... expected output ...\n")
		fmt.Fprintf(output, "}\n\n")
		return
	}

	protoVal := reflect.ValueOf(repProto)
	if protoVal.Kind() == reflect.Ptr {
		protoVal = protoVal.Elem() // Dereference to get the struct value
	}
	protoType := protoVal.Type()

	for i := 0; i < protoType.NumField(); i++ {
		typeField := protoType.Field(i)
		valueField := protoVal.Field(i)

		// Skip unexported fields (like state, sizeCache, unknownFields)
		if !typeField.IsExported() {
			continue
		}

		if valueField.Kind() == reflect.Ptr {
			if valueField.IsNil() {
				// This optional field was not set, so we don't include it
				continue
			}
			value := valueField.Elem() // Dereference the pointer

			switch value.Kind() {
			case reflect.Uint32:
				fmt.Fprintf(output, "			%s: uint32Pointer(%v),\n", typeField.Name, value.Uint())
			case reflect.Uint64:
				fmt.Fprintf(output, "			%s: uint64Pointer(%v),\n", typeField.Name, value.Uint())
			case reflect.Int32:
				fmt.Fprintf(output, "			%s: int32Pointer(%v),\n", typeField.Name, value.Int())
			case reflect.String:
				fmt.Fprintf(output, "			%s: stringPointer(%q),\n", typeField.Name, value.String())
			case reflect.Float32:
				fmt.Fprintf(output, "			%s: float32Pointer(%v),\n", typeField.Name, float32(value.Float()))
			case reflect.Float64:
				fmt.Fprintf(output, "			%s: float64Pointer(%v),\n", typeField.Name, value.Float())
			default:
				fmt.Fprintf(output, "			// WARNING: Unsupported pointer type for %s: %s\n", typeField.Name, value.Kind())
			}
		} else if valueField.Kind() == reflect.Slice {
			if valueField.IsNil() || valueField.Len() == 0 {
				continue
			}

			fmt.Fprintf(output, "			%s: []*protobuf.CEconItemPreviewDataBlock_Sticker{\n", typeField.Name)
			for j := 0; j < valueField.Len(); j++ {
				stickerVal := valueField.Index(j)
				if stickerVal.Kind() == reflect.Ptr && !stickerVal.IsNil() {
					fmt.Fprintf(output, "				&protobuf.CEconItemPreviewDataBlock_Sticker{\n")
					stickerElem := stickerVal.Elem()
					stickerType := stickerElem.Type()

					for k := 0; k < stickerType.NumField(); k++ {
						stickerField := stickerType.Field(k)
						stickerValue := stickerElem.Field(k)

						if !stickerField.IsExported() {
							continue
						}

						if stickerValue.Kind() == reflect.Ptr && !stickerValue.IsNil() {
							actualStickerValue := stickerValue.Elem()
							// --- ADDED FLOAT32 AND FLOAT64 CASES FOR NESTED STICKER FIELDS ---
							switch actualStickerValue.Kind() {
							case reflect.Uint32:
								fmt.Fprintf(output, "					%s: uint32Pointer(%v),\n", stickerField.Name, actualStickerValue.Uint())
							case reflect.Float32: // New: Handle float32 in sticker
								fmt.Fprintf(output, "					%s: float32Pointer(%v),\n", stickerField.Name, float32(actualStickerValue.Float()))
							case reflect.Float64: // New: Handle float64 in sticker (if any)
								fmt.Fprintf(output, "					%s: float64Pointer(%v),\n", stickerField.Name, actualStickerValue.Float())
							case reflect.String: // New: Handle string in sticker (if any)
								fmt.Fprintf(output, "					%s: stringPointer(%q),\n", stickerField.Name, actualStickerValue.String())
							// Add more cases for other sticker field types if needed
							default:
								fmt.Fprintf(output, "					// WARNING: Unsupported sticker field type for %s: %s\n", stickerField.Name, actualStickerValue.Kind())
							}
						} else if stickerValue.Kind() == reflect.String {
							// Handle direct string fields if they are not pointers but optional/omitempty
							// (though Protobuf usually uses pointers for optional scalars)
							if stickerValue.String() != "" { // Check if not empty
								fmt.Fprintf(output, "					%s: %q,\n", stickerField.Name, stickerValue.String())
							}
						}
					}
					fmt.Fprintf(output, "				},\n")
				}
			}
			fmt.Fprintf(output, "			},\n")

		} else {
			fmt.Fprintf(output, "		// WARNING: Field %s has unexpected kind: %s\n", typeField.Name, valueField.Kind())
		}
	}
	fmt.Fprintf(output, "		}\n")

	fmt.Fprintf(output, "	expectedItem = &item.Item{}\n")
	fmt.Fprintf(output, "	expectedItem.PopulateProto(input)\n")
	fmt.Fprintf(output, "	tests[\"%s\"] = protoTestCase{\n", name)
	fmt.Fprintf(output, "		input: input,\n")
	fmt.Fprintf(output, "		expectedItem: expectedItem,\n")
	fmt.Fprintf(output, "		expectedError: nil,\n")
	fmt.Fprintf(output, "	}\n\n")
}
