CREATE TABLE items (
  -- uint64 can not be stored as an integer due to bit width, so will store as 
  -- text so numbers won't loose precision
  M TEXT NOT NULL,
  D TEXT NOT NULL,
  S TEXT NOT NULL,
  Accountid INTEGER,
  Itemid TEXT NOT NULL,
  Defindex INTEGER,
  Paintindex INTEGER,
  Rarity INTEGER,
  Quality INTEGER,
  Paintwear INTEGER,
  Paintseed INTEGER,
  Killeaterscoretype INTEGER,
  Killeatervalue INTEGER,
  Customname TEXT,
  Inventory INTEGER,
  Origin INTEGER,
  Questid INTEGER,
  Dropreason INTEGER,
  Musicindex INTEGER,
  Entindex INTEGER,
  Petindex INTEGER,
  PRIMARY KEY (M, A, D, S)
);

CREATE TABLE mods (
  Modid INTEGER PRIMARY KEY,
  Slot INTEGER,
  StickerId INTEGER,
  Wear REAL,
  Scale REAL,
  Rotation REAL,
  TintId INTEGER,
  OffsetX REAL,
  OffsetY REAL,
  OffsetZ REAL,
  Pattern INTEGER,
  HighlightReel INTEGER
);

CREATE TABLE modsStickers (
  Modid INTEGER NOT NULL,
  M TEXT NOT NULL,
  D TEXT NOT NULL,
  S TEXT NOT NULL,
  Itemid TEXT NOT NULL,
  PRIMARY KEY (Modid, M, D, S, Itemid),
  FOREIGN KEY (Modid) REFERENCES stickers (Modid),
  FOREIGN KEY (M) REFERENCES items (M),
  FOREIGN KEY (D) REFERENCES items (D),
  FOREIGN KEY (S) REFERENCES items (S),
  FOREIGN KEY (Itemid) REFERENCES items (Itemid)
);

CREATE TABLE modsChains (
  Modid INTEGER NOT NULL,
  M TEXT NOT NULL,
  D TEXT NOT NULL,
  S TEXT NOT NULL,
  Itemid TEXT NOT NULL,
  PRIMARY KEY (Modid, M, D, S, Itemid),
  FOREIGN KEY (Modid) REFERENCES stickers (Modid),
  FOREIGN KEY (M) REFERENCES items (M),
  FOREIGN KEY (D) REFERENCES items (D),
  FOREIGN KEY (S) REFERENCES items (S),
  FOREIGN KEY (Itemid) REFERENCES items (Itemid)
);
