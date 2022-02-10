Значит так ебана в рот чтоб запустить хуяришь `./run.sh` в терминал. Все уже скомпилено, вшито, свистелки, перделки, batteries included.

Раскладка в файле commands, ща сюда ебану текущую версию но может она уже не текущая а поменялась. Мне что делать что ли нехуй еще ридми обновлять. Ебись оно конем

```BtnSouth:         undoCmd,
	BtnEast:          {uinput.KeyBackspace},
	BtnNorth:         {uinput.KeySpace},
	BtnWest:          {uinput.KeyEnter},
	BtnC:             NoAction,
	BtnZ:             NoAction,
	BtnLeftTrigger:   SwitchToTyping,
	BtnLeftTrigger2:  {RightMouse},
	BtnRightTrigger:  SwitchLang,
	BtnRightTrigger2: {LeftMouse},
	BtnSelect:        {uinput.KeyLeftmeta},
	BtnStart:         {uinput.KeyEsc},
	BtnMode:          NoAction,
	BtnLeftThumb:     copyCmd,
	BtnRightThumb:    pasteCmd,
	BtnDPadUp:        {uinput.KeyUp},
	BtnDPadDown:      {uinput.KeyDown},
	BtnDPadLeft:      {uinput.KeyLeft},
	BtnDPadRight:     {uinput.KeyRight},
```

Раскладка для печати в файле `layout.csv`

Switching language may not work on your system because go fuck yourself

App supports controller's connection and disconnection at runtime.

На руском копи паст не работает не ебу почему, issue что ли создай и PR чтоб исправить.

Xbox is likely to work as well. Testing is required.
