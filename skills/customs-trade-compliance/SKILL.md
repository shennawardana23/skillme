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

Every entry rests on a small set of documents — commercial invoice, packing list, certificate of origin, bill of lading/air waybill, ISF 10+2 (US), Entry Summary — each with its own required data elements and penalty exposure if wrong. Full per-document requirements: [references/documentation-and-regional-reference.md](references/documentation-and-regional-reference.md) — load when preparing or reviewing a document set.

### Incoterms 2020

Incoterms are contractual terms, not law — they govern cost/risk/responsibility transfer and must be explicitly incorporated into the contract. They do not transfer title (a separate matter of the sale contract), and customs valuation still follows the importing jurisdiction's own rules regardless of which term is used. The term-by-term breakdown (EXW, FCA, CPT/CIP, DAP, DDP) is in [references/documentation-and-regional-reference.md](references/documentation-and-regional-reference.md) — load it when selecting or reviewing an Incoterm.

### Duty Optimization

Every FTA has a product-specific rule of origin (tariff shift, regional value content, or net cost, depending on the agreement); beyond FTAs, Foreign Trade Zones, temporary import bonds/ATA Carnets, and duty drawback each offer a distinct duty-savings mechanism with its own eligibility and filing-window rules. Program-by-program detail: [references/documentation-and-regional-reference.md](references/documentation-and-regional-reference.md) — load when structuring a specific claim.

### Restricted Party Screening

Screening must cover every party in the transaction — buyer, seller, consignee, end user, freight forwarder, banks, intermediate consignees — against the applicable mandatory lists (list names in the reference below).

Red flags warranting enhanced due diligence: reluctance to provide end-use information, unusual routing through free ports, willingness to pay cash for high-value goods, delivery to a forwarder/trading company with no clear end user, product capability exceeding the stated application, no business background in the product type.

The large majority of screening hits are false positives. Adjudicate on exact vs. partial name match, address correlation, date of birth, country nexus, and alias analysis — document the rationale for every hit, since regulators will ask for it during an audit. List names and jurisdiction-specific screening notes: [references/documentation-and-regional-reference.md](references/documentation-and-regional-reference.md).

### Regional Specialties

US CBP, EU Customs Union, UK post-Brexit, and China each layer their own authorizations, tariff structures, and certification regimes on the general framework above. Jurisdiction-by-jurisdiction detail: [references/documentation-and-regional-reference.md](references/documentation-and-regional-reference.md) — load when working a specific jurisdiction.

### Penalties and Compliance

US penalty exposure scales sharply with culpability — negligence, gross negligence, and fraud carry very different multipliers and mitigation options — and **prior disclosure**, filed before the government opens a formal investigation or pre-penalty notice, is the single most powerful way to cap exposure. Penalty-tier detail and record-keeping requirements: [references/penalties-and-escalation.md](references/penalties-and-escalation.md) — load when assessing exposure on a known violation.

## Decision Frameworks

Step-by-step procedures for classification decision logic, FTA
qualification analysis, valuation method selection, and screening hit
assessment are in
[references/decision-frameworks.md](references/decision-frameworks.md)
— load it when actually executing one of these decisions.

## Escalation Protocols

The full escalation trigger table and chain (detention/seizure,
restricted-party hits, penalty exposure, self-disclosure decisions) is
in [references/penalties-and-escalation.md](references/penalties-and-escalation.md)
— load it when routing an active incident.

## Gotchas

- Incoterms do not transfer title to goods — title passes under the sale contract and applicable law, a separate question from who bears risk and cost under the Incoterm.
- FOB is technically the wrong term for containerized ocean freight (risk transfers at the container yard, not the ship's rail) — seeing FOB used for a containerized shipment is a quick signal to check the rest of the documentation for similar looseness.
- Prior disclosure only works if filed before the government opens a formal investigation or issues a pre-penalty notice — filing after either event forfeits most of its mitigation value.
- Roughly the large majority of restricted-party screening hits are false positives — treating every hit as a true positive (over-blocking) and treating every hit as noise (under-investigating) are both real failure modes; adjudicate each one on the documented criteria.
- A binding ruling or classification opinion on an analogous product is persuasive but not automatically binding on a different importer's identical good — verify the facts match closely enough before relying on it as precedent.

## Real-world grounding

The Harmonized System is maintained by the World Customs Organization and adopted by essentially every trading nation, which is why its 6-digit core and GRI 1-6 interpretive rules are the genuine international legal backbone this skill's classification section is built on. Customs valuation here follows the WTO Agreement on Customs Valuation (based on GATT Article VII), the actual multilateral treaty establishing the hierarchical transaction-value-first approach used by US, EU, and most other customs authorities. Transshipment schemes that route goods through a third country to evade anti-dumping/countervailing duty orders are pursued in the US through real EAPA (Enforce and Protect Act) evasion investigations, illustrating why "substantial transformation" tests exist in origin determination.

## Verification

- [ ] Classification was derived by walking GRI 1 through 6 in order, not asserted from the product name
- [ ] Every transaction party was screened against the relevant denied-party lists before shipment
- [ ] A restricted-party hit was documented with match-quality rationale before being cleared or escalated
- [ ] FTA claims trace non-originating materials through the BOM against the specific product rule of origin
- [ ] Any known violation was assessed for prior disclosure before a formal investigation could begin
