// Package similarity finds skills whose descriptions overlap enough to be
// worth a human's attention during the "search the catalog first" pre-flight
// check CONTRIBUTING.md already asks for — currently a purely manual grep.
//
// This uses TF-IDF-weighted cosine similarity over each skill's name +
// description, computed in plain Go with no third-party dependency and no
// model/API call — consistent with this repo's stated invariant for
// cmd/smeval (skill-catalog-authoring/SKILL.md: "plain Go, no third-party
// dependency"). A semantic-embedding approach (calling a model to embed each
// description) would catch paraphrased overlap this lexical approach misses,
// at the cost of a real API key and per-run spend; that's a genuine
// trade-off, not an oversight, and needs its own explicit buy-in the way the
// live-eval-in-CI cost decision did — this package is deliberately the free
// path.
package similarity

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

// stopWords are removed before scoring — without this, "use", "the", and
// "when" (present in nearly every skill's description by convention: "Use
// when...") would dominate every pair's similarity score regardless of
// actual topic overlap, since they carry no discriminative signal at all.
var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "or": true, "the": true, "this": true,
	"that": true, "to": true, "of": true, "in": true, "on": true, "for": true,
	"is": true, "are": true, "be": true, "it": true, "its": true, "with": true,
	"use": true, "used": true, "uses": true, "using": true, "when": true,
	"skill": true, "should": true, "user": true, "asks": true,
	"otherwise": true, "not": true, "than": true, "into": true, "from": true,
	"as": true, "by": true, "at": true, "also": true, "any": true, "such": true,
}

var tokenPattern = regexp.MustCompile(`[a-z0-9]+`)

func tokenize(text string) []string {
	var out []string
	for _, tok := range tokenPattern.FindAllString(strings.ToLower(text), -1) {
		if len(tok) < 3 || stopWords[tok] {
			continue
		}
		out = append(out, tok)
	}
	return out
}

// Corpus is a TF-IDF model built from a fixed set of named documents (one
// per skill: its name + description).
type Corpus struct {
	docs  map[string]map[string]float64 // skill name -> term -> TF-IDF weight
	names []string
	idf   map[string]float64
}

// BuildCorpus computes TF-IDF vectors for every document in docs (skill name
// -> raw text). Term frequency is normalized by document length so a longer
// description doesn't automatically score "more similar" to everything just
// from raw term repetition.
func BuildCorpus(docs map[string]string) *Corpus {
	tokenized := make(map[string][]string, len(docs))
	df := make(map[string]int)
	for name, text := range docs {
		toks := tokenize(text)
		tokenized[name] = toks
		seen := make(map[string]bool, len(toks))
		for _, t := range toks {
			if !seen[t] {
				seen[t] = true
				df[t]++
			}
		}
	}

	n := float64(len(docs))
	idf := make(map[string]float64, len(df))
	for term, count := range df {
		idf[term] = math.Log(n / float64(count))
	}

	c := &Corpus{
		docs: make(map[string]map[string]float64, len(docs)),
		idf:  idf,
	}
	for name, toks := range tokenized {
		tf := make(map[string]int, len(toks))
		for _, t := range toks {
			tf[t]++
		}
		vec := make(map[string]float64, len(tf))
		total := float64(len(toks))
		for t, count := range tf {
			if total == 0 {
				continue
			}
			vec[t] = (float64(count) / total) * idf[t]
		}
		c.docs[name] = vec
		c.names = append(c.names, name)
	}
	sort.Strings(c.names)
	return c
}

// Names returns every document name in the corpus, sorted.
func (c *Corpus) Names() []string { return c.names }

// Similarity returns the cosine similarity between two documents' TF-IDF
// vectors, in [0, 1] (both vectors are non-negative, so cosine similarity
// can't go negative here). 0 if either name is unknown or either vector is
// entirely zero (no scored terms — e.g. a description that's only stop
// words, which shouldn't happen in practice but must not divide by zero).
func (c *Corpus) Similarity(a, b string) float64 {
	va, ok := c.docs[a]
	if !ok {
		return 0
	}
	vb, ok := c.docs[b]
	if !ok {
		return 0
	}

	var dot, normA, normB float64
	for t, wa := range va {
		normA += wa * wa
		if wb, ok := vb[t]; ok {
			dot += wa * wb
		}
	}
	for _, wb := range vb {
		normB += wb * wb
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Pair is one scored document pair, A < B lexicographically so a pair is
// never reported twice in both orders.
type Pair struct {
	A, B  string
	Score float64
}

// TopPairs returns every pair scoring at or above minScore, highest first.
// Comparing all pairs is O(n^2) but n is the number of skills in a catalog
// (low hundreds), not documents in general — trivially fast at this scale.
func (c *Corpus) TopPairs(minScore float64) []Pair {
	var pairs []Pair
	for i := 0; i < len(c.names); i++ {
		for j := i + 1; j < len(c.names); j++ {
			a, b := c.names[i], c.names[j]
			score := c.Similarity(a, b)
			if score >= minScore {
				pairs = append(pairs, Pair{A: a, B: b, Score: score})
			}
		}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Score > pairs[j].Score })
	return pairs
}
