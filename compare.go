package main
import(
	"./packages/readfiles"
	"fmt"
)
/*
func Search(b [] int, a int) bool {
	for _, x := range b {
		if x == a {
			return true
		}
	}
	return false
}
*/

func MinIntSlice (sl [] int) [] int  {
	var m int
	var id int
	for i, e := range sl {
		if i==0 || e < m {
			m = e
			id = i
		}
	}
	return [] int {m, id}
}

func main() {
	// Initialization.
	m0 := map[string] [] int {}
	m1 := map[string] [] int {}
	m2 := map[string] [] int {}
	m3 := map[string] [] int {}

	data := ReadFiles.ReadJSON("PosSARs_3.json")
	fmt.Println("Length of data: ", len(data))

	// Find position of match and mismatch.
	pos := map[string] [] int {}

	for name, v := range data {
		n := len(v)
		var NM []int
		var DR []int
		var LC []int
		for i := 0; i < n; i++ {
			for dr, val := range v[i] {
				for lc, dif := range val {
					nm := dif
					LC = append(LC, lc)
					NM = append(NM, nm)
					DR = append(DR, dr)
				}
			}
		}
		min := MinIntSlice(NM)
		pos[name] = []int{DR[min[1]], LC[min[1]], min[0]}
	}

	for name, vl := range pos {

		if vl[2] == 0 {
			m0[name] = vl
		} else if vl[2] == 1 {
			m1[name] = vl
		} else if vl[2] == 2 {
			m2[name] = vl
		} else if vl[2] == 3 {
			m3[name] = vl
		}
	}

	fmt.Println("Length of match: ", len(m0))
	fmt.Println("Length of mis 1: ", len(m1))
	fmt.Println("Length of mis 2: ", len(m2))
	fmt.Println("Length of mis 3: ", len(m3))

	ReadFiles.WritePos(m0, "m0.json")
	ReadFiles.WritePos(m1, "m1.json")
	ReadFiles.WritePos(m2, "m2.json")
	ReadFiles.WritePos(m3, "m3.json")
}