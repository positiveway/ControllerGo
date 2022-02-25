cd ./build/linux
rm -rf ./*
cp -f ~/CLionProjects/ControllerRust/target/release/controllerRust ./
cd ../../src
go build -o ../build/linux/controllerGo
cd ..
chmod +x ./run.sh
./run.sh