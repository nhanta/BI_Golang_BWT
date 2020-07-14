package ReadFiles

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func WriteText(text string, name string) {
	f, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(text)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ReadText(path string) string{
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


	b, err := ioutil.ReadAll(file)
	return string(b)
}