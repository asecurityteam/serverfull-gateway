package serverfullgw

//go:generate mockgen -destination mock_signer_test.go -package serverfullgw github.com/asecurityteam/serverfull-gateway/pkg Signer
//go:generate mockgen -destination mock_roundtripper_test.go -package serverfullgw net/http RoundTripper
