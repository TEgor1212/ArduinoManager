package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Greeting returns a greeting for the given name.
func Greeting(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// SumRange returns the sum of numbers from `from` to `to` inclusive.
func SumRange(from, to int) int {
	sum := 0
	for i := from; i <= to; i++ {
		sum += i
	}
	return sum
}

// wordRank ranks how well a name matches a query:
// 0 = prefix/word-boundary match, 1 = contains, 2 = no match.
func wordRank(name, query string) int {
	name = strings.ToLower(name)
	query = strings.ToLower(query)
	if query == "" || !strings.Contains(name, query) {
		return 2
	}
	for i := 0; ; {
		idx := strings.Index(name[i:], query)
		if idx < 0 {
			break
		}
		idx += i
		if idx == 0 || !isAlphaNum(name[idx-1]) {
			return 0
		}
		if idx+len(query) >= len(name) {
			return 1
		}
		i = idx + 1
	}
	return 1
}

// isAlphaNum returns true if the byte is an ASCII letter or digit.
func isAlphaNum(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= '0' && b <= '9'
}

// guessBoardNames guesses vendor and architecture from a GitHub repo full name.
// Returns (vendor, arch).
func guessBoardNames(fullName string) (string, string) {
	known := map[string][2]string{
		"esp8266/Arduino":                    {"esp8266", "esp8266"},
		"espressif/arduino-esp32":            {"espressif", "esp32"},
		"espressif/esp32":                    {"espressif", "esp32"},
		"adafruit/Adafruit_nRF52_Arduino":    {"adafruit", "nrf52"},
		"arduino/ArduinoCore-avr":            {"arduino", "avr"},
		"arduino/ArduinoCore-megaavr":        {"arduino", "megaavr"},
		"arduino/ArduinoCore-samd":           {"arduino", "samd"},
		"arduino/ArduinoCore-mbed":           {"arduino", "mbed"},
		"arduino/ArduinoCore-nRF528x-mbedos": {"arduino", "nrf52840"},
		"arduino/ArduinoCore-renesas":        {"arduino", "renesas"},
		"stm32duino/Arduino_Core_STM32":      {"stm32duino", "stm32"},
		"SpenceKonde/ATTinyCore":             {"ATTinyCore", "avr"},
		"SpenceKonde/megaTinyCore":           {"megaTinyCore", "megaavr"},
	}
	if x, ok := known[fullName]; ok {
		return x[0], x[1]
	}
	parts := strings.SplitN(fullName, "/", 2)
	owner, repo := parts[0], parts[len(parts)-1]
	arch := strings.ToLower(repo)
	for _, prefix := range []string{"arduino-", "arduino_", "arduinocore-", "arduinocore_"} {
		arch = strings.TrimPrefix(arch, prefix)
	}
	arch = strings.TrimSuffix(arch, "_core")
	arch = strings.TrimSuffix(arch, "-core")
	if arch == "" {
		arch = "core"
	}
	return strings.ToLower(owner), arch
}

// maxFloat32 returns the larger of two float32 values.
func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// parseKeyValueFile reads a key=value or key:value file and returns the value.
func parseKeyValueFile(path, key string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, key+":") && !strings.HasPrefix(line, key+"=") {
			continue
		}
		v := line[strings.IndexAny(line, ":=")+1:]
		v = strings.TrimSpace(strings.Trim(v, "\"'"))
		return v
	}
	return ""
}

// libraryVersion reads the version from an Arduino library.properties file.
func libraryVersion(dir string) string {
	f, err := os.Open(filepath.Join(dir, "library.properties"))
	if err != nil {
		return "IDE 1 only"
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "version=") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "version="))
			if v != "" {
				return "v" + v
			}
		}
	}
	return "vunknown"
}
