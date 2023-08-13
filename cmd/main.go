package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Saturn745/OBS-Shuffle-Media/helper"
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/fatih/color"
	ffprobe "github.com/vansante/go-ffprobe"
)

type Config struct {
	Obs   ObsConfig   `json:"obs"`
	Media MediaConfig `json:"media"`
}

type ObsConfig struct {
	Host        string `json:"host"`
	Password    string `json:"password"`
	MediaSource string `json:"media-source"`
	TitleSource string `json:"title-source"`
}

type MediaConfig struct {
	Sources []string `json:"sources"`
	Font    string   `json:"font"`
}

var AppConfig Config

func loadConfig() {
	// Read config.json file
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error reading config file:", err)
		return
	}

	// Parse JSON into Config struct
	err = json.Unmarshal(file, &AppConfig)
	if err != nil {
		log.Fatal("Error parsing config file:", err)
		return
	}
}

var videoFiles []string
var OBSClient goobs.Client

func main() {
	loadConfig()

	color.Green("[OBS] Attempting to connect to OBS...")

	client, err := goobs.New(AppConfig.Obs.Host, goobs.WithPassword(AppConfig.Obs.Password))
	if err != nil {
		log.Fatal(err)
	}
	OBSClient = *client
	defer client.Disconnect()

	color.Green("[OBS] Connected!")

	version, err := client.General.GetVersion()
	if err != nil {
		log.Fatal(err)
	}
	color.Blue("OBS Studio version: %s\n", version.ObsVersion)
	color.Blue("Websocket server version: %s\n", version.ObsWebSocketVersion)

	color.Green("Scanning video directories... If you have lots of video's this may take a bit")
	for i := 0; i < len(AppConfig.Media.Sources); i++ {
		files := helper.ScanDirectory(AppConfig.Media.Sources[i])
		videoFiles = append(videoFiles, files...)
	}
	color.Green("Found %d videos!", len(videoFiles))
	playRandomVideo()
}

func playRandomVideo() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := random.Intn(len(videoFiles))
	randomVideo := videoFiles[randomIndex]
	data, err := ffprobe.GetProbeData(randomVideo, 120000*time.Millisecond)
	if err != nil {
		log.Fatal("Error getting media files data")
		return
	}
	color.Blue("Playing video file: %s", randomVideo)
	helper.GenerateImageWithWrappedText(AppConfig.Media.Font, helper.GetFileNameFromPath(randomVideo), "temp/title.png")
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	OBSClient.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{InputName: AppConfig.Obs.MediaSource, InputSettings: map[string]interface{}{"local_file": randomVideo}})
	if AppConfig.Obs.TitleSource != "" {
		OBSClient.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{InputName: AppConfig.Obs.TitleSource, InputSettings: map[string]interface{}{"image_file": currentPath + "/temp/title.png"}})
	}
	time.Sleep(time.Duration(data.Format.Duration().Seconds()) * time.Second)
	color.Blue("%s finished playing after %s seconds", randomVideo, data.Format.Duration().Seconds())
	playRandomVideo()
}
