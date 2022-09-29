package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nschimek/nba-scraper/cmd"
	"github.com/nschimek/nba-scraper/core"
	"github.com/spf13/viper"
)

func main() {
	core.SetupViper()
	if viper.GetBool("serverless") {
		lambda.Start(LambdaHandler)
	} else {
		cmd.Execute()
	}
}

func LambdaHandler() (string, error) {
	if err := cmd.Execute(); err != nil {
		return "", err
	} else {
		return "success", nil
	}
}
