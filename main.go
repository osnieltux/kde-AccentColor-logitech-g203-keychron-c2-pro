package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/godbus/dbus/v5"
)

const DEFAULT_COLOR = "0,255,0"
const FALLBACK_COLOR = "00FF00"

// Keyboard settings
const BASEDIR = "/sys/class/hidraw/"

var KEYBOARD_PATH = ""
var vendorID = "3434"
var productID = "0526"
var productInterface = "1.1"

const KEYBOARD_RANDOM_BRIGHTNESS = true
const KEYBOARD_DEFAULT_BRIGHTNESS = 100
const CHANGE_KEYBOARD_COLOR = true

var CURRENT_COLOR = FALLBACK_COLOR

const DEBUG = false

var DEVICE_NAME = "cheering-viscacha" //"hollering-marmot"

func main() {
	err := checkDependencies()
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

	// Changing keyboard
	if CHANGE_KEYBOARD_COLOR {
		KEYBOARD_PATH, err = detectKeyboard()
		if err != nil {
			log.Println(err)
			return
		}
		if DEBUG {
			log.Printf("KEYBOARD_PATH: %s", KEYBOARD_PATH)
		}
		if KEYBOARD_PATH == "" {
			if DEBUG {
				log.Println("No keyboard detected.")
			}
			return
		}
		if err := setColor(color); err != nil {
			log.Println(err)
		}

		// Changing brightness
		if err := setBrightness(getBrightness()); err != nil {
			log.Println(err)
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

func detectKeyboard() (string, error) {
	files, err := os.ReadDir(BASEDIR)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %v", BASEDIR, err)
	}
	pathToCheck := ""
	for _, file := range files {
		pathToCheck = BASEDIR + file.Name()
		// os.Readlink devuelve el destino del enlace, que puede ser una ruta relativa.
		relativeDest, err := os.Readlink(pathToCheck)
		if err != nil {
			log.Printf("Error al leer el enlace symbolic %s: %v", pathToCheck, err)
			continue
		}
		absolutePath := filepath.Clean(filepath.Join(filepath.Dir(pathToCheck), relativeDest))

		valueToCheck := strings.Split(absolutePath, "/")
		device := strings.ReplaceAll(valueToCheck[9], ".", ":")
		deviceParts := strings.Split(device, ":")

		deviceVendorID := deviceParts[1]
		deviceProductID := deviceParts[2]
		deviceInterface := strings.Split(valueToCheck[8], ":")[1]

		//device
		if (deviceVendorID == vendorID) && (deviceProductID == productID) && (deviceInterface == productInterface) {
			// fmt.Println("AA", deviceParts, deviceInterface)
			// fmt.Printf("Enlace: %s -> %s : %s %s\n", pathToCheck, absolutePath, device, valueToCheck[8])
			return filepath.Clean(filepath.Join("/dev/", file.Name())), nil
		}
	}
	return "", nil
}

func setColor(color string) error {
	var r, g, b int
	_, err := fmt.Sscanf(color, "%d,%d,%d", &r, &g, &b)
	if err != nil {
		return err
	}

	brightnessPercent := 100 // 0-100%

	hi, lo := RGBtoHSVBytes(r, g, b)
	brightness := PercentToByte(brightnessPercent)

	packet := []byte{
		0x07,
		0x03,
		0x04,       // set COLOR
		hi,         //
		lo,         //
		brightness, //
		0x00,
		0x00,
	}

	// string to "\xHH\xHH..."
	var cmdStr string
	for _, b := range packet {
		cmdStr += fmt.Sprintf("\\x%02x", b)
	}

	cmd := fmt.Sprintf("echo -ne \"%s\" > %v", cmdStr, KEYBOARD_PATH)
	if DEBUG {
		log.Println("cmd:", cmd)
	}

	c := exec.Command("/bin/sh", "-c", cmd)
	if err := c.Run(); err != nil {
		return err
	}

	if DEBUG {
		log.Println("keyboard color changed")
	}

	return nil
}

func RGBtoHSVBytes(r, g, b int) (hi byte, lo byte) {
	if r < 0 {
		r = 0
	}
	if r > 255 {
		r = 255
	}
	if g < 0 {
		g = 0
	}
	if g > 255 {
		g = 255
	}
	if b < 0 {
		b = 0
	}
	if b > 255 {
		b = 255
	}

	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := rf
	if gf > max {
		max = gf
	}
	if bf > max {
		max = bf
	}

	min := rf
	if gf < min {
		min = gf
	}
	if bf < min {
		min = bf
	}

	delta := max - min
	vFloat := max

	var hFloat float64
	if delta == 0 {
		hFloat = 0
	} else if max == rf {
		hFloat = 60.0 * ((gf - bf) / delta)
	} else if max == gf {
		hFloat = 60.0 * (2.0 + (bf-rf)/delta)
	} else {
		hFloat = 60.0 * (4.0 + (rf-gf)/delta)
	}

	if hFloat < 0 {
		hFloat += 360
	}

	// Escalar a 0-255
	hi = byte(hFloat * 255 / 360)
	lo = byte(vFloat * 255) // Value max brightness
	return
}

// Transform percentage 0-100 a byte 0-255
func PercentToByte(percent int) byte {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	return byte(percent * 255 / 100)
}

func setBrightness(percentToSet int) error {
	b := byte(percentToSet * 255 / 100)
	brightness := fmt.Sprintf("0x%02X", b)

	if DEBUG {
		log.Println("setBrightness:", brightness, percentToSet, "%")
	}
	packet := []byte{
		0x07,
		0x03,
		0x01, // Set to brightness
		b,    // how much
		0x00,
		0x00,
		0x00,
		0x00,
	}

	var cmdStr string
	for _, b := range packet {
		cmdStr += fmt.Sprintf("\\x%02x", b)
	}

	cmd := fmt.Sprintf("echo -ne \"%s\" > %v", cmdStr, KEYBOARD_PATH)
	if DEBUG {
		log.Println("cmd set brightness:", cmd)
	}

	c := exec.Command("/bin/sh", "-c", cmd)
	if err := c.Run(); err != nil {
		return err
	}
	if DEBUG {
		log.Println("brightness changed")
	}

	return nil
}

func getBrightness() int {
	if DEBUG {
		log.Println("get Brightness")
	}
	if !KEYBOARD_RANDOM_BRIGHTNESS {
		return KEYBOARD_DEFAULT_BRIGHTNESS
	}
	// generate random value between 1 to 100
	return min(rand.Intn(100)+1, 100)
}
