// ACAB - Another Cloud Auditing Binary
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// checkAWSRead checks an AWS bucket for public read access.
func auditAWS(bucketName string) {
	ctx := context.Background()
	sess := session.New(&aws.Config{
		Credentials: credentials.AnonymousCredentials,
	})
	fmt.Println("getting bucket region")
	region, err := s3manager.GetBucketRegion(ctx, sess, bucketName, "us-west-2")
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			fmt.Fprintf(os.Stderr, "unable to find bucket %s's region not found\n", bucketName)
		}
		return
	}
	fmt.Printf("Bucket %s is in %s region\n", bucketName, region)

	svc := s3.New(sess, &aws.Config{
		Region: &region,
	})

	fmt.Println("getting bucket ACL")
	result, err := svc.GetBucketAcl(&s3.GetBucketAclInput{Bucket: &bucketName})
	if err != nil {
		log.Printf(err.Error())
	}

	fmt.Println("Owner:", *result.Owner.DisplayName)
	fmt.Println("")
	fmt.Println("Grants")

	for _, g := range result.Grants {
		// If we add a canned ACL, the name is nil
		if g.Grantee.DisplayName == nil {
			fmt.Println("  Grantee:    EVERYONE")
		} else {
			fmt.Println("  Grantee:   ", *g.Grantee.DisplayName)
		}

		fmt.Println("  Type:      ", *g.Grantee.Type)
		fmt.Println("  Permission:", *g.Permission)
		fmt.Println("")
	}

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(2),
	}

	obj, err := svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(obj)
}

func usage() {
	fmt.Print(""+
		"ACAB - Another Cloud Auditing Binary\n"+
		"Usage:\n"+
		os.Args[0], " [-h] <bucket name>\n"+
		"\t-h\t print this help screen\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" {
		usage()
	}

	bucketName := os.Args[1]
	auditAWS(bucketName)
}
