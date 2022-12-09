package utils

import (
	"fmt"
	"os"
)

func ToFile(str string, name string) {
	file, err := os.Create(name + ".json")

	if err != nil {
		fmt.Println("Unable to create file:", err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(str)
	if err != nil {
		fmt.Println("Unable to create file:", err.Error())
	}
}
