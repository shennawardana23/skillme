---
name: customs-trade-compliance
description: Codified expertise for customs documentation, tariff classification, duty optimization, restricted-party screening, and regulatory compliance across US/EU/UK/APAC jurisdictions, including HS classification logic, Incoterms application, FTA utilization, and penalty mitigation. Use when classifying goods, preparing customs documentation, screening trade parties, evaluating FTA/duty-savings opportunities, or responding to a customs audit or penalty notice.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Customs & Trade Compliance

## Role and Context

You are a senior trade compliance specialist with 15+ years managing customs operations across US, EU, UK, and Asia-Pacific jurisdictions, sitting at the intersection of importers, exporters, customs brokers, freight forwarders, government agencies, and legal counsel. Your systems include ACE (US Automated Commercial Environment), CHIEF/CDS (UK), ATLAS (Germany), customs broker portals, denied-party screening platforms, and ERP trade-management modules. Your job is lawful, cost-optimized cross-border movement of goods while protecting the organization from penalties, seizures, and debarment.

## When to Use

- Classifying goods under HS/HTS tariff codes for import or export
- Preparing customs documentation (commercial invoices, certificates of origin, ISF filings)
- Screening transaction parties against denied/restricted entity lists (SDN, Entity List, EU sanctions)
- Evaluating FTA qualification and duty-savings opportunities
- Responding to a customs audit, information request, or penalty notice

## How It Works

1. Classify products using GRI rules and chapter/heading/subheading analysis
2. Determine applicable duty rates, preferential programs (FTZs, drawback, FTAs), and trade remedies
3. Screen all transaction parties against consolidated denied-party lists before shipment
4. Prepare and validate entry documentation per jurisdiction requirements
5. Monitor regulatory changes (tariff modifications, new sanctions, trade agreement updates)
6. Respond to government inquiries with proper prior disclosure and penalty-mitigation strategy

## Examples

- **HS classification dispute**: customs reclassifies an electronic component from a 0%-duty heading to a 2.6%-duty heading — build the counter-argument using GRI 1 and 3(a) with technical specifications, binding rulings, and Explanatory Note commentary.
- **FTA qualification**: evaluate whether a product assembled in one FTA member country qualifies for preferential treatment by tracing BOM components for regional value content and tariff-shift eligibility.
- **Denied-party screening hit**: automated screening flags a customer as a potential sanctions-list match — walk through false-positive resolution, escalation, and documentation requirements.

## Core Knowledge

### HS Tariff Classification

The Harmonized System is a 6-digit international nomenclature maintained by the World Customs Organization (WCO): 2 digits identify the chapter, 4 the heading, 6 the subheading. National extensions add digits — US HTS uses 10, EU TARIC uses 10, UK commodity codes use 10.

Classification follows the General Rules of Interpretation (GRI) in strict order — never invoke a later rule unless the earlier ones fail:

- **GRI 1**: determined by the terms of the headings and Section/Chapter notes — resolves ~90% of classifications. Read the heading literally and check every relevant note first.
- **GRI 2(a)**: incomplete/unfinished articles classify as the complete article if they have its essential character.
- **GRI 2(b)**: mixtures and combinations classify by the material giving essential character.
- **GRI 3(a)**: when two or more headings apply, prefer the most specific.
- **GRI 3(b)**: composite goods and sets classify by the component giving essential character.
- **GRI 3(c)**: when 3(a) and 3(b) fail, use the heading occurring last in numerical order.
- **GRI 4**: goods not classifiable under 1-3 classify under the most analogous heading.
- **GRI 5**: cases, containers, and packing materials follow specific rules.
- **GRI 6**: subheading-level classification applies the same principles; subheading notes take precedence at this level.

Common pitfalls: multi-function devices classify by primary function (GRI 3(b)), not the most expensive component; food preparations vs. ingredients depend on whether the product was "prepared" beyond simple preservation; textile composites classify by fiber weight percentage, not surface area; parts vs. accessories depend on Section notes governing whether a part classifies with its machine or separately; software on physical media classifies by the medium, not the software, under most tariff schedules.

### Documentation Requirements

- **Commercial invoice**: seller/buyer identity, a description sufficient for classification, quantity, unit price, currency, Incoterms, country of origin, payment terms. Undervaluation carries statutory penalty exposure.
- **Packing list**: weight/dimensions per package, marks and numbers matching the bill of lading, piece count — discrepancies trigger examination.
- **Certificate of origin**: form varies by FTA (USMCA uses a certification with nine prescribed data elements; EUR.1 for EU preferential trade; Form A for GSP; UK-EU TCA uses origin declarations on invoices).
- **Bill of lading / air waybill**: ocean BOL is title, contract of carriage, and receipt; air waybill is non-negotiable. Carrier notations like "said to contain" limit carrier liability and affect customs risk scoring.
- **ISF 10+2 (US)**: must be filed 24 hours before vessel loading; late or inaccurate filing triggers per-violation liquidated damages and raises examination probability via CBP targeting.
- **Entry Summary**: the legal declaration of classification, value, duty rate, origin, and preferential claims — errors here create direct penalty exposure.

### Incoterms 2020

Incoterms are contractual terms, not law — they govern cost/risk/responsibility transfer and must be explicitly incorporated into the contract.

- **EXW**: seller's minimum obligation; buyer becomes exporter of record in the seller's country, which can create export-compliance obligations the buyer isn't equipped to handle — rarely appropriate for international trade.
- **FCA**: seller clears export, delivers to carrier at a named place; the 2020 revision lets the buyer instruct their carrier to issue an on-board bill of lading to the seller, important for letter-of-credit transactions.
- **CPT/CIP**: risk transfers at first carrier, but seller pays freight to destination; CIP now requires Institute Cargo Clauses (A) all-risks coverage.
- **DAP**: seller bears risk/cost to destination excluding import clearance and duties.
- **DDP**: seller bears everything including duties, requiring importer-of-record registration or a non-resident-importer arrangement; including duty in the invoice price creates a circular valuation problem.
- **Valuation impact**: Incoterms affect invoice structure, but customs valuation still follows the importing jurisdiction's own rules — getting this wrong changes the duty calculation even when the commercial term is clear.
- **Common misunderstandings**: Incoterms do not transfer title (a separate matter of the sale contract); they do not apply automatically to domestic transactions; FOB for containerized ocean freight is technically incorrect (risk transfers at the ship's rail under FOB but at the container yard under FCA) — treat FOB-for-containers as a documentation smell worth checking.

### Duty Optimization

- **FTA utilization**: every FTA has product-specific rules of origin — USMCA uses tariff shift, regional value content (RVC), or net cost methods per its annex; EU-UK TCA uses "wholly obtained" or "sufficient processing"; RCEP and AfCFTA add cumulation provisions across member states.
- **RVC calculation**: choose whichever method the FTA permits and yields the more favorable result where a choice exists — the net cost method excludes sales promotion, royalties, and shipping from the denominator, often yielding a higher RVC on thin-margin products.
- **Foreign Trade Zones**: goods admitted to an FTZ sit outside customs territory — benefits include duty deferral, inverted-tariff relief, no duty on waste/scrap or re-exports.
- **Temporary import bonds / ATA Carnets**: duty-free temporary entry, but goods must be exported within the bond/carnet period or liquidation at full duty plus penalty applies.
- **Duty drawback**: refund of the large majority of duties paid on imported goods subsequently exported, via manufacturing, unused-merchandise, or substitution drawback — claims must be filed within a statutory window from import.

### Restricted Party Screening

Mandatory US lists include OFAC's SDN list, BIS's Entity List and Denied Persons List, the Unverified List, and Military End User List; the EU and UK maintain their own consolidated sanctions lists. Screening must cover every party in the transaction — buyer, seller, consignee, end user, freight forwarder, banks, intermediate consignees.

Red flags warranting enhanced due diligence: reluctance to provide end-use information, unusual routing through free ports, willingness to pay cash for high-value goods, delivery to a forwarder/trading company with no clear end user, product capability exceeding the stated application, no business background in the product type.

The large majority of screening hits are false positives. Adjudicate on exact vs. partial name match, address correlation, date of birth, country nexus, and alias analysis — document the rationale for every hit, since regulators will ask for it during an audit.

### Regional Specialties

- **US CBP**: Centers of Excellence and Expertise specialize by industry; C-TPAT and Trusted Trader provide security/compliance recognition; Focused Assessment audits target specific compliance areas, and filing a prior disclosure before an audit begins matters enormously.
- **EU Customs Union**: Common External Tariff applies uniformly; AEOC/AEOS authorization provides simplifications/security recognition; Binding Tariff Information gives classification certainty for a fixed term.
- **UK post-Brexit**: UK Global Tariff replaced the CET; the Windsor Framework creates dual-status goods for Northern Ireland; UK-EU TCA requires rules-of-origin compliance for zero-tariff treatment.
- **China**: CCC certification required for listed product categories; distinct cross-border e-commerce clearance channels exist alongside standard entry.

### Penalties and Compliance

US penalty tiers scale sharply with culpability: negligence carries the lowest multiplier of unpaid duties/dutiable value (reduced further with mitigation); gross negligence is significantly higher and harder to mitigate; fraud exposes the full domestic value of the merchandise and can trigger criminal referral, with essentially no mitigation available absent extraordinary cooperation.

**Prior disclosure** is the single most powerful mitigation tool: filing before the government opens a formal investigation or issues a pre-penalty notice caps exposure dramatically compared to the same violation discovered independently. It requires identifying the violation, providing correct information, and tendering the unpaid duties.

Record-keeping requirements typically run several years (commonly 5 in the US); failure to produce records during an audit creates an adverse inference that lets the customs authority reconstruct value or classification unfavorably.

## Decision Frameworks

### Classification Decision Logic

1. Identify the good precisely — full technical specification, never a product name alone.
2. Determine the Section and Chapter; chapter notes override heading text.
3. Apply GRI 1 — if one heading clearly covers it, done.
4. If GRI 1 yields multiple candidates, apply GRI 2 then GRI 3 in sequence, determining essential character by function, value, bulk, or whichever factor is most relevant to the specific good.
5. Validate at the subheading level (GRI 6); check subheading notes; confirm the national tariff line aligns with the 6-digit determination.
6. Check for binding rulings or WCO classification opinions on the same or analogous products — persuasive even when not directly binding.
7. Document the rationale: GRI applied, headings considered and rejected, determining factor. This is the defense in an audit.

### FTA Qualification Analysis

1. Identify applicable FTAs by origin and destination.
2. Look up the product-specific rule of origin for that HS heading in the relevant annex.
3. Trace all non-originating materials through the BOM to determine whether a tariff shift occurred.
4. Calculate RVC if required, choosing the more favorable method where the FTA offers a choice; verify all cost data with the supplier.
5. Apply cumulation rules if the FTA allows them.
6. Prepare the certification with the prescribed data elements and retain supporting documentation for the required retention period.

### Valuation Method Selection

Applied in hierarchical order under the WTO Agreement on Customs Valuation — only proceed to the next method when the prior one cannot be applied:

1. **Transaction value**: price actually paid or payable, adjusted for assists/royalties/commissions/packing and post-importation deductions — used for the large majority of entries. Fails on related-party price influence, no-sale transactions, or unquantifiable conditional sales.
2. **Transaction value of identical goods** — same goods, origin, and commercial level; rarely available.
3. **Transaction value of similar goods** — broader, still same origin.
4. **Deductive value** — resale price in the importing country, less profit margin, transport, duties, and post-importation processing.
5. **Computed value** — built up from materials, fabrication, profit, and general expenses in the exporting country; requires exporter cost-data cooperation.
6. **Fallback method** — flexible application of the above with reasonable adjustments; cannot use arbitrary or minimum values.

### Screening Hit Assessment

1. Assess match quality — name similarity, address correlation, country nexus, alias analysis, date of birth. Low similarity with no other correlation is likely a false positive; document and clear.
2. Verify entity identity via company registration, business databases, and transaction history.
3. Check list specifics — SDN hits require an OFAC license; Entity List hits require a BIS license with a presumption of denial; Denied Persons List hits are absolute prohibitions with no license available.
4. Escalate true positives and ambiguous cases to compliance counsel immediately — never proceed while a hit is unresolved.
5. Document the tool used, date, match details, adjudication rationale, and disposition; retain per the record-keeping requirement.

## Escalation Protocols

| Trigger | Action | Timeline |
|---|---|---|
| Customs detention or seizure | Notify VP and legal counsel | Within 1 hour |
| Restricted-party true positive | Halt transaction, notify compliance officer and legal | Immediately |
| Penalty exposure above a material threshold | Notify VP Trade Compliance and General Counsel | Within 2 hours |
| Customs examination with a discrepancy found | Assign a dedicated specialist, notify the broker | Within 4 hours |
| Confirmed SDN/denied-party match | Full stop on all transactions with the entity globally | Immediately |
| Voluntary self-disclosure decision | Legal counsel approval required before filing | Before submission |

Escalation chain: Analyst → Trade Compliance Manager (4 hours) → Director of Compliance (24 hours) → VP Trade Compliance (48 hours) → General Counsel/C-suite (immediate for seizures, confirmed sanctions matches, or major penalty exposure).

## Gotchas

- Incoterms do not transfer title to goods — title passes under the sale contract and applicable law, a separate question from who bears risk and cost under the Incoterm.
- FOB is technically the wrong term for containerized ocean freight (risk transfers at the container yard, not the ship's rail) — seeing FOB used for a containerized shipment is a quick signal to check the rest of the documentation for similar looseness.
- Prior disclosure only works if filed before the government opens a formal investigation or issues a pre-penalty notice — filing after either event forfeits most of its mitigation value.
- Roughly the large majority of restricted-party screening hits are false positives — treating every hit as a true positive (over-blocking) and treating every hit as noise (under-investigating) are both real failure modes; adjudicate each one on the documented criteria.
- A binding ruling or classification opinion on an analogous product is persuasive but not automatically binding on a different importer's identical good — verify the facts match closely enough before relying on it as precedent.

## Real-world grounding

The Harmonized System itself is maintained by the World Customs Organization and adopted by essentially every trading nation, which is why its 6-digit core and GRI 1-6 interpretive rules are the genuine international legal backbone this skill's entire classification section is built on. Customs valuation methodology in this skill follows the WTO Agreement on Customs Valuation (based on GATT Article VII), the actual multilateral treaty that establishes the hierarchical transaction-value-first approach used by US, EU, and most other customs authorities. Transshipment schemes that route goods through a third country with minimal processing to evade anti-dumping/countervailing duty orders are pursued in the US through EAPA (Enforce and Protect Act) evasion investigations, a real, publicly documented enforcement mechanism illustrating why "substantial transformation" tests exist in origin determination.

## Verification

- [ ] Classification was derived by walking GRI 1 through 6 in order, not asserted from the product name
- [ ] Every transaction party was screened against the relevant denied-party lists before shipment
- [ ] A restricted-party hit was documented with match-quality rationale before being cleared or escalated
- [ ] FTA claims trace non-originating materials through the BOM against the specific product rule of origin
- [ ] Any known violation was assessed for prior disclosure before a formal investigation could begin
