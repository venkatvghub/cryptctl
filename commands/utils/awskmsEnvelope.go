package utils

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type EnvelopeProvider struct {
	kmsClient      *kms.Client
	dynamoDbClient *dynamodb.Client
}

type localStackEndpointResolver struct {
	endpoint string
}

func createEnvelopeProvider(cfg aws.Config, KmsKeyAlias, DynamoDbTable string) (*EnvelopeProvider, error) {
	// Create a KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Create a DynamoDB client
	dynamoDBClient := dynamodb.NewFromConfig(cfg)

	_, err := checkDynamoDBTableExists(dynamoDBClient, DynamoDbTable)
	if err != nil {
		return nil, err
	}

	_, err = checkKMSKeyAliasExists(kmsClient, KmsKeyAlias)
	if err != nil {
		return nil, err
	}

	return &EnvelopeProvider{
		kmsClient:      kmsClient,
		dynamoDbClient: dynamoDBClient,
	}, nil
}

func (r *localStackEndpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL: r.endpoint,
	}, nil
}

func initAwsKmsEnvelopeLocalStack(KmsKeyAlias, DynamoDbTable string) (*EnvelopeProvider, error) {
	localStackEndpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	if localStackEndpoint == "" {
		return nil, errors.New("LOCALSTACK_ENDPOINT environment variable is not set")
	}

	// Load the AWS configuration with the LocalStack endpoint
	/*cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(
		aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: localStackEndpoint,
			}, nil
		}),
	))*/
	// Load the AWS configuration with the LocalStack endpoint
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(
		&localStackEndpointResolver{endpoint: localStackEndpoint},
	))

	if err != nil {
		return nil, err
	}
	return createEnvelopeProvider(cfg, KmsKeyAlias, DynamoDbTable)
}

func initAwsKmsEnvelopeAws(KmsKeyAlias, DynamoDbTable string) (*EnvelopeProvider, error) {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	return createEnvelopeProvider(cfg, KmsKeyAlias, DynamoDbTable)
}

func InitAwsKmsEnvelope(KmsKeyAlias, DynamoDbTable string, Debug bool) error {
	if Debug {
		_, err := initAwsKmsEnvelopeLocalStack(KmsKeyAlias, DynamoDbTable)
		if err != nil {
			return err
		}
	} else {
		_, err := initAwsKmsEnvelopeAws(KmsKeyAlias, DynamoDbTable)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkDynamoDBTableExists(client *dynamodb.Client, tableName string) (bool, error) {
	_, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		// Check if the error is because the table doesn't exist
		var notFoundErr *dynamodbtypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func checkKMSKeyAliasExists(client *kms.Client, aliasName string) (bool, error) {
	input := &kms.ListAliasesInput{
		Limit: aws.Int32(100),
	}
	paginator := kms.NewListAliasesPaginator(client, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return false, err
		}
		for _, alias := range output.Aliases {
			if aws.ToString(alias.AliasName) == aliasName {
				return true, nil
			}
		}
	}
	return false, nil
}
