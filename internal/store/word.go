package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
)

// DB access for Words

const (
	// QuizMode keeps track of the quiz mode setting
	QuizModeAny string = "any"
	QuizModeWellKnown string = "well_known"
	QuizModeLeastKnown string = "least_known"
)

// Word is a single Spanish/English/Ukraine noun pair.
type Word struct {
	ID          int64   // 1-based ind
	Spanish     string  // Spanish Word
	English     string  // English Word
	Ukrainian   string  // Ukraine Word
	Gender      string  // TODO use/implement/delete?
	EnToEsScore float64 // how well the user knows english word from spanish
	EnToUkScore float64 // how well the user knows english word from ukrainian
	EsToEnScore float64 // how well the user knows spanish word from english
	EsToUkScore float64 // how well the user knows spanish word from ukrainian
	UkToEsScore float64 // how well the user knows ukrainian word from spanish
	UkToEnScore float64 // how well the user knows ukrainian word from english
}

// GetTextFromLang returns the word text given a lang ('es', 'en', 'uk')
func (w Word) GetTextFromLang(lang string) string {
	switch lang {
	case "es":
		return w.Spanish
	case "en":
		return w.English
	case "uk":
		return w.Ukrainian
	}
	return "Unknown word" // Shouldn't get here
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

// GetWordsForQuestion gets 10 words for a quiz question, excluding the excludedWordIDs (previously shown words)
func (s *Store) GetWordsForQuestion (ctx context.Context, profile string, quizMode string, lang string,
									excludedWordIDs []int64) ([]Word, error) {
	 var sb strings.Builder
	 sb.WriteString(`
	 SELECT w.id, w.spanish, w.english, w.ukrainian
		FROM words w
		INNER JOIN rankings r ON w.id = r.word_id
		INNER JOIN profiles p ON p.id = r.profile_id
		WHERE p.name = ?
	 `)
	 args := []any{profile}

	if len(excludedWordIDs) > 0 {
		placeholders := make([]string, len(excludedWordIDs))
		for i, id := range excludedWordIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		sb.WriteString(" AND w.id NOT IN (" + strings.Join(placeholders, ",") + ") ")
	}

	switch quizMode {
	case QuizModeAny:
		sb.WriteString(" ORDER BY RANDOM() LIMIT ?")
	case QuizModeLeastKnown:
		sb.WriteString(fmt.Sprintf(" ORDER BY %s ASC, RANDOM() LIMIT ?", lang))
	case QuizModeWellKnown:
		sb.WriteString(fmt.Sprintf(" ORDER BY %s DESC, RANDOM() LIMIT ?", lang))
	}

	log.DebugSQL("SQL: %s", sb.String())

	args = append(args, 10) // LIMIT 10 words from which we can select a question

	rows, err := s.db.QueryContext(ctx, sb.String(), args...)

	if err != nil {
		return nil, fmt.Errorf("store: get words for question: %w", err)
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Spanish, &w.English, &w.Ukrainian); err != nil {
			return nil, fmt.Errorf("store: get words for question: %w", err)
		}
		words = append(words, w)
	}
	return words, rows.Err()
}