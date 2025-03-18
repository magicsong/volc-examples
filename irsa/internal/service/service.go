package service

import (
	"context"
	"fmt"
	"os"

	"github.com/magicsong/volc-examples/irsa/pkg/display"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
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
	TosClient *tos.ClientV2
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

// 修正 init 方法，使用 getCredentials 获取凭证
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
	// 获取 AK、SK 用于 TOS 客户端
	v, err := creds.Get()
	if err != nil {
		return fmt.Errorf("failed to get credentials: %v", err)
	}
	akValue := v.AccessKeyID
	skValue := v.SecretAccessKey

	tosOpts := []tos.ClientOption{
		tos.WithRegion("cn-beijing"),
	}

	tosOpts = append(tosOpts, tos.WithCredentials(tos.NewStaticCredentials(akValue, skValue)))

	s.TosClient, err = tos.NewClientV2("tos-cn-beijing.volces.com", tosOpts...)
	if err != nil {
		return fmt.Errorf("failed to create tos client: %v", err)
	}

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
	v, err := s.KmsClient.GetSecretValue(&kms.GetSecretValueInput{
		SecretName: volcengine.String("hello"),
	})
	if err != nil {
		return fmt.Errorf("failed to get secret value: %w", err)
	}
	display.PrintAsJSON(v)
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

// ListBuckets
func (s *Service) ListBuckets() error {
	resp, err := s.TosClient.ListBuckets(context.TODO(), &tos.ListBucketsInput{})
	if err != nil {
		return err
	}
	display.PrintAsJSON(resp)
	return nil
}
