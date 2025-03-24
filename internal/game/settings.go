package game

import (
	"errors"
	"github.com/goccy/go-yaml"
	"log/slog"
	"os"
)

// Settings struct to persist global game settings. Does not aim to be compatible with the old config file format.
// Instead, we serialize to a YAML file. We do reuse the original setting names, however.
var Settings = struct {
	MouseSensitivity int  `yaml:"mouse_sensitivity"`
	SfxVolume        int  `yaml:"sfx_volume"`
	MusicVolume      int  `yaml:"music_volume"`
	ShowMessages     bool `yaml:"show_messages"`

	KeyRight       int `yaml:"key_right"`
	KeyLeft        int `yaml:"key_left"`
	KeyUp          int `yaml:"key_up"`
	KeyDown        int `yaml:"key_down"`
	KeyStrafeLeft  int `yaml:"key_strafeleft"`
	KeyStrafeRight int `yaml:"key_straferight"`

	KeyFire   int `yaml:"key_fire"`
	KeyUse    int `yaml:"key_use"`
	KeyStrafe int `yaml:"key_strafe"`
	KeySpeed  int `yaml:"key_speed"`

	UseMouse           bool `yaml:"use_mouse"`
	MouseButtonFire    int  `yaml:"mouseb_fire"`
	MouseButtonStrafe  int  `yaml:"mouseb_strafe"`
	MouseButtonForward int  `yaml:"mouseb_forward"`

	UseJoystick          bool `yaml:"use_joystick"`
	JoystickButtonFire   int  `yaml:"joyb_fire"`
	JoystickButtonStrafe int  `yaml:"joyb_strafe"`
	JoystickButtonUse    int  `yaml:"joyb_use"`
	JoystickButtonSpeed  int  `yaml:"joyb_speed"`

	ScreenBlocks int `yaml:"screenblocks"`
	DetailLevel  int `yaml:"detaillevel"`

	SoundChannels int `yaml:"snd_channels"`

	UseGamma bool `yaml:"use_gamma"`

	ChatMacro0 string `yaml:"chat_macro0"`
	ChatMacro1 string `yaml:"chat_macro1"`
	ChatMacro2 string `yaml:"chat_macro2"`
	ChatMacro3 string `yaml:"chat_macro3"`
	ChatMacro4 string `yaml:"chat_macro4"`
	ChatMacro5 string `yaml:"chat_macro5"`
	ChatMacro6 string `yaml:"chat_macro6"`
	ChatMacro7 string `yaml:"chat_macro7"`
	ChatMacro8 string `yaml:"chat_macro8"`
	ChatMacro9 string `yaml:"chat_macro9"`
}{}

// TODO: Move to a separate file
const (
	KEY_RIGHTARROW = 0xae
	KEY_LEFTARROW  = 0xac
	KEY_UPARROW    = 0xad
	KEY_DOWNARROW  = 0xaf
	KEY_RSHIFT     = 0x80 + 0x36
	KEY_RCTRL      = 0x80 + 0x1d
	KEY_RALT       = 0x80 + 0x38
)

// TODO: Move to a separate file for english strings and provide separate french strings
const (
	HUSTR_CHATMACRO0 = "I'm ready to kick butt!"
	HUSTR_CHATMACRO1 = "I'm OK."
	HUSTR_CHATMACRO2 = "I'm not looking too good!"
	HUSTR_CHATMACRO3 = "Help!"
	HUSTR_CHATMACRO4 = "You suck!"
	HUSTR_CHATMACRO5 = "Next time, scumbag..."
	HUSTR_CHATMACRO6 = "Come here!"
	HUSTR_CHATMACRO7 = "I'll take care of it."
	HUSTR_CHATMACRO8 = "Yes"
	HUSTR_CHATMACRO9 = "No"
)

func initDefaults() {
	Settings.MouseSensitivity = 5
	Settings.SfxVolume = 8
	Settings.MusicVolume = 8
	Settings.ShowMessages = true

	Settings.KeyRight = KEY_RIGHTARROW
	Settings.KeyLeft = KEY_LEFTARROW
	Settings.KeyUp = KEY_UPARROW
	Settings.KeyDown = KEY_DOWNARROW
	Settings.KeyStrafeLeft = ','
	Settings.KeyStrafeRight = '.'

	Settings.KeyFire = KEY_RCTRL
	Settings.KeyUse = ' '
	Settings.KeyStrafe = KEY_RALT
	Settings.KeySpeed = KEY_RSHIFT

	Settings.UseMouse = true
	Settings.MouseButtonFire = 0
	Settings.MouseButtonStrafe = 1
	Settings.MouseButtonForward = 2

	Settings.UseJoystick = false
	Settings.JoystickButtonFire = 0
	Settings.JoystickButtonStrafe = 1
	Settings.JoystickButtonUse = 3
	Settings.JoystickButtonSpeed = 2

	Settings.ScreenBlocks = 9
	Settings.DetailLevel = 0

	Settings.SoundChannels = 3

	Settings.UseGamma = false

	Settings.ChatMacro0 = HUSTR_CHATMACRO0
	Settings.ChatMacro1 = HUSTR_CHATMACRO1
	Settings.ChatMacro2 = HUSTR_CHATMACRO2
	Settings.ChatMacro3 = HUSTR_CHATMACRO3
	Settings.ChatMacro4 = HUSTR_CHATMACRO4
	Settings.ChatMacro5 = HUSTR_CHATMACRO5
	Settings.ChatMacro6 = HUSTR_CHATMACRO6
	Settings.ChatMacro7 = HUSTR_CHATMACRO7
	Settings.ChatMacro8 = HUSTR_CHATMACRO8
	Settings.ChatMacro9 = HUSTR_CHATMACRO9
}

func InitSettings(configPath string) {
	// start with defaults
	initDefaults()

	// read the config file if exists to override defaults with persisted user settings
	file, err := os.Open(configPath)
	if errors.Is(err, os.ErrNotExist) {
		slog.Debug("No config file found. Using default settings.", "configPath", configPath)
		return
	}

	// otherwise read the file to load the settings
	if err != nil {
		slog.Error("Failed to open config file", "err", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&Settings); err != nil {
		slog.Error("Failed to decode config file", "err", err)
		os.Exit(1)
	}
	slog.Debug("Loaded settings from config file", "configPath", configPath)
	slog.Debug("Settings", "settings", Settings)
}

func SaveSettings(configPath string) {
	file, err := os.Create(configPath)
	if err != nil {
		slog.Error("Failed to create config file. Settings not saved.", "err", err)
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(Settings); err != nil {
		slog.Error("Failed to encode settings. Settings not saved.", "err", err)
		return
	}
	slog.Debug("Saved settings to config file", "configPath", configPath)
}
