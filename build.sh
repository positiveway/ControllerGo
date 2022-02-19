cd ./build/linux
rm ./ControllerGo
rm ./ControllerRust
cp -f ~/CLionProjects/ControllerRust/target/release/ControllerRust ./
cd ../../src
go build -o ../build/linux/ControllerGo
cd ..
chmod +x ./run.sh
./run.sh