rm ./ControllerGo
cp -f ~/CLionProjects/ControllerRust/target/release/ControllerRust ./
go build
sudo ./getLocale.sh
chmod +x ./run.sh
./run.sh