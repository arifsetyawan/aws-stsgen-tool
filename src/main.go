package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
)

func main() {

	var theRes map[string]interface{}

	command, value, errArgs := argsh()
	if errArgs != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", errArgs)
		os.Exit(1)
	}

	switch command {
	case "install":
		installSerialNumber(value)
	case "set":
		result, err := requestStsCredential(value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		json.Unmarshal(result.([]byte), &theRes)

		credentials := theRes["Credentials"].(map[string]interface{})
		accessKey := credentials["AccessKeyId"].(string)
		secretKey := credentials["SecretAccessKey"].(string)
		sessionToken := credentials["SessionToken"].(string)

		profile := "cred-mfa"

		fmt.Println(credentials)

		configureLocalCredential(accessKey, secretKey, sessionToken, profile)
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

func requestStsCredential(token string) (interface{}, error) {
	binary, lookErr := exec.LookPath("aws")
	if lookErr != nil {
		return nil, lookErr
	}

	serial_id := readSerialNumber()

	cmd := exec.Command(binary, "sts", "get-session-token", "--serial-number", serial_id, "--token-code", token, "--duration-seconds", "129600", "--profile", "wallex")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func configureLocalCredential(accessKey string, secretKey string, sessionToken string, profile string) bool {

	binary, lookErr := exec.LookPath("aws")
	if lookErr != nil {
		return false
	}
	exec.Command(binary, "configure", "set", "aws_access_key_id", accessKey, "--profile", profile).CombinedOutput()
	exec.Command(binary, "configure", "set", "aws_secret_access_key", secretKey, "--profile", profile).CombinedOutput()
	exec.Command(binary, "configure", "set", "aws_session_token", sessionToken, "--profile", profile).CombinedOutput()

	return true

}

func installSerialNumber(value string) bool {

	usr, err := user.Current()
	if err != nil {
		return false
	}

	d1 := []byte(value)
	f, err := os.Create(usr.HomeDir+"/.awsstsgen/.serial_id")
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

	dat, err := ioutil.ReadFile(usr.HomeDir+"/.awsstsgen/.serial_id")
	fmt.Print(string(dat))

	return string(dat)
}