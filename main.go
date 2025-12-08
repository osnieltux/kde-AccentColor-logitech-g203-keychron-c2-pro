package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/godbus/dbus/v5"
)



const DEFAULT_COLOR = "0,255,0"
const FALLBACK_COLOR = "00FF00"
var CURRENT_COLOR = FALLBACK_COLOR

const DEBUG = false

var DEVICE_NAME = "hollering-marmot"

func main() {
	err:= checkDependencies() 
	if err != nil {
		log.Fatal(err)
	}

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("Error connecting to session bus: %v", err)
	}
	defer conn.Close()

	call := conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch", 0,
		"type='signal',interface='org.kde.kconfig.notify'",
	)
	if call.Err != nil {
		log.Fatalf("AddMatch error: %v", call.Err)
	}

	signalChan := make(chan *dbus.Signal, 10)
	conn.Signal(signalChan)

	// Initial setup
	onConfigChanged()

	for sig := range signalChan {
		if sig.Name == "org.kde.kconfig.notify.ConfigChanged" {
			if DEBUG {
				log.Printf("ConfigChanged")
			}
			onConfigChanged()
		}
	}
}

func onConfigChanged() {
	color, err := getAccentColor()
	if err != nil {
		log.Printf("Error getting accent color: %v", err)
		// set default color
		color = DEFAULT_COLOR
	} else {
		if DEBUG {
			log.Printf("AccentColor retrieved: %s", color)
		}
	}

	if color != CURRENT_COLOR {
		CURRENT_COLOR = transformColor(color)
		setMouseAccentColor()
		if DEBUG {
			log.Printf("Accent color changed to: %s", CURRENT_COLOR)
		}
	} else {
		if DEBUG {
			log.Println("Accent color has not changed.")
		}
	}

}

func setMouseAccentColor() {
	cmd := exec.Command("ratbagctl", DEVICE_NAME, "led", "0", "set", "color", CURRENT_COLOR)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func transformColor(inputColor string) string {
	// transform r,g,b to hex
	var r, g, b int
	_, err := fmt.Sscanf(inputColor, "%d,%d,%d", &r, &g, &b)
	if err != nil {
		if DEBUG {
			log.Printf("Error parsing RGB color: %v , inputColor: %v, using FALLBACK_COLOR: %v", err, inputColor, FALLBACK_COLOR)
		}
		return FALLBACK_COLOR
	}
	return fmt.Sprintf("%02X%02X%02X", r, g, b)
}

func getAccentColor() (string, error) {
	out, err := exec.Command("kreadconfig6", "--file", "kdeglobals", "--group", "General", "--key", "AccentColor").Output()
	if err != nil {
		return "", err
	}

	s := string(out)

	if s == "" {
		return "", fmt.Errorf("empty AccentColor value")
	}

	return s, nil
}

func checkDependencies() error {
	_, err := exec.LookPath("kreadconfig6")
	if err != nil {
		return fmt.Errorf("Missing dependencies, kreadconfig6 not found in PATH")
		
	}
	_, err = exec.LookPath("ratbagctl")
	if err != nil {
		return fmt.Errorf("Missing dependencies, ratbagctl not found in PATH")
	}
	return nil
}