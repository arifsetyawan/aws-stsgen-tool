# aws-stsgen-tool

[![Go Report Card](https://goreportcard.com/badge/github.com/arifsetyawan/aws-stsgen-tool)](https://goreportcard.com/report/github.com/arifsetyawan/aws-stsgen-tool)

Amazon Web Service STS Helper to Local Credential File. 

This small humble program is tends to help the user that have required to access their aws cli under the AWS  Security Token Service (STS). What this program do is : 

- help perform aws sts command.
- save the aws sts result to .aws/credential file under targeted profile that defined when configure this program.

## Binnary
In this git repository, I have built the binary and upload it under `/bin`. You can choose what operating system you are working on Darwin or Linux. If You have doubts about the binary, you can read through the code and build it on your own.

## Install
Before moving forward. Please makesure you already have aws-cli and base setup your aws credential and aws profile. After the base aws-cli setup is complete You need to setup or built config for this program by run `awsstsgen install`.
```
./bin/darwin/awsstsgen install
```

or you can set it manually in `$HOME/.awsstsgen/config.json` and provide this value : 
```
{
 "base-profile": "",
 "mfa-arn": "",
 "target-profile": ""
}
```
`mfa-arn` example is `arn:aws:iam::<your_aws_account_id>:mfa/<your_username>`

it will prompt for configuration of `your mfa arn`, `your current base credential name in ~/.aws/credential` 

## Set
```
./bin/darwin/awsstsgen set
```
set is the command to help you set your sts to `~/.aws/credential` file under the target-profile name. It will prompt for your mfa token. 
