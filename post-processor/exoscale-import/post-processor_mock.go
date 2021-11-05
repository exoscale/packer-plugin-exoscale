package exoscaleimport

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	egoscale "github.com/exoscale/egoscale/v2"
	"github.com/stretchr/testify/mock"
)

type exoscaleClientMock struct {
	mock.Mock
}

func (m *exoscaleClientMock) CopyTemplate(
	ctx context.Context,
	zone string,
	template *egoscale.Template,
	dst string,
) (*egoscale.Template, error) {
	args := m.Called(ctx, zone, template, dst)
	return args.Get(0).(*egoscale.Template), args.Error(1)
}

func (m *exoscaleClientMock) DeleteTemplate(ctx context.Context, zone string, template *egoscale.Template) error {
	args := m.Called(ctx, zone, template)
	return args.Error(0)
}

func (m *exoscaleClientMock) RegisterTemplate(
	ctx context.Context,
	zone string,
	template *egoscale.Template,
) (*egoscale.Template, error) {
	args := m.Called(ctx, zone, template)
	return args.Get(0).(*egoscale.Template), args.Error(1)
}

type s3ClientMock struct {
	mock.Mock
}

func (m *s3ClientMock) PutObject(
	ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.Options),
) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func (m *s3ClientMock) UploadPart(
	context.Context,
	*s3.UploadPartInput,
	...func(*s3.Options),
) (*s3.UploadPartOutput, error) {
	panic("not implemented")
}

func (m *s3ClientMock) CreateMultipartUpload(
	context.Context,
	*s3.CreateMultipartUploadInput,
	...func(*s3.Options),
) (*s3.CreateMultipartUploadOutput, error) {
	panic("not implemented")
}

func (m *s3ClientMock) CompleteMultipartUpload(
	context.Context,
	*s3.CompleteMultipartUploadInput,
	...func(*s3.Options),
) (*s3.CompleteMultipartUploadOutput, error) {
	panic("not implemented")
}

func (m *s3ClientMock) AbortMultipartUpload(
	context.Context,
	*s3.AbortMultipartUploadInput,
	...func(*s3.Options),
) (*s3.AbortMultipartUploadOutput, error) {
	panic("not implemented")
}

func (m *s3ClientMock) DeleteObject(
	ctx context.Context,
	params *s3.DeleteObjectInput,
	optFns ...func(*s3.Options),
) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
}
