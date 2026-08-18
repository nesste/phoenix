package main

import (
	"math"
	"sort"
)

const (
	bootstrapReplicates = 10000
	monteCarloFlips     = 100000
)

type pairedValue struct {
	FamilyID   string
	CaseID     string
	Repetition int
	X          float64
	Y          float64
	Complete   bool
}

type inference struct {
	Method     string
	Families   int
	Pairs      int
	Point      float64
	LowerBound float64
	UpperBound float64
	PWorse     *float64
}

func inferDifference(values []pairedValue) inference {
	families := pairedFamilies(values)
	if len(families) >= 20 {
		distribution := hierarchicalBootstrap(values, analysisSeed, bootstrapReplicates)
		return inference{
			Method: "10000-replicate paired hierarchical bootstrap", Families: len(families), Pairs: len(values),
			Point: meanPairDifference(values), LowerBound: quantile(distribution, 0.05),
			UpperBound: quantile(distribution, 0.95),
		}
	}
	familyMeans := unweightedFamilyMeans(values)
	null := signFlipDistribution(familyMeans, analysisSeed)
	point := mean(familyMeans)
	critical := quantile(null, 0.95)
	pWorse := lowerTail(null, point, len(familyMeans) > 16)
	method := "exact family-mean sign-flip permutation"
	if len(familyMeans) > 16 {
		method = "100000-draw seeded family-mean sign-flip permutation"
	}
	return inference{
		Method: method, Families: len(familyMeans), Pairs: len(values), Point: point,
		LowerBound: point - critical, UpperBound: point + critical, PWorse: &pWorse,
	}
}

func pairedValues(observations []observation, intervention, comparator string, include func(observation) bool, value func(observation) float64) []pairedValue {
	type pair struct {
		x *observation
		y *observation
	}
	pairs := map[string]pair{}
	for index := range observations {
		item := &observations[index]
		if !include(*item) || (item.Arm != intervention && item.Arm != comparator) {
			continue
		}
		key := item.CaseID + "\x00" + integerString(item.Repetition)
		entry := pairs[key]
		if item.Arm == intervention {
			entry.x = item
		} else {
			entry.y = item
		}
		pairs[key] = entry
	}
	values := make([]pairedValue, 0, len(pairs))
	for _, entry := range pairs {
		if entry.x == nil || entry.y == nil {
			continue
		}
		values = append(values, pairedValue{
			FamilyID: entry.x.FamilyID, CaseID: entry.x.CaseID, Repetition: entry.x.Repetition,
			X: value(*entry.x), Y: value(*entry.y), Complete: entry.x.CompleteCase && entry.y.CompleteCase,
		})
	}
	sort.Slice(values, func(left, right int) bool {
		if values[left].CaseID != values[right].CaseID {
			return values[left].CaseID < values[right].CaseID
		}
		return values[left].Repetition < values[right].Repetition
	})
	return values
}

func pairedFamilies(values []pairedValue) map[string][]pairedValue {
	families := map[string][]pairedValue{}
	for _, value := range values {
		families[value.FamilyID] = append(families[value.FamilyID], value)
	}
	return families
}

func unweightedFamilyMeans(values []pairedValue) []float64 {
	families := pairedFamilies(values)
	ids := sortedKeys(families)
	result := make([]float64, 0, len(ids))
	for _, familyID := range ids {
		cases := map[string][]float64{}
		for _, value := range families[familyID] {
			cases[value.CaseID] = append(cases[value.CaseID], value.X-value.Y)
		}
		caseMeans := make([]float64, 0, len(cases))
		for _, differences := range cases {
			caseMeans = append(caseMeans, mean(differences))
		}
		result = append(result, mean(caseMeans))
	}
	return result
}

func hierarchicalBootstrap(values []pairedValue, seed uint64, replicates int) []float64 {
	families := pairedFamilies(values)
	ids := sortedKeys(families)
	random := splitMix64{state: seed}
	distribution := make([]float64, replicates)
	for replicate := range replicates {
		var sampled []float64
		for range ids {
			familyID := ids[random.uniform(uint64(len(ids)))]
			cases := familyCases(families[familyID])
			caseIDs := sortedKeys(cases)
			for range caseIDs {
				caseID := caseIDs[random.uniform(uint64(len(caseIDs)))]
				for _, value := range cases[caseID] {
					sampled = append(sampled, value.X-value.Y)
				}
			}
		}
		distribution[replicate] = mean(sampled)
	}
	return distribution
}

func familyCases(values []pairedValue) map[string][]pairedValue {
	cases := map[string][]pairedValue{}
	for _, value := range values {
		cases[value.CaseID] = append(cases[value.CaseID], value)
	}
	return cases
}

func signFlipDistribution(familyMeans []float64, seed uint64) []float64 {
	if len(familyMeans) <= 16 {
		count := 1 << len(familyMeans)
		result := make([]float64, count)
		for mask := range count {
			var total float64
			for index, value := range familyMeans {
				if mask&(1<<index) == 0 {
					total -= value
				} else {
					total += value
				}
			}
			result[mask] = total / float64(len(familyMeans))
		}
		return result
	}
	random := splitMix64{state: seed}
	result := make([]float64, monteCarloFlips)
	for draw := range monteCarloFlips {
		var total float64
		for _, value := range familyMeans {
			if random.next()&1 == 0 {
				total -= value
			} else {
				total += value
			}
		}
		result[draw] = total / float64(len(familyMeans))
	}
	return result
}

func pairedSuccessRatio(values []pairedValue, successX, successY int, capImbalanced bool) ratioReport {
	result := ratioReport{Pairs: len(values)}
	if successX < 10 || successY < 10 {
		result.Reason = "fewer than 10 successful trials in at least one compared arm"
		return result
	}
	if capImbalanced {
		result.Reason = "cap-hit rates differ by more than 0.02"
		return result
	}
	if len(values) == 0 {
		result.Reason = "no pairing keys where both arms succeeded"
		return result
	}
	families := pairedFamilies(values)
	pointLogs := familyLogRatios(families)
	point := math.Exp(mean(pointLogs))
	upper := 0.0
	if len(families) >= 20 {
		upper = ratioBootstrapUpper(values, analysisSeed, bootstrapReplicates)
		result.Method = "10000-replicate paired hierarchical bootstrap of family log ratios"
	} else {
		null := signFlipDistribution(pointLogs, analysisSeed)
		upper = math.Exp(mean(pointLogs) + quantile(null, 0.95))
		result.Method = "exact family-log-ratio sign-flip permutation"
		if len(families) > 16 {
			result.Method = "100000-draw seeded family-log-ratio sign-flip permutation"
		}
	}
	result.Determinate = true
	result.Ratio = &point
	result.UpperBound = &upper
	return result
}

func familyLogRatios(families map[string][]pairedValue) []float64 {
	ids := sortedKeys(families)
	logs := make([]float64, 0, len(ids))
	for _, familyID := range ids {
		var left, right float64
		for _, value := range families[familyID] {
			left += value.X
			right += value.Y
		}
		if left == 0 {
			left += 0.5
		}
		if right == 0 {
			right += 0.5
		}
		logs = append(logs, math.Log(left/right))
	}
	return logs
}

func ratioBootstrapUpper(values []pairedValue, seed uint64, replicates int) float64 {
	families := pairedFamilies(values)
	ids := sortedKeys(families)
	random := splitMix64{state: seed}
	distribution := make([]float64, replicates)
	for replicate := range replicates {
		logs := make([]float64, 0, len(ids))
		for range ids {
			familyID := ids[random.uniform(uint64(len(ids)))]
			cases := familyCases(families[familyID])
			caseIDs := sortedKeys(cases)
			var left, right float64
			for range caseIDs {
				caseID := caseIDs[random.uniform(uint64(len(caseIDs)))]
				for _, value := range cases[caseID] {
					left += value.X
					right += value.Y
				}
			}
			if left == 0 {
				left += 0.5
			}
			if right == 0 {
				right += 0.5
			}
			logs = append(logs, math.Log(left/right))
		}
		distribution[replicate] = math.Exp(mean(logs))
	}
	return quantile(distribution, 0.95)
}

func sensitivity(values []pairedValue) sensitivityReport {
	complete := make([]pairedValue, 0, len(values))
	for _, value := range values {
		if value.Complete {
			complete = append(complete, value)
		}
	}
	var completePoint *float64
	if len(complete) > 0 {
		value := meanPairDifference(complete)
		completePoint = &value
	}
	caseDifferences := map[string][]float64{}
	for _, value := range values {
		caseDifferences[value.CaseID] = append(caseDifferences[value.CaseID], value.X-value.Y)
	}
	caseMeans := make([]float64, 0, len(caseDifferences))
	for _, differences := range caseDifferences {
		caseMeans = append(caseMeans, mean(differences))
	}
	return sensitivityReport{
		CompleteCase: completePoint, OneVotePerCase: mean(caseMeans),
		OneVotePerFamily: mean(unweightedFamilyMeans(values)),
	}
}

func meanPairDifference(values []pairedValue) float64 {
	differences := make([]float64, len(values))
	for index, value := range values {
		differences[index] = value.X - value.Y
	}
	return mean(differences)
}

func quantile(values []float64, probability float64) float64 {
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	if len(ordered) == 0 {
		return math.NaN()
	}
	index := int(math.Ceil(probability*float64(len(ordered)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(ordered) {
		index = len(ordered) - 1
	}
	return ordered[index]
}

func lowerTail(distribution []float64, observed float64, monteCarlo bool) float64 {
	count := 0
	for _, value := range distribution {
		if value <= observed+1e-12 {
			count++
		}
	}
	if monteCarlo {
		return float64(count+1) / float64(len(distribution)+1)
	}
	return float64(count) / float64(len(distribution))
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	var total float64
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func integerString(value int) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	index := len(digits)
	for value > 0 {
		index--
		digits[index] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[index:])
}

type splitMix64 struct{ state uint64 }

func (random *splitMix64) next() uint64 {
	random.state += 0x9e3779b97f4a7c15
	value := random.state
	value = (value ^ (value >> 30)) * 0xbf58476d1ce4e5b9
	value = (value ^ (value >> 27)) * 0x94d049bb133111eb
	return value ^ (value >> 31)
}

func (random *splitMix64) uniform(bound uint64) uint64 {
	threshold := (0 - bound) % bound
	for {
		value := random.next()
		if value >= threshold {
			return value % bound
		}
	}
}
