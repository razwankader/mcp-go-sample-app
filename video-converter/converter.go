package videoconverter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type QualityPreset struct {
	CRF    string
	Preset string
}

type VideoConverter struct{}

var QUALITY_PRESETS = map[string]QualityPreset{
	"low":    {CRF: "28", Preset: "fast"},
	"medium": {CRF: "23", Preset: "medium"},
	"high":   {CRF: "18", Preset: "slow"},
}

var SUPPORTED_FORMATS = map[string]bool{
	"webm": true,
	"mkv":  true,
	"avi":  true,
	"mov":  true,
	"gif":  true,
}

func ValidateInput(inputPath string) error {
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file not found: %s", inputPath)
	}

	if !strings.HasSuffix(strings.ToLower(inputPath), ".mp4") {
		return errors.New("input file must be an MP4 file")
	}

	return nil
}

func GenerateOutputPath(inputPath, format string) string {
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	return fmt.Sprintf("%s.%s", base, strings.ToLower(format))
}

func BuildFFmpegCommand(inputPath, outputPath, format string) ([]string, error) {
	preset := QUALITY_PRESETS["medium"]

	cmd := []string{"ffmpeg", "-i", inputPath, "-y"}

	format = strings.ToLower(format)

	if format == "gif" {
		cmd = append(cmd,
			"-vf", "fps=15,scale=480:-1:flags=lanczos",
			"-c:v", "gif",
			outputPath,
		)
	} else if SUPPORTED_FORMATS[format] {
		cmd = append(cmd,
			"-c:v", "libx264",
			"-preset", preset.Preset,
			"-crf", preset.CRF,
			"-c:a", "aac",
			"-b:a", "128k",
			outputPath,
		)
	} else {
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}

	return cmd, nil
}

func Convert(inputPath, format string) (string, error) {
	if err := ValidateInput(inputPath); err != nil {
		return "", err
	}

	outputPath := GenerateOutputPath(inputPath, format)

	cmdArgs, err := BuildFFmpegCommand(inputPath, outputPath, format)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", errors.New("ffmpeg not found. ensure it is installed and in PATH")
		}
		return "", fmt.Errorf("ffmpeg conversion failed: %s", string(output))
	}

	return fmt.Sprintf("Successfully converted %s to %s", inputPath, outputPath), nil
}
