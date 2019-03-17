package exocomp

import (
	"regexp"
	"testing"
)

func TestMatchFiles(t *testing.T) {
	tests := []struct {
		file        string
		shouldMatch bool
	}{
		{"Dockerfile", true},
		{"main.go", false},
		{"docker-compose.yaml", true},
		{"docker-compose.yml", true},
		{"deploy.yaml", true},
		{"deploy.yml", true},
		{".gitlab-ci.yml", true},
		{"index.php", false},
	}
	var matcher = regexp.MustCompile(fileWithImages)

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			if result := matcher.MatchString(tt.file); result != tt.shouldMatch {
				t.Errorf("matcher.MatchString(%v) = %v, wanted %v", tt.file, result, tt.shouldMatch)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	tests := []struct {
		line        string
		shouldMatch bool
		image       string
	}{
		{"# comment\n", false, ""},
		{"FROM nginx:1.13.7-alpine", true, "nginx:1.13.7-alpine"},
		{"FROM nginx:1.13.7-alpine\n", true, "nginx:1.13.7-alpine"},
		{"FROM nginx:1.13.7-alpine \n", true, "nginx:1.13.7-alpine"},
		{"FROM ubuntu AS weather", true, "ubuntu"},
		{"FROM php:7.3 as fibonacci", true, "php:7.3"},
	}
	var matcher = regexp.MustCompile(linesFrom)

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := matcher.MatchString(tt.line)
			if result != tt.shouldMatch {
				t.Errorf("matcher.MatchString(%v) = %v, wanted %v", tt.line, result, tt.shouldMatch)
			}
			if tt.shouldMatch == false {
				return
			}
			if m := matcher.FindStringSubmatch(tt.line); m[1] != tt.image {
				t.Errorf("image(%v) = %v, wanted %v", tt.line, m[1], tt.image)
			}

		})
	}
}
