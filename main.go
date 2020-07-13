package main

import (
	"./packages/memusage"
	"./packages/readfiles"
	"bytes"
	"fmt"
	"reflect"
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
	firstOccurrence map[byte] int) map[string] map[int] int{

	N := len(pattern)
	startPositionsChannel := make(chan [] int, 100)
	startPositions := map[string] map[int] int {}
	SP := map[int] int {} // Map of start position.
	// Declare a channel to store TopBottomUpdate.
	TopBottomUpdateChannel := make(chan [] int, len(bwt) +  2)
	posChannel := make(chan [] int, 100) // Channel of positions.
	TopBottomUpdateChannel <- [] int {initialTop, initialBottom, 0, N}


	gs := runtime.NumCPU()

	var wt sync.WaitGroup
	wt.Add(gs)

	for g := 0; g < gs; g++ {

		go func() {

			for tbt := range TopBottomUpdateChannel {

				//  Note that if we receive a value of a channel, the value will be removed from the channel.
				// Initial set up for loop.
				top := tbt[0]
				bottom := tbt[1]
				initialMismatch := tbt[2]
				id := tbt[3]
				// Find unique symbols in an string interval.
				aprSymbolAlphabet := UniqueGene(alphabet, bwt[top:bottom+1])

				// If pattern != "", we continue loop.
				if id != 0 {
					pattern := pattern[0:id]
					symbol := pattern[len(pattern)-1]

					// Because the max approximate symbol is 5, I assign goroutines to the length of aprSymbolAlphabet.
					goroutines1 := len(aprSymbolAlphabet)
					var wg sync.WaitGroup

					wg.Add(goroutines1)

					for g := 0; g < goroutines1; g++ {
						if reflect.DeepEqual(tbt, [] int {0, 0, 0, 0}) == true {
							TopBottomUpdateChannel <- [] int {0, 0, 0, 0}
							break
						}
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
									newBottom = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newBottom+1) - 1
								}
								// Update new top and new bottom.
								TopBottomUpdateChannel <- []int{newTop, newBottom, mismatch, id - 1}
							}
							wg.Done()
						}(g)
					}
					wg.Wait()

				} else {
					// Else we break loop with result.
					if reflect.DeepEqual(tbt, [] int {0, 0, 0, 0}) == false {
						posChannel <- tbt
					}
				}
				if len(TopBottomUpdateChannel) == 0 {
					for i:= 0; i < gs; i++ {
						TopBottomUpdateChannel <- [] int {0, 0, 0, 0}
					}
					break
				}
			}
			wt.Done()
		}()
	}
	wt.Wait()
	if len(posChannel) != 0 {
		// Close position channel to use for range channel.
		close(posChannel)
		var wgr sync.WaitGroup
		grs := runtime.NumCPU()
		wgr.Add(grs)

		for g := 0; g < grs; g++ {

			go func() {
				// Find all position of pattern matching string.
				for x := range posChannel {
					top := x[0]
					bottom := x[1]
					NM := x[2] // Number of mismatch.

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
						startPositionsChannel <- [] int {partialSuffixArray[lastIndex]+ct, NM}
					}
				}
				wgr.Done()
			}()
		}
		wgr.Wait()
		close(startPositionsChannel)
		for x := range startPositionsChannel {
			SP[x[0]] = x[1]
		}
		startPositions[sraName] = SP
	}
	return startPositions
}

// Get all position of approximate pattern in the string with concurrency.
func GetApproximatePatternMatchingWithCheckPointsC (bwt string, sra [][] string, psa map[int] int,
	                                                countDict map[int] map[byte] int,
	                                                firstOccurrence map[byte] int,
	                                                d int) map[string] map[int] int {

	runtime.GC()
	initialTop := 0
	initialBottom := len(bwt) - 1
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	totalPatterns := len(sra)
	goroutines := runtime.NumCPU()
	lastGoroutine := goroutines - 1
	stride := totalPatterns/goroutines

	startingPositions := map[string] map[int] int {}
	ch := make(chan map[string] map[int] int, totalPatterns)

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
					                                                 psa, countDict, firstOccurrence)

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
	fmt.Println("Memory of matching: ")
	MemUsage.PrintMemUsage()
	return startingPositions
}

func main(){

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

	fmt.Println("Information of Computing System")
	fmt.Println("OS\t\t", runtime.GOOS)
	fmt.Println("ARCH \t\t", runtime.GOARCH)
	fmt.Println("CPUs \t", runtime.NumCPU())
	fmt.Println("")

	bwt := ReadFiles.ReadText("BWT_Ecoli.txt")
	sra := ReadFiles.ReadSra("Sra_Ecoli.fasta")
	psa := ReadFiles.ReadPSA("PSA_Ecoli1.json")
	cd := ReadFiles.ReadCountDict("Count_Ecoli1.json")
	firstColumn := GetFirstColumn(bwt)
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	firstOccurrence := GetFirstOccurrence(alphabet, firstColumn)
	d := 4

	fmt.Println("")
	fmt.Println("Length of genome: ", len(bwt))
	fmt.Println("Number of Sequence Read Archive: ", len(sra))

	fmt.Println("")
	fmt.Println("******************************** IMPLEMENTING ********************************")
	fmt.Println("")
	fmt.Println("The algorithm runtime may take a long time, please wait...")
	fmt.Println("")

	sta := time.Now()
	pos := GetApproximatePatternMatchingWithCheckPointsC (bwt, sra, psa, cd, firstOccurrence, d)
	fmt.Println("Runtime of Concurrency: ------", time.Since(sta), "------")
	fmt.Println("Number of matching pattern: ", len(pos))
	ReadFiles.WritetoJSON(pos, "positions_0.json")

	fmt.Println("")
	fmt.Println("********************************** FINISHED **********************************")
	fmt.Println("")
}
