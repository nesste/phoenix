package main

import (
	"fmt"
	"sort"
)

const (
	phase1ScheduleSeed = uint64(20260817)
	phase1Repetitions  = 3
)

var phase1Arms = []string{"A", "B", "C", "D", "E"}

type launchSchedule struct {
	V           int             `json:"v"`
	Tranche     string          `json:"tranche"`
	Seed        uint64          `json:"seed"`
	Repetitions int             `json:"repetitions"`
	Arms        []string        `json:"arms"`
	Entries     []scheduleEntry `json:"entries"`
}

type scheduleEntry struct {
	LaunchIndex  int    `json:"launch_index"`
	FamilyBlock  int    `json:"family_block"`
	PairingIndex int    `json:"pairing_index"`
	FamilyID     string `json:"family_id"`
	CaseID       string `json:"case_id"`
	Repetition   int    `json:"repetition"`
	Arm          string `json:"arm"`
}

type scheduleCase struct {
	CaseID   string
	FamilyID string
}

type pairingKey struct {
	caseID     string
	repetition int
}

func generateSchedule(tranche string, cases []scheduleCase, seed uint64, repetitions int, arms []string) (launchSchedule, error) {
	if tranche == "" || len(cases) == 0 || repetitions <= 0 || len(arms) < 2 {
		return launchSchedule{}, fmt.Errorf("schedule requires a tranche, cases, repetitions, and at least two arms")
	}
	families, familyIDs, err := groupScheduleCases(cases)
	if err != nil {
		return launchSchedule{}, err
	}
	sequences, err := williamsSequences(arms)
	if err != nil {
		return launchSchedule{}, err
	}
	random := splitMix64{state: seed}
	shuffleStrings(familyIDs, &random)
	schedule := launchSchedule{
		V: 1, Tranche: tranche, Seed: seed, Repetitions: repetitions,
		Arms: append([]string(nil), arms...), Entries: []scheduleEntry{},
	}
	pairingIndex := 0
	for familyBlock, familyID := range familyIDs {
		caseIDs := append([]string(nil), families[familyID]...)
		shuffleStrings(caseIDs, &random)
		pairings := familyPairings(caseIDs, repetitions)
		shufflePairings(pairings, &random)
		for localIndex, pairing := range pairings {
			sequence := sequences[(localIndex+familyBlock)%len(sequences)]
			for _, arm := range sequence {
				schedule.Entries = append(schedule.Entries, scheduleEntry{
					LaunchIndex: len(schedule.Entries), FamilyBlock: familyBlock, PairingIndex: pairingIndex,
					FamilyID: familyID, CaseID: pairing.caseID, Repetition: pairing.repetition, Arm: arm,
				})
			}
			pairingIndex++
		}
	}
	return schedule, nil
}

func groupScheduleCases(cases []scheduleCase) (map[string][]string, []string, error) {
	families := map[string][]string{}
	seen := map[string]struct{}{}
	for _, item := range cases {
		if item.CaseID == "" || item.FamilyID == "" {
			return nil, nil, fmt.Errorf("schedule case and family identifiers are required")
		}
		if _, duplicate := seen[item.CaseID]; duplicate {
			return nil, nil, fmt.Errorf("schedule case %q is duplicated", item.CaseID)
		}
		seen[item.CaseID] = struct{}{}
		families[item.FamilyID] = append(families[item.FamilyID], item.CaseID)
	}
	familyIDs := make([]string, 0, len(families))
	for familyID := range families {
		familyIDs = append(familyIDs, familyID)
		sort.Strings(families[familyID])
	}
	sort.Strings(familyIDs)
	return families, familyIDs, nil
}

func familyPairings(caseIDs []string, repetitions int) []pairingKey {
	pairings := make([]pairingKey, 0, len(caseIDs)*repetitions)
	for _, caseID := range caseIDs {
		for repetition := 0; repetition < repetitions; repetition++ {
			pairings = append(pairings, pairingKey{caseID: caseID, repetition: repetition})
		}
	}
	return pairings
}

func williamsSequences(arms []string) ([][]string, error) {
	seen := map[string]struct{}{}
	for _, arm := range arms {
		if arm == "" {
			return nil, fmt.Errorf("williams arm names are required")
		}
		if _, duplicate := seen[arm]; duplicate {
			return nil, fmt.Errorf("williams arm %q is duplicated", arm)
		}
		seen[arm] = struct{}{}
	}
	base := make([]int, len(arms))
	base[0] = 0
	for index := 1; index < len(base); index++ {
		if index%2 == 1 {
			base[index] = (index + 1) / 2
		} else {
			base[index] = len(base) - index/2
		}
	}
	rows := make([][]string, 0, len(arms)*2)
	for shift := range arms {
		row := make([]string, len(arms))
		for index, value := range base {
			row[index] = arms[(value+shift)%len(arms)]
		}
		rows = append(rows, row)
	}
	if len(arms)%2 == 1 {
		for index := 0; index < len(arms); index++ {
			rows = append(rows, reverseStrings(rows[index]))
		}
	}
	return rows, nil
}

func reverseStrings(values []string) []string {
	reversed := make([]string, len(values))
	for index := range values {
		reversed[len(values)-1-index] = values[index]
	}
	return reversed
}

type splitMix64 struct {
	state uint64
}

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

func shuffleStrings(values []string, random *splitMix64) {
	for index := len(values) - 1; index > 0; index-- {
		other := int(random.uniform(uint64(index + 1)))
		values[index], values[other] = values[other], values[index]
	}
}

func shufflePairings(values []pairingKey, random *splitMix64) {
	for index := len(values) - 1; index > 0; index-- {
		other := int(random.uniform(uint64(index + 1)))
		values[index], values[other] = values[other], values[index]
	}
}
