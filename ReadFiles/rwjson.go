package ReadFiles
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func WritetoJSON(dataStr map[string] [] map[int] map[int] int, path string) {
	// Marshal the map into a JSON string.
	empData, err := json.Marshal(dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// write to JSON file

	jsonFile, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(empData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())
}

// Read JSON file
func ReadJSON(path string) map[string] [] map[int] map[int] int {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ", jsonFile.Name())
	// defer the closing of our jsonFile so that we can parse it later on

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string] [] map[int] map[int] int
	json.Unmarshal(byteValue, &result)

	return result
}

func WriteCountDict(dataStr map[int] map[byte] int, path string) {
	// Marshal the map into a JSON string.
	empData, err := json.Marshal(dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// write to JSON file

	jsonFile, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(empData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())
}

// Read JSON file
func ReadCountDict(path string) map[int] map[byte] int {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ", jsonFile.Name())
	// defer the closing of our jsonFile so that we can parse it later on

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[int] map[byte] int
	json.Unmarshal(byteValue, &result)

	return result
}

func WritePSA(dataStr map[int] int, path string) {
	// Marshal the map into a JSON string.
	empData, err := json.Marshal(dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// write to JSON file

	jsonFile, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(empData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())
}

// Read JSON file
func ReadPSA(path string) map[int] int {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened ", jsonFile.Name())
	// defer the closing of our jsonFile so that we can parse it later on

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[int] int
	json.Unmarshal(byteValue, &result)

	return result
}

func WritePos(dataStr [] string, path string) {
	// Marshal the map into a JSON string.
	empData, err := json.Marshal(dataStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// write to JSON file

	jsonFile, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(empData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())
}