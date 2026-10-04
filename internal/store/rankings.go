package store

import (
	"context"
	"fmt"

	"github.com/alexbudy/go_spanish_rewrite/internal/log"
)

/* rankings table DB queries */

// ResetRankingsForWord resets the rankings for a word, lang1_to_lang2
func (s *Store)	ResetRankingsForWord(ctx context.Context, wordId int64, lang1 string, lang2 string, profileId int64) error {
	resetQry := "UPDATE rankings SET %s = 0.0 WHERE word_id = ? AND profile_id = ?"

	_, err := s.db.ExecContext(ctx, fmt.Sprintf(resetQry, lang1+"_to_"+lang2), wordId, profileId)
	if err != nil {
		return fmt.Errorf("store: reset rankings: %w", err)
	}
	return nil
}

// UpdateRankingForWord adjusts the ranking for a word and profile in the given language direction.
func (s *Store) UpdateRankingForWord(ctx context.Context, wordID int64, profileID int64,
	lang1 string, lang2 string, adjustment float64) error {
	rankingToUpdate := lang1 + "_to_" + lang2
	updateQry := fmt.Sprintf(
		"UPDATE rankings SET %s = %s + ? WHERE word_id = ? AND profile_id = ?",
		rankingToUpdate,
		rankingToUpdate,
	)
	if _, err := s.db.ExecContext(ctx, updateQry, adjustment, wordID, profileID); err != nil {
		return fmt.Errorf("store: update ranking: %w", err)
	}
	log.DebugSQL(updateQry+" for word %d and profile %d with adjustment %f", wordID, profileID, adjustment)
	return nil
}