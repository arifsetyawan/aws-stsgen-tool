package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
)

func main() {

	var theRes map[string]interface{}

	command, errArgs := argsh()
	if errArgs != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", errArgs)
		os.Exit(1)
	}

	switch command {
	case "install":
		install()
	case "set":

		// Get Current User
		usr, err := user.Current()
		if err != nil {
			fmt.Errorf("%v", err)
		}

		viper.SetConfigName("config") // name of config file (without extension)
		viper.AddConfigPath(usr.HomeDir+"/.awsstsgen/")
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		var mfaToken string

		// Get token prompt
		fmt.Print("Please input current MFA token: ")
		_, err = fmt.Scanln(&mfaToken)
		if err != nil {
			panic("Require MFA Token")
		}

		result, err := requestStsCredential(mfaToken, viper.Get("mfa-arn").(string), viper.Get("base-profile").(string))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		json.Unmarshal(result.([]byte), &theRes)

		credentials := theRes["Credentials"].(map[string]interface{})
		accessKey := credentials["AccessKeyId"].(string)
		secretKey := credentials["SecretAccessKey"].(string)
		sessionToken := credentials["SessionToken"].(string)

		profile := viper.Get("target-profile").(string)

		fmt.Println("accessKey: "+accessKey)
		fmt.Println("secretKey: "+secretKey)
		fmt.Println("sessionToken: "+sessionToken)

		if configureLocalCredential(accessKey, secretKey, sessionToken, profile) == true {
			fmt.Println("\nCREDENTIAL UPDATED")
		}
	}

}

func argsh() (string, error) {

	if len(os.Args) < 2 {
		return "", fmt.Errorf("%s", "Incomplete Command")
	}

	command := os.Args[1:][0]
	return command, nil
}

func requestStsCredential(token string, serial string, baseProfile string) (interface{}, error) {
	binary, lookErr := exec.LookPath("aws")
	if lookErr != nil {
		return nil, lookErr
	}

	cmd := exec.Command(binary, "sts", "get-session-token", "--serial-number", serial, "--token-code", token, "--duration-seconds", "129600", "--profile", baseProfile)
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
	_, err := exec.Command(binary, "configure", "set", "aws_access_key_id", accessKey, "--profile", profile).CombinedOutput()
	_, err = exec.Command(binary, "configure", "set", "aws_secret_access_key", secretKey, "--profile", profile).CombinedOutput()
	_, err = exec.Command(binary, "configure", "set", "aws_session_token", sessionToken, "--profile", profile).CombinedOutput()
	if err != nil {
		fmt.Errorf("Error Execute aws configure ", err)
		return false
	}

	return true
}

func install() {

	var resultArn, resultBaseProfile, resultTargetProfile string

	fmt.Print("Your mfa arn : ")
	_, err := fmt.Scanln(&resultArn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Print("Your base credential name in ~/.aws/credential (default): ")
	_, err = fmt.Scanln(&resultBaseProfile)
	if err != nil {
		resultBaseProfile = "default"
	}

	fmt.Print("Your target credential name (default-mfa): ")
	_, err = fmt.Scanln(&resultTargetProfile)
	if err != nil {
		resultTargetProfile = "default-mfa"
	}

	preparedConfig := map[string]interface{}{
		"mfa-arn": resultArn,
		"base-profile": resultBaseProfile,
		"target-profile": resultTargetProfile,
	}

	preparedConfigData, _ := json.MarshalIndent(preparedConfig, "", " ")

	// Get Current User
	usr, err := user.Current()
	if err != nil {
		fmt.Errorf("%v", err)
	}

	// Checking if the .awsstsgen is exist, otherwise create it.
	if _, err := os.Stat(usr.HomeDir+"/.awsstsgen/"); os.IsNotExist(err) {
		_ = os.Mkdir(usr.HomeDir+"/.awsstsgen", 0777)
	}

	// save config file
	errSaveConfig := ioutil.WriteFile(usr.HomeDir+"/.awsstsgen/config.json", preparedConfigData, 0644)
	if errSaveConfig != nil {
		fmt.Errorf("%v", errSaveConfig)
		panic(errSaveConfig)
	}

	fmt.Println("\n=======================================================================")
	fmt.Println("Install complete!\nyou can check your config at ~/.awsstsgen/config.json")

}