package serverfullgw

import (
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

// Signer represents the AWS request signing process.
type Signer interface {
	Sign(r *http.Request, body io.ReadSeeker) error
}

// AWSSigner wraps the AWS specific signing process.
type AWSSigner struct {
	Session *session.Session
	Signer  *v4.Signer
}

// Sign a request using the AWS algorithm.
func (s *AWSSigner) Sign(r *http.Request, body io.ReadSeeker) error {
	_, err := s.Signer.Sign(r, body, "lambda", *s.Session.Config.Region, time.Now())
	return err
}
