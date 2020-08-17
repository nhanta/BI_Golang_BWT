package main

import (
	"./packages/readfiles"
	"./packages/memusage"
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strconv"
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
func GetFirstOccurrence(firstColumn string) map[byte] int {

	m := map[byte] int {}
	al := [] byte {'A', 'T', 'G', 'C'}
	var w sync.WaitGroup
	var mu sync.Mutex
	w.Add(4)

	for g := 0; g < 4; g++ {
		go func(g int) {
			id := bytes.IndexByte([] byte(firstColumn),  al[g])
			mu.Lock()
			m[al[g]] = id
			mu.Unlock()
			w.Done()
		}(g)
	}
	w.Wait()
	m['$'] = 0
	return m
}

// Count a symbol in string.
func GetCountSymbol(bwt string, symbol string, index int) int {
	return bytes.Count([] byte(bwt[0 : (index + 1)]), [] byte(symbol))
}

//Calculate differences with c = k = 1.
func CalculateDMaxTime(pattern string, firstOccurrence map[byte] int,
	countDict map[int] map[byte] int, l, N  int) map[int] int {
	top := 0
	bottom := l - 1

	z := 0
	D := map[int] int {}
	for i := 0; i < N; i++ {
		symbol := pattern[i]
		top = firstOccurrence[symbol] + countDict[top - 1][symbol]
		bottom = firstOccurrence[symbol] + countDict[bottom][symbol] - 1
		if top > bottom {
			top = 0
			bottom = l - 1
			z += 1
		}
		D[i] = z
	}
	return D
}

//Calculate differences.
func CalculateD(bwt, pattern string, firstOccurrence map[byte] int,
	countDict map[int] map[byte] int, l, N int) map[int] int {
	top := 0
	bottom := l - 1

	z := 0
	D := map[int] int {}

	for i := 0; i < N; i++ {
		symbol := pattern[i]
		if _, ok := countDict[top - 1]; ok {
			top = firstOccurrence[symbol] + countDict[top - 1][symbol]
		} else {
			top = firstOccurrence[symbol] + GetCountSymbol(bwt, string(symbol), top - 1)
		}

		if _, ok := countDict[bottom]; ok {
			bottom = firstOccurrence[symbol] + countDict[bottom][symbol] - 1
		} else {
			bottom = firstOccurrence[symbol] + GetCountSymbol(bwt, string(symbol), bottom) - 1
		}

		if top > bottom {
			top = 0
			bottom = l - 1
			z += 1
		}
		D[i] = z
	}
	return D
}

// Get an approximate pattern with concurrency.
func GetSubApproximateStringWithCheckPointsC(bwt string,
	                                         pattern string, sraName string, initialTop, initialBottom int,
	                                         partialSuffixArray map [int] int, D map[int] int, dfr, N int,
	                                         countDict map[int] map[byte] int,
	                                         firstOccurrence map[byte] int) map[string] map[int] int{

	startPositionsChannel := make(chan [] int, 100)
	startPositions := map[string] map[int] int {}
	SP := map[int] int {} // Map of start position.
	// Declare a channel to store TopBottomUpdate.
	TopBottomUpdateChannel := make(chan [] int, len(bwt) +  2)
	posChannel := make(chan [] int, 100) // Channel of positions.
	TopBottomUpdateChannel <- [] int {initialTop, initialBottom, dfr, N}

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
				d := tbt[2]
				id := tbt[3]
				// Find unique symbols in an string interval.
				aprSymbolAlphabet := [] byte {'A', 'T', 'G', 'C'}

				// If pattern != "", we continue loop.
				if id != 0 {
					pattern := pattern[0:id]
					symbol := pattern[id - 1]

					// Because the max approximate symbol is 5, I assign goroutines to the length of aprSymbolAlphabet.
					goroutines := 4
					var wg sync.WaitGroup

					wg.Add(goroutines)

					for g := 0; g < goroutines; g++ {
						if reflect.DeepEqual(tbt, [] int {0, 0, 0, 0}) == true {
							TopBottomUpdateChannel <- [] int {0, 0, 0, 0}
							break
						}
						go func(g int) {
							// Apply loops for every new symbols in the bwt interval.
							aprSymbol := aprSymbolAlphabet[g]
							nD := d
							newTop := top
							newBottom := bottom

							if nD >= D[id - 1] {
								// Compare symbol of pattern to reference string.
								if aprSymbol != symbol {
									// If symbol is not the same letter of string we reduce differences.
									nD -= 1
								}
								// Find new top and new bottom by checking whether they are in check point arrays or not.
								if _, ok := countDict[newTop - 1]; ok {
									newTop = firstOccurrence[aprSymbol] + countDict[newTop - 1][aprSymbol]
								} else {
									newTop = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newTop - 1)
								}

								if _, ok := countDict[newBottom]; ok {
									newBottom = firstOccurrence[aprSymbol] + countDict[newBottom][aprSymbol] - 1
								} else {
									newBottom = firstOccurrence[aprSymbol] + GetCountSymbol(bwt, string(aprSymbol), newBottom) - 1
								}
								if newTop <= newBottom && nD >= 0 {
									// Update new top and new bottom.
									TopBottomUpdateChannel <- []int{newTop, newBottom, nD, id - 1}
								}
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
					for i:= 0; i <= gs; i++ {
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
					NM := dfr - x[2] // Number of mismatches or gaps.

					// Implement loop if last index is not in partial suffix arrays.
					for j := top; j <= bottom; j++ {
						ct := 0
						lastIndex := j
						for _, ok := partialSuffixArray[lastIndex]; ok == false; {
							ct += 1
							symbol := bwt[lastIndex]
							if _, ok := countDict[lastIndex - 1]; ok {
								lastIndex = firstOccurrence[symbol] + countDict[lastIndex - 1][symbol]
							} else {
								lastIndex = firstOccurrence[symbol] + GetCountSymbol(bwt, string(symbol), lastIndex - 1)
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

// Get an approximate pattern with concurrency.
func GetSubApproximateStringWithCheckPointsMaxTime(bwt string,
	pattern string, sraName string, initialTop, initialBottom int,
	partialSuffixArray map [int] int, D map[int] int, dfr, N int,
	countDict map[int] map[byte] int,
	firstOccurrence map[byte] int) map[string] map[int] int{

	startPositionsChannel := make(chan [] int, 100)
	startPositions := map[string] map[int] int {}
	SP := map[int] int {} // Map of start position.
	// Declare a channel to store TopBottomUpdate.
	TopBottomUpdateChannel := make(chan [] int, len(bwt) +  2)
	posChannel := make(chan [] int, 100) // Channel of positions.
	TopBottomUpdateChannel <- [] int {initialTop, initialBottom, dfr, N}

	gs := runtime.NumCPU()

	var wt sync.WaitGroup

	wt.Add(gs)

	for g := 0; g < gs; g++ {

		go func() {

			for tbt := range TopBottomUpdateChannel{
				//fmt.Println("after: ", len(TopBottomUpdateChannel))
				//  Note that if we receive a value of a channel, the value will be removed from the channel.
				// Initial set up for loop.
				top := tbt[0]
				bottom := tbt[1]
				d := tbt[2]
				id := tbt[3]
				// Find unique symbols in an string interval.
				aprSymbolAlphabet := [] byte {'A', 'T', 'G', 'C'}

				// If pattern != "", we continue loop.
				if id != 0 {
					pattern := pattern[0 : id]
					symbol := pattern[id - 1]

					// Because the max approximate symbol is 4, I assign goroutines to the length of aprSymbolAlphabet.
					goroutines := 4
					var wg sync.WaitGroup

					wg.Add(goroutines)

					for g := 0; g < goroutines; g++ {
						if reflect.DeepEqual(tbt, [] int {0, 0, 0, 0}) == true {
							TopBottomUpdateChannel <- [] int {0, 0, 0, 0}
							break
						}
						go func(g int) {
							// Apply loops for every new symbols in the bwt interval.
							aprSymbol := aprSymbolAlphabet[g]
							nD := d
							newTop := top
							newBottom := bottom

							if nD >= D[id - 1] {
								// Compare symbol of pattern to reference string.
								if aprSymbol != symbol {
									// If symbol is not the same letter of string we reduce differences.
									nD -= 1
								}
								// Find new top and new bottom.
								newTop = firstOccurrence[aprSymbol] + countDict[newTop - 1][aprSymbol]
								newBottom = firstOccurrence[aprSymbol] + countDict[newBottom][aprSymbol] - 1
								if newTop <= newBottom && nD >= 0 {
									// Update new top and new bottom.
									TopBottomUpdateChannel <- []int{newTop, newBottom, nD, id - 1}
								}
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
					for i:= 0; i <= gs ; i++ {
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
					NM := dfr - x[2] // Number of mismatches of gaps.

					// Implement loop if last index is not in partial suffix arrays.
					for j := top; j <= bottom; j++ {
						// We need to plus count for backtrack.
						startPositionsChannel <- [] int {partialSuffixArray[j], NM}
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
func GetApproximatePatternMatchingWithCheckPointsC (bwt, rBWT string, sra [][] string, psa map[int] int,
	                                                countDict map[int] map[byte] int, rCD map[int] map[byte] int,
	                                                firstOccurrence map[byte] int,
	                                                c, k, d int,
	                                                ) map[string] [] map[int] map[int] int {

	initialTop := 0
	l := len(bwt)
	initialBottom := l - 1
	totalPatterns := len(sra)
	goroutines := runtime.NumCPU()
	lastGoroutine := goroutines - 1
	stride := totalPatterns/goroutines

	startingPositions := map[string] map[int] map[int] int {}
	rStartingPositions := map[string] map[int] map[int] int {}
	start := map[string] [] map[int] map[int] int {}
	ch := make(chan map[string] map[int] int, totalPatterns)
	rCh := make(chan map[string] map[int] int, totalPatterns)

	// Collect sra names.
	var name [] string
	for i := 0; i < totalPatterns; i++{
		name = append(name, sra[i][0])
	}

	var wg sync.WaitGroup
	wg.Add(goroutines)

	if c == 1 && k == 1 {
		for g := 0; g < goroutines; g++ {
			go func(g int) {
				start := g * stride
				end := start + stride
				if g == lastGoroutine {
					end = totalPatterns
				}
				for i := start; i < end; i++{

					sraName := sra[i][0]
					pattern1 := sra[i][1]
					pattern2 := sra[i][2]
					N := len(pattern1)
					D1 := CalculateDMaxTime(pattern1, firstOccurrence, rCD, l, N)
					D2 := CalculateDMaxTime(pattern2, firstOccurrence, rCD, l, N)
					var w sync.WaitGroup
					w.Add(2)

					go func() {
						pos1 := GetSubApproximateStringWithCheckPointsMaxTime(bwt, pattern1, sraName,
							initialTop, initialBottom,
							psa, D1, d, N, countDict, firstOccurrence)
						ch <- pos1
						w.Done()
					}()
					go func() {
						pos2 := GetSubApproximateStringWithCheckPointsMaxTime(bwt, pattern2, sraName,
							initialTop, initialBottom,
							psa, D2, d, N, countDict, firstOccurrence)
						rCh <- pos2
						w.Done()
					}()
					w.Wait()
				}
				wg.Done()
			}(g)
		}
		wg.Wait()
	} else {
		for g := 0; g < goroutines; g++ {
			go func(g int) {
				start := g * stride
				end := start + stride
				if g == lastGoroutine {
					end = totalPatterns
				}
				for i := start; i < end; i++{

					sraName := sra[i][0]
					pattern1 := sra[i][1]
					pattern2 := sra[i][2]
					N := len(pattern1)
					D1 := CalculateD(rBWT, pattern1, firstOccurrence, rCD, l, N)
					D2 := CalculateD(rBWT, pattern2, firstOccurrence, rCD, l, N)

					var w sync.WaitGroup
					w.Add(2)

					go func() {
						pos1 := GetSubApproximateStringWithCheckPointsC(bwt, pattern1, sraName,
							initialTop, initialBottom,
							psa, D1, d, N, countDict, firstOccurrence)
						ch <- pos1
						w.Done()
					}()
					go func() {
						pos2 := GetSubApproximateStringWithCheckPointsC(bwt, pattern2, sraName,
							initialTop, initialBottom,
							psa, D2, d, N, countDict, firstOccurrence)
						rCh <- pos2
						w.Done()
					}()
					w.Wait()
				}
				wg.Done()
			}(g)
		}
		wg.Wait()
	}

	close(ch)
	close(rCh)
	// Get positions of pros.
	for pos := range ch {
		for key, value := range pos {
			startingPositions[key] = map[int] map[int] int {
				0 : value,
			}
		}
	}
	//Get positions of reverse.
	for pos := range rCh {
		for key, value := range pos {
			rStartingPositions[key] = map[int] map[int] int {
				16 : value,
			}
		}
	}

	// Get all positions.
	for i := 0; i < totalPatterns; i++ {
		key := name[i]
		var value [] map[int] map [int] int
		if val, ok := startingPositions[key]; ok {
			value = append(value, val)
		}
		if v, ok := rStartingPositions[key]; ok {
			value = append(value, v)
		}
		if len(value) != 0 {
			start[key] = value
		}
	}
	return start
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

	runtime.GC()
	bwt := ReadFiles.ReadText("BWT_SARs_CoV_2.txt")
	rBWT := ReadFiles.ReadText("rBWT_SARs_CoV_2.txt")
	sra := ReadFiles.ReadSra("Sra_SARs_CoV_2.fasta")

	var c, k, d int

	fmt.Print("Step of Partial Suffix Arrays, Enter c: ")
	_, err := fmt.Scanf("%d\n", &c)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("Step of Check Point Arrays, Enter k: ")
	_, er := fmt.Scanf("%d\n", &k)
	if er != nil {
		fmt.Println(er)
	}

	fmt.Print("Number of difference, Enter d: ")
	_, ee := fmt.Scanf("%d\n", &d)
	if ee != nil {
		fmt.Println(ee)
	}

	psaPath := "PSA_SARs_CoV_2_" + strconv.Itoa(c) + ".json"
	countPath := "Count_SARs_CoV_2_" + strconv.Itoa(k)+ ".json"
	rCountPath := "rCount_SARs_CoV_2_" + strconv.Itoa(k)+ ".json"

	psa := ReadFiles.ReadPSA(psaPath)
	cd := ReadFiles.ReadCountDict(countPath)
	rCD := ReadFiles.ReadCountDict(rCountPath)
	firstColumn := GetFirstColumn(bwt)
	firstOccurrence := GetFirstOccurrence(firstColumn)

	fmt.Println("")
	fmt.Println("Length of genome: ", len(bwt))
	fmt.Println("Number of Sequence Read Archive: ", len(sra))

	fmt.Println("")
	fmt.Println("******************************** IMPLEMENTING ********************************")
	fmt.Println("")
	fmt.Println("The algorithm runtime may take a long time, please wait...")
	fmt.Println("")

	sta := time.Now()
	pos := GetApproximatePatternMatchingWithCheckPointsC (bwt, rBWT, sra, psa, cd, rCD, firstOccurrence, c, k, d)
	fmt.Println("Runtime of Concurrency: ------", time.Since(sta), "------")
	fmt.Println("Number of matching pattern: ", len(pos))
	posPath := "PosSARs_" + strconv.Itoa(d) + ".json"
	ReadFiles.WritetoJSON(pos, posPath)

	fmt.Println("Memory of matching: ")
	MemUsage.PrintMemUsage()

	fmt.Println("")
	fmt.Println("********************************** FINISHED **********************************")
	fmt.Println("")
}
