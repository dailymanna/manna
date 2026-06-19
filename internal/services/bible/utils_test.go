package bible

import (
	"slices"
	"testing"
)

var testRefs = []CrossReference{
	{ID: 1, ToBook: "John", ToChapter: 3, ToVerseStart: 16, Votes: 50},
	{ID: 2, ToBook: "Romans", ToChapter: 8, ToVerseStart: 28, Votes: 120},
	{ID: 3, ToBook: "Genesis", ToChapter: 1, ToVerseStart: 1, Votes: 80},
	{ID: 4, ToBook: "John", ToChapter: 1, ToVerseStart: 1, Votes: 30},
}

func TestSortedCrossRefs_ByVotesDesc(t *testing.T) {
	got := SortedCrossRefs(slices.Values(testRefs), ByVotesDesc)

	wantVotes := []int{120, 80, 50, 30}
	for i, ref := range got {
		if ref.Votes != wantVotes[i] {
			t.Errorf("position %d: got votes %d, want %d", i, ref.Votes, wantVotes[i])
		}
	}
}

func TestSortedCrossRefs_ByTargetLocation(t *testing.T) {
	got := SortedCrossRefs(slices.Values(testRefs), ByTargetLocation)

	// Genesis < John < Romans alphabetically; within John: ch1 before ch3
	wantIDs := []int{3, 4, 1, 2}
	for i, ref := range got {
		if ref.ID != wantIDs[i] {
			t.Errorf("position %d: got ID %d, want %d", i, ref.ID, wantIDs[i])
		}
	}
}

func TestSortedCrossRefs_ComposedComparator(t *testing.T) {
	refs := []CrossReference{
		{ID: 1, ToBook: "John", ToChapter: 3, ToVerseStart: 16, Votes: 80},
		{ID: 2, ToBook: "Romans", ToChapter: 8, ToVerseStart: 28, Votes: 80},
		{ID: 3, ToBook: "Genesis", ToChapter: 1, ToVerseStart: 1, Votes: 120},
	}

	// Sort by votes desc, then location as tiebreaker.
	got := SortedCrossRefs(slices.Values(refs), func(a, b CrossReference) int {
		if n := ByVotesDesc(a, b); n != 0 {
			return n
		}
		return ByTargetLocation(a, b)
	})

	// Genesis (120) first; then John and Romans tied on votes — John < Romans alphabetically.
	wantIDs := []int{3, 1, 2}
	for i, ref := range got {
		if ref.ID != wantIDs[i] {
			t.Errorf("position %d: got ID %d, want %d", i, ref.ID, wantIDs[i])
		}
	}
}
