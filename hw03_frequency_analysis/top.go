package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordFreq struct {
	name string
	freq int
}

func Top10(str string) []string {
	sliceWord := strings.Fields(str)

	freqMap := make(map[string]int)

	for i := 0; i < len(sliceWord); i++ {
		freqMap[sliceWord[i]] += 1
	}

	wordFreqSlice := make([]wordFreq, 0, len(freqMap))
	for k, v := range freqMap {
		wordFreqSlice = append(wordFreqSlice, wordFreq{k, v})
	}
	sort.Slice(wordFreqSlice, func(i, j int) bool {
		if wordFreqSlice[i].freq > wordFreqSlice[j].freq {
			return true
		} else if wordFreqSlice[i].freq == wordFreqSlice[j].freq {
			if wordFreqSlice[i].name < wordFreqSlice[j].name {
				return true
			} else {
				return false
			}
		} else {
			return false
		}

	})

	wordSlice := make([]string, 0, len(freqMap))
	for _, val := range wordFreqSlice {
		wordSlice = append(wordSlice, val.name)
	}

	if len(wordSlice) >= 10 {
		return wordSlice[0:10]
	} else {
		return wordSlice
	}

}
