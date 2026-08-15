---
name: dependency-and-license-policy
description: Use when the user asks to "add a dependency", "vet a license", "check if this package is safe to use", "can we use this library", "add this to go.mod/package.json", or is preparing a proprietary/closed-source product that bundles open-source code. Guides license-compatibility review before a dependency is pulled in, distinguishing copyleft obligations (GPL/AGPL/LGPL) from permissive ones (MIT/Apache-2.0/BSD) and flagging when legal review is required.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Dependency and License Policy

Every dependency you add is a legal commitment, not just a technical one.
Vet the license before the import, not after a customer's legal team asks
about it during due diligence.

## The core distinction: permissive vs. copyleft

**Permissive** licenses (MIT, BSD-2/3-Clause, Apache-2.0, ISC) let you use,
modify, and ship the code inside proprietary software with no obligation to
release your own source. Apache-2.0 additionally grants an explicit patent
license — prefer it over BSD/MIT when patent exposure matters.

**Copyleft** licenses require derivative works to be distributed under the
same license — this is the mechanism, not a matter of intent:

- **GPL (v2/v3)**: "strong" copyleft. Linking GPL code into your binary
  (statically or, for GPLv3's interpretation, in ways that create one
  combined work) generally obligates you to release your own source under
  GPL if you distribute the binary. This is why companies shipping
  proprietary software avoid GPL dependencies in anything that ships to
  customers.
- **LGPL**: "weak" copyleft — a library under LGPL can be *dynamically
  linked* by proprietary code without that obligation propagating, as long
  as the LGPL component itself stays swappable/replaceable. Static linking
  or bundling in ways that prevent replacement can still trigger the
  obligation — this is the most commonly misunderstood distinction in the
  whole license landscape.
- **AGPL**: extends GPL's obligation to network use — merely running AGPL
  code as a network service (SaaS) that users interact with over a network,
  without ever distributing the binary, still triggers the source-release
  obligation. This is why AGPL is often an outright ban at companies running
  SaaS products, even though GPL itself might be tolerable there.

## Procedure: vetting a new dependency

1. **Identify the license.** Check `LICENSE`/`COPYING` in the package repo,
   or `go.mod`'s module page on pkg.go.dev, or `npm view <pkg> license`.
   Don't trust a README's one-line claim — read the actual license file;
   READMEs are sometimes stale relative to a license change.
2. **Classify it**: permissive, weak copyleft (LGPL/MPL), strong copyleft
   (GPL/AGPL), or "custom/unclear" (source-available, dual-licensed,
   no-license-file). Custom and no-license-file cases require legal review
   before merge — never guess.
3. **Check how it's consumed**: is it linked into a binary you distribute,
   run as a subprocess/separate service, or only used at build/test time
   (dev dependency)? A GPL *build tool* invoked as a subprocess and never
   linked into your shipped artifact is a fundamentally different exposure
   than a GPL *library* imported into your codebase.
4. **Check for license changes.** Some widely used projects have relicensed
   mid-life (e.g., from a permissive license to a source-available or
   Business Source License) — pin and re-vet on every major version bump,
   don't assume the license you checked once still applies.
5. **Escalate, don't decide unilaterally**, when you hit: AGPL anywhere,
   GPL in something distributed to customers, a custom/non-OSI license, or
   a "dual license" you don't fully understand. Flag to legal/engineering
   leadership with the specific clause in question.

## Gotchas

- "MIT-licensed" doesn't guarantee a patent grant — Apache-2.0 does. If
  patent risk from the dependency's authors matters to your business,
  Apache-2.0 gives you more protection than MIT for equivalent freedom.
- A transitive dependency can carry a different, more restrictive license
  than the direct dependency you chose — license scanning tools must walk
  the full dependency graph (`go mod graph`, `npm ls --all`), not just
  top-level `go.mod`/`package.json` entries.
- Dynamically linking LGPL code satisfies the "swappable" requirement in
  most interpretations; statically linking it into a single compiled Go
  binary is far more likely to be scrutinized as creating one combined
  work — Go's static-linking-by-default model makes this a live issue for
  Go services in particular, more so than for languages that dynamically
  link by convention.
- "Source-available" licenses (e.g., BSL, SSPL, Elastic License) look like
  open source at a glance but are not OSI-approved and often carry
  commercial-use restrictions — treat any license you don't recognize by
  name as "custom," not as a variant of MIT/Apache.
- License compatibility is not symmetric: you can combine an MIT dependency
  into a GPL-licensed project, but not the reverse — check the direction of
  the obligation relative to what *you* are shipping.

## Real-world grounding

The GPL/LGPL/AGPL distinctions and the "linking" test that determines when
copyleft obligations propagate are documented directly by the Free
Software Foundation's own GPL FAQ (gnu.org/licenses/gpl-faq.html) — the
canonical, publicly available source for how these licenses are actually
meant to be interpreted, and the reference point most corporate
open-source-usage policies cite when drawing the line between "safe to use
freely" and "requires legal sign-off."

## Verification

- [ ] Every new dependency's license is identified from its actual license
      file, not assumed from the README
- [ ] Transitive dependencies were scanned, not just direct ones
- [ ] Any GPL/AGPL/custom license found triggers escalation, not a
      unilateral decision
- [ ] The consumption model (linked binary vs. subprocess vs. dev-only) was
      considered before accepting a copyleft dependency
- [ ] License re-checked on major version bumps, not just at first add
