//go:generate mockgen -package=mock -destination=s3uploader.go github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface UploaderAPI
//go:generate mockgen -package=mock -destination=s3api.go github.com/aws/aws-sdk-go/service/s3/s3iface S3API
//go:generate mockgen -package=mock -destination=upload.go -source=../upload.go

package mock

///usr/local/Cellar/go/1.13/libexec/src/mime/multipart
