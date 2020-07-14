package main
import(
	"./packages/readfiles"
	"fmt"
)

func Search(b [] int, a int) bool {
	for _, x := range b {
		if x == a {
			return true
		}
	}
	return false
}

func main() {
	// Initialization.
	var m [] string
	var m1 [] string
	var m2 [] string
	var m3 [] string

	data := ReadFiles.ReadJSON("positions_3.json")

	// Find position of match and mismatch.
	pos := map[string] map[int] int {}
	for k, v := range data {
		n := len(v)
		m := map[int] int {}

		for i := 0; i < n; i++{
			for _, val := range v[i] {
				for a, b := range val {
					m[a] = b
				}
			}
		}
		pos[k] = m
	}

	for key, value := range pos {
		var vl [] int
		for _, v := range value {
			vl = append(vl, v)
		}
		if Search(vl, 0) == true {
			m = append(m, key[1:])
		} else if Search(vl, 1) == true {
			m1 = append(m1, key[1:])
		} else if Search(vl, 2) == true {
			m2 = append(m2, key[1:])
		} else if Search(vl, 3) == true {
			m3 = append(m3, key[1:])
		}
	}

	fmt.Println("Length of match: ", len(m))
	fmt.Println("Length of mis 1: ", len(m1))
	fmt.Println("Length of mis 2: ", len(m2))
	fmt.Println("Length of mis 3: ", len(m3))

	ReadFiles.WritePos(m, "m.json")
	ReadFiles.WritePos(m1, "m1.json")
	ReadFiles.WritePos(m2, "m2.json")
	ReadFiles.WritePos(m3, "m3.json")

}