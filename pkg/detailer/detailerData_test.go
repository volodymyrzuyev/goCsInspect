package detailer

import (
	"github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/common/types"
)

func getTestCases() map[string]protoTestCase {
	tests := make(map[string]protoTestCase)
	var input *protobuf.CEconItemPreviewDataBlock
	// https://steamcommunity.com/market/listings/730/StatTrak%E2%84%A2%20AWP%20%7C%20Printstream%20%28Factory%20New%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M659329991369725526A44816452463D640747018096259789
	// itemid:44816452463 defindex:9 paintindex:1206 rarity:6 quality:9 paintwear:1024300654 paintseed:953 killeaterscoretype:0 killeatervalue:54 stickers:{slot:0 sticker_id:4983 rotation:-90 offset_x:0.022147745 offset_y:0.0029743314} stickers:{slot:3 sticker_id:6016 rotation:-42 offset_x:0.20414282 offset_y:0.07944468} stickers:{slot:3 sticker_id:7284 offset_x:0.07931936 offset_y:0.06176564} stickers:{slot:3 sticker_id:6694 offset_x:0.033233702 offset_y:0.064612746} stickers:{slot:3 sticker_id:7312 offset_x:-0.0047016516 offset_y:0.06625104} inventory:40 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:             uint64Pointer(44816452463),
		Defindex:           uint32Pointer(9),
		Paintindex:         uint32Pointer(1206),
		Rarity:             uint32Pointer(6),
		Quality:            uint32Pointer(9),
		Paintwear:          uint32Pointer(1024300654),
		Paintseed:          uint32Pointer(953),
		Killeaterscoretype: uint32Pointer(0),
		Killeatervalue:     uint32Pointer(54),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(4983),
				Rotation:  float32Pointer(-90),
				OffsetX:   float32Pointer(0.022147745),
				OffsetY:   float32Pointer(0.0029743314),
			},
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(3),
				StickerId: uint32Pointer(6016),
				Rotation:  float32Pointer(-42),
				OffsetX:   float32Pointer(0.20414282),
				OffsetY:   float32Pointer(0.07944468),
			},
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(3),
				StickerId: uint32Pointer(7284),
				OffsetX:   float32Pointer(0.07931936),
				OffsetY:   float32Pointer(0.06176564),
			},
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(3),
				StickerId: uint32Pointer(6694),
				OffsetX:   float32Pointer(0.033233702),
				OffsetY:   float32Pointer(0.064612746),
			},
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(3),
				StickerId: uint32Pointer(7312),
				OffsetX:   float32Pointer(-0.0047016516),
				OffsetY:   float32Pointer(0.06625104),
			},
		},
		Inventory: uint32Pointer(40),
		Origin:    uint32Pointer(8),
	}
	tests["Skin StatTrak Stickers"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.0345672890543938,
			MinFloat:     "0",
			MaxFloat:     "1",
			ItemName:     "Printstream",
			QualityName:  "StatTrak™",
			WeaponType:   "AWP",
			RarityName:   "Covert",
			WearName:     "Factory New",
			FullItemName: "StatTrak™ AWP | Printstream (Factory New)",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "stockh2021_team_evl_holo",
					Material: "Holo",
					Name:     "Sticker | Evil Geniuses (Holo) | Stockholm 2021",
				},
				types.Modification{
					Proto:    input.Stickers[1],
					CodeName: "rio2022_team_mouz_gold",
					Material: "Gold",
					Name:     "Sticker | MOUZ (Gold) | Rio 2022",
				},
				types.Modification{
					Proto:    input.Stickers[2],
					CodeName: "cph2024_team_cplx_holo",
					Material: "Holo",
					Name:     "Sticker | Complexity Gaming (Holo) | Copenhagen 2024",
				},
				types.Modification{
					Proto:    input.Stickers[3],
					CodeName: "paris2023_team_cplx_holo",
					Material: "Holo",
					Name:     "Sticker | Complexity Gaming (Holo) | Paris 2023",
				},
				types.Modification{
					Proto:    input.Stickers[4],
					CodeName: "cph2024_team_gl_holo",
					Material: "Holo",
					Name:     "Sticker | GamerLegion (Holo) | Copenhagen 2024",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/%E2%98%85%20StatTrak%E2%84%A2%20Bayonet%20%7C%20Gamma%20Doppler%20%28Factory%20New%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M657078191632766063A24514899689D17025473724194038752
	// itemid:24514899689 defindex:500 paintindex:569 rarity:6 quality:3 paintwear:1017758561 paintseed:750 killeaterscoretype:0 killeatervalue:740 inventory:90 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:             uint64Pointer(24514899689),
		Defindex:           uint32Pointer(500),
		Paintindex:         uint32Pointer(569),
		Rarity:             uint32Pointer(6),
		Quality:            uint32Pointer(3),
		Paintwear:          uint32Pointer(1017758561),
		Paintseed:          uint32Pointer(750),
		Killeaterscoretype: uint32Pointer(0),
		Killeatervalue:     uint32Pointer(740),
		Inventory:          uint32Pointer(90),
		Origin:             uint32Pointer(8),
	}
	tests["Knife StatTrack"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.0207230467349291,
			MinFloat:     "0",
			MaxFloat:     "0.08",
			ItemName:     "Gamma Doppler",
			QualityName:  "★",
			WeaponType:   "Bayonet",
			RarityName:   "Covert",
			WearName:     "Factory New",
			FullItemName: "★ StatTrak™ Bayonet | Gamma Doppler (Factory New)",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sticker%20%7C%20BIG%20%7C%20Krakow%202017
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M636811354732916948A44663924163D2360593803736594085
	// itemid:44663924163 defindex:1209 paintindex:0 rarity:3 quality:4 stickers:{slot:0 sticker_id:2103} inventory:3221225475 origin:2
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44663924163),
		Defindex:   uint32Pointer(1209),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(3),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(2103),
			},
		},
		Inventory: uint32Pointer(3221225475),
		Origin:    uint32Pointer(2),
	}
	tests["Sticker Paper"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "BIG | Krakow 2017",
			QualityName:  "Unique",
			WeaponType:   "Sticker",
			RarityName:   "High Grade",
			WearName:     "",
			FullItemName: "Sticker | BIG | Krakow 2017",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "krakow2017_team_big",
					Material: "Paper",
					Name:     "Sticker | BIG | Krakow 2017",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Freeze%20%28Lenticular%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M658204091488801075A45029291629D11962765263242622181
	// itemid:45029291629 defindex:1209 paintindex:0 rarity:6 quality:4 stickers:{slot:0 sticker_id:5955} inventory:3221225482 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(45029291629),
		Defindex:   uint32Pointer(1209),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(5955),
			},
		},
		Inventory: uint32Pointer(3221225482),
		Origin:    uint32Pointer(8),
	}
	tests["Sticker Lenticular"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Freeze (Lenticular)",
			QualityName:  "Unique",
			WeaponType:   "Sticker",
			RarityName:   "Extraordinary",
			WearName:     "",
			FullItemName: "Sticker | Freeze (Lenticular)",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "csgo10_freeze_lenticular",
					Material: "Lenticular",
					Name:     "Sticker | Freeze (Lenticular)",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Alyx%20Pin
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M666085390839497849A43786144675D2747665804181709680
	// itemid:43786144675 defindex:6134 paintindex:0 rarity:6 quality:4 inventory:49 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(43786144675),
		Defindex:   uint32Pointer(6134),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Inventory:  uint32Pointer(49),
		Origin:     uint32Pointer(8),
	}
	tests["Pin"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Alyx Pin",
			QualityName:  "Unique",
			WeaponType:   "Pin",
			RarityName:   "Extraordinary",
			WearName:     "",
			FullItemName: "Alyx Pin",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/%E2%98%85%20Bayonet
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M666085390842291119A45045640461D5362143296623681792
	// itemid:45045640461 defindex:500 paintindex:0 rarity:6 quality:3 paintwear:1052449268 paintseed:41 inventory:3221225482 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(45045640461),
		Defindex:   uint32Pointer(500),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(3),
		Paintwear:  uint32Pointer(1052449268),
		Paintseed:  uint32Pointer(41),
		Inventory:  uint32Pointer(3221225482),
		Origin:     uint32Pointer(8),
	}
	tests["Knife Vanila"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.36543238162994385,
			MinFloat:     "0.06",
			MaxFloat:     "0.8",
			QualityName:  "★",
			WeaponType:   "Bayonet",
			RarityName:   "Covert",
			WearName:     "Field-Tested",
			FullItemName: "★ Bayonet",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/%E2%98%85%20StatTrak%E2%84%A2%20Bayonet
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M660453062733259347A44410422468D776985712673648068
	// itemid:44410422468 defindex:500 paintindex:0 rarity:6 quality:3 paintwear:1040853480 paintseed:336 killeaterscoretype:0 killeatervalue:0 inventory:18 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:             uint64Pointer(44410422468),
		Defindex:           uint32Pointer(500),
		Paintindex:         uint32Pointer(0),
		Rarity:             uint32Pointer(6),
		Quality:            uint32Pointer(3),
		Paintwear:          uint32Pointer(1040853480),
		Paintseed:          uint32Pointer(336),
		Killeaterscoretype: uint32Pointer(0),
		Killeatervalue:     uint32Pointer(0),
		Inventory:          uint32Pointer(18),
		Origin:             uint32Pointer(8),
	}
	tests["Knife Vanila StatTrak"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.1349254846572876,
			MinFloat:     "0.06",
			MaxFloat:     "0.8",
			QualityName:  "★",
			WeaponType:   "Bayonet",
			RarityName:   "Covert",
			WearName:     "Minimal Wear",
			FullItemName: "★ StatTrak™ Bayonet",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sticker%20%7C%20somebody%20%28Gold%29%20%7C%20Shanghai%202024
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M648070992403718122A45045428824D14025976628708375625
	// itemid:45045428824 defindex:1209 paintindex:0 rarity:6 quality:4 stickers:{slot:0 sticker_id:8504} inventory:7 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(45045428824),
		Defindex:   uint32Pointer(1209),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(8504),
			},
		},
		Inventory: uint32Pointer(7),
		Origin:    uint32Pointer(8),
	}
	tests["Sticker Gold"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "somebody (Gold) | Shanghai 2024",
			QualityName:  "Unique",
			WeaponType:   "Sticker",
			RarityName:   "Extraordinary",
			WearName:     "",
			FullItemName: "Sticker | somebody (Gold) | Shanghai 2024",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "sha2024_signature_somebody_4_gold",
					Material: "Gold",
					Name:     "Sticker | somebody (Gold) | Shanghai 2024",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Boom%20%28Foil%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M661581516166896654A44745305965D11845634205072979149
	// itemid:44745305965 defindex:1209 paintindex:0 rarity:5 quality:4 stickers:{slot:0 sticker_id:976} inventory:11 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44745305965),
		Defindex:   uint32Pointer(1209),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(5),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(976),
			},
		},
		Inventory: uint32Pointer(11),
		Origin:    uint32Pointer(8),
	}
	tests["Sticker Foil"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Boom (Foil)",
			QualityName:  "Unique",
			WeaponType:   "Sticker",
			RarityName:   "Exotic",
			WearName:     "",
			FullItemName: "Sticker | Boom (Foil)",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "slid3_boom_foil",
					Material: "Foil",
					Name:     "Sticker | Boom (Foil)",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sir%20Bloody%20Miami%20Darryl%20%7C%20The%20Professionals
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M657077916547623203A44924191790D2903522233577573364
	// itemid:44924191790 defindex:4726 paintindex:0 rarity:6 quality:4 inventory:95 origin:23
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44924191790),
		Defindex:   uint32Pointer(4726),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Inventory:  uint32Pointer(95),
		Origin:     uint32Pointer(23),
	}
	tests["Agent"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Sir Bloody Miami Darryl | The Professionals",
			QualityName:  "Unique",
			WeaponType:   "Agent",
			RarityName:   "Master",
			WearName:     "",
			FullItemName: "Sir Bloody Miami Darryl | The Professionals",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sir%20Bloody%20Miami%20Darryl%20%7C%20The%20Professionals
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M653697663271621789A44360148598D433416971386187216
	// itemid:44360148598 defindex:4726 paintindex:0 rarity:6 quality:4 stickers:{slot:1 sticker_id:4560 scale:1.0963718} inventory:100 origin:23
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44360148598),
		Defindex:   uint32Pointer(4726),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(1),
				StickerId: uint32Pointer(4560),
				Scale:     float32Pointer(1.0963718),
			},
		},
		Inventory: uint32Pointer(100),
		Origin:    uint32Pointer(23),
	}
	tests["Agent Patch"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Sir Bloody Miami Darryl | The Professionals",
			QualityName:  "Unique",
			WeaponType:   "Agent",
			RarityName:   "Master",
			WearName:     "",
			FullItemName: "Sir Bloody Miami Darryl | The Professionals",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "patch_wildfire",
					Name:     "Patch | Wildfire",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Charm%20%7C%20Baby%20Karat%20T
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M658204091494906705A45048018784D10296998617952951029
	// itemid:45048018784 defindex:1355 paintindex:0 rarity:6 quality:4 inventory:3221225482 origin:23 keychains:{slot:0 sticker_id:30 pattern:30299}
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(45048018784),
		Defindex:   uint32Pointer(1355),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Inventory:  uint32Pointer(3221225482),
		Origin:     uint32Pointer(23),
		Keychains: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(30),
				Pattern:   uint32Pointer(30299),
			},
		},
	}
	tests["Keychain"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Baby Karat T",
			QualityName:  "Unique",
			WeaponType:   "Charm",
			RarityName:   "Extraordinary",
			WearName:     "",
			FullItemName: "Charm | Baby Karat T",
			Keychains: []types.Modification{
				types.Modification{
					Proto:    input.Keychains[0],
					CodeName: "kc_wpn_tknife_gold",
					Name:     "Charm | Baby Karat T",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/AWP%20%7C%20Printstream%20%28Factory%20New%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M645819192553877985A44825267139D16736572647846215177
	// itemid:44825267139 defindex:9 paintindex:1206 rarity:6 quality:4 paintwear:1027353092 paintseed:966 inventory:153 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44825267139),
		Defindex:   uint32Pointer(9),
		Paintindex: uint32Pointer(1206),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Paintwear:  uint32Pointer(1027353092),
		Paintseed:  uint32Pointer(966),
		Inventory:  uint32Pointer(153),
		Origin:     uint32Pointer(8),
	}
	tests["Skin"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.045938506722450256,
			MinFloat:     "0",
			MaxFloat:     "1",
			ItemName:     "Printstream",
			QualityName:  "Unique",
			WeaponType:   "AWP",
			RarityName:   "Covert",
			WearName:     "Factory New",
			FullItemName: "AWP | Printstream (Factory New)",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/AK-47%20%7C%20Vulcan%20%28Factory%20New%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M650322517130496214A44417931234D11961637775872526773
	// itemid:44417931234 defindex:7 paintindex:302 rarity:6 quality:4 paintwear:1028819163 paintseed:387 stickers:{slot:0 sticker_id:666 wear:0.23460516} stickers:{slot:1 sticker_id:1338} stickers:{slot:2 sticker_id:140 wear:0.20677716} stickers:{slot:3 sticker_id:260} inventory:48 origin:8 keychains:{slot:0 sticker_id:1 offset_x:24.17532 offset_y:0.31357658 offset_z:2.541658 pattern:86000}
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44417931234),
		Defindex:   uint32Pointer(7),
		Paintindex: uint32Pointer(302),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(4),
		Paintwear:  uint32Pointer(1028819163),
		Paintseed:  uint32Pointer(387),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(666),
				Wear:      float32Pointer(0.23460516),
			},
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(1),
				StickerId: uint32Pointer(1338),
			},
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(2),
				StickerId: uint32Pointer(140),
				Wear:      float32Pointer(0.20677716),
			},
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(3),
				StickerId: uint32Pointer(260),
			},
		},
		Inventory: uint32Pointer(48),
		Origin:    uint32Pointer(8),
		Keychains: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(1),
				OffsetX:   float32Pointer(24.17532),
				OffsetY:   float32Pointer(0.31357658),
				OffsetZ:   float32Pointer(2.541658),
				Pattern:   uint32Pointer(86000),
			},
		},
	}
	tests["Skin Stickers Keychain"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.05140004679560661,
			MinFloat:     "0",
			MaxFloat:     "0.9",
			ItemName:     "Vulcan",
			QualityName:  "Unique",
			WeaponType:   "AK-47",
			RarityName:   "Covert",
			WearName:     "Factory New",
			FullItemName: "AK-47 | Vulcan (Factory New)",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "eslcologne2015_team_clg_foil",
					Material: "Foil",
					Name:     "Sticker | Counter Logic Gaming (Foil) | Cologne 2015",
				},
				types.Modification{
					Proto:    input.Stickers[1],
					CodeName: "cologne2016_team_liq_holo",
					Material: "Holo",
					Name:     "Sticker | Team Liquid (Holo) | Cologne 2016",
				},
				types.Modification{
					Proto:    input.Stickers[2],
					CodeName: "cologne2014_epsilonesports",
					Material: "Paper",
					Name:     "Sticker | Epsilon eSports | Cologne 2014",
				},
				types.Modification{
					Proto:    input.Stickers[3],
					CodeName: "drugwarveteran",
					Material: "Paper",
					Name:     "Sticker | Drug War Veteran",
				},
			},
			Keychains: []types.Modification{
				types.Modification{
					Proto:    input.Keychains[0],
					CodeName: "kc_missinglink_ava",
					Name:     "Charm | Lil' Ava",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/%E2%98%85%20Bloodhound%20Gloves%20%7C%20Bronzed%20%28Factory%20New%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M6664853361691228811A39592993407D14448091259802652413
	// itemid:39592993407 defindex:5027 paintindex:10008 rarity:6 quality:3 paintwear:1031772734 paintseed:896 inventory:3221225482 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(39592993407),
		Defindex:   uint32Pointer(5027),
		Paintindex: uint32Pointer(10008),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(3),
		Paintwear:  uint32Pointer(1031772734),
		Paintseed:  uint32Pointer(896),
		Inventory:  uint32Pointer(3221225482),
		Origin:     uint32Pointer(8),
	}
	tests["Gloves"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.06240295618772507,
			MinFloat:     "0.06",
			MaxFloat:     "0.8",
			ItemName:     "Bronzed",
			QualityName:  "★",
			WeaponType:   "Bloodhound Gloves",
			RarityName:   "Covert",
			WearName:     "Factory New",
			FullItemName: "★ Bloodhound Gloves | Bronzed (Factory New)",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Lotus%20%28Glitter%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M643567392820418477A42466192195D9975096757687244901
	// itemid:42466192195 defindex:1209 paintindex:0 rarity:4 quality:4 stickers:{slot:0 sticker_id:7251} inventory:85 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(42466192195),
		Defindex:   uint32Pointer(1209),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(4),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(7251),
			},
		},
		Inventory: uint32Pointer(85),
		Origin:    uint32Pointer(8),
	}
	tests["Sticker Glitter"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Lotus (Glitter)",
			QualityName:  "Unique",
			WeaponType:   "Sticker",
			RarityName:   "Remarkable",
			WearName:     "",
			FullItemName: "Sticker | Lotus (Glitter)",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "glitter_lotus",
					Material: "Glitter",
					Name:     "Sticker | Lotus (Glitter)",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Patch%20%7C%20Elder%20God
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M644693292653772584A23430518132D3314979815666773509
	// itemid:23430518132 defindex:4609 paintindex:0 rarity:5 quality:4 stickers:{slot:0 sticker_id:4948} inventory:268 origin:23
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(23430518132),
		Defindex:   uint32Pointer(4609),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(5),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(4948),
			},
		},
		Inventory: uint32Pointer(268),
		Origin:    uint32Pointer(23),
	}
	tests["Patch"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Elder God",
			QualityName:  "Unique",
			WeaponType:   "Patch",
			RarityName:   "Exotic",
			WearName:     "",
			FullItemName: "Patch | Elder God",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "patch_op11_cthulhu",
					Name:     "Patch | Elder God",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/StatTrak%E2%84%A2%20AWP%20%7C%20Printstream%20%28Factory%20New%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M649195897111665513A44158728090D11837336287321839237
	// itemid:44158728090 defindex:9 paintindex:1206 rarity:6 quality:9 paintwear:1027111819 paintseed:93 killeaterscoretype:0 killeatervalue:54 inventory:253 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:             uint64Pointer(44158728090),
		Defindex:           uint32Pointer(9),
		Paintindex:         uint32Pointer(1206),
		Rarity:             uint32Pointer(6),
		Quality:            uint32Pointer(9),
		Paintwear:          uint32Pointer(1027111819),
		Paintseed:          uint32Pointer(93),
		Killeaterscoretype: uint32Pointer(0),
		Killeatervalue:     uint32Pointer(54),
		Inventory:          uint32Pointer(253),
		Origin:             uint32Pointer(8),
	}
	tests["Skin StatTrak"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.045039694756269455,
			MinFloat:     "0",
			MaxFloat:     "1",
			ItemName:     "Printstream",
			QualityName:  "StatTrak™",
			WeaponType:   "AWP",
			RarityName:   "Covert",
			WearName:     "Factory New",
			FullItemName: "StatTrak™ AWP | Printstream (Factory New)",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/%E2%98%85%20Bayonet%20%7C%20Gamma%20Doppler%20%28Factory%20New%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M653700491952159799A44851210636D9387414943281127756
	// itemid:44851210636 defindex:500 paintindex:571 rarity:6 quality:3 paintwear:1023716016 paintseed:86 inventory:15 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44851210636),
		Defindex:   uint32Pointer(500),
		Paintindex: uint32Pointer(571),
		Rarity:     uint32Pointer(6),
		Quality:    uint32Pointer(3),
		Paintwear:  uint32Pointer(1023716016),
		Paintseed:  uint32Pointer(86),
		Inventory:  uint32Pointer(15),
		Origin:     uint32Pointer(8),
	}
	tests["Knife"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0.03238934278488159,
			MinFloat:     "0",
			MaxFloat:     "0.08",
			ItemName:     "Gamma Doppler",
			QualityName:  "★",
			WeaponType:   "Bayonet",
			RarityName:   "Covert",
			WearName:     "Factory New",
			FullItemName: "★ Bayonet | Gamma Doppler (Factory New)",
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sticker%20%7C%20Lit%20%28Holo%29
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M650322792241404634A44809037288D922326464949352065
	// itemid:44809037288 defindex:1209 paintindex:0 rarity:4 quality:4 stickers:{slot:0 sticker_id:7248} inventory:159 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44809037288),
		Defindex:   uint32Pointer(1209),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(4),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(7248),
			},
		},
		Inventory: uint32Pointer(159),
		Origin:    uint32Pointer(8),
	}
	tests["Sticker Holo"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Lit (Holo)",
			QualityName:  "Unique",
			WeaponType:   "Sticker",
			RarityName:   "Remarkable",
			WearName:     "",
			FullItemName: "Sticker | Lit (Holo)",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "holo_lit",
					Material: "Holo",
					Name:     "Sticker | Lit (Holo)",
				},
			},
		},
		expectedError: nil,
	}

	// https://steamcommunity.com/market/listings/730/Sealed%20Graffiti%20%7C%20Drug%20War%20Veteran
	// steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20M643567392817240427A44225554008D11566448795219559360
	// itemid:44225554008 defindex:1348 paintindex:0 rarity:5 quality:4 stickers:{slot:0 sticker_id:1655} inventory:6 origin:8
	input = &protobuf.CEconItemPreviewDataBlock{
		Itemid:     uint64Pointer(44225554008),
		Defindex:   uint32Pointer(1348),
		Paintindex: uint32Pointer(0),
		Rarity:     uint32Pointer(5),
		Quality:    uint32Pointer(4),
		Stickers: []*protobuf.CEconItemPreviewDataBlock_Sticker{
			&protobuf.CEconItemPreviewDataBlock_Sticker{
				Slot:      uint32Pointer(0),
				StickerId: uint32Pointer(1655),
			},
		},
		Inventory: uint32Pointer(6),
		Origin:    uint32Pointer(8),
	}
	tests["Graffiti"] = protoTestCase{
		input: input,
		expectedItem: &types.Item{
			Proto:        input,
			FloatValue:   0,
			ItemName:     "Drug War Veteran",
			QualityName:  "Unique",
			WeaponType:   "Sealed Graffiti",
			RarityName:   "Exotic",
			WearName:     "",
			FullItemName: "Sealed Graffiti | Drug War Veteran",
			Stickers: []types.Modification{
				types.Modification{
					Proto:    input.Stickers[0],
					CodeName: "spray_drugwarveteran",
					Name:     "Sealed Graffiti | Drug War Veteran",
				},
			},
		},
		expectedError: nil,
	}

	return tests
}
