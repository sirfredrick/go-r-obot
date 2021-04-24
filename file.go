package main

import (
"io/ioutil"
"fmt"
)

func main() {
	ioutil.WriteFile("data.txt", []byte(fmt.Sprint(336)), 0666)

}
