sudo rm -rf ~/go/pkg/mod/github.com/positiveway
sudo rm -rf ~/go/pkg/mod/cache/download/github.com/positiveway

cd ../src
go get -u github.com/positiveway/gofuncs@master