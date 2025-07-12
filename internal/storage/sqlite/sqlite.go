package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	proto "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/volodymyrzuyev/goCsInspect/internal/storage/sqlite/sql/sqlc"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common"
	"github.com/volodymyrzuyev/goCsInspect/pkg/common/errors"
	"github.com/volodymyrzuyev/goCsInspect/pkg/storage"
	t "github.com/volodymyrzuyev/goCsInspect/pkg/types"
)

type Sqlite struct {
	dbPath string
	db     *sql.DB
	q      *sqlc.Queries
}

func NewSQLiteStore(dbPath string) storage.Storage {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	migrateDB(dbPath)

	q := sqlc.New(db)

	return &Sqlite{
		dbPath: dbPath,
		db:     db,
		q:      q,
	}
}

func getDBProto(params t.InspectParameters, proto *proto.CEconItemPreviewDataBlock) sqlc.InsertItemParams {
	return sqlc.InsertItemParams{
		M:                  fmt.Sprintf("%v", params.M),
		D:                  fmt.Sprintf("%v", params.D),
		S:                  fmt.Sprintf("%v", params.S),
		Accountid:          common.NullInt64Uint32Ptr(proto.Accountid),
		Itemid:             fmt.Sprintf("%v", params.A),
		Defindex:           common.NullInt64Uint32Ptr(proto.Defindex),
		Paintindex:         common.NullInt64Uint32Ptr(proto.Paintindex),
		Rarity:             common.NullInt64Uint32Ptr(proto.Rarity),
		Quality:            common.NullInt64Uint32Ptr(proto.Quality),
		Paintwear:          common.NullInt64Uint32Ptr(proto.Paintwear),
		Paintseed:          common.NullInt64Uint32Ptr(proto.Paintseed),
		Killeaterscoretype: common.NullInt64Uint32Ptr(proto.Killeaterscoretype),
		Killeatervalue:     common.NullInt64Uint32Ptr(proto.Killeatervalue),
		Customname:         common.NullStringStringPtr(proto.Customname),
		Inventory:          common.NullInt64Uint32Ptr(proto.Inventory),
		Origin:             common.NullInt64Uint32Ptr(proto.Origin),
		Questid:            common.NullInt64Uint32Ptr(proto.Questid),
		Dropreason:         common.NullInt64Uint32Ptr(proto.Dropreason),
		Musicindex:         common.NullInt64Uint32Ptr(proto.Musicindex),
		Entindex:           common.NullInt64Int32Ptr(proto.Entindex),
		Petindex:           common.NullInt64Uint32Ptr(proto.Petindex),
	}

}

func getDBModProto(protos []*proto.CEconItemPreviewDataBlock_Sticker) []sqlc.InsertModParams {
	mods := make([]sqlc.InsertModParams, len(protos))

	for i, proto := range protos {
		mods[i] = sqlc.InsertModParams{
			Slot:          common.NullInt64Uint32Ptr(proto.Slot),
			Stickerid:     common.NullInt64Uint32Ptr(proto.StickerId),
			Wear:          common.NullFloat64Float32Ptr(proto.Wear),
			Scale:         common.NullFloat64Float32Ptr(proto.Scale),
			Rotation:      common.NullFloat64Float32Ptr(proto.Rotation),
			Tintid:        common.NullInt64Uint32Ptr(proto.TintId),
			Offsetx:       common.NullFloat64Float32Ptr(proto.OffsetX),
			Offsety:       common.NullFloat64Float32Ptr(proto.OffsetY),
			Offsetz:       common.NullFloat64Float32Ptr(proto.OffsetZ),
			Pattern:       common.NullInt64Uint32Ptr(proto.Pattern),
			Highlightreel: common.NullInt64Uint32Ptr(proto.Pattern),
		}
	}

	return mods
}

func (s *Sqlite) StoreItem(ctx context.Context, params t.InspectParameters, proto *proto.CEconItemPreviewDataBlock) error {
	dbItem := getDBProto(params, proto)
	stickers := getDBModProto(proto.GetStickers())
	chains := getDBModProto(proto.GetKeychains())
	modJoinTableStorer := sqlc.InsertModStickerParams{
		M:      dbItem.M,
		D:      dbItem.D,
		S:      dbItem.S,
		Itemid: dbItem.Itemid,
	}

	dbTX, err := s.db.Begin()
	if err != nil {
		slog.Error("Error when beggining db tx", "error", err)
		return errors.ErrDB
	}
	defer dbTX.Rollback()

	tx := s.q.WithTx(dbTX)

	if err = tx.InsertItem(ctx, dbItem); err != nil {
		slog.Error("Error when adding proto to db", "error", err)
		return errors.ErrInsertItem
	}

	for _, sticker := range stickers {
		modId, err := tx.InsertMod(ctx, sticker)
		if err != nil {
			slog.Error("Error when inserting sticker", "error", err)
			return errors.ErrInsertSticker
		}
		modJoinTableStorer.Modid = modId

		err = tx.InsertModSticker(ctx, modJoinTableStorer)
		if err != nil {
			slog.Error("Error when inserting sticker", "error", err)
			return errors.ErrInsertSticker
		}
	}

	for _, chain := range chains {
		modId, err := tx.InsertMod(ctx, chain)
		if err != nil {
			slog.Error("Error when inserting keychain", "error", err)
			return errors.ErrInsertKeychain
		}
		modJoinTableStorer.Modid = modId

		err = tx.InsertModChain(ctx, sqlc.InsertModChainParams(modJoinTableStorer))
		if err != nil {
			slog.Error("Error when inserting keychain", "error", err)
			return errors.ErrInsertKeychain
		}
	}

	return dbTX.Commit()
}

func converDBItemToProto(i sqlc.Item) (*proto.CEconItemPreviewDataBlock, error) {
	itemID, err := strconv.ParseUint(i.Itemid, 10, 64)
	if err != nil {
		return nil, err
	}

	item := &proto.CEconItemPreviewDataBlock{
		Accountid:          common.Uint32PtrNullInt64(i.Accountid),
		Itemid:             common.Uint64Pointer(itemID),
		Defindex:           common.Uint32PtrNullInt64(i.Defindex),
		Paintindex:         common.Uint32PtrNullInt64(i.Paintindex),
		Rarity:             common.Uint32PtrNullInt64(i.Rarity),
		Quality:            common.Uint32PtrNullInt64(i.Quality),
		Paintwear:          common.Uint32PtrNullInt64(i.Paintwear),
		Paintseed:          common.Uint32PtrNullInt64(i.Paintseed),
		Killeaterscoretype: common.Uint32PtrNullInt64(i.Killeaterscoretype),
		Killeatervalue:     common.Uint32PtrNullInt64(i.Killeatervalue),
		Customname:         common.StringPtrNullString(i.Customname),
		Inventory:          common.Uint32PtrNullInt64(i.Inventory),
		Origin:             common.Uint32PtrNullInt64(i.Origin),
		Questid:            common.Uint32PtrNullInt64(i.Questid),
		Dropreason:         common.Uint32PtrNullInt64(i.Dropreason),
		Musicindex:         common.Uint32PtrNullInt64(i.Musicindex),
		Entindex:           common.Int32PtNullInt64(i.Entindex),
		Petindex:           common.Uint32PtrNullInt64(i.Petindex),
	}

	return item, nil
}

func parseDbModToProto(mods []sqlc.Mod) []*proto.CEconItemPreviewDataBlock_Sticker {
	protos := make([]*proto.CEconItemPreviewDataBlock_Sticker, len(mods))

	for i, mod := range mods {
		protos[i] = &proto.CEconItemPreviewDataBlock_Sticker{
			Slot:          common.Uint32PtrNullInt64(mod.Slot),
			StickerId:     common.Uint32PtrNullInt64(mod.Stickerid),
			Wear:          common.Float32PtrNullFloat(mod.Wear),
			Scale:         common.Float32PtrNullFloat(mod.Scale),
			Rotation:      common.Float32PtrNullFloat(mod.Rotation),
			TintId:        common.Uint32PtrNullInt64(mod.Tintid),
			OffsetX:       common.Float32PtrNullFloat(mod.Offsetx),
			OffsetY:       common.Float32PtrNullFloat(mod.Offsety),
			OffsetZ:       common.Float32PtrNullFloat(mod.Offsetz),
			Pattern:       common.Uint32PtrNullInt64(mod.Pattern),
			HighlightReel: common.Uint32PtrNullInt64(mod.Highlightreel),
		}
	}

	return protos
}

func assembleItem(dbItem sqlc.Item, chains, stickers []sqlc.Mod) (*proto.CEconItemPreviewDataBlock, error) {
	item, err := converDBItemToProto(dbItem)
	if err != nil {
		slog.Error("Got error when getting dbItem", "error", err)
		return nil, errors.ErrFetchItem
	}

	item.Stickers = parseDbModToProto(stickers)
	item.Keychains = parseDbModToProto(chains)

	return item, nil
}

func (s *Sqlite) GetItem(ctx context.Context, params t.InspectParameters) (*proto.CEconItemPreviewDataBlock, error) {
	itemDbParams := sqlc.GetItemParams{
		M:      fmt.Sprintf("%v", params.M),
		D:      fmt.Sprintf("%v", params.D),
		S:      fmt.Sprintf("%v", params.S),
		Itemid: fmt.Sprintf("%v", params.A),
	}

	dbItem, err := s.q.GetItem(ctx, itemDbParams)
	if err != nil {
		slog.Error("Got error when getting dbItem", "error", err)
		return nil, errors.ErrFetchItem
	}

	stickers, err := s.q.GetStickers(ctx, sqlc.GetStickersParams(itemDbParams))
	if err != nil {
		slog.Error("Got error when getting stickers", "error", err)
		return nil, errors.ErrFetchSticker
	}

	chains, err := s.q.GetChains(ctx, sqlc.GetChainsParams(itemDbParams))
	if err != nil {
		slog.Error("Got error when getting keychain", "error", err)
		return nil, errors.ErrFetchKeychain
	}

	return assembleItem(dbItem, stickers, chains)
}
