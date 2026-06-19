package bible

import (
	"cmp"
	"iter"
	"slices"
)

// SortedCrossRefs returns a sorted slice from any CrossReference sequence.
// Use slices.Values(refs) to get an iter.Seq from a plain slice.
func SortedCrossRefs(seq iter.Seq[CrossReference], compare func(a, b CrossReference) int) []CrossReference {
	return slices.SortedFunc(seq, compare)
}

// ByVotesDesc orders CrossReferences by votes descending (most popular first).
func ByBook(a, b CrossReference) int {
	return cmp.Compare(b.Votes, a.Votes)
}

// ByVotesDesc orders CrossReferences by votes descending (most popular first).
func ByVotesDesc(a, b CrossReference) int {
	return cmp.Compare(b.Votes, a.Votes)
}

// ByTargetLocation orders CrossReferences by canonical Bible location
// (ToBook → ToChapter → ToVerseStart).
func ByTargetLocation(a, b CrossReference) int {
	if n := cmp.Compare(a.ToBook, b.ToBook); n != 0 {
		return n
	}
	if n := cmp.Compare(a.ToChapter, b.ToChapter); n != 0 {
		return n
	}
	return cmp.Compare(a.ToVerseStart, b.ToVerseStart)
}
