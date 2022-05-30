# go-aws-lambda-twitter-bot
Run a Twitter bot serverlessly on Amazon Web Services

## deployment instructions

1. Install AWS CLI (e.g. by running `install-aws-cli.sh`)
1. Install terraform
1. Clone this project
1. Run `deploy.sh`

## development instructions

1. Create `~/.aws/credentials` and `~/.aws/config` as described in the [official guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)
1. cd `source/dev`
1. `go run .`