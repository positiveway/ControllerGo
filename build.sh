cd ./Build
rm ./ControllerGo
cp -f ~/CLionProjects/ControllerRust/target/release/ControllerRust ./
cd ../src
go build -o ../Build/ControllerGo
chmod +x ./getLocale.sh
cd ..
chmod +x ./run.sh
./run.sh