-- name: GetItem :one
SELECT
*
FROM items
WHERE ItemID = ?;

-- name: GetModifiers :many
SELECT
*
FROM Modifiers
WHERE ItemID = ?;

-- name: InsertItem :exec
INSERT INTO items (
    AccountID,
    ItemID,
    DefIndex,
    PaintIndex,
    Rarity,
    Quality,
    Paintwear,
    Paintseed,
    Killeaterscoretype,
    Killeatervalue,
    Customname,
    Inventory,
    Origin,
    Questid,
    Dropreason,
    Musicindex,
    Entindex,
    Petindex,
    ParamD,
    ParamM,
    ParamS,
    FloatValue,
    MaxFloat,
    MinFloat,
    WeaponType,
    ItemName,
    RarityName,
    QualityName,
    OriginName,
    WearName,
    MarketHashName,
    LastUpdated
) VALUES (
    ?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14, ?15, ?16,
    ?17, ?18, ?19, ?20, ?21, ?22, ?23, ?24, ?25, ?26, ?27, ?28, ?29, ?30, ?31,
    ?32
);

-- name: InsertModifier :exec
INSERT INTO Modifiers(
	ItemID,
	ModifierID,
	ModifierType,
	Slot,
	Wear,
	Scale,
	Rotation,
	TintId,
	OffsetX,
	OffsetY,
	OffsetZ,
	Pattern,
	MarketHashName
)
VALUES ( 
	?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13
);


-- name: DeleteItem :exec
DELETE FROM items
WHERE ItemID = ?;

-- name: DeleteModifiers :exec
DELETE FROM Modifiers
WHERE ItemID = ?;
