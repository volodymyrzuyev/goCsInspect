CREATE TABLE items (
  M INTEGER,
  D INTEGER,
  S INTEGER,
  Accountid INTEGER,
  Itemid INTEGER,
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

CREATE TABLE stickers (
  Modid INTEGER PRIMARY KEY,
  Slot INTEGER NOT NULL,
  StickerId INTEGER NOT NULL,
  Wear REAL NOT NULL,
  Scale REAL NOT NULL,
  Rotation REAL NOT NULL,
  TintId INTEGER NOT NULL,
  OffsetX REAL NOT NULL,
  OffsetY REAL NOT NULL,
  OffsetZ REAL NOT NULL,
  Pattern INTEGER NOT NULL,
  HighlightReel INTEGER NOT NULL
);

CREATE TABLE modsStickers (
  Modid INTEGER,
  M INTEGER,
  D INTEGER,
  S INTEGER,
  Itemid INTEGER,
  PRIMARY KEY (Modid, M, D, S, Itemid),
  FOREIGN KEY (Modid) REFERENCES stickers (Modid),
  FOREIGN KEY (M) REFERENCES items (M),
  FOREIGN KEY (D) REFERENCES items (D),
  FOREIGN KEY (S) REFERENCES items (S),
  FOREIGN KEY (Itemid) REFERENCES items (Itemid)
);

CREATE TABLE modsChains (
  Modid INTEGER,
  M INTEGER,
  D INTEGER,
  S INTEGER,
  Itemid INTEGER,
  PRIMARY KEY (Modid, M, D, S, Itemid),
  FOREIGN KEY (Modid) REFERENCES stickers (Modid),
  FOREIGN KEY (M) REFERENCES items (M),
  FOREIGN KEY (D) REFERENCES items (D),
  FOREIGN KEY (S) REFERENCES items (S),
  FOREIGN KEY (Itemid) REFERENCES items (Itemid)
);
