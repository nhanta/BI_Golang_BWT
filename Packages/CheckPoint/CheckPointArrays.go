package CheckPoint

import (
	"bytes"
	"runtime"
	"sync"
)

// Make slice of integer.
func MakeRange(min, max, k int) [] int {
	x := make([] int, 1 + (max - min)/k)
	for i := range x {
		x[i] = min + i * k
	}
	return x
}

// Count a symbol in string.
func GetCountSymbol(bwt string, symbol string, index int) int {
	return bytes.Count([] byte(bwt[0 : (index + 1)]), [] byte(symbol))
}

// Construct check point arrays.
func GetCheckPointArrays (alphabet [] byte, bwt string, l int, c int) map[int] map[byte] int {

	indexCheckPoint := MakeRange(0, l - 1, c) // Make integer slice.
	total:= len(indexCheckPoint)
	countDict := make(chan map[int] map[byte] int, total)
	C := map[int] map[byte] int {}

	goroutines := runtime.NumCPU()
	lastGoroutine := goroutines - 1
	stride := total/goroutines

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(g int) {
			start := g * stride
			end := start + stride
			if g == lastGoroutine {
				end = total
			}
			for i := start; i < end; i++{

				id := indexCheckPoint[i]

				gs := 5
				var wgr sync.WaitGroup
				var mu sync.Mutex
				wgr.Add(gs)

				// Count each letter in string interval.
				counter := map[byte] int {}
				for j := 0; j < gs; j++ {
					go func(j int) {
						mu.Lock()
						key := alphabet[j]
						value := GetCountSymbol(bwt, string(key), id)
						counter[key] = value
						wgr.Done()
						mu.Unlock()
					}(j)
				}
				wgr.Wait()
				countDict <- map[int] map[byte] int {
					id: counter,
				}
			}
			wg.Done()
		}(g)
	}
	wg.Wait()
	close(countDict)
	for x := range countDict {
		for key, value := range x {
			C[key] = value
		}
	}
	return C
}
