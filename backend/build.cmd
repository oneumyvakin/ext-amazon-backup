go get -u "github.com/aws/aws-sdk-go/aws"
go get -u "github.com/aws/aws-sdk-go/aws/awserr"
go get -u "github.com/aws/aws-sdk-go/service/s3"
go get -u "github.com/aws/aws-sdk-go/aws/session"
go get -u "github.com/aws/aws-sdk-go/service/s3/s3manager"

set PKGNAME=backup-amazon
set LOCALPATH=%~dp0

mklink /J "%GOPATH%\src\%PKGNAME%" "%LOCALPATH%"


set GOOS=windows
set GOARCH=386
go build -o ../sbin/%PKGNAME%.exe %PKGNAME%


set GOOS=linux
set GOARCH=386
go build -o ../sbin/%PKGNAME%.386 %PKGNAME%


rmdir "%GOPATH%\src\%PKGNAME%"