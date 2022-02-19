sudo apt-get install -y unzip libsdl2-dev libdrm-dev libhidapi-dev libusb-1.0-0 libusb-1.0-0-dev libevdev-dev
rm ./build/linux/ControllerRust
rm -rf ./tmp
mkdir tmp
cd ./tmp
wget -O ControllerRust.zip https://github.com/positiveway/ControllerRust/archive/refs/heads/master.zip
unzip ControllerRust.zip
cd ./ControllerRust-master
chmod +x ./install.sh
./install.sh
cp ./target/release/ControllerRust ../../build/linux/
cd ../../
rm -rf ./tmp
chmod +x ./run.sh