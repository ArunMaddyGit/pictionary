package words

import "testing"

func TestGetRandomWords_Count(t *testing.T) {
	got := GetRandomWords(3)
	if len(got) != 3 {
		t.Fatalf("len(GetRandomWords(3)) = %d, want 3", len(got))
	}
}

func TestGetRandomWords_NoDuplicates(t *testing.T) {
	got := GetRandomWords(10)
	seen := make(map[string]struct{}, len(got))
	for _, w := range got {
		if _, exists := seen[w]; exists {
			t.Fatalf("duplicate word found: %q", w)
		}
		seen[w] = struct{}{}
	}
}

func TestGetRandomWords_FromWordBank(t *testing.T) {
	got := GetRandomWords(12)

	bank := GetWordBank()
	all := make(map[string]struct{})
	for _, words := range bank {
		for _, w := range words {
			all[w] = struct{}{}
		}
	}

	for _, w := range got {
		if _, ok := all[w]; !ok {
			t.Fatalf("word %q is not in word bank", w)
		}
	}
}

func TestIsDuplicateInHistory_True(t *testing.T) {
	history := []string{"elephant", "bicycle", "swimming"}
	if !IsDuplicateInHistory("bicycle", history) {
		t.Fatal("expected duplicate word")
	}
}

func TestIsDuplicateInHistory_False(t *testing.T) {
	history := []string{"elephant", "bicycle", "swimming"}
	if IsDuplicateInHistory("kangaroo", history) {
		t.Fatal("expected non-duplicate word")
	}
}

func TestGetWordBank_IsCopy(t *testing.T) {
	copyBank := GetWordBank()
	copyBank["animals"][0] = "mutated"
	copyBank["new-category"] = []string{"x"}

	original := GetWordBank()
	if original["animals"][0] == "mutated" {
		t.Fatal("original WordBank should not be affected by copy mutation")
	}
	if _, exists := original["new-category"]; exists {
		t.Fatal("original WordBank should not gain new categories from copy mutation")
	}
}
