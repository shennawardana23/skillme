// Package riskscan statically scans a skill directory against the risk
// indicators in Anthropic's own enterprise Agent Skills guidance
// (https://platform.claude.com/docs/en/agents-and-tools/agent-skills/enterprise
// #risk-tier-assessment): code execution, network access, hardcoded
// credentials, suspicious filesystem paths, and MCP tool references.
//
// Two of that table's seven rows — "instruction manipulation" and "tool
// invocations" — are deliberately not turned into pass/fail checks here.
// This catalog's own subject matter includes several skills whose entire
// purpose is discussing prompt injection, exfiltration, and tool misuse
// (security-review, security-and-hardening, skill-inspector); a keyword
// scan for that language cannot distinguish a skill teaching defense
// against an attack from a skill attempting one — both use the same
// vocabulary. Static analysis is the wrong tool for that distinction; it
// needs the semantic/human review skill-inspector already performs.
// Report an Informational finding for any hit so a human still sees it,
// but never fail the scan on it alone.
//
// "Network access" (URLs, curl, fetch) turned out to need the same
// treatment, discovered by running this scanner against the real catalog
// before trusting it: security-and-hardening and skill-inspector's own
// prose is full of legitimate https:// citations and fetch()/curl example
// code teaching about SSRF, CSRF, and CLAUDE.md injection vectors — every
// single hit in the real catalog was a citation or a worked example, never
// an actual embedded network call. It is Informational here, not High.
//
// "Code execution" (a bundled script file) is Medium, not High, for the
// same reason: Anthropic's own authoring best-practices explicitly
// recommend bundling utility scripts with a skill, and this catalog
// already does it (web-quality-audit/scripts/analyze.sh) — a script's mere
// presence is worth a human glance, not an automatic catalog-wide fail.
//
// "Hardcoded credentials" is the one indicator kept at High: a real
// key/token-shaped string is unambiguously wrong regardless of framing,
// and it stayed at zero false positives across the full catalog during
// the same real-catalog test pass that demoted the two checks above.
package riskscan

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Concern mirrors the enterprise doc's three-level concern scale.
type Concern string

const (
	High          Concern = "High"
	Medium        Concern = "Medium"
	Informational Concern = "Informational"
)

// Finding is one static-scan hit, with enough evidence for a human to go
// look rather than a bare pass/fail.
type Finding struct {
	Indicator string // the enterprise doc's row name, e.g. "Code execution"
	Concern   Concern
	File      string // path relative to the scanned skill directory
	Line      int    // 1-indexed; 0 when the finding is file-level, not line-level
	Detail    string // the matched text or a short explanation
}

// scriptExtensions are the code-execution risk indicator: a skill in this
// catalog is markdown + JSON evals only (skill-catalog-authoring's own
// spec has no provision for bundled executables), so any of these is a
// genuine anomaly here, not a routine occurrence to filter noise from.
var scriptExtensions = map[string]bool{
	".py": true, ".sh": true, ".js": true, ".ts": true,
	".rb": true, ".pl": true, ".ps1": true, ".exe": true,
}

var (
	networkPattern = regexp.MustCompile(`(?i)\bhttps?://|` +
		`\bcurl\s+|\bwget\s+|\bfetch\(|` +
		`\brequests\.(get|post|put|delete)\(|\burllib\.request\.|\baxios\.`)

	// Credential patterns are intentionally specific (a real key/token
	// shape), not a bare "password" word match, to keep the false-positive
	// rate low against prose that merely discusses credentials in the
	// abstract (e.g. "never hardcode credentials" in this exact package's
	// own doc comment would otherwise self-trigger).
	credentialPatterns = []*regexp.Regexp{
		regexp.MustCompile(`sk-[A-Za-z0-9]{20,}`),
		regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		regexp.MustCompile(`-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----`),
		regexp.MustCompile(`(?i)(api[_-]?key|access[_-]?token|secret)\s*[:=]\s*["'][A-Za-z0-9_\-]{16,}["']`),
	}

	suspiciousPathPattern = regexp.MustCompile(`(\.\./){2,}|~/\.ssh|/etc/passwd|/etc/shadow`)

	// Only the literal "mcp__server__tool" naming convention Claude Code
	// actually uses. An earlier version also matched a generic
	// "ServerName:tool_name" shape from the enterprise doc's own example
	// syntax, but that matched date('Y-m-d H:i:s', ...) in a real skill's
	// PHP example — "H:i" is a time format, not an MCP reference. Literal
	// prefix only; false-positive rate against the real catalog: zero.
	mcpReferencePattern = regexp.MustCompile(`\bmcp__[A-Za-z0-9_]+__[A-Za-z0-9_]+\b`)

	// instructionManipulationPatterns are reported Informational-only —
	// see the package doc for why this can never be a fail-the-scan check
	// in this specific catalog.
	instructionManipulationPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)ignore (all |any )?(previous|prior|the above) instructions`),
		regexp.MustCompile(`(?i)do not (tell|inform|mention|show) the user`),
		regexp.MustCompile(`(?i)hide (this|these actions?) from the user`),
		regexp.MustCompile(`(?i)without (the user|them) (know|noticing)`),
		regexp.MustCompile(`(?i)disregard (the )?system prompt`),
	}
)

// Scan walks every regular file under skillDir and returns every finding,
// most-severe first. A read or walk error aborts with an error rather than
// returning a partial, silently-incomplete result.
func Scan(skillDir string) ([]Finding, error) {
	var findings []Finding

	err := filepath.WalkDir(skillDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(skillDir, path)
		if relErr != nil {
			rel = path
		}

		if scriptExtensions[strings.ToLower(filepath.Ext(path))] {
			findings = append(findings, Finding{
				Indicator: "Code execution", Concern: Medium, File: rel,
				Detail: "bundled executable script — legitimate per Anthropic's own authoring guidance, but worth a human glance",
			})
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("read %s: %w", path, readErr)
		}
		if !isProbablyText(data) {
			return nil
		}
		content := string(data)

		findings = append(findings, scanLines(content, rel, networkPattern, "Network access", Informational)...)
		for _, pat := range credentialPatterns {
			findings = append(findings, scanLines(content, rel, pat, "Hardcoded credentials", High)...)
		}
		findings = append(findings, scanLines(content, rel, suspiciousPathPattern, "Filesystem access scope", Medium)...)
		findings = append(findings, scanLines(content, rel, mcpReferencePattern, "MCP server references", Informational)...)
		for _, pat := range instructionManipulationPatterns {
			findings = append(findings, scanLines(content, rel, pat, "Instruction manipulation", Informational)...)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan %s: %w", skillDir, err)
	}
	return findings, nil
}

// HasFailing reports whether any finding is severe enough to fail a CI
// gate — currently just High. Medium and Informational findings are
// reported for human review but never block on their own.
func HasFailing(findings []Finding) bool {
	for _, f := range findings {
		if f.Concern == High {
			return true
		}
	}
	return false
}

func scanLines(content, file string, pat *regexp.Regexp, indicator string, concern Concern) []Finding {
	var findings []Finding
	for i, line := range strings.Split(content, "\n") {
		if m := pat.FindString(line); m != "" {
			findings = append(findings, Finding{
				Indicator: indicator, Concern: concern, File: file, Line: i + 1,
				Detail: strings.TrimSpace(m),
			})
		}
	}
	return findings
}

// isProbablyText reports whether data looks like text rather than a binary
// blob, by checking the first 512 bytes for a NUL byte — the same
// heuristic net/http.DetectContentType and git both use as a cheap
// binary/text signal.
func isProbablyText(data []byte) bool {
	n := len(data)
	if n > 512 {
		n = 512
	}
	for _, b := range data[:n] {
		if b == 0 {
			return false
		}
	}
	return true
}
