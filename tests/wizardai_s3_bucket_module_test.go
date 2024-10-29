package test

import (
	"fmt"
	"log"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go/aws"
	awsSesson "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestS3BucketModule(t *testing.T) {
	t.Parallel()

	// Define input variables
	bucketName := "test"
	environment := "development"
	region := "us-west-2"

	terraformOptions := &terraform.Options{
		TerraformDir: "../wizardai_s3_bucket_module",

		Vars: map[string]interface{}{
			"name":        bucketName,
			"environment": environment,
			"region":      region,
		},

		// Prevent Terraform from storing the plan output in your working directory
		NoColor: true,
	}

	// Run `terraform init` and `terraform apply`, fail the test if there are any errors
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Run tests
	ValidateBucketNamingConvention(t, terraformOptions, bucketName, environment)
	ValidateBucketEncryptionAtRest(t, terraformOptions, region)
	ValidateBucketEncryptionInTransit(t, terraformOptions, region)
	ValidateBucketPublicAccessBlocked(t, terraformOptions, region)
	ValidateBucketVersioning(t, terraformOptions, region)
}

// Test that the bucket name adheres to the naming convention
func ValidateBucketNamingConvention(t *testing.T, terraformOptions *terraform.Options, name string, env string) {
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")
	expectedName := "wizardai-" + name + "-" + env

	assert.Equal(t, expectedName, bucketName, "Bucket name does not follow naming convention")
}

// Test that the bucket has encryption at rest enabled
func ValidateBucketEncryptionAtRest(t *testing.T, terraformOptions *terraform.Options, region string) {
	/// terratest doesn't support GetBucketEncryption so we inplement our own by AWS Go SDK
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")
	sess, err := awsSesson.NewSession(&awsSDK.Config{
		Region: awsSDK.String(region),
	})
	svc := awsS3.New(sess)

	input := &awsS3.GetBucketEncryptionInput{
		Bucket: awsSDK.String(bucketName),
	}
	result, err := svc.GetBucketEncryption(input)
	assert.NoError(t, err)

	for _, rule := range result.ServerSideEncryptionConfiguration.Rules {
		assert.NotNil(t, rule.ApplyServerSideEncryptionByDefault)
		assert.NotNil(t, rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		assert.Equal(t, "aws:kms", *rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	}
}

// Test that the public access is blocked
func ValidateBucketPublicAccessBlocked(t *testing.T, terraformOptions *terraform.Options, region string) {
	/// terratest doesn't support GetBucketEncryption so we inplement our own by AWS Go SDK
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")
	sess, err := awsSesson.NewSession(&awsSDK.Config{
		Region: awsSDK.String(region),
	})
	svc := awsS3.New(sess)
	// Get Bucket ACL
	aclResult, err := svc.GetBucketAcl(&s3.GetBucketAclInput{
		Bucket: awsSDK.String(bucketName),
	})
	if err != nil {
		log.Fatalf("Unable to get bucket ACL: %v", err)
	}

	// Check if any grants allow public access
	for _, grant := range aclResult.Grants {
		if *grant.Grantee.Type == "Group" && *grant.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
			fmt.Println("Bucket is publicly accessible via ACL")
		}
	}

	// Get Bucket Policy Status
	policyStatusResult, err := svc.GetBucketPolicyStatus(&s3.GetBucketPolicyStatusInput{
		Bucket: awsSDK.String(bucketName),
	})
	if err != nil {
		log.Fatalf("Unable to get bucket policy status: %v", err)
	}

	if *policyStatusResult.PolicyStatus.IsPublic {
		fmt.Println("Bucket is publicly accessible via policy")
	} else {
		fmt.Println("Bucket is not publicly accessible")
	}
}

// Test that the bucket has encryption in transit enforced by the bucket policy
func ValidateBucketEncryptionInTransit(t *testing.T, terraformOptions *terraform.Options, region string) {
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")
	bucketPolicy := aws.GetS3BucketPolicy(t, region, bucketName)

	// Assert that the bucket policy denies any requests without SecureTransport
	assert.Contains(t, bucketPolicy, `"aws:SecureTransport":"false"`, "Bucket policy does not enforce encryption in transit")
}

// Test that the bucket has versioning enabled
func ValidateBucketVersioning(t *testing.T, terraformOptions *terraform.Options, region string) {
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")
	versioningStatus := aws.GetS3BucketVersioning(t, region, bucketName)

	// Assert that the bucket policy denies any requests without SecureTransport
	assert.Equal(t, versioningStatus, `Enabled`, "Bucket Versioning is not enabled")
}
