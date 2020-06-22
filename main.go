package main

import (
	"./packages/partialsuffixarrays"
	"./packages/texttobwt"
	"./packages/readfiles"
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Sort string.
type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Len() int{
	return len(s)
}
func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort string BWT to receive first column.
func GetFirstColumn (bwt string) string {
	s := [] rune(bwt)
	sort.Sort(sortRunes(s))
	return string(s)
}

// Get first occurrence.
func GetFirstOccurrence(alphabet [] byte, firstColumn string) map[byte] int {

	m := map[byte] int {}
	for i:= range alphabet {
		id := bytes.IndexByte([] byte(firstColumn),  alphabet[i])
		m[alphabet[i]] = id
	}
	return m
}

// Count a symbol in string.
func GetCountSymbol(bwt string, symbol string, index int) int {
	return bytes.Count([] byte(bwt[0 : index]), [] byte(symbol))
}

// Make slice of integer.
func MakeRange(min, max, k int) [] int {
	x := make([] int, (max - min + 1)/k)
	for i := range x {
		x[i] = min + i * k
	}
	return x
}

// Construct check point arrays.
func GetCheckPointArrays (alphabet [] byte, bwt string, l int, c int) map[int] map[byte] int {
	// Make integer slice.
	indexCheckPoint := MakeRange(0, l, c)
	countDict := map[int] map[byte] int {}

	// Make a check point arrays.
	for i := range indexCheckPoint {
		counter := map[byte] int {}
		id := indexCheckPoint[i]
		instanceColumn := bwt[0:id]

		// Count each letter in string interval.
		for j := range alphabet {
			key := alphabet[j]
			value := GetCountSymbol(instanceColumn, string(key), id)
			counter[key] = value
		}
		countDict[id] = counter
	}

	return countDict
}

// Find unique symbol in a string.
func UniqueGene(alphabet [] byte, str string) [] byte{
	var unique [] byte
	for i := range alphabet {
		/* Only find index symbol in a string by byte,
		not for integer because byte < 256.
		 */
		if bytes.IndexByte([]byte(str), alphabet[i]) != -1 {
			unique = append(unique, alphabet[i])
		}
	}
	return unique
}

// Insert an item into a certain position of an array.
func InsertArr(arr [][] int, element [] int, index int) [][] int{
	arr = append(arr, [] int {})
	copy(arr[index + 1 :], arr[index :])
	arr[index] = element
	return arr
}

// Get an approximate pattern with check point array.
func GetSubApproximateStringWithCheckPointsNC(bwt string, alphabet [] byte,
	                                          pattern string, d int, initialTop, initialBottom int,
	                                          partialSuffixArray map [int] int,
	                                          countDict map[int] map[byte] int,
	                                          firstOccurrence map[byte] int, startPositions [] int) [] int{

	// Declare a TopBottomUpdate.
	TopBottomUpdate := [] [] int {{initialTop, initialBottom, 0, 0}}

	// Declare a  channel to store sequential label
	sequentialLabel := make(chan int, 2)
	sequentialLabel <- 0

	for len(TopBottomUpdate) != 0 {

		// Initial set up for loop.
		tbt := TopBottomUpdate[0]
		top := tbt[0]
		bottom := tbt[1]
		initialMismatch := tbt[2]
		label := tbt[3]

		// Find unique symbols in an string interval.
		aprSymbolAlphabet := UniqueGene(alphabet, bwt[top : bottom + 1])

		// Remove the first item in TopBottomUpdate.
		TopBottomUpdate = TopBottomUpdate [1:]

		// If loop is out of previous interval, we remove the last symbol of pattern.
		if label != <- sequentialLabel{
			pattern = pattern[0 : len(pattern) - 1]}
		// If pattern != "", we continue loop.
		if pattern != "" {

			// Back symbol of the pattern.
			symbol := pattern[len(pattern) - 1]

			// Apply loops for every new symbols in the bwt interval.
			for i := range aprSymbolAlphabet {

				aprSymbol := aprSymbolAlphabet[i]
				count := label + 1
				mismatch := initialMismatch
				newTop := top
				newBottom := bottom

					// Compare symbol of pattern to reference string.
					if aprSymbol != symbol {
						/* If symbol is not the same letter of string
						   we increase the mismatch.
						 */
						mismatch += 1
					}

					if mismatch <= d {

						// Find new top and new bottom by checking whether they are in check point arrays or not.
						if _, ok := countDict[newTop]; ok {
							newTop = firstOccurrence[aprSymbol] + countDict[newTop][aprSymbol]
						} else {
							newTop = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newTop)
						}

						if _, ok := countDict[newBottom + 1]; ok  {
							newBottom = firstOccurrence[aprSymbol] + countDict[newBottom + 1][aprSymbol] - 1
						} else {
							newBottom = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newBottom + 1) - 1
						}
						// Update new top and new bottom.
						TopBottomUpdate = append(TopBottomUpdate, [] int {newTop, newBottom, mismatch, count})
					}
			}
			sequentialLabel <- label
		// Else we break loop with result.
		} else {
			// Remove the last item of TopBottomUpdate.
			NewTopBottomUpdate := append(TopBottomUpdate, tbt)

			// Find all position of pattern matching string.
			for i := range NewTopBottomUpdate {
				x := NewTopBottomUpdate[i]
				top := x[0]
				bottom := x[1]

				// Implement loop if last index is not in partial suffix arrays.
				for j := top; j <= bottom; j++ {
					ct := 0
					lastIndex := j
					for _, ok := partialSuffixArray[lastIndex]; ok == false; {
						ct += 1
						symbol := bwt[lastIndex]
						if _, ok := countDict[lastIndex]; ok {
							lastIndex = firstOccurrence[symbol] + countDict[lastIndex][symbol]
						} else {
							lastIndex = firstOccurrence[symbol] + GetCountSymbol(bwt, string(symbol), lastIndex)
						}
						if _, ok := partialSuffixArray[lastIndex]; ok {
							break
						}
					}
					// We need to plus count for backtrack.
					startPositions = append(startPositions, partialSuffixArray[lastIndex] + ct)
				}
			}
			return startPositions
		}
	}
	return [] int {}
}

// Get an approximate pattern with concurrency.
func GetSubApproximateStringWithCheckPointsC(bwt string, alphabet [] byte,
	                                         pattern string, d int, initialTop, initialBottom int,
	                                         partialSuffixArray map [int] int,
	                                         countDict map[int] map[byte] int,
	                                         firstOccurrence map[byte] int, startPositions [] int) [] int{

	// Declare a channel to store TopBottomUpdate.
	TopBottomUpdateChannel := make(chan [] int, len(bwt) +  1)
	TopBottomUpdateChannel <- [] int {initialTop, initialBottom, 0, 0}

	// Storage sequential label in a channel.
	previousLabel := make(chan int, 2)
	previousLabel <- 0

	for len(TopBottomUpdateChannel) != 0 {

		//  Note that if we receive a value of a channel, the value will be removed from the channel.
		// Initial set up for loop.
		tbt := <- TopBottomUpdateChannel
		top := tbt[0]
		bottom := tbt[1]
		initialMismatch := tbt[2]
		label := tbt[3]

		// Find unique symbols in an string interval.
		aprSymbolAlphabet := UniqueGene(alphabet, bwt[top : bottom + 1])

		// If loop is out of previous interval, we remove the last symbol of pattern.
		// Back symbol of the pattern.
		if label != <- previousLabel {
			pattern = pattern[0 : len(pattern) - 1]
		}

		// If pattern != "", we continue loop.
		if pattern != "" {

			symbol := pattern[len(pattern) - 1]

			// Because the max approximate symbol is 5, I assign goroutines to the length of aprSymbolAlphabet.
			goroutines1 := len(aprSymbolAlphabet)

			var wg sync.WaitGroup

			wg.Add(goroutines1)

			for g := 0; g < goroutines1; g++ {
				go func(g int) {
					// Apply loops for every new symbols in the bwt interval.
					aprSymbol := aprSymbolAlphabet[g]
					count := label + 1
					mismatch := initialMismatch
					newTop := top
					newBottom := bottom

					// Compare symbol of pattern to reference string.
					if aprSymbol != symbol {
						// If symbol is not the same letter of string we increase the mismatch.
						mismatch += 1
					}

					if mismatch <= d {

						// Find new top and new bottom by checking whether they are in check point arrays or not.
						if _, ok := countDict[newTop]; ok {
							newTop = firstOccurrence[aprSymbol] + countDict[newTop][aprSymbol]
						} else {
							newTop = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newTop)
						}

						if _, ok := countDict[newBottom+1]; ok {
							newBottom = firstOccurrence[aprSymbol] + countDict[newBottom+1][aprSymbol] - 1
						} else {
							newBottom = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newBottom+1) - 1
						}
						// Update new top and new bottom.
						TopBottomUpdateChannel <- [] int {newTop, newBottom, mismatch, count}
					}
					wg.Done()
				}(g)
			}
			wg.Wait()
			previousLabel <- label
		} else {
			// Else we break loop with result.
			NewTopBottomUpdate := [][] int {tbt}

			// Close top bottom update channel to use for range channel.
			close(TopBottomUpdateChannel)

			for i := range TopBottomUpdateChannel {
				NewTopBottomUpdate = append(NewTopBottomUpdate, i)
			}
			// Find all position of pattern matching string.
			for i := range NewTopBottomUpdate {
				x := NewTopBottomUpdate[i]
				top := x[0]
				bottom := x[1]

				// Implement loop if last index is not in partial suffix arrays.
				for j := top; j <= bottom; j++ {
					ct := 0
					lastIndex := j
					for _, ok := partialSuffixArray[lastIndex]; ok == false; {
						ct += 1
						symbol := bwt[lastIndex]
						if _, ok := countDict[lastIndex]; ok {
							lastIndex = firstOccurrence[symbol] + countDict[lastIndex][symbol]
						} else {
							lastIndex = firstOccurrence[symbol] + GetCountSymbol(bwt, string(symbol), lastIndex)
						}
						if _, ok := partialSuffixArray[lastIndex]; ok {
							break
						}
					}
					// We need to plus count for backtrack.
					startPositions = append(startPositions, partialSuffixArray[lastIndex]+ct)
				}
			}
			return startPositions
			}
		}
	return [] int {}
}

// Get all position of approximate pattern in the string.
func GetApproximatePatternMatchingWithCheckPointsNC (text string, patterns [] string, c, k, d int) [][] int {

	// Initial set up for each pattern.
	var startingPositions [][] int

	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	bwt := texttobwt.GetBWT(text)
	l := len(bwt)
	initialTop := 0
	initialBottom := len(bwt) - 1

	countDict := GetCheckPointArrays (alphabet, bwt, l, c)
	firstColumn := GetFirstColumn(bwt)
	firstOccurrence := GetFirstOccurrence(alphabet, firstColumn)

	psa := partialsuffixarrays.ConstructPartialSuffixArrayNoneConcurrency(text, k)

	for i := range patterns {
		pattern := patterns[i]
		positions := GetSubApproximateStringWithCheckPointsNC(bwt, alphabet, pattern, d,
			                                                  initialTop, initialBottom,
			                                                  psa, countDict,
			                                                  firstOccurrence , [] int {})
		startingPositions = append(startingPositions, positions)
	}
	return startingPositions
}

// Get all position of approximate pattern in the string with concurrency.
func GetApproximatePatternMatchingWithCheckPointsC (text string, patterns [] string,
	                                                c, k, d int) [][] int {

	// Initial set up for each pattern.
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	bwt := texttobwt.GetBWT(text)
	l := len(bwt)
	initialTop := 0
	initialBottom := len(bwt) - 1

	countDict := GetCheckPointArrays (alphabet, bwt, l, c)
	firstColumn := GetFirstColumn(bwt)
	firstOccurrence := GetFirstOccurrence(alphabet, firstColumn)

	psa := partialsuffixarrays.ConstructPartialSuffixArrayNoneConcurrency(text, k)

	totalPatterns := len(patterns)
	goroutines := runtime.NumCPU()
	lastGoroutine := goroutines - 1
	stride := totalPatterns/goroutines

	var startingPositions [][] int
	ch := make(chan [] int, totalPatterns)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(g int) {
			start := g * stride
			end := start + stride
			if g == lastGoroutine {
				end = totalPatterns
			}
			for i := start; i < end; i++{
				pattern := patterns[i]
				positions := GetSubApproximateStringWithCheckPointsC(bwt, alphabet, pattern, d,
					                                                 initialTop, initialBottom,
					                                                 psa, countDict,
					                                                 firstOccurrence , [] int {})
				/*positions := GetSubApproximateStringWithCheckPointsNC(bwt, alphabet, pattern, d,
					initialTop, initialBottom,
					psa, countDict,
					firstOccurrence , [] int {})*/
				ch <- positions
			}
			wg.Done()
		}(g)
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		startingPositions = append(startingPositions, i)
	}
	return startingPositions
}

func main(){
	//bwt := "TCAAAAAAACTCCGCTGTACGGCTATCACCTATTGGTTTTCGGTGTTGTTCGGTTTGGCG$GGTCTGCTCCGCGTCCCGTCCTGCAGCCGTGTCGGATTTCTTGCTTTCGCTGCCATCCCCTCTTATTATGTAGGTAATGATCATAGCGCCACCCGGGGTTTACCTGACTCGCCTCCGATCGACTAGTGCGCTATCCTACCGTGTATTCCACATCATTCTCGCGCGCCCTCGTTTGCATCTAAACTCATGCAAATAGTTATCATAGGGAAGTACGTCCCTCCACCGACTGCCCGTGAGCGCCAAACTCCTGACTCATTTGGTTAGATCTTTTTTTTCACTGAATGTGGTATAACTTGTGATGGGGGCCACCCGGCTTCTTTTTACCGTTTCTGGGGAGCCTTAGTATTAAAAAAAAATAATAAAAAAAGGTAAGGTTAGTACGGTCAAAGTTTAAACATACATTATTTCGCCTAACATTCATAGGGATTTACTATTTTTTGTTACTTTTAGGCCCAGTGCTGATACGTCCGCCTTGGGCAAACGTGGCCTCGTTTTCCGTTTTTGCTTAGCAGGCTGGTCGAAAACAAACCTAGGATGGGGATCTACTAGGACAGTTACGGCCCCCGCTGTTTCAACGCTATGACGCGCCTGCCTAAGCTACTATTGCCCCCATCACGGTCCGTCCAGGATAAGTTCCCAGCCCCGTACTTTTGGCGCCCTATCATGTCTTTTTGGCCGGCTCCGGACCAGGAAGAGTGGGGCCGTGCCAGGTCCAGCCTGGTGGTTACCGCCCGCGATCCTAACGACGTTACCGGGGAACATCCGGGCGTAGTGTTGTTAGACTCGCAAATACATTGCCTACGAGTGACAAAAAAGAATGTGCGTGCCAGTGGTTAACAGACGCGTATCCCCAGGCAGCCCGTATCGTGCGGGCGACCCGGCCAGCGCGGTTATTGATTATCAAACCCGGACTGGGGCGCTGCTAGAGTTGCTTTCGACTCGTCCTACCTTGATAGTCTTACGGGGGCTCTCCGGTTGAAGCCCGGAGCTCTCAAACGCACTCGGACTTTTAATCCTTTCCATCACCCCTGCCTTTCAGGTCATCTCTCTTTCTACCGCCACAACTTAATCCCTGGGATCACCGATCCTTTGAGCCGATTTTTTAAAAGTGCTGGTGGCTTTCTCTTGAATCATTTCGAAAAAATGCTCTTGCACCTTTCCTCGGCCGGCTTATGCTCATCTTGCCTCTTGTCTCTGTCTTACCTAACCAGGTGTCAACTGCTTCCTGCGTCCTGAAAAACCGGCGTTTCTTGGCAGATCACGTCGTTTTCACCCCAGCACATAAAGAGAACCAGCCGCCTTTACCAGGGTTCCGACCTCTCCAGCCTTCCCTTGGTTCAGGTTGGCTCAATTGCGCTCGAACGGATCATCGCAGGTGATGTTGTGCATTGTGGCCCCCCTCTGAAATTCCCTCTTTTATGCGTGTACGAGGGGTCGCTATGTTATTGAGTAAGTAAGCAGTACATAATCACATAACCAGTGCAAGACCCCCGGGACGACGGAGTCTGTGTGTAAGGGCCTCCCGCGCGGGATGGGGACGTCCTGACGTTGGTGTGATCGAATGCCCCACTTTTTGCCCCACCCGTGCACCCTTTTGGAGCGGTCTGAAGGTCTCAACTTAATCCTGTGCACGCGTTTAACTAGAATGCCGTATTCCGGGAATCGGGCTCCCGTTTTAGGCCGCCTTTTTCCCACGCGTGTCGCTTTTCACCCTTCGACGTATCTTTCCCCTGCCCTCTATGACGGGATCTAGCGTATTGAGTCTGCGCAACGATTCTTGTCCCACACCACGAATCTGTTAGCAACGTCCCGGTGTGTGCAAGTCCCGGTTACGCGGGGTCTGGGACTTCCTGCATTGCAGTCTGCGTGAGTTACCGTGCTTGTCCCGGAGCATAGCAGACTGGAGTCTGAGCGTAGTGGAGGGACACCGCCTCGTGTGATCTTGCAGCCCTTGCGCCTTAGTTAGGCTCGGCTTACAATCTCCCTCTATTCGTCCCTCCCCGCCCTCCATCGCTTAGTTCTGTTTTCGTTTCCTCGCCCCGTCGGTCGCTGACTTTCTCGTAAAAAATCATAGGCACTGAATGATCATCGATTCTCCAGCAAGGCCCGTACATAACTTGACTGAAGCCTCAGTAAGCGTCCCTACTACGCCTGCCGGGCACATTCCGCATCCGAAGGACACTTCATATAATAGCCCCCCGGCCAGCAGGTTTCCGAGGGCTGATCGGTCCACTCGGCCCGACGAGGTCCGAGTTTAGCCACATATAGTGAGAGCCGGGCCAAGGAAAATGTATCTGAGGCGAGTGTGCAGGCACTGAATGGGATACTGTTCACATTACCCAGTCCGAGGGTCTAGGGGGGGTACCCTTAGACGATCCACCATCCTCTCTTGATGGCCGCGGGCGCGTTGTCGAACCGTGTAGATTCTTCGGCCGGGGGGGGGGCTTTGTAGGGCTCCATGAGAGGGCGTTTGAGATCGACCTGCGGCTGCTCGAAAACCGAAGAGTATGGGGCCGAAGAACACAACAGTCAGCCCATATCTCGGTGCTTAAAACCTGCATCATACTCGGGGGGAGGCACGTCCCCCTGCTTGTCGGGATGCGTGTGCCTTGACCCTTTGTGAACCTCTTTTGACCATCCGGGCGGTCAAACGCCCTTTTCATATATTTTTCTAGACCAATCGCTCGGGGGTTTCCCGAGATCCGCTCGAGCGGCTTTCTTGCGTCCGTACATATGGCCCGAATTTCATGGTGGGTGGTCCCGGCGTGGGCGTACGGAAACCGTCCCGGGGTGGTTGCTCCATCGGCGGGATAGGGCTTGAGTTATAAAATATCAGTTCACCAGGCCTAGGTACGGGCGCTCGCGAGACCGAAAAGAACGGCAGACGATTAGCAACCCACAAAGCGAGGCTGCCCTTAGGGAACGCACCGCGGGGAAACAGAAGATGAATCTACTACACCAGGACCCTACAATACCAGCGGGGATCCACTGGAACCAATCCCGCTCTCTGACGTCGTAAAGACGCACAGACGAGGCAGCGCCTCTATGACATGTACGTCAGGATTAGTTTTTTTTTTGTG"
	//bwt := "ATTC$AA" // BATATA$
	/*text := "ACATGCTACTTT$"
	patterns := [] string {"ATT", "GCC", "GCTA", "TATT"}
	d := 1*/
	//pattern := patterns[0]

	words, err := ReadFiles.ScanWords("dataset_304_10 (4).txt")
	if err != nil {
		panic(err)
	}
	var patterns [] string

	text := words[0] + "$"
	lenData := len(words)

	for i := 1; i < lenData - 1; i++ {
		pattern := words[i]
		patterns = append(patterns, pattern)
	}

	d := 2

	/*bwt := texttobwt.GetBWT(text)
	l := len(bwt)
	initialtop := 0
	initialbottom := len(bwt) - 1

	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	firstColumn := GetFirstColumn(bwt)

	// fmt.Println(firstColumn)
	firstOccurence := GetFirstOccurrence(alphabet, firstColumn)
	// fmt.Println(firstOccurence)
	// fmt.Println(getCountSymbol(bwt, "A", 6))
	// fmt.Println(makeRange(2, 23, 3))

	countDict := GetCheckPointArrays (alphabet, bwt, l, 1)
	//fmt.Println("index check point: ", indexCheckPoint)
	//fmt.Println(countDict)
	// fmt.Println(uniqueGene(alphabet, bwt))

	psa := partialsuffixarrays.ConstructPartialSuffixArrayNoneConcurrency(text, 2)
	//fmt.Println(psa)
	var sai [] byte
	for k := range psa {
		sai = append(sai, byte(k))
	}*/
	// fmt.Println("Suffix Array Index: ", sai)

	/*pos := GetSubApproximateStringWithCheckPointsNC(bwt, alphabet, pattern, d, initialtop, initialbottom,
		psa,  countDict,
		firstOccurence , [] int {})
	fmt.Println(pos)

	start := time.Now()
	positions := GetSubApproximateStringWithCheckPointsC(bwt, alphabet, pattern, d, initialtop, initialbottom,
		                                               psa,  countDict,
		                                               firstOccurence , [] int {})
	fmt.Println("None Concurrency: ", time.Since(start))
	fmt.Println(positions)*/

	start := time.Now()
	positions := GetApproximatePatternMatchingWithCheckPointsNC (text, patterns, 100, 100, d)
	fmt.Println("None Concurrency: ", time.Since(start))
	fmt.Println(len(positions))

	sta := time.Now()
	pos := GetApproximatePatternMatchingWithCheckPointsC (text, patterns, 100, 100, d)
	fmt.Println("Concurrency: ", time.Since(sta))
	fmt.Println(len(pos))

	/*
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	bwt := texttobwt.GetBWT(text)
	l := len(bwt)
	indexCheckPoint, countDict := GetCheckPointArrays (alphabet, bwt, l, 100)
	fmt.Println(len(indexCheckPoint))
	fmt.Println(len(countDict))*/

}
