package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/magicsong/volc-examples/irsa/pkg/display"
	"github.com/volcengine/volcengine-go-sdk/service/kms"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

type Service struct {
	VpcClient *vpc.VPC
	KmsClient kms.KMSAPI
	region    string
	config    *volcengine.Config
}

// NewService creates a new Service instance.
func NewService() (*Service, error) {
	s := &Service{
		region: getEnvWithDefault("VOLCENGINE_REGION", "cn-beijing"),
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize service: %v", err)
	}

	return s, nil
}

func (s *Service) init() error {
	// 获取凭证
	creds := getCredentials()

	// 创建配置并存储在服务结构体中
	s.config = volcengine.NewConfig().
		WithCredentialsChainVerboseErrors(true).
		WithCredentials(creds).
		WithRegion(s.region)
	sess, err := session.NewSession(s.config)
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	s.VpcClient = vpc.New(sess)
	s.KmsClient = kms.New(sess)

	return nil
}

// getCredentials retrieves credentials from OpenID Connect provider
// when running in a Kubernetes pod with IAM role
func getCredentials() *credentials.Credentials {
	return credentials.NewCredentials(credentials.NewOIDCCredentialsProviderFromEnv())
	// return credentials.NewChainCredentials([]credentials.Provider{
	// 	credentials.NewOIDCCredentialsProviderFromEnv(),
	// 	credentials.NewEnvCredentials().GetProvider()})
}

func (s *Service) DoSomething() error {
	secretName := getEnvWithDefault("KMS_SECRET_NAME", "default-secret")
	akKey := getEnvWithDefault("KMS_AK_KEY", "accessKey")
	skKey := getEnvWithDefault("KMS_SK_KEY", "secretKey")
	typeVar := getEnvWithDefault("CLOUD_TYPE", "default")

	v, err := s.KmsClient.GetSecretValue(&kms.GetSecretValueInput{
		SecretName: volcengine.String(secretName),
	})
	if err != nil {
		return fmt.Errorf("failed to get secret value: %w", err)
	}

	// Parse the secret value as JSON
	var secretData map[string]string
	if err := json.Unmarshal([]byte(*v.SecretValue), &secretData); err != nil {
		return fmt.Errorf("failed to parse secret value as JSON: %w", err)
	}

	// Retrieve AK and SK from the parsed JSON
	ak, akExists := secretData[akKey]
	sk, skExists := secretData[skKey]
	if !akExists || !skExists {
		return fmt.Errorf("missing required keys in secret: %s, %s", akKey, skKey)
	}

	if typeVar == "aws" {
		// Create AWS credentials file
		awsCreds := fmt.Sprintf("[default]\naws_access_key_id = %s\naws_secret_access_key = %s\n", ak, sk)
		awsCredsPath := os.ExpandEnv("$HOME/.aws/credentials")
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(awsCreds, 0700); err != nil {
			return fmt.Errorf("failed to create AWS credentials directory: %w", err)
		}
		if err := os.WriteFile(awsCredsPath, []byte(awsCreds), 0600); err != nil {
			return fmt.Errorf("failed to write AWS credentials file: %w", err)
		}
	}

	display.PrintAsJSON(map[string]string{
		"AccessKey": ak,
		"SecretKey": sk,
	})
	return nil
}

// Helper function to get environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
