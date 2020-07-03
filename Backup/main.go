package main

import (
	"./packages/partialsuffixarrays"
	"./packages/texttobwt"
	//"./packages/readfiles"
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

// Get an approximate pattern with concurrency.
func GetSubApproximateStringWithCheckPointsC(bwt string, alphabet [] byte,
	                                         pattern string, sraName string, d int, initialTop, initialBottom int,
	                                         partialSuffixArray map [int] int,
	                                         countDict map[int] map[byte] int,
	                                         firstOccurrence map[byte] int) map[string] []int{

	N := len(pattern)
	// Declare a channel to store TopBottomUpdate.
	TopBottomUpdateChannel := make(chan [] int, len(bwt) +  2)
	TopBottomUpdateChannel <- [] int {initialTop, initialBottom, 0, N}

	// Storage sequential label in a channel.
	//previousLabel := make(chan int, 2)
	//previousLabel <- 0

	for len(TopBottomUpdateChannel) != 0 {

		//  Note that if we receive a value of a channel, the value will be removed from the channel.
		// Initial set up for loop.
		tbt := <- TopBottomUpdateChannel
		top := tbt[0]
		bottom := tbt[1]
		initialMismatch := tbt[2]
		id := tbt[3]
		pattern = pattern[0 : id]

		// Find unique symbols in an string interval.
		aprSymbolAlphabet := UniqueGene(alphabet, bwt[top : bottom + 1])

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
							newBottom = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newBottom + 1) - 1
						}
						// Update new top and new bottom.
						TopBottomUpdateChannel <- [] int {newTop, newBottom, mismatch, id - 1}
					}
					wg.Done()
				}(g)
			}
			wg.Wait()
		} else {
			// Else we break loop with result.
			TopBottomUpdateChannel <- tbt

			// Close top bottom update channel to use for range channel.
			close(TopBottomUpdateChannel)
			startPositionsChannel := make(chan int, 100)

			startPositions := map[string] []int {}
			var SP [] int

			var wgr sync.WaitGroup
			gs := runtime.NumCPU()
			wgr.Add(gs)

			for g := 0; g < gs; g++ {

				go func() {
					// Find all position of pattern matching string.
					for x := range TopBottomUpdateChannel {
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
							startPositionsChannel <- partialSuffixArray[lastIndex]+ct
						}
					}
					wgr.Done()
				}()
			}
			wgr.Wait()
			close(startPositionsChannel)
			for i := range startPositionsChannel {
				SP = append(SP, i)
			}
			startPositions[sraName] = SP
			return startPositions
		}
	}
	return map [string] [] int {}
}

// Get all position of approximate pattern in the string with concurrency.
func GetApproximatePatternMatchingWithCheckPointsC (text string, sra [][] string,
	                                                c, k, d int) map[string] [] int {

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

	totalPatterns := len(sra)
	goroutines := runtime.NumCPU()
	lastGoroutine := goroutines - 1
	stride := totalPatterns/goroutines

	startingPositions := map[string] [] int {}
	ch := make(chan map[string] []int, totalPatterns)

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

				sraName := sra[i][0]
				pattern := sra[i][1]

				positions := GetSubApproximateStringWithCheckPointsC(bwt, alphabet, pattern, sraName, d,
					                                                 initialTop, initialBottom,
					                                                 psa, countDict,
					                                                 firstOccurrence)
				ch <- positions
			}
			wg.Done()
		}(g)
	}
	wg.Wait()
	close(ch)
	for pos := range ch {
		for key, value := range pos {
			startingPositions[key] = value
		}
	}
	return startingPositions
}

func main(){
	/*
	fmt.Println("")
	fmt.Println("|============================================================================|")
	fmt.Println("|                                                                            |")
	fmt.Println("|       -----     BWT - MULTIPLE APPROXIMATE PATTERN MATCHING    -----       |")
	fmt.Println("|                                                                            |")
	fmt.Println("|============================================================================|")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("********************************* INPUT DATA *********************************")
	fmt.Println("")
	fmt.Println("")

	seq, title := ReadFiles.ReadRefseq("Refseq_Ecoli.fasta")
	fmt.Println(title)
	fmt.Println("")
	fmt.Println("Length of genome: ", len(seq))
	
	sra := ReadFiles.ReadSra("Sra_Ecoli.fasta")
	fmt.Println("Number of Sequence Read Archive: ", len(sra))
	sra = sra[0:17]
	seq += "$"

	d := 1
	c := 100
	k := 100

	fmt.Println("")
	fmt.Println("Maximum of mismatch d = ", d)
	fmt.Println("Check point arrays c = ", c)
	fmt.Println("Partial suffix arrays k = ", k)

	fmt.Println("")
	fmt.Println("")
	fmt.Println("******************************** IMPLEMENTING ********************************")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("The algorithm runtime may take a long time, please wait...")
	fmt.Println("")
	fmt.Println("")
	sta := time.Now()
	pos := GetApproximatePatternMatchingWithCheckPointsC (seq, sra, c, k, d)
	fmt.Println("Runtime of Concurrency: ------", time.Since(sta), "------")
	fmt.Println("Number of matching pattern: ", len(pos))
	ReadFiles.WritetoJSON(pos, "positions.json")

	fmt.Println("")
	fmt.Println("")
	fmt.Println("********************************** FINISHED **********************************")
	fmt.Println("")*/

	text := "ACATGCTACTTT$"
	patterns := [][] string {{"pattern 0", "ATT"}, {"pattern 1", "GCC"}, {"pattern 2", "GCTA"}, {"pattern 3", "TATT"} }
	d := 1
	sta := time.Now()
	pos := GetApproximatePatternMatchingWithCheckPointsC (text, patterns, 100, 100, d)
	fmt.Println("Concurrency: ", time.Since(sta))
	fmt.Println(pos)
}
