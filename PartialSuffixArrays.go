package PSA

import(
	"sort"
	"sync"
)

type Slice struct {
	sort.StringSlice
	idx [] int
}

func (s Slice) Swap(i, j int) {
	s.StringSlice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n [] string) *Slice {
	s := &Slice{StringSlice: sort.StringSlice(n), idx: make([]int, len(n))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}
/*func makeRange(min, max int) [] int {
	x := make([] int, max - min + 1)
	for i := range x {
		x[i] = min + i
	}
	return x
}*/

func constructPartialSuffixArrayConcurrency(goroutines int, text string, c int) map [int] int {

	l := len(text)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var textList []string

	psa := map[int]int{}
	KeyValueMap := map[int]int{}
	SuffixMapCh := make(chan map[int]int, l)
	KeyValueCh := make(chan map[int]int, l)

	// Creating suffix array
	for i := 0; i < l; i++ {
		textList = append(textList, string(text[i:]))
	}

	s := NewSlice(textList)
	sort.Sort(s)
	suffixArr := s.idx

	for i := 0; i < l; i++ {
		SuffixMapCh <- map[int]int{
			i: suffixArr[i],
		}
	}
	close(SuffixMapCh)

	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			for SuffixMap := range SuffixMapCh {
				for key, v := range SuffixMap {
					KeyValueCh <- map[int]int{
						v: key,
					}

				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(KeyValueCh)

	// order := makeRange(0, l)
	k := (l - 1) / c

	for KeyValue := range KeyValueCh {
		for v, key := range KeyValue {
			mu.Lock()
			KeyValueMap[v] = key
			mu.Unlock()
		}
	}

	for i := 0; i <= k; i++ {
		SAValue := i * c
		id := KeyValueMap[SAValue]
		psa[id] = SAValue


	}
	return psa
}

// Construct Partial Suffix Array with Concurrency.

func constructPartialSuffixArrayNoneConcurrency(text string, c int) map [int] int{

	l := len(text)
	var textList [] string
	KeyValueMap := map[int] int {}
	suffixMap := map[int] int {}
	psa := map[int] int {}

	// Creating suffix array
	for i := 0; i < l; i++ {
		textList = append(textList, string(text[i:]))
	}

	s := NewSlice(textList)
	sort.Sort(s)
	suffixArr := s.idx

	for i:= 0; i < l; i++ {
		suffixMap[i] = suffixArr[i]
	}
	for key, v := range suffixMap {
		KeyValueMap[v] = key
	}

	// order := makeRange(0, l)
	k := (l-1)/c
	for i := 0; i <= k; i++ {
		SAValue := i * c
		id :=  KeyValueMap[SAValue]
		psa[id] = SAValue
	}

	return psa
}