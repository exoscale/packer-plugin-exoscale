package exoscaleimport

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepUploadImage struct {
	postProcessor *PostProcessor
}

func (s *stepUploadImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		ui       = state.Get("ui").(packer.Ui)
		artifact = state.Get("artifact").(packer.Artifact)

		imageFile  = artifact.Files()[0]
		bucketFile = filepath.Base(imageFile)
	)

	ui.Say("Uploading template file")

	f, err := os.Open(imageFile)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to read template file: %v", err))
		return multistep.ActionHalt
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to read template file: %v", err))
		return multistep.ActionHalt
	}

	// For tracking image file upload progress
	pf := ui.TrackProgress(imageFile, 0, fileInfo.Size(), f)
	defer pf.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		ui.Error(fmt.Sprintf("Unable to compute template file checksum: %v", err))
		return multistep.ActionHalt
	}
	if _, err := f.Seek(0, 0); err != nil {
		ui.Error(fmt.Sprintf("Unable to compute template file checksum: %v", err))
		return multistep.ActionHalt
	}

	output, err := s3manager.
		NewUploader(s.postProcessor.sos).
		Upload(ctx,
			&s3.PutObjectInput{
				Bucket:     aws.String(s.postProcessor.config.ImageBucket),
				Key:        aws.String(bucketFile),
				Body:       pf,
				ContentMD5: aws.String(base64.StdEncoding.EncodeToString(hash.Sum(nil))),
				ACL:        s3types.ObjectCannedACLPublicRead,
			})
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to upload template file: %v", err))
		return multistep.ActionHalt
	}

	state.Put("image_url", output.Location)
	state.Put("image_checksum", fmt.Sprintf("%x", hash.Sum(nil)))

	return multistep.ActionContinue
}

func (s *stepUploadImage) Cleanup(state multistep.StateBag) {
	var (
		ui       = state.Get("ui").(packer.Ui)
		artifact = state.Get("artifact").(packer.Artifact)

		imageFile  = artifact.Files()[0]
		bucketFile = filepath.Base(imageFile)
	)

	if s.postProcessor.config.SkipClean {
		return
	}

	ui.Say("Deleting uploaded template file")

	_, err := s.postProcessor.sos.DeleteObject(context.Background(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.postProcessor.config.ImageBucket),
			Key:    aws.String(bucketFile),
		})
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to delete template file: %v", err))
	}
}
