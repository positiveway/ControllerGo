cd ./Build/linux
rm ./ControllerGo
rm ./ControllerRust
cp -f ~/CLionProjects/ControllerRust/target/release/ControllerRust ./
cd ../../src
go build -o ../Build/linux/ControllerGo
cd ..
chmod +x ./run.sh
./run.sh