module example.com/deploy

go 1.17

replace example.com/tweet => ../tweet

require (
	example.com/tweet v0.0.0-00010101000000-000000000000
	github.com/nirasan/go-oauth-pkce-code-verifier v0.0.0-20220510032225-4f9f17eaec4c
	golang.org/x/oauth2 v0.0.0-20220524215830-622c5d57e401
)

require (
	github.com/aws/aws-sdk-go-v2 v1.16.4 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.10 // indirect
	github.com/aws/smithy-go v1.11.2 // indirect
	github.com/g8rswimmer/go-twitter/v2 v2.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220325170049-de3da57026de // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
