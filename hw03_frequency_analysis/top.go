package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordFreqType struct {
	name string
	freq int
}

const numberOutputElements = 10

func Top10(str string) []string {
	if len(str) == 0 {
		return nil
	}

	sliceWord := strings.Fields(str)

	if len(sliceWord) == 0 {
		return nil
	}

	freqMap := make(map[string]int)

	for i := 0; i < len(sliceWord); i++ {
		freqMap[sliceWord[i]]++
	}

	wordFreqSlice := make([]wordFreqType, 0, len(freqMap))
	for k, v := range freqMap {
		wordFreqSlice = append(wordFreqSlice, wordFreqType{k, v})
	}
	sort.Slice(wordFreqSlice, func(i, j int) bool {
		if wordFreqSlice[i].freq > wordFreqSlice[j].freq {
			return true
		} else if wordFreqSlice[i].freq == wordFreqSlice[j].freq {
			if wordFreqSlice[i].name < wordFreqSlice[j].name {
				return true
			}
		}
		return false
	})

	if len(wordFreqSlice) >= numberOutputElements {
		wordFreqSlice = wordFreqSlice[:numberOutputElements]
	}

	wordSlice := make([]string, 0, len(wordFreqSlice))
	for _, val := range wordFreqSlice {
		wordSlice = append(wordSlice, val.name)
	}

	return wordSlice
}
