package store

import (
	"context"
	"fmt"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
)

// DB access for Words

// Word is a single Spanish/English/Ukraine noun pair.
type Word struct {
	ID          int64   // 1-based ind
	Spanish     string  // Spanish Word
	English     string  // English Word
	Ukrainian     string  // Ukraine Word
	Gender      string  // TODO use/implement/delete?
	EnToEsScore float64 // how well the user knows english word from spanish
	EnToUkScore float64 // how well the user knows english word from ukrainian
	EsToEnScore float64 // how well the user knows spanish word from english
	EsToUkScore float64 // how well the user knows spanish word from ukrainian
	UkToEsScore float64 // how well the user knows ukrainian word from spanish
	UkToEnScore float64 // how well the user knows ukrainian word from english
}


// GetAllWords returns all words associated with a profile.
func (s *Store) GetAllWords(ctx context.Context, profileName string) ([]Word, error) {
	var words []Word

	selQry := `SELECT w.id, spanish, english, ukrainian, gender_uk, en_to_es, en_to_uk, es_to_en, es_to_uk, uk_to_es, uk_to_en  FROM words w
			INNER JOIN rankings r ON w.id = r.word_id
			INNER JOIN profiles p ON p.id = r.profile_id
			WHERE p.name = ?`
	rows, err := s.db.QueryContext(ctx, selQry, profileName)
	if err != nil {
		return nil, fmt.Errorf("store: get words: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Spanish, &w.English, &w.Ukrainian, &w.Gender, &w.EnToEsScore, &w.EnToUkScore, &w.EsToEnScore, &w.EsToUkScore, &w.UkToEsScore, &w.UkToEnScore); err != nil {
			return nil, fmt.Errorf("store: get words: %w", err)
		}
		words = append(words, w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: get words: rows: %w", err)
	}
	
	log.DebugSQL("Successfully queried all words for profile %s, total count: %d", profileName, len(words))
	return words, nil

}