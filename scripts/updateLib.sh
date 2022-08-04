username="positiveway"
repository="gofuncs"

#sudo rm -rf ~/go/pkg/mod/github.com/$username
#sudo rm -rf ~/go/pkg/mod/cache/download/github.com/$username

cd ../src
commitID=$(git ls-remote https://github.com/$username/$repository.git HEAD | awk '{print substr($1, 0, 7)}')
go get github.com/$username/$repository@$commitID