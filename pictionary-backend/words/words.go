package words

import (
	"math/rand"
	"slices"
	"strings"
	"time"
)

// WordBank groups available words by category.
var WordBank = map[string][]string{
	"animals": {
		"elephant", "giraffe", "penguin", "dolphin", "kangaroo",
		"tiger", "lion", "zebra", "octopus", "rabbit",
		"panda", "flamingo", "crocodile", "hedgehog", "koala",
	},
	"objects": {
		"umbrella", "telescope", "bicycle", "backpack", "lantern",
		"clock", "camera", "hammer", "toothbrush", "suitcase",
		"helmet", "keyboard", "microscope", "notebook", "ladder",
	},
	"actions": {
		"swimming", "juggling", "climbing", "painting", "sleeping",
		"dancing", "laughing", "running", "jumping", "cooking",
		"reading", "driving", "whistling", "skiing", "gardening",
	},
}

// GetRandomWords picks count unique random words from the combined word pool.
func GetRandomWords(count int) []string {
	if count <= 0 {
		return []string{}
	}

	pool := allWords()
	if len(pool) == 0 {
		return []string{}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	if count >= len(pool) {
		return pool
	}
	return pool[:count]
}

// GetWordBank returns a deep copy so callers cannot mutate the global WordBank.
func GetWordBank() map[string][]string {
	out := make(map[string][]string, len(WordBank))
	for category, words := range WordBank {
		out[category] = slices.Clone(words)
	}
	return out
}

// IsDuplicateInHistory reports whether word already appears in history.
func IsDuplicateInHistory(word string, history []string) bool {
	for _, w := range history {
		if strings.EqualFold(w, word) {
			return true
		}
	}
	return false
}

func allWords() []string {
	total := 0
	for _, words := range WordBank {
		total += len(words)
	}

	out := make([]string, 0, total)
	for _, words := range WordBank {
		out = append(out, words...)
	}
	return out
}
