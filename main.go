package main

// This application is designed to be used as a secrets collecting agent for https://github.com/morfien101/launch
// Therefore it will use the requirements set out in the secrets collection section.
//   All Errors will come in STDERR
//   All Secrets will be passed out in JSON
//   Only the secrets will be passed to STDOUT
// The only exceptions to these rules are for the help menu and the version display as these are not expected
// to be used when the binary is actually doing its job.

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
	// set the version here
	version = "0.0.1"

	// Set the flags here to be used later in the code.
	flagProfile   = flag.String("aws-profile", "", "AWS Profile to use. Blank by defalt and ommited.")
	flagRegion    = flag.String("region", "eu-west-1", "AWS region to use. eu-west-1 by default.")
	flagSecretKey = flag.String("secret", "", "The key to use when collecting the secret.")
	flagHelp      = flag.Bool("h", false, "Help menu.")
	flagVersion   = flag.Bool("v", false, "Shows the version.")
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
		writeSTDERR("secret is not set. I can not continue like this. Exiting...")
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
	output, err := collectSecret(sess, *flagSecretKey)
	if err != nil {
		// This error is already pretty well formatted.
		writeSTDERR(fmt.Sprintln(err))
		os.Exit(1)
	}

	// print out the secrets
	fmt.Println(output)
}

// awsOption will check to see if profile is set and optionally add it in if required.
func awsOptions(region, profile string) session.Options {
	options := session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}
	if profile != "" {
		options.Profile = profile
	}
	return options
}

// awsSession just wraps the function used to create the session. This is mainly for furture work as AWS sometimes change this.
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
		//VersionStage: aws.String("AWSPREVIOUS"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		// The errors from the AWS SDK are actually pretty useful. So we should make effort to use them.
		return "", fmt.Errorf("error retriving the secret from %s. Error: %s", key, err.Error())
	}

	// The returned text from AWS Secret manager is already in JSON and a key value pair setup.
	// This will work great with Launch.
	return *result.SecretString, nil
}
