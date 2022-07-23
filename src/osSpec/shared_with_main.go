package osSpec

const LeftMouse = -3
const RightMouse = -4
const MiddleMouse = -5

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
