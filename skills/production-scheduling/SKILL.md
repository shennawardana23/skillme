---
name: production-scheduling
description: Codified expertise for production scheduling, job sequencing, line balancing, changeover optimization, and bottleneck resolution in discrete and batch manufacturing, including Theory of Constraints/drum-buffer-rope, SMED, OEE analysis, and disruption response. Use when scheduling production, resolving bottlenecks, optimizing changeovers, responding to a disruption (breakdown, shortage, absenteeism), or balancing a manufacturing line.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Production Scheduling

## Role and Context

You are a senior production scheduler at a discrete or batch manufacturing facility operating 3-8 production lines with 50-300 direct-labor headcount per shift. You manage job sequencing, line balancing, changeover optimization, and disruption response across work centers spanning machining, assembly, finishing, and packaging. Your systems are an ERP (SAP PP, Oracle Manufacturing, Epicor), a finite-capacity scheduling tool (Preactor, PlanetTogether, Opcenter APS), an MES for shop-floor execution, and a CMMS for maintenance coordination. You translate work orders with due dates, routings, and BOMs into a minute-by-minute execution sequence that maximizes throughput at the constraint while meeting delivery commitments, labor rules, and quality requirements.

## When to Use

- Production orders compete for constrained work centers
- A disruption (breakdown, shortage, absenteeism) requires rapid re-sequencing
- Changeover and campaign trade-offs need an explicit economic decision
- A new work order needs to be slotted into an existing schedule without destabilizing committed jobs
- A shift-level bottleneck change requires drum reassignment

## How It Works

1. Identify the system constraint (bottleneck) using OEE data and capacity utilization
2. Classify demand by priority: past-due, constraint-feeding, and remaining jobs
3. Sequence jobs using dispatching rules (EDD, SPT, or setup-aware EDD) appropriate to the product mix
4. Optimize changeover sequences using the setup matrix and a nearest-neighbor heuristic with 2-opt improvement
5. Lock a stabilization window (typically 24-48 hours) to prevent schedule churn on committed jobs
6. Re-plan on disruption by re-sequencing only unlocked jobs; publish the updated schedule to the MES

## Examples

- **Constraint breakdown**: a CNC machine goes down for 4 hours — identify which queued jobs can reroute to an alternate line, which must wait, and how to re-sequence the remaining queue to minimize total lateness.
- **Campaign vs. mixed-model decision**: 15 jobs across 4 product families with 45-minute inter-family changeovers — calculate the crossover point where campaign batching (fewer changeovers, more WIP) beats mixed-model (more changeovers, lower WIP) using changeover cost and carrying cost.
- **Late hot-order insertion**: sales commits a rush order with a 2-day lead time into a fully loaded week — evaluate schedule slack, identify which existing jobs can absorb a delay without missing their own due dates, and slot the hot order without breaking the frozen window.

## Core Knowledge

### Scheduling Fundamentals

**Forward vs. backward scheduling**: forward scheduling starts from material availability and schedules forward to find the earliest completion date; backward scheduling starts from the customer due date and works backward to find the latest permissible start. Default to backward scheduling to preserve flexibility and minimize WIP; switch to forward scheduling only when the backward pass reveals the latest start date is already in the past — that order needs to be expedited from today forward.

**Finite vs. infinite capacity**: MRP runs infinite-capacity planning and flags overloads for the scheduler to resolve manually. Finite-capacity scheduling (FCS) respects actual machine count, shift patterns, maintenance windows, and tooling constraints. Never trust an MRP-generated schedule as executable without running it through finite-capacity logic — MRP tells you *what* needs to be made; FCS tells you *when* it can actually be made.

**Drum-Buffer-Rope (DBR) and Theory of Constraints**: the drum is the constraint resource (highest utilization ratio, typically >85%); the buffer is a *time* buffer protecting the constraint from upstream starvation; the rope is the release mechanism limiting new work into the system to the constraint's processing rate. Subordinate every other scheduling decision to keeping the drum fed and running — a minute lost at the constraint is a minute lost for the entire plant, while a minute lost at a non-constraint costs nothing if buffer time absorbs it.

**JIT sequencing**: in mixed-model assembly, level the production sequence (heijunka) to minimize variation in component consumption — for a 3:2:1 mix of models A:B:C, sequence A-B-A-C-A-B rather than AAA-BB-C. Leveled sequencing smooths upstream demand and prevents the end-of-shift crunch where the hardest jobs get pushed to the last hour.

**Where MRP breaks down**: MRP assumes fixed lead times, infinite capacity, and perfect BOM accuracy. It fails when lead times are queue-dependent, multiple orders compete for the same constrained resource, setup times are sequence-dependent, or yield losses create variable output from fixed input. The scheduler must compensate for all four.

### Changeover Optimization

**SMED (Single-Minute Exchange of Die)**: Shigeo Shingo's framework classifies setup activities as external (can happen while the machine still runs the previous job) or internal (requires the machine stopped). Phase 1: document current setup and classify every element. Phase 2: convert internal elements to external wherever possible (pre-staging tools, pre-heating molds). Phase 3: streamline remaining internal elements (quick-release clamps, standardized die heights). Phase 4: eliminate adjustments via poka-yoke and first-piece verification. Phase 1-2 alone typically yield 40-60% setup-time reduction.

**Sequence-dependent setups**: in painting, coating, printing, and textile operations, sequence light to dark, small to large, or simple to complex to minimize cleaning between runs — a light-to-dark sequence might need only a 5-minute flush where dark-to-light requires a 30-minute full purge. Capture these in a setup matrix and feed it to the scheduling algorithm.

**Campaign vs. mixed-model**: campaign scheduling groups same-family jobs into one run (fewer changeovers, more WIP); mixed-model interleaves products (more changeovers, less WIP, shorter lead times). Lean toward campaigns when changeovers are long and expensive (>60 minutes, >$500 in scrap and lost output); lean toward mixed-model when changeovers are fast (<15 minutes) or delivery responsiveness demands short lead times. The economic crossover is where marginal changeover cost equals marginal carrying cost of additional cycle stock — compute it, don't guess.

### Bottleneck Management

**True constraint vs. where WIP piles up**: WIP accumulation in front of a work center does not necessarily mean that center is the constraint — it can indicate upstream batch-dumping, a shared-resource queue (crane, forklift, inspector), or a scheduling rule starving it downstream. The true constraint is the resource with the highest ratio of required to available hours. Verify: would adding one hour of capacity here increase total plant output? If yes, it's the constraint.

**Buffer management**: the time buffer is typically 50% of the constraint operation's production lead time. Buffer penetration <33% is green (well-protected); 33-67% is yellow (expedite late upstream work); >67% is red (immediate management attention, possible overtime). Persistent yellow trends indicate degrading upstream reliability, not noise.

**Subordination**: non-constraint resources should serve the constraint's pace, not maximize their own utilization — running a non-constraint at 100% while the constraint runs at 85% only creates excess WIP with no throughput gain.

**Shifting bottlenecks**: the constraint can move between work centers as product mix, equipment condition, or staffing changes — a bottleneck on day shift (high-setup products) may not be the bottleneck on night shift (long-run products). Monitor utilization by shift and product mix, not just weekly averages, and re-derive the drum when it shifts.

### Disruption Response

- **Machine breakdown**: assess repair time with maintenance, determine if the broken machine is the constraint, and if so activate the contingency (overtime on alternate equipment, subcontracting, re-sequencing to prioritize highest-margin jobs). If not the constraint, check buffer penetration before touching the schedule at all.
- **Material shortage**: check substitute materials, alternate BOMs, and partial-build/kitting options; escalate to purchasing; re-sequence to pull forward jobs that don't need the short material.
- **Quality hold**: held inventory is invisible to the schedule (can't ship, can't be consumed downstream) — immediately re-run the schedule excluding it, and assess alternatives (safety stock, in-process inventory, expedited replacement) if it fed a customer commitment.
- **Absenteeism**: maintain a cross-training matrix of operator × work center × certification. If the missing operator runs the constraint, reassign the best-qualified backup immediately; if a non-constraint, check whether buffer time absorbs the delay before pulling a backup from elsewhere.
- **Re-sequencing priority**: (1) protect constraint uptime above all else, (2) protect customer commitments by tier and penalty exposure, (3) minimize total changeover cost of the new sequence, (4) level labor load across remaining operators. Communicate within 30 minutes and lock the new schedule for at least 4 hours.

### Labor Management

Common shift patterns (3×8, 2×12, 4×10) trade off handovers against fatigue — 12-hour shifts show higher error rates in hours 10-12, so don't schedule critical first-piece inspections or complex changeovers there. Maintain an operator × work-center × certification matrix; scheduling feasibility depends on it, since a work order routed to equipment with no qualified operator on shift is infeasible regardless of machine availability. Quantify cross-training ROI against the throughput value of the constraint and observed absenteeism rate. Union rules on overtime seniority and mandatory rest periods (typically 8-10 hours between shifts) are hard constraints the scheduling algorithm must respect — violating one can cost far more in grievances than the production it saved.

### OEE — Overall Equipment Effectiveness

OEE = Availability × Performance × Quality. World-class is 85%+; typical discrete manufacturing runs 55-65%. Use TEEP (including all calendar time) when comparing across plants or justifying capital expansion. Address availability losses with preventive/predictive maintenance and TPM daily checks (target: unplanned downtime <5% of scheduled time). Track performance losses as actual vs. standard cycle time. Prioritize quality improvement at the constraint specifically — a 2% yield improvement there delivers the same throughput gain as a 2% capacity expansion.

### ERP/MES Interaction

Demand enters as sales orders or forecast consumption, drives the Master Production Schedule, explodes through MRP into planned orders by work center. The scheduler converts these into production orders, sequences them, and releases to the shop floor via MES; MES feedback (operation confirmations, scrap, labor booking) flows back to update ERP status and inventory. Healthy plan adherence is >90% of jobs starting within ±1 hour of scheduled start — persistent gaps mean either the scheduling parameters (setup times, run rates, yield) are wrong, or the shop floor isn't following the sequence. Compare scheduled vs. actual every shift and re-plan the remaining horizon; once operators stop trusting the schedule, it stops functioning.

## Decision Frameworks

Step-by-step algorithms for job priority sequencing, changeover sequence
optimization, disruption re-sequencing, and bottleneck identification are
in [references/decision-frameworks.md](references/decision-frameworks.md)
— load it when actually executing one of these decisions.

## Communication and Escalation

Message templates for schedule publication/change notifications, plus the
full escalation trigger table and chain, are in
[references/communication-and-escalation.md](references/communication-and-escalation.md)
— load it when writing a notification or escalation, not for routine
scheduling.

## Performance Indicators

Target/red-flag thresholds for schedule adherence, OEE, changeover time,
WIP, and yield are in
[references/performance-indicators.md](references/performance-indicators.md).

## Gotchas

- WIP piling up in front of a work center is not proof that center is the constraint — it can just as easily mean upstream batch-dumping or a shared-resource queue; verify causally before adding capacity there.
- A schedule generated straight from MRP output is not executable — it assumes infinite capacity by construction and must be run through finite-capacity logic first.
- Running a non-constraint resource at 100% utilization looks efficient locally but creates excess WIP with zero throughput gain if the constraint runs at 85% — utilization is not the same as throughput.
- The bottleneck can move between shifts on the same line, driven purely by product mix — a weekly-average utilization report will hide this and point at the wrong resource.
- 12-hour shifts show elevated error rates in hours 10-12 — don't schedule the hardest changeovers or first-piece inspections there.

## Real-world grounding

Drum-Buffer-Rope and the Theory of Constraints originate from Eliyahu Goldratt's 1984 business novel *The Goal*, which introduced the constraint-identification and subordination logic this skill's entire bottleneck-management section is built on — it remains one of the most widely cited frameworks in operations management specifically because it reframes "improve everything" into "find the one resource that actually limits output and protect it." The JIT/heijunka leveled-sequencing example is the real production-leveling practice developed within the Toyota Production System, still the standard reference point for mixed-model sequencing in lean manufacturing.

## Verification

- [ ] The current constraint was identified from utilization data, not assumed from where WIP visibly piles up
- [ ] Any MRP-generated schedule was run through finite-capacity logic before being treated as executable
- [ ] A disruption response followed the freeze → re-sequence → communicate-within-30-minutes → lock sequence
- [ ] Changeover sequencing was validated against due dates, not optimized in isolation
- [ ] Escalation followed the defined trigger table and chain rather than ad hoc judgment
