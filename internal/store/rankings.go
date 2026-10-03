package store

import (
	"context"
	"fmt"
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