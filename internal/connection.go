package internal

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewDBConnection established and returned dynamoDB connection.
func NewDBConnection() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	return svc
}
