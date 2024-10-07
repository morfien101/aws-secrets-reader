package main

// This application is designed to be used as a secrets collecting agent for https://github.com/morfien101/launch
// Therefore it will use the requirements set out in the secrets collection section.
//   All Errors will come in STDERR
//   All Secrets will be passed out in JSON
//   Only the secrets will be passed to STDOUT
// The only exceptions to these rules are for the help menu and the version display as these are not expected
// to be used when the binary is actually doing its job.

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	// set the version here
	// The build version is expected to be set by the build script.
	version = "development"

	// Set the flags here to be used later in the code.
	flagProfile     = flag.String("aws-profile", "", "AWS Profile to use. Blank by default and omitted.")
	flagRegion      = flag.String("region", "eu-west-1", "AWS region to use. eu-west-1 by default.")
	flagSecretKey   = flag.String("secret", "", "The key to use when collecting the secret.")
	flagUpperCase   = flag.Bool("upper-case", false, "Attempt to uppercase all the returned keys")
	flagPrependKeys = flag.String("prepend-with", "", "Prepend the returned keys with given string. Upper casing happens after this is applied.")
	flagHelp        = flag.Bool("h", false, "Help menu.")
	flagVersion     = flag.Bool("v", false, "Shows the version.")
	flagFormat      = flag.String("format", "json", "The format to return the secrets in. Default is JSON. Supported values are: json, yaml, env and shell_export")
)

func main() {
	flag.Parse()
	if *flagHelp {
		// Show a nice help menu then exit
		fmt.Println("Collects secrets from AWS Secrets manager.")
		fmt.Println("version:", version)
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// If there is no secret path we can no continue.
	if *flagSecretKey == "" {
		writeSTDERR(fmt.Sprintln("secret is not set. I can not continue like this. Exiting..."))
		os.Exit(1)
	}

	// At this point we are ready to go get some secrets.
	// Obviously there can still be issues but we have everything we need to try now.
	// Setup the session
	sess, err := awsSession(awsOptions(*flagRegion, *flagProfile))
	if err != nil {
		writeSTDERR(fmt.Sprintf("there was an error creating the AWS Session. Error: %s", err))
		os.Exit(1)
	}
	// Attempt to collect secret
	rawSecret, err := collectSecret(sess, *flagSecretKey)
	if err != nil {
		// This error is already pretty well formatted.
		writeSTDERR(fmt.Sprintln(err))
		os.Exit(1)
	}

	// PostProcessing
	// Append prefix: Some people like to prepend the secrets keys with a tag to know which ones they collected. Like AWS_SM_<value>
	// Other reasons might be to allow automatic collection from applications if they have the correct keys.
	secretMap, err := postProcess(rawSecret, *flagPrependKeys, *flagUpperCase)
	if err != nil {
		writeSTDERR(fmt.Sprintln(err))
		os.Exit(1)
	}
	formattedSecret, err := format(secretMap, *flagFormat)
	if err != nil {
		writeSTDERR(fmt.Sprintln(err))
		os.Exit(1)
	}

	// print out the secrets
	fmt.Println(formattedSecret)
}

// awsOption will check to see if profile is set and optionally add it in if required.
func awsOptions(region, profile string) session.Options {
	options := session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}
	if profile != "" {
		// We need to tell go to load the configuration files as users might be sourcing profiles from roles.
		os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
		options.Profile = profile
	}
	return options
}

// awsSession just wraps the function used to create the session. This is mainly for future work as AWS sometimes change this.
func awsSession(options session.Options) (*session.Session, error) {
	return session.NewSessionWithOptions(options)
}

// wrtieSTDERR is like an easy to use function to send data to STDERR.
func writeSTDERR(s string) {
	os.Stderr.WriteString(s)
}

// collectSecret will go off and collect the secrets from AWS and return the string
// and an error if there was one.
func collectSecret(awsSession *session.Session, key string) (string, error) {

	svc := secretsmanager.New(awsSession)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		// The errors from the AWS SDK are actually pretty useful. So we should make effort to use them.
		return "", fmt.Errorf("error retrieving the secret from %s. Error: %s", key, err.Error())
	}

	// The returned text from AWS Secret manager is already in JSON and a key value pair setup.
	// This will work great with Launch.
	return *result.SecretString, nil
}

func postProcess(input string, prepend string, upperCase bool) (map[string]string, error) {

	// First we need to convert out input in to a map so we can edit it.
	editMe := map[string]string{}
	err := json.Unmarshal([]byte(input), &editMe)
	if err != nil {
		return map[string]string{}, fmt.Errorf("there was an error reading the collected secrets before editing. Error: %s", err)
	}

	// To append we need to actually create new keys with the appended string and then remove the old key
	if prepend != "" {
		newMap := map[string]string{}
		for key, value := range editMe {
			newMap[fmt.Sprintf("%s%s", prepend, key)] = value
		}
		editMe = newMap
	}

	// Now cycle through the map again and upper case everything if required.
	// ToUpper returns s with all Unicode letters mapped to their upper case.
	// https://golang.org/pkg/strings/#ToUpper
	if upperCase {
		newMap := map[string]string{}
		for key, value := range editMe {
			newMap[strings.ToUpper(key)] = value
		}
		editMe = newMap
	}

	return editMe, nil
}

func format(output map[string]string, format string) (string, error) {
	switch format {
	case "json":
		jsonOutput, err := json.Marshal(output)
		if err != nil {
			return "", fmt.Errorf("there was an error while converting the updated values back to JSON. Error: %s", err)
		}
		return string(jsonOutput), nil
	case "yaml":
		yamlOutput, err := yaml.Marshal(output)
		if err != nil {
			return "", fmt.Errorf("there was an error while converting the updated values back to YAML. Error: %s", err)
		}
		return string(yamlOutput), nil
	case "env":
		secrets := ""
		for key, value := range output {
			secrets += fmt.Sprintf("%s=%s\n", key, value)
		}
		data, err := godotenv.Unmarshal(secrets)
		if err != nil {
			return "", fmt.Errorf("there was an error while converting the updated values back to dotenv. Error: %s", err)
		}
		content, err := godotenv.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("there was an error while converting the updated values back to dotenv. Error: %s", err)
		}
		return content, nil
	case "shell_export":
		secrets := ""
		for key, value := range output {
			secrets += fmt.Sprintf("%s=%s\n", key, value)
		}
		data, err := godotenv.Unmarshal(secrets)
		if err != nil {
			return "", fmt.Errorf("there was an error while converting the updated values back to dotenv. Error: %s", err)
		}
		exportedSecrets := ""
		for key, value := range data {
			exportedSecrets += fmt.Sprintf("export %s=%s\n", key, value)
		}
		return exportedSecrets, nil
	default:
		return "", fmt.Errorf("the format %s is not supported", format)
	}
}
