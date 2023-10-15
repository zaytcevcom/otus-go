package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	Text string
	Rate int
}

func Top10(str string) []string {
	str = strings.TrimSpace(str)

	if str == "" {
		return make([]string, 0)
	}

	arr := getWordsRate(strings.Fields(str))

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].Rate != arr[j].Rate {
			return arr[i].Rate > arr[j].Rate
		}

		return arr[i].Text < arr[j].Text
	})

	result := make([]string, 10)
	i := 0

	for _, v := range arr {
		result[i] = v.Text

		i++

		if i == 10 {
			break
		}
	}

	return result
}

func getWordsRate(words []string) []Word {
	wordMap := make(map[string]int, len(words))

	for _, word := range words {
		wordMap[word]++
	}

	result := make([]Word, len(wordMap))

	i := 0

	for text, count := range wordMap {
		result[i] = Word{Text: text, Rate: count}
		i++
	}

	return result
}
