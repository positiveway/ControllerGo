rm -rf ./build/linux/*
cd src
go build -o ../build/linux/ControllerGo
cd ..
chmod +x ./run.sh
./run.sh