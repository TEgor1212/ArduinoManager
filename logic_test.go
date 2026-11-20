package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGreeting(t *testing.T) {
	got := Greeting("Fyne")
	want := "Hello, Fyne!"
	if got != want {
		t.Errorf("Greeting() = %q, want %q", got, want)
	}
}

func TestSumRange(t *testing.T) {
	tests := []struct {
		from, to, expected int
	}{
		{1, 10, 55},
		{0, 0, 0},
		{-3, 3, 0},
		{1, 100, 5050},
	}
	for _, tt := range tests {
		got := SumRange(tt.from, tt.to)
		if got != tt.expected {
			t.Errorf("SumRange(%d, %d) = %d, want %d", tt.from, tt.to, got, tt.expected)
		}
	}
}

func TestWordRank(t *testing.T) {
	tests := []struct {
		name, query string
		expected    int
	}{
		{"Arduino Library", "arduino", 0},
		{"MyLibrary", "mylibrary", 0},
		{"SomeLib", "lib", 1},
		{"Hello", "xyz", 2},
		{"Test", "", 2},
	}
	for _, tt := range tests {
		got := wordRank(tt.name, tt.query)
		if got != tt.expected {
			t.Errorf("wordRank(%q, %q) = %d, want %d", tt.name, tt.query, got, tt.expected)
		}
	}
}

func TestIsAlphaNum(t *testing.T) {
	if !isAlphaNum(byte('a')) {
		t.Errorf("isAlphaNum('a') should be true")
	}
	if !isAlphaNum(byte('z')) {
		t.Errorf("isAlphaNum('z') should be true")
	}
	if !isAlphaNum(byte('0')) {
		t.Errorf("isAlphaNum('0') should be true")
	}
	if !isAlphaNum(byte('9')) {
		t.Errorf("isAlphaNum('9') should be true")
	}
	if isAlphaNum(byte('-')) {
		t.Errorf("isAlphaNum('-') should be false")
	}
	if isAlphaNum(byte(' ')) {
		t.Errorf("isAlphaNum(' ') should be false")
	}
}

func TestGuessBoardNames(t *testing.T) {
	vendor, arch := guessBoardNames("esp8266/Arduino")
	if vendor != "esp8266" || arch != "esp8266" {
		t.Errorf("guessBoardNames(esp8266/Arduino) = (%q, %q), want (esp8266, esp8266)", vendor, arch)
	}

	vendor, arch = guessBoardNames("arduino/ArduinoCore-avr")
	if vendor != "arduino" || arch != "avr" {
		t.Errorf("guessBoardNames(arduino/ArduinoCore-avr) = (%q, %q), want (arduino, avr)", vendor, arch)
	}

	vendor, arch = guessBoardNames("espressif/arduino-esp32")
	if vendor != "espressif" || arch != "esp32" {
		t.Errorf("guessBoardNames(espressif/arduino-esp32) = (%q, %q), want (espressif, esp32)", vendor, arch)
	}

	// Unknown repo — heuristic fallback.
	vendor, arch = guessBoardNames("myuser/MyBoardCore")
	if vendor == "" || arch == "" {
		t.Errorf("guessBoardNames(myuser/MyBoardCore) should return a guess, got (%q, %q)", vendor, arch)
	}
}

func TestMaxFloat32(t *testing.T) {
	if maxFloat32(1.0, 2.0) != 2.0 {
		t.Errorf("maxFloat32(1.0, 2.0) = %f, want 2.0", maxFloat32(1.0, 2.0))
	}
	if maxFloat32(5.5, 3.2) != 5.5 {
		t.Errorf("maxFloat32(5.5, 3.2) = %f, want 5.5", maxFloat32(5.5, 3.2))
	}
	if maxFloat32(0.0, 0.0) != 0.0 {
		t.Errorf("maxFloat32(0.0, 0.0) = %f, want 0.0", maxFloat32(0.0, 0.0))
	}
}

func TestParseKeyValueFile(t *testing.T) {
	// Create a temporary config file.
	tmpDir, err := os.MkdirTemp("", "arduinomanager-test-")
	if err != nil {
		t.Errorf("failed to create temp dir: %v", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	cfgPath := filepath.Join(tmpDir, "test.conf")
	content := `# Arduino CLI config
sketchbook.path = /home/user/Arduino
board_manager.url = https://downloads.arduino.cc
`
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Errorf("failed to write test config: %v", err)
		return
	}

	got := parseKeyValueFile(cfgPath, "sketchbook.path")
	want := "/home/user/Arduino"
	if got != want {
		t.Errorf("parseKeyValueFile(sketchbook.path) = %q, want %q", got, want)
	}

	got = parseKeyValueFile(cfgPath, "nonexistent")
	if got != "" {
		t.Errorf("parseKeyValueFile(nonexistent) = %q, want empty", got)
	}

	got = parseKeyValueFile("/nonexistent/path", "key")
	if got != "" {
		t.Errorf("parseKeyValueFile on missing file should return empty, got %q", got)
	}
}

func TestLibraryVersion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arduinomanager-test-")
	if err != nil {
		t.Errorf("failed to create temp dir: %v", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	// Test: missing library.properties returns "IDE 1 only".
	got := libraryVersion(tmpDir)
	if got != "IDE 1 only" {
		t.Errorf("libraryVersion on empty dir = %q, want 'IDE 1 only'", got)
	}

	// Test: valid library.properties.
	propPath := filepath.Join(tmpDir, "library.properties")
	content := `name=My Test Lib
version=2.1.0
author=Test
`
	if err := os.WriteFile(propPath, []byte(content), 0o644); err != nil {
		t.Errorf("failed to write library.properties: %v", err)
		return
	}
	got = libraryVersion(tmpDir)
	want := "v2.1.0"
	if got != want {
		t.Errorf("libraryVersion() = %q, want %q", got, want)
	}
}