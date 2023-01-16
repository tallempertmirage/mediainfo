package mediainfo

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

var mediainfoBinary = flag.String("mediainfo-bin", "mediainfo", "the path to the mediainfo binary if it is not in the system $PATH")

type MediaInfoJson struct {
	Media struct {
		Track []Track `json:"track"`
	} `json:"media"`
}

type Track struct {
	Type                  string `json:"@type"`
	ID                    string `json:"ID"`
	VideoCount            string `json:"VideoCount,omitempty"`
	AudioCount            string `json:"AudioCount,omitempty"`
	FileExtension         string `json:"FileExtension,omitempty"`
	Format                string `json:"Format"`
	FileSize              string `json:"FileSize,omitempty"`
	Duration              string `json:"Duration"`
	OverallBitRateMode    string `json:"OverallBitRate_Mode,omitempty"`
	OverallBitRate        string `json:"OverallBitRate,omitempty"`
	FrameRate             string `json:"FrameRate"`
	FrameCount            string `json:"FrameCount,omitempty"`
	FileModifiedDate      string `json:"File_Modified_Date,omitempty"`
	FileModifiedDateLocal string `json:"File_Modified_Date_Local,omitempty"`
	Extra                 struct {
		OverallBitRatePrecisionMin string `json:"OverallBitRate_Precision_Min"`
		OverallBitRatePrecisionMax string `json:"OverallBitRate_Precision_Max"`
	} `json:"extra,omitempty"`
	StreamOrder                    string `json:"StreamOrder,omitempty"`
	MenuID                         string `json:"MenuID,omitempty"`
	FormatProfile                  string `json:"Format_Profile,omitempty"`
	FormatLevel                    string `json:"Format_Level,omitempty"`
	FormatSettingsCABAC            string `json:"Format_Settings_CABAC,omitempty"`
	FormatSettingsRefFrames        string `json:"Format_Settings_RefFrames,omitempty"`
	CodecID                        string `json:"CodecID,omitempty"`
	BitRateMode                    string `json:"BitRate_Mode,omitempty"`
	BitRateNominal                 string `json:"BitRate_Nominal,omitempty"`
	BitRateMaximum                 string `json:"BitRate_Maximum,omitempty"`
	Width                          string `json:"Width,omitempty"`
	Height                         string `json:"Height,omitempty"`
	StoredHeight                   string `json:"Stored_Height,omitempty"`
	SampledWidth                   string `json:"Sampled_Width,omitempty"`
	SampledHeight                  string `json:"Sampled_Height,omitempty"`
	PixelAspectRatio               string `json:"PixelAspectRatio,omitempty"`
	DisplayAspectRatio             string `json:"DisplayAspectRatio,omitempty"`
	ColorSpace                     string `json:"ColorSpace,omitempty"`
	ChromaSubsampling              string `json:"ChromaSubsampling,omitempty"`
	BitDepth                       string `json:"BitDepth,omitempty"`
	ScanType                       string `json:"ScanType,omitempty"`
	Delay                          string `json:"Delay,omitempty"`
	EncodedLibrary                 string `json:"Encoded_Library,omitempty"`
	EncodedLibraryName             string `json:"Encoded_Library_Name,omitempty"`
	EncodedLibraryVersion          string `json:"Encoded_Library_Version,omitempty"`
	EncodedLibrarySettings         string `json:"Encoded_Library_Settings,omitempty"`
	BufferSize                     string `json:"BufferSize,omitempty"`
	ColourDescriptionPresent       string `json:"colour_description_present,omitempty"`
	ColourDescriptionPresentSource string `json:"colour_description_present_Source,omitempty"`
	ColourRange                    string `json:"colour_range,omitempty"`
	ColourRangeSource              string `json:"colour_range_Source,omitempty"`
	ColourPrimaries                string `json:"colour_primaries,omitempty"`
	ColourPrimariesSource          string `json:"colour_primaries_Source,omitempty"`
	TransferCharacteristics        string `json:"transfer_characteristics,omitempty"`
	TransferCharacteristicsSource  string `json:"transfer_characteristics_Source,omitempty"`
	MatrixCoefficients             string `json:"matrix_coefficients,omitempty"`
	MatrixCoefficientsSource       string `json:"matrix_coefficients_Source,omitempty"`
	FormatVersion                  string `json:"Format_Version,omitempty"`
	FormatAdditionalFeatures       string `json:"Format_AdditionalFeatures,omitempty"`
	MuxingMode                     string `json:"MuxingMode,omitempty"`
	Channels                       string `json:"Channels,omitempty"`
	ChannelPositions               string `json:"ChannelPositions,omitempty"`
	ChannelLayout                  string `json:"ChannelLayout,omitempty"`
	SamplesPerFrame                string `json:"SamplesPerFrame,omitempty"`
	SamplingRate                   string `json:"SamplingRate,omitempty"`
	SamplingCount                  string `json:"SamplingCount,omitempty"`
	CompressionMode                string `json:"Compression_Mode,omitempty"`
	DelaySource                    string `json:"Delay_Source,omitempty"`
	Language                       string `json:"Language,omitempty"`
}

type MediaInfo struct {
	General []Track
	Video   []Track
	Audio   []Track
}

func IsInstalled() bool {
	cmd := exec.Command(*mediainfoBinary)
	err := cmd.Run()
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") ||
			strings.HasSuffix(err.Error(), "executable file not found in %PATH%") ||
			strings.HasSuffix(err.Error(), "executable file not found in $PATH") {
			return false
		} else if strings.HasPrefix(err.Error(), "exit status 255") {
			return true
		}
	}
	return true
}

func (info MediaInfo) IsMedia() bool {
	return info.Video[0].Duration != "" && info.Audio[0].Duration != ""
}

func GetMediaInfo(fname string, outputType string) (MediaInfo, error) {
	info := MediaInfoJson{}
	mediaInfo := MediaInfo{}

	if !IsInstalled() {
		return mediaInfo, fmt.Errorf("Must install mediainfo")
	}
	out, err := exec.Command(*mediainfoBinary, "--Output=JSON", "-f", fname).Output()

	if err != nil {
		return mediaInfo, err
	}

	if err := json.Unmarshal(out, &info); err != nil {
		return mediaInfo, err
	}

	for _, t := range info.Media.Track {
		if t.Type == "Video" {
			mediaInfo.Video = append(mediaInfo.Video, t)
		}
		if t.Type == "Audio" {
			mediaInfo.Audio = append(mediaInfo.Audio, t)
		}
		if t.Type == "General" {
			mediaInfo.General = append(mediaInfo.General, t)
		}
	}

	if !mediaInfo.IsMedia() {
		return mediaInfo, errors.New("The media file is invalid: no duration")
	}

	return mediaInfo, nil
}
