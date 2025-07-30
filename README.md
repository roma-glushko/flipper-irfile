# Flipper IR File Parser

🐬 Flipper Zero Infra Red (IR) File Parser. 

You can use to load, process and save [Flipper Zero IR signal and IR library files](https://github.com/jamisonderek/flipper-zero-tutorials/wiki/Infrared):

```text
Filetype: IR library file
Version: 1
#
name: Power
type: parsed
protocol: SIRC
address: 01 00 00 00
command: 15 00 00 00
#
name: Power
type: raw
frequency: 38000
duty_cycle: 0.330000
data: 1617 4604 1559 1537 1560 1537 1560 4661 1533 33422 1613 4607 1566 1530 1556 1540 1536 4685 1539
```

## Install

```sh
go get github.com/roma-glushko/flipperirfile
```

## Usage

```go
package main

import (
	"fmt"
    "io"
    "os"

	irfile "github.com/roma-glushko/flipper-irfile"
)

func main() {
	f, err := os.Open("tv.ir")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	rawData, err := io.ReadAll(f)

	if err != nil {
		panic(err)
	}

	// Load IR signals
	sigLib, err := irfile.Unmarshal(rawData)

	// Process IR signals
	fmt.Println(fmt.Sprintf("Filetype: %s", sigLib.FileType))
	fmt.Println(fmt.Sprintf("Version: %d", sigLib.Version))
	fmt.Println(fmt.Sprintf("Signals: %d", len(sigLib.Signals)))

	parsedSignals := make([]irfile.Signal, 0, 20)

	for _, s := range sigLib.Signals {
		if s.Type != irfile.SignalTypeParsed {
			continue
		}

		parsedSignals = append(parsedSignals, s)
	}

	sigLib.Signals = parsedSignals

	// Save IR signals
	rawLib, err = irfile.Marshal(sigLib)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile("parsed.ir", rawLib, 0644)

	if err != nil {
		panic(err)
	}
}
```

### Generate IR Signals

```golang
package main

import (
	"fmt"
	"os"
	
	irfile "github.com/roma-glushko/flipper-irfile"
)

func main() {
    lib := irfile.SignalLib{
        Filetype: irfile.FiletypeSignalFile,
        Version:  "1",
    }

	addr := 0x05
	cmdCodeMin := 0x4000
	cmdCodeMax := 0x40FF
    
    for cmdCode := cmdCodeMin; cmdCode <= cmdCodeMax; cmdCode++ {
        sig := irfile.Signal{
            Name:     fmt.Sprintf("Cmd %d", cmdCode),
            Type:     irfile.SignalTypeParsed,
            Protocol: irfile.ProtocolNECExt,
            Address:  addr,
            Command:  cmd,
        }
		
        lib.Signals = append(lib.Signals, sig)
    }
        
    data, err := Marshal(lib)
    
    if err != nil {
        panic(err)
    }
    
    err = os.WriteFile("gen.ir", data, 0644)
    
    if err != nil {
        panic(err)
    }
}   
```

## Credits

Made with ❤️ by [Roman Glushko](https://github.com/roma-glushko), Apache License 2.0.
