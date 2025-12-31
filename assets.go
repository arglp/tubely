package main

import (
	"fmt"
	"os"
	"errors"
	"path/filepath"
	"strings"
	"crypto/rand"
	"encoding/base64"
	"os/exec"
	"bytes"
	"encoding/json"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(mediaType string) string {
	ext := mediaTypeToExt(mediaType)
	random := make([]byte, 32)
	rand.Read(random)
	fileName := base64.RawURLEncoding.EncodeToString(random)
	return fmt.Sprintf("%s%s", fileName, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func getVideoAspectRatio (filePath string) (string, error) {

	type stream struct {
		Width int `json:"width"`
		Height int `json:"height"`
	}
		
	type videoInformation struct {
		Streams []stream `json:"streams"`
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	videoInf := videoInformation{}
	err = json.Unmarshal(buf.Bytes(), &videoInf)
	width := videoInf.Streams[0].Width
	height := videoInf.Streams[0].Height
	if width == 0 || height == 0 {
		return "", errors.New("couldn't find video data")
	}

	div := float64(width) / float64(height)

	if div >= 1.7 && div <= 1.8 {
		return "16:9", nil
	}
	if div >= 0.5 && div <= 0.6 {
		return "9:16", nil
	}
	return "other", nil
}