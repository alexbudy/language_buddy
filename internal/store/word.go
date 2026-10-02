package store

// DB access for Words

// Word is a single Spanish/English/Ukraine noun pair.
type Word struct {
	ID          int64   // 1-based ind
	Spanish     string  // Spanish Word
	English     string  // English Word
	Ukraine     string  // Ukraine Word
	Gender      string  // TODO use/implement/delete?
	EnToEsScore float64 // how well the user knows english word from spanish
	EnToUkScore float64 // how well the user knows english word from ukrainian
	EsToEnScore float64 // how well the user knows spanish word from english
	EsToUkScore float64 // how well the user knows spanish word from ukrainian
	UkToEsScore float64 // how well the user knows ukrainian word from spanish
	UkToEnScore float64 // how well the user knows ukrainian word from english
}
