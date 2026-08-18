package main

import (
	"fmt"
	"math"
)

type armPair struct{ intervention, comparator string }

var comparedArms = []armPair{{"C", "A"}, {"C", "B"}, {"C", "D"}, {"D", "E"}}

func analyze(summary scheduledSummary, observations []observation, incomplete int) analysisReport {
	report := analysisReport{
		V: 1, Tranche: summary.Tranche, AnalysisSeed: analysisSeed,
		ScheduleDigest: summary.ScheduleDigest, WorldBuild: summary.WorldBuild,
		RuntimeVersion: summary.Configuration.RuntimeVersion, Model: summary.Configuration.Model,
		GraderDigest: summary.Configuration.GraderDigest, ArmBDocument: summary.Configuration.ArmBDocumentDigest,
		Status: "analysis_complete", StatusReasons: []string{}, IncompletePairings: incomplete,
		Arms: armReports(observations), Measures: measureReport(observations),
	}
	globalReasons := indeterminateReasons(summary, report.Arms, incomplete)
	report.StatusReasons = append(report.StatusReasons, globalReasons...)
	report.Claims = []claimReport{
		directClaim(observations, report.Arms, len(globalReasons) > 0),
		headlineClaim(observations, report.Arms, len(globalReasons) > 0),
		frontierClaim(observations, report.Arms, len(globalReasons) > 0),
		teachingClaim(observations, report.Arms, len(globalReasons) > 0),
	}
	if len(globalReasons) > 0 {
		report.Status = "indeterminate"
		return report
	}
	report.Status = report.Claims[1].Decision
	for _, claim := range report.Claims {
		if claim.ID != "headline_value" && (claim.Decision == "reject_surface" || claim.Decision == "remove_frontier" || claim.Decision == "replace_teaching_refusals") {
			report.Status = "component_rejected"
			report.StatusReasons = append(report.StatusReasons, claim.ID+": "+claim.Reason)
		}
	}
	return report
}

func armReports(observations []observation) map[string]armReport {
	reports := map[string]armReport{}
	for _, arm := range phase1Arms {
		reports[arm] = armReport{}
	}
	for _, item := range observations {
		report := reports[item.Arm]
		report.Assigned++
		if item.ITTSuccess {
			report.Successes++
			report.MeanTokensToSuccess += item.Tokens
			report.MeanActsToSuccess += float64(item.Acts)
		} else {
			report.MeanFailedTrialTokens += item.Tokens
			report.MeanFailedTrialActs += float64(item.Acts)
		}
		if item.Unresolved {
			report.Unresolved++
		}
		if item.CapHit {
			report.CapHits++
		}
		report.TotalTokens += item.Tokens
		report.TotalCostUSD += item.CostUSD
		report.TotalTurns += item.Turns
		report.WallTimeMS += item.WallTimeMS
		reports[item.Arm] = report
	}
	for arm, report := range reports {
		if report.Assigned > 0 {
			report.SuccessRate = float64(report.Successes) / float64(report.Assigned)
			report.UnresolvedRate = float64(report.Unresolved) / float64(report.Assigned)
			report.CapHitRate = float64(report.CapHits) / float64(report.Assigned)
			report.MeanWallTimeMS = report.WallTimeMS / float64(report.Assigned)
		}
		if report.Successes > 0 {
			report.MeanTokensToSuccess /= float64(report.Successes)
			report.MeanActsToSuccess /= float64(report.Successes)
		}
		failures := report.Assigned - report.Successes
		if failures > 0 {
			report.MeanFailedTrialTokens /= float64(failures)
			report.MeanFailedTrialActs /= float64(failures)
		}
		reports[arm] = report
	}
	return reports
}

func indeterminateReasons(summary scheduledSummary, arms map[string]armReport, incomplete int) []string {
	var reasons []string
	if summary.Status == "indeterminate" || summary.StopReason != "" {
		reasons = append(reasons, "scheduled run stopped or was marked indeterminate")
	}
	if incomplete > 0 {
		reasons = append(reasons, fmt.Sprintf("%d incomplete pairing keys were discarded and counted unresolved", incomplete))
	}
	for _, arm := range phase1Arms {
		if arms[arm].UnresolvedRate > 0.05 {
			reasons = append(reasons, fmt.Sprintf("arm %s unresolved rate %.6f exceeds 0.05", arm, arms[arm].UnresolvedRate))
		}
	}
	for _, pair := range comparedArms {
		difference := math.Abs(arms[pair.intervention].UnresolvedRate - arms[pair.comparator].UnresolvedRate)
		if difference > 0.02 {
			reasons = append(reasons, fmt.Sprintf("%s-%s unresolved-rate imbalance %.6f exceeds 0.02", pair.intervention, pair.comparator, difference))
		}
	}
	return reasons
}

func directClaim(observations []observation, arms map[string]armReport, globallyIndeterminate bool) claimReport {
	values := successPairs(observations, "C", "A", classFilter("direct"))
	inference := inferDifference(values)
	claim := baseClaim("direct_no_tax", "C", "A", "direct", "ITT task success", inference, values)
	if globallyIndeterminate {
		claim.Decision, claim.Reason = "indeterminate", "tranche-level missing-data or stopping rule"
	} else if inference.PWorse < 0.05 {
		claim.Decision, claim.Reason = "reject_surface", "one-sided test finds C worse than A at p < 0.05"
	} else {
		claim.Decision, claim.Reason = "no_harm_signal", "harm test did not reject; this is not a noninferiority claim"
	}
	_ = arms
	return claim
}

func headlineClaim(observations []observation, arms map[string]armReport, globallyIndeterminate bool) claimReport {
	all := func(observation) bool { return true }
	values := successPairs(observations, "C", "B", all)
	inference := inferDifference(values)
	claim := baseClaim("headline_value", "C", "B", "all eight classes", "ITT task success", inference, values)
	deadEndValues := pairedValues(observations, "C", "B", all, func(item observation) float64 { return boolFloat(item.DeadEnd) })
	deadEndInference := inferDifference(deadEndValues)
	claim.DeadEndDifference = inferencePointer(deadEndInference)
	capValues := pairedValues(observations, "C", "B", all, func(item observation) float64 { return boolFloat(item.TimeoutOrCap) })
	capInference := inferDifference(capValues)
	claim.CapHitDifference = inferencePointer(capInference)
	pairedTokens := pairedValues(observations, "C", "B", func(item observation) bool { return item.ITTSuccess }, func(item observation) float64 { return item.Tokens })
	capImbalanced := math.Abs(arms["C"].CapHitRate-arms["B"].CapHitRate) > 0.02
	ratio := pairedSuccessRatio(pairedTokens, arms["C"].Successes, arms["B"].Successes, capImbalanced)
	claim.CostRatio = &ratio
	if globallyIndeterminate {
		claim.Decision, claim.Reason = "indeterminate", "tranche-level missing-data or stopping rule"
		return claim
	}
	if inference.LowerBound > 0 {
		claim.Decision, claim.Reason = "capability_pass", "one-sided 95% lower bound for C-B ITT success is greater than zero"
		return claim
	}
	efficiency := inference.LowerBound > -0.05 && ratio.Determinate && ratio.UpperBound != nil && *ratio.UpperBound <= 0.80 && deadEndInference.UpperBound < 0.05 && capInference.UpperBound < 0.05
	if efficiency {
		claim.Decision, claim.Reason = "efficiency_only", "capability did not pass; all four precommitted efficiency-only conditions hold"
	} else if capImbalanced {
		claim.Decision, claim.Reason = "indeterminate", "cap-hit imbalance makes the required paired-success token ratio indeterminate"
	} else {
		claim.Decision, claim.Reason = "fail", "capability lower bound is not above zero and the efficiency-only conjunction does not hold"
	}
	return claim
}

func frontierClaim(observations []observation, arms map[string]armReport, globallyIndeterminate bool) claimReport {
	filter := func(item observation) bool {
		return item.Class == "cascade" || item.Class == "far_discovery" || item.Class == "recovery" || item.Class == "temptation"
	}
	values := successPairs(observations, "C", "D", filter)
	inference := inferDifference(values)
	claim := baseClaim("frontier_value", "C", "D", "cascade, far_discovery, recovery, and temptation", "ITT task success", inference, values)
	if globallyIndeterminate {
		claim.Decision, claim.Reason = "indeterminate", "tranche-level missing-data or stopping rule"
	} else if inference.Point <= 0 || inference.PWorse < 0.05 {
		claim.Decision, claim.Reason = "remove_frontier", "C-D point estimate is nonpositive or the one-sided harm test rejects"
	} else {
		claim.Decision, claim.Reason = "retain_frontier", "C-D point estimate is positive and the one-sided harm test does not reject"
	}
	_ = arms
	return claim
}

func teachingClaim(observations []observation, arms map[string]armReport, globallyIndeterminate bool) claimReport {
	filter := classFilter("recovery")
	recoveryValues := pairedValues(observations, "D", "E", filter, func(item observation) float64 { return boolFloat(item.Recovery) })
	recoveryInference := inferDifference(recoveryValues)
	claim := baseClaim("teaching_refusal_value", "D", "E", "recovery", "triggered recovery with successful grade", recoveryInference, recoveryValues)
	downstreamValues := successPairs(observations, "D", "E", filter)
	downstream := inferDifference(downstreamValues)
	claim.DownstreamSuccess = inferencePointer(downstream)
	if globallyIndeterminate {
		claim.Decision, claim.Reason = "indeterminate", "tranche-level missing-data or stopping rule"
	} else if recoveryInference.Point <= 0 || recoveryInference.PWorse < 0.05 || downstream.LowerBound < -0.05 {
		claim.Decision, claim.Reason = "replace_teaching_refusals", "recovery effect is nonpositive, harm test rejects, or downstream ITT lower bound is below -0.05"
	} else {
		claim.Decision, claim.Reason = "retain_teaching_refusals", "recovery effect is positive without a harm signal or downstream guardrail failure"
	}
	_ = arms
	return claim
}

func baseClaim(id, intervention, comparator, population, endpoint string, result inference, values []pairedValue) claimReport {
	lower, upper, pWorse := result.LowerBound, result.UpperBound, result.PWorse
	return claimReport{
		ID: id, Intervention: intervention, Comparator: comparator, Population: population, Endpoint: endpoint,
		Families: result.Families, Pairs: result.Pairs, Method: result.Method, PointEstimate: result.Point,
		LowerBound: &lower, UpperBound: &upper, PWorse: &pWorse, Sensitivity: sensitivity(values),
	}
}

func inferencePointer(result inference) *inferenceReport {
	lower, upper, pWorse := result.LowerBound, result.UpperBound, result.PWorse
	return &inferenceReport{Point: result.Point, LowerBound: &lower, UpperBound: &upper, PWorse: &pWorse}
}

func successPairs(observations []observation, intervention, comparator string, include func(observation) bool) []pairedValue {
	return pairedValues(observations, intervention, comparator, include, func(item observation) float64 { return boolFloat(item.ITTSuccess) })
}

func classFilter(class string) func(observation) bool {
	return func(item observation) bool { return item.Class == class }
}

func boolFloat(value bool) float64 {
	if value {
		return 1
	}
	return 0
}

func measureReport(observations []observation) descriptiveMeasures {
	var acts, wrong, deadEnds, capHits, frontierShown, frontierTaken, frontierPassed int
	for _, item := range observations {
		acts += item.Acts
		wrong += item.WrongActs
		if item.DeadEnd {
			deadEnds++
		}
		if item.TimeoutOrCap {
			capHits++
		}
		frontierShown += item.FrontierShown
		frontierTaken += item.FrontierTaken
		frontierPassed += item.FrontierPassed
	}
	report := descriptiveMeasures{}
	if acts > 0 {
		report.WrongVerbRate = float64(wrong) / float64(acts)
	}
	if len(observations) > 0 {
		report.DeadEndRate = float64(deadEnds) / float64(len(observations))
		report.TimeoutOrCapHitRate = float64(capHits) / float64(len(observations))
	}
	if frontierShown > 0 {
		take := float64(frontierTaken) / float64(frontierShown)
		passed := float64(frontierPassed) / float64(frontierShown)
		report.FrontierTakeRate = &take
		report.FrontierTakeAndSucceedRate = &passed
	}
	return report
}
