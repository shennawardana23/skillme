# Decision Frameworks

Step-by-step algorithms for the four recurring production-scheduling
decisions. Load this when actually executing one of these, not for
general scheduling questions.

## Job Priority Sequencing

1. Any job past-due or about to miss its due date? Schedule those first, ordered by penalty exposure (contractual > reputational > internal KPI).
2. Any job feeding a constraint whose buffer is yellow or red? Schedule those next.
3. Among the rest, apply the dispatching rule fit for the mix: EDD for high-variety short-run (minimizes maximum lateness); SPT for long-run few-product (minimizes average flow time/WIP); setup-aware EDD (swap adjacent jobs when it saves >30 minutes of setup without a due-date miss) for mixed sequence-dependent-setup environments.
4. Tie-break on customer tier, then margin.

## Changeover Sequence Optimization

1. Build the setup matrix (changeover time and cost for every product-pair transition).
2. Identify mandatory sequence constraints (allergen cross-contamination, hazmat sequencing) — these are non-negotiable, not optimizable.
3. Apply nearest-neighbor heuristic for a feasible baseline sequence.
4. Improve with 2-opt swaps, keeping any swap that reduces total changeover time without violating a due date.
5. Validate against due dates last — due-date compliance always trumps changeover optimization.

## Disruption Re-Sequencing

1. Assess the impact window and whether the disrupted resource is the constraint.
2. Freeze committed work (in-process or within 2 hours of start) unless physically impossible to continue.
3. Re-sequence remaining jobs with the job-priority framework, using updated availability.
4. Communicate the revised schedule within 30 minutes.
5. Lock it for at least 4 hours — constant re-sequencing creates more chaos than the original disruption.

## Bottleneck Identification

1. Pull utilization by work center over the trailing 2 weeks, by shift, not averaged.
2. Rank by load-hours/available-hours ratio; the top center is the suspected constraint.
3. Verify causally: would one added hour of capacity here raise total output?
4. Check for shifting patterns across shifts or product mix; if the top center changes, schedule the constraint per-shift, not on a weekly average.
5. Distinguish true constraints from artificial ones — a center overloaded only because upstream batch-dumps into it needs the upstream release rate fixed, not added downstream capacity.
