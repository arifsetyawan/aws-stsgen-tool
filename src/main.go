package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

func main() {

	command, value, errArgs := argsh()
	if errArgs != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", errArgs)
		os.Exit(1)
	}

	switch command {
	case "install":
		installSerialNumber(value)
	}

}


func argsh() (string, string, error) {
	if len(os.Args) < 3 {
		return "", "", fmt.Errorf("%s", "Incomplete Command")
	}
	command := os.Args[1:][0]
	value := os.Args[1:][1]

	fmt.Println(command);
	fmt.Println(value);

	return command, value, nil
}


func installSerialNumber(value string) bool {

	usr, err := user.Current()
	if err != nil {
		return false
	}

	d1 := []byte(value)
	f, err := os.Create(usr.HomeDir+"/.walstsgen/.serial_id")
	defer f.Close()
	n2, err := f.Write(d1)
	fmt.Printf("wrote %d bytes\n", n2)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func readSerialNumber() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	dat, err := ioutil.ReadFile(usr.HomeDir+"/.walstsgen/.serial_id")
	fmt.Print(string(dat))

	return string(dat)
}