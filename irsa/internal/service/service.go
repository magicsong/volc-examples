package service

import (
	"context"
	"fmt"
	"os"

	"github.com/magicsong/volc-examples/irsa/pkg/display"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

type Service struct {
	VpcClient *vpc.VPC
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



// getCredentials retrieves credentials from OpenID Connect provider
// when running in a Kubernetes pod with IAM role
func getCredentials() *credentials.Credentials {
    // Check if running in a pod with IRSA configuration
    tokenFile := os.Getenv("VOLCENGINE_OIDC_TOKEN_FILE")
    roleArn := os.Getenv("VOLCENGINE_ROLE_TRN")

    if tokenFile != "" && roleArn != "" {
        // Running with IRSA, use web identity token
        service := &Service{}
        creds, err := service.assumeRoleWithOIDC()
        if err != nil {
            fmt.Printf("Warning: Failed to assume role with OIDC: %v\n", err)
        } else {
            return creds
        }
    }

    // Fallback to static credentials for local development
    ak := os.Getenv("VOLCENGINE_ACCESS_KEY")
    sk := os.Getenv("VOLCENGINE_SECRET_KEY")
    
    if ak != "" && sk != "" {
        return credentials.NewStaticCredentials(ak, sk, "")
    }
    
    // Return empty credentials if no method available
    // This will cause authentication to fail appropriately
    return credentials.NewStaticCredentials("", "", "")
}

// 修正 init 方法，使用 getCredentials 获取凭证
func (s *Service) init() error {
    // 获取凭证
    creds := getCredentials()

    // 创建配置并存储在服务结构体中
    s.config = volcengine.NewConfig().
        WithCredentials(creds).
        WithRegion(s.region)

    sess, err := session.NewSession(s.config)
    if err != nil {
        return fmt.Errorf("failed to create session: %v", err)
    }

    s.VpcClient = vpc.New(sess)
    
    // 获取 AK、SK 用于 TOS 客户端
    akValue, skValue, token := creds.Get()
    
    tosOpts := []tos.ClientOption{
        tos.WithRegion("cn-beijing"),
    }
    
    if token != "" {
        tosOpts = append(tosOpts, tos.WithCredentials(tos.NewSessionCredentials(akValue, skValue, token)))
    } else {
        tosOpts = append(tosOpts, tos.WithCredentials(tos.NewStaticCredentials(akValue, skValue)))
    }
    
    s.TosClient, err = tos.NewClientV2("tos-cn-beijing.volces.com", tosOpts...)
    if err != nil {
        return fmt.Errorf("failed to create tos client: %v", err)
    }
    
    return nil
}

// getCredentials retrieves credentials from OpenID Connect provider
// when running in a Kubernetes pod with IAM role
func getCredentials() *credentials.Credentials {
	// Check if running in a pod with IRSA configuration
	tokenFile := os.Getenv("VOLCENGINE_OIDC_TOKEN_FILE")
	roleArn := os.Getenv("VOLCENGINE_ROLE_TRN")

	if tokenFile != "" && roleArn != "" {
		// Running with IRSA, use web identity token
		return 
	}

	// Fallback to static credentials for local development
	ak := os.Getenv("VOLCENGINE_ACCESS_KEY")
	sk := os.Getenv("VOLCENGINE_SECRET_KEY")
	
	if ak != "" && sk != "" {
		return credentials.NewStaticCredentials(ak, sk, "")
	}
	
	// Return empty credentials if no method available
	// This will cause authentication to fail appropriately
	return credentials.NewStaticCredentials("", "", "")
}

// init initializes the service clients
func (s *Service) init() error {
	// Get credentials from environment variables
	// Create the config and store it in the service struct
	s.config = volcengine.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(ak, sk, "")).
		WithRegion(s.region)

	sess, err := session.NewSession(s.config)
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	s.VpcClient = vpc.New(sess)
	s.TosClient, err = tos.NewClientV2("tos-cn-beijing.volces.com", tos.WithCredentials(tos.NewStaticCredentials(ak, sk)),tos.WithRegion("cn-beijing"))
	if err != nil {
		return fmt.Errorf("failed to create tos client: %v", err)
	}
	return nil
}

func (s *Service) DoSomething() error {
	return s.ListBuckets()
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
