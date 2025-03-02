package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	gt "github.com/volodymyrzuyev/goCsInspect/cmd/globalTypes"
	custSql "github.com/volodymyrzuyev/goCsInspect/cmd/storage/db/sqlc"
)

type db struct {
	q *custSql.Queries
	d *sql.DB
}

func (d db) GetItem(itemId int64) (gt.Item, error) {
	item, err := d.q.GetItem(context.Background(), itemId)
	if err != nil {
		fmt.Println("err here")
		return gt.Item{}, fmt.Errorf("De?")
	}

	mods, err := d.q.GetModifiers(context.Background(), itemId)
	if err != nil {
		fmt.Println("err there")
		return gt.Item{}, fmt.Errorf("SHo?")
	}

	stickers, chains := getStickersAndChains(mods)

	itemStruct := gt.Item{
		AccountID:          int(item.Accountid.Int64),
		ItemID:             int(item.Itemid),
		DefIndex:           int(item.Defindex.Int64),
		PaintIndex:         int(item.Paintindex.Int64),
		Rarity:             int(item.Rarity.Int64),
		Quality:            int(item.Quality.Int64),
		Paintwear:          int(item.Paintwear.Int64),
		Paintseed:          int(item.Paintseed.Int64),
		Killeaterscoretype: int(item.Killeaterscoretype.Int64),
		Killeatervalue:     int(item.Killeatervalue.Int64),
		Customname:         item.Customname.String,
		Stickers:           stickers,
		Inventory:          int(item.Inventory.Int64),
		Origin:             int(item.Origin.Int64),
		Questid:            int(item.Questid.Int64),
		Dropreason:         int(item.Dropreason.Int64),
		Musicindex:         int(item.Musicindex.Int64),
		Entindex:           int(item.Entindex.Int64),
		Petindex:           int(item.Petindex.Int64),
		Keychains:          chains,
		ParamD:             int(item.Paramd.Int64),
		ParamM:             int(item.Paramm.Int64),
		ParamS:             int(item.Params.Int64),
		FloatValue:         item.Floatvalue.Float64,
		MaxFloat:           item.Maxfloat.Float64,
		MinFloat:           item.Minfloat.Float64,
		WeaponType:         item.Weapontype.String,
		ItemName:           item.Itemname.String,
		RarityName:         item.Rarityname.String,
		QualityName:        item.Qualityname.String,
		OriginName:         item.Originname.String,
		WearName:           item.Wearname.String,
		MarketHashName:     item.Markethashname.String,
		LastModified:       time.Unix(item.Lastupdated.Int64, 0),
	}

	return itemStruct, nil
}

func (d db) InsertItem(item gt.Item) error {
	itemArgs := custSql.InsertItemParams{
		Accountid:          getNullInt64(int64(item.AccountID)),
		Itemid:             int64(item.ItemID),
		Defindex:           getNullInt64(int64(item.DefIndex)),
		Paintindex:         getNullInt64(int64(item.PaintIndex)),
		Rarity:             getNullInt64(int64(item.Rarity)),
		Quality:            getNullInt64(int64(item.Quality)),
		Paintwear:          getNullInt64(int64(item.Paintwear)),
		Paintseed:          getNullInt64(int64(item.Paintseed)),
		Killeaterscoretype: getNullInt64(int64(item.Killeaterscoretype)),
		Killeatervalue:     getNullInt64(int64(item.Killeatervalue)),
		Customname:         getNullString(item.Customname),
		Inventory:          getNullInt64(int64(item.Inventory)),
		Origin:             getNullInt64(int64(item.Origin)),
		Questid:            getNullInt64(int64(item.Questid)),
		Dropreason:         getNullInt64(int64(item.Dropreason)),
		Musicindex:         getNullInt64(int64(item.Musicindex)),
		Entindex:           getNullInt64(int64(item.Entindex)),
		Petindex:           getNullInt64(int64(item.Petindex)),
		Paramd:             getNullInt64(int64(item.ParamD)),
		Paramm:             getNullInt64(int64(item.ParamM)),
		Params:             getNullInt64(int64(item.ParamS)),
		Floatvalue:         getNullFloat64(item.FloatValue),
		Maxfloat:           getNullFloat64(item.MaxFloat),
		Minfloat:           getNullFloat64(item.MinFloat),
		Weapontype:         getNullString(item.WeaponType),
		Itemname:           getNullString(item.ItemName),
		Rarityname:         getNullString(item.RarityName),
		Qualityname:        getNullString(item.QualityName),
		Originname:         getNullString(item.OriginName),
		Wearname:           getNullString(item.WearName),
		Markethashname:     getNullString(item.MarketHashName),
		Lastupdated:        getNullInt64(item.LastModified.Unix()),
	}

	tx, err := d.d.Begin()
	if err != nil {
		return DbError
	}
	defer tx.Rollback()

	qtx := d.q.WithTx(tx)

	err = qtx.InsertItem(context.Background(), itemArgs)
	if err != nil {
		fmt.Println(err)
		return DbError
	}

	for _, m := range item.Stickers {
		err = qtx.InsertModifier(context.Background(), modToInsertModArgs(stickersEnum, m, int64(item.ItemID)))
		if err != nil {
			fmt.Println("tam")
			return DbError
		}
	}

	for _, m := range item.Keychains {
		err = qtx.InsertModifier(context.Background(), modToInsertModArgs(chainEnum, m, int64(item.ItemID)))
		if err != nil {
			fmt.Println("tut")
			return DbError
		}
	}

	return tx.Commit()
}

func migrateUp(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	migrationsDir := filepath.Join(currentDir, "cmd", "storage", "db", "sqlc", "sql")

	migrationsPath := fmt.Sprintf("file://%s", migrationsDir)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath, // Path to migration files
		"sqlite3",
		driver, // Driver instance
	)
	if err != nil {
		return err
	}

	// Run migrations up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func InitDB(dbPath string) (db, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return db{}, err
	}
	q := custSql.New(conn)

	err = migrateUp(conn)
	if err != nil {
		fmt.Println(err)
		return db{}, err
	}

	return db{q, conn}, nil
}

const (
	stickersEnum = "sticker"
	chainEnum    = "chain"
)

func modToInsertModArgs(enum string, mod gt.Modifier, itemID int64) custSql.InsertModifierParams {
	args := custSql.InsertModifierParams{
		Itemid:         itemID,
		Modifierid:     getNullInt64(int64(mod.ModifierId)),
		Modifiertype:   enum,
		Slot:           getNullInt64(int64(mod.Slot)),
		Wear:           getNullFloat64(mod.Wear),
		Scale:          getNullFloat64(mod.Scale),
		Rotation:       getNullFloat64(mod.Rotation),
		Tintid:         getNullInt64(int64(mod.TintId)),
		Offsetx:        getNullFloat64(mod.OffsetX),
		Offsety:        getNullFloat64(mod.OffsetY),
		Offsetz:        getNullFloat64(mod.OffsetZ),
		Pattern:        getNullInt64(int64(mod.Pattern)),
		Markethashname: getNullString(mod.MarketHashName),
	}

	return args
}

func getStickersAndChains(mods []custSql.Modifier) (stickers []gt.Modifier, chains []gt.Modifier) {
	for _, m := range mods {

		newItem := gt.Modifier{
			Slot:           int(m.Slot.Int64),
			ModifierId:     int(m.Modifierid.Int64),
			Wear:           m.Wear.Float64,
			Scale:          m.Scale.Float64,
			Rotation:       m.Rotation.Float64,
			TintId:         int(m.Tintid.Int64),
			OffsetX:        m.Offsetx.Float64,
			OffsetY:        m.Offsety.Float64,
			OffsetZ:        m.Offsetz.Float64,
			Pattern:        int(m.Pattern.Int64),
			MarketHashName: m.Markethashname.String,
		}
		switch m.Modifiertype {
		case stickersEnum:
			stickers = append(stickers, newItem)
		case chainEnum:
			chains = append(chains, newItem)
		}
	}

	return
}

func getNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

func getNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func getNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: true}
}
