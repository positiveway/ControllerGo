cd ./Build
rm ./ControllerGo
rm ./ControllerRust
cp -f ~/CLionProjects/ControllerRust/target/release/ControllerRust ./
cd ../src
go build -o ../Build/ControllerGo
chmod +x ./getLocale.sh
cd ..
chmod +x ./run.sh
./run.sh