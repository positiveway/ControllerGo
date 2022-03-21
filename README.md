## Functionality
 - High precision mouse emulation using controller's stick or touchpad 
 - Unique typing mechanics that allow typing using only touchpads or sticks 
 - Custom layouts for commands and typing
 - Different actions for holding and shortly pressing a key
 - Remmaping of any key or trigger

## Supported controllers
 - Steam Controller
 - PlayStation Dualshock
 - Xbox

## Supported OS
 - Linux (requires **sudo**)
	 - Wayland
	 - X11
 - Windows (support is comming)

## Typing

For typing use both sticks, there are 8 states for each stick:
 - Right
 - UpRight 
 - Up 
 - UpLeft 
 - Left 
 - DownLeft 	
 - Down 	
 - DownRight

Which gives 8x8 = **64** possible combinations to assign your keys

Layouts are stored in `Layouts` folder

## How to run

> ./run.sh

Build instructions:
 - Linux: 
	 - `./build.sh`
