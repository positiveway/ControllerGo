rm -rf ./build/linux/*
cd src
go build -o ../build/linux/controllerGo
cd ..
chmod +x ./run.sh
./run.sh