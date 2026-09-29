CREATE TABLE IF NOT EXISTS words (
  id INTEGER PRIMARY KEY,
  spanish TEXT NOT NULL,
  english TEXT NOT NULL,
  ukrainian TEXT NOT NULL,
  gender_uk TEXT,
  UNIQUE(english)
);

CREATE TABLE IF NOT EXISTS profiles (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  enable_speech BOOLEAN NOT NULL DEFAULT TRUE,
  default_num_questions INT NOT NULL DEFAULT 10,
  default_num_answers INT NOT NULL DEFAULT 4,
  lang1 TEXT NOT NULL DEFAULT 'en',
  lang2 TEXT NOT NULL DEFAULT 'es',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL,
  UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS rankings (
  profile_id INTEGER NOT NULL,
  word_id INTEGER NOT NULL,
  en_to_es REAL NOT NULL DEFAULT 0,
  es_to_en REAL NOT NULL DEFAULT 0,
  en_to_uk REAL NOT NULL DEFAULT 0,
  uk_to_en REAL NOT NULL DEFAULT 0,
  es_to_uk REAL NOT NULL DEFAULT 0,
  uk_to_es REAL NOT NULL DEFAULT 0,
  FOREIGN KEY(profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
  FOREIGN KEY(word_id) REFERENCES words(id) ON DELETE CASCADE,
  UNIQUE(profile_id, word_id)
);
