module example.com/deploy

go 1.18

replace example.com/tweet => ../tweet

require (
	example.com/tweet v0.0.0-00010101000000-000000000000
	github.com/nirasan/go-oauth-pkce-code-verifier v0.0.0-20220510032225-4f9f17eaec4c
	golang.org/x/oauth2 v0.0.0-20220524215830-622c5d57e401
)

require (
	cloud.google.com/go v0.65.0 // indirect
	github.com/aws/aws-lambda-go v1.24.0 // indirect
	github.com/aws/aws-sdk-go v1.40.42 // indirect
	github.com/g8rswimmer/go-twitter/v2 v2.0.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/api v0.30.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
