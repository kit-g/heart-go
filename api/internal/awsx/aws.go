package awsx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	env "heart/internal/config"
	"time"
)

var (
	S3       *s3.Client
	Env      env.AwsConfig
	events   *scheduler.Client
	s3Signer *s3.PresignClient
)

func Init(c env.AwsConfig) error {
	Env = c
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(c.AwsRegion),
	)

	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	events = scheduler.NewFromConfig(cfg)
	S3 = s3.NewFromConfig(cfg)
	s3Signer = s3.NewPresignClient(S3)

	return nil
}

func GeneratePresignedPostURL(
	ctx context.Context,
	bucket string,
	key string,
	contentType string,
	tagging string,
) (*s3.PresignedPostRequest, error) {
	input := s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		// Don't set Tagging here - we'll handle it manually
	}

	request, err := s3Signer.PresignPostObject(
		ctx,
		&input,
		func(options *s3.PresignPostOptions) {
			options.Expires = 15 * time.Minute
			options.Conditions = []interface{}{
				[]interface{}{"content-length-range", minContentLength, maxContentLength},
				map[string]string{"Content-Type": contentType},
				map[string]string{"tagging": tagging},
			}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	if request.Values == nil {
		request.Values = make(map[string]string)
	}
	request.Values["tagging"] = tagging
	request.Values["Content-Type"] = contentType

	return request, nil
}

func DeleteObject(ctx context.Context, bucket string, key string) (*s3.DeleteObjectOutput, error) {
	options := s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	return S3.DeleteObject(ctx, &options)
}

func CreateAccountDeletionSchedule(ctx context.Context, userId string) (*time.Time, *string, error) {
	scheduleName := fmt.Sprintf("account-deletion-%s", userId)
	when := time.Now().UTC().AddDate(0, 0, Env.AccountDeletionOffset)
	desc := fmt.Sprintf("Deletes user %s account after %d days", userId, Env.AccountDeletionOffset)

	payload, err := json.Marshal(
		map[string]any{
			"Event": "AccountDeletion",
			"Payload": map[string]string{
				"user_id": userId,
			},
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal input payload: %w", err)
	}

	input := scheduler.CreateScheduleInput{
		ActionAfterCompletion: types.ActionAfterCompletionDelete,
		Description:           aws.String(desc),
		FlexibleTimeWindow: &types.FlexibleTimeWindow{
			Mode: types.FlexibleTimeWindowModeOff,
		},
		GroupName:          aws.String(Env.ScheduleGroup),
		Name:               aws.String(scheduleName),
		ScheduleExpression: aws.String(fmt.Sprintf("at(%s)", when.Format("2006-01-02T15:04:05"))),
		State:              types.ScheduleStateEnabled,
		Target: &types.Target{
			Arn:     aws.String(Env.BackgroundFunctionArn),
			Input:   aws.String(string(payload)),
			RoleArn: aws.String(Env.BackgroundFunctionRole),
		},
	}
	out, err := events.CreateSchedule(ctx, &input)

	if err != nil {
		var conflictErr *types.ConflictException
		if errors.As(err, &conflictErr) {
			// schedule already exists, treat as success
			return &when, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return &when, out.ScheduleArn, nil
}

const (
	minContentLength = 128
	maxContentLength = 31_457_280 // 30 MB max
)
