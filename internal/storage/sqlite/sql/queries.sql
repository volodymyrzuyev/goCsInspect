-- name: InsertItem :exec
INSERT INTO items (
    M,
    D,
    S,
    Accountid,
    Itemid,
    Defindex,
    Paintindex,
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
    Petindex
)
VALUES (
    ?1, 
    ?2, 
    ?3, 
    ?4, 
    ?5, 
    ?6, 
    ?7, 
    ?8, 
    ?9, 
    ?10, 
    ?11, 
    ?12, 
    ?13, 
    ?14, 
    ?15, 
    ?16,
    ?17, 
    ?18, 
    ?19, 
    ?20, 
    ?21
);

-- name: InsertMod :one
INSERT INTO mods (
    Slot,
    StickerId,
    Wear,
    Scale,
    Rotation,
    TintId,
    OffsetX,
    OffsetY,
    OffsetZ,
    Pattern,
    HighlightReel
)
VALUES (
    ?1, 
    ?2, 
    ?3, 
    ?4, 
    ?5, 
    ?6, 
    ?7, 
    ?8, 
    ?9, 
    ?10, 
    ?11
) 
RETURNING Modid;

-- name: InsertModSticker :exec
INSERT INTO modsStickers (
  Modid,
  M,
  D,
  S,
  Itemid
)
VALUES (
    ?1, 
    ?2, 
    ?3, 
    ?4, 
    ?5
);

-- name: InsertModChain :exec
INSERT INTO modsChains (
  Modid,
  M,
  D,
  S,
  Itemid
)
VALUES (
    ?1, 
    ?2, 
    ?3, 
    ?4, 
    ?5
);

-- name: GetItem :one
SELECT * FROM items WHERE
    M = ?1
    AND D = ?2
    AND S = ?3
    AND Itemid = ?4;

-- name: GetStickers :many
SELECT mods.* FROM mods, modsStickers s WHERE 
    s.M = ?1
    AND s.D = ?2
    AND s.S = ?3
    AND s.Itemid = ?4
    AND mods.Modid = s.Modid
    ORDER BY mods.Modid;

-- name: GetChains :many
SELECT mods.* FROM mods, modsChains s WHERE 
    s.M = ?1
    AND s.D = ?2
    AND s.S = ?3
    AND s.Itemid = ?4
    AND mods.Modid = s.Modid
    ORDER BY mods.Modid;
