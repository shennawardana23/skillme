package similarity

import "testing"

// TestSimilarity_RelatedDescriptionsScoreHigherThanUnrelated proves the
// core ranking property: two descriptions sharing real, specific
// vocabulary (not just stop words) score higher than two unrelated ones,
// and a description compared with itself scores 1 (or extremely close,
// modulo floating point).
func TestSimilarity_RelatedDescriptionsScoreHigherThanUnrelated(t *testing.T) {
	c := BuildCorpus(map[string]string{
		"security-review":        "Reviews code changes for security vulnerabilities, injection risks, missing authorization checks",
		"security-and-hardening": "Threat-models trust boundaries and designs OWASP Top 10 prevention, authorization checks, injection defenses",
		"seo":                    "Optimizes meta tags, sitemaps, structured data, and crawlability for search engine ranking",
	})

	related := c.Similarity("security-review", "security-and-hardening")
	unrelated := c.Similarity("security-review", "seo")

	if related <= unrelated {
		t.Fatalf("related score %.3f should be greater than unrelated score %.3f", related, unrelated)
	}
	if related <= 0 {
		t.Fatalf("related score %.3f should be positive — both descriptions share real vocabulary (injection, authorization, checks)", related)
	}

	self := c.Similarity("security-review", "security-review")
	if self < 0.99 {
		t.Fatalf("a document compared with itself should score ~1.0, got %.3f", self)
	}
}

func TestSimilarity_UnknownNameReturnsZero(t *testing.T) {
	c := BuildCorpus(map[string]string{"x": "some description text here"})
	if got := c.Similarity("x", "does-not-exist"); got != 0 {
		t.Fatalf("Similarity with an unknown name = %v, want 0", got)
	}
}

func TestTopPairs_RespectsMinScoreAndOrdering(t *testing.T) {
	c := BuildCorpus(map[string]string{
		"a": "reviews code changes for security vulnerabilities and injection risks",
		"b": "reviews code changes for security vulnerabilities and injection defenses",
		"c": "optimizes meta tags and sitemaps for search engine ranking",
	})

	// a and b share 6 of 7 tokens, differing only in the last word — but
	// TF-IDF weighs each doc's one unique word (df=1, higher idf) more
	// heavily per-token than the 6 shared words (df=2, lower idf), so the
	// actual cosine similarity lands around 0.45, not near 1.0 as a naive
	// "6/7 words match" intuition would suggest. 0.4 is chosen from that
	// real computed value, not assumed — this is the same lesson as
	// evalspec.Lint's own history: verify a threshold against real output
	// before trusting it, don't guess a round number.
	pairs := c.TopPairs(0.4)
	if len(pairs) != 1 {
		t.Fatalf("TopPairs(0.4) = %v, want exactly the (a,b) pair", pairs)
	}
	if !(pairs[0].A == "a" && pairs[0].B == "b") {
		t.Fatalf("TopPairs(0.4)[0] = %+v, want the a/b pair", pairs[0])
	}

	all := c.TopPairs(0)
	if len(all) != 3 {
		t.Fatalf("TopPairs(0) over 3 docs should return all 3 pairs, got %d", len(all))
	}
	for i := 1; i < len(all); i++ {
		if all[i-1].Score < all[i].Score {
			t.Fatalf("TopPairs is not sorted descending: %+v then %+v", all[i-1], all[i])
		}
	}
}

// TestStopWordsDoNotInflateSimilarity proves the specific failure mode this
// package exists to avoid: two totally unrelated skills whose descriptions
// both follow this catalog's "Use when the user asks..." convention must
// not score as similar just from sharing that boilerplate.
func TestStopWordsDoNotInflateSimilarity(t *testing.T) {
	c := BuildCorpus(map[string]string{
		"a": "Use when the user asks to review a pull request for security vulnerabilities",
		"b": "Use when the user asks to optimize a database query for performance",
	})
	if got := c.Similarity("a", "b"); got > 0.15 {
		t.Fatalf("Similarity = %.3f, want near 0 — these share only boilerplate phrasing, not real topic overlap", got)
	}
}
