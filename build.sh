rm ./ControllerGo
cp -f ~/CLionProjects/ControllerRust/target/release/ControllerRust ./
go build
chmod +x ./run.sh
./run.sh