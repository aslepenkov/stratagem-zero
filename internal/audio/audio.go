package audio

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Sound int

const (
	SoundKeypress Sound = iota
	SoundSuccess
	SoundFail
	SoundLaunch
)

type AudioPlayer interface {
	Play(Sound)
	Close()
}

type NoopAudioPlayer struct{}

func (n *NoopAudioPlayer) Play(s Sound) {}
func (n *NoopAudioPlayer) Close()     {}

// LinuxAudioPlayer plays WAV audio using system players (pw-play, paplay, aplay).
type LinuxAudioPlayer struct {
	playerBin string
	soundMap  map[Sound]string
	tempDir   string
	enabled   bool
	mu        sync.Mutex
}

func NewLinuxAudioPlayer(disabled bool, soundFS embed.FS) AudioPlayer {
	if disabled {
		return &NoopAudioPlayer{}
	}

	// Detect available CLI player
	var playerBin string
	for _, bin := range []string{"pw-play", "paplay", "aplay"} {
		if path, err := exec.LookPath(bin); err == nil {
			playerBin = path
			break
		}
	}

	if playerBin == "" {
		return &NoopAudioPlayer{}
	}

	tempDir, err := os.MkdirTemp("", "stratagem-zero-audio-*")
	if err != nil {
		return &NoopAudioPlayer{}
	}

	soundFiles := map[Sound]string{
		SoundKeypress: "key.wav",
		SoundSuccess:  "success.wav",
		SoundFail:     "fail.wav",
		SoundLaunch:   "launch.wav",
	}

	soundMap := make(map[Sound]string)

	entries, _ := soundFS.ReadDir(".")
	for _, entry := range entries {
		_ = entry
	}

	for sound, name := range soundFiles {
		// Read either directly or from subfolder if present
		data, err := soundFS.ReadFile(name)
		if err != nil {
			data, err = soundFS.ReadFile("sounds/" + name)
		}
		if err != nil {
			// Find match in embedded FS
			data = findAndReadFile(soundFS, name)
		}
		if len(data) > 0 {
			tmpFilePath := filepath.Join(tempDir, name)
			if err := os.WriteFile(tmpFilePath, data, 0600); err == nil {
				soundMap[sound] = tmpFilePath
			}
		}
	}

	return &LinuxAudioPlayer{
		playerBin: playerBin,
		soundMap:  soundMap,
		tempDir:   tempDir,
		enabled:   true,
	}
}

func findAndReadFile(soundFS embed.FS, target string) []byte {
	var found []byte
	var walk func(dir string)
	walk = func(dir string) {
		entries, err := soundFS.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			path := e.Name()
			if dir != "." {
				path = dir + "/" + e.Name()
			}
			if e.IsDir() {
				walk(path)
			} else if strings.HasSuffix(path, target) {
				b, err := soundFS.ReadFile(path)
				if err == nil {
					found = b
					return
				}
			}
		}
	}
	walk(".")
	return found
}

func (l *LinuxAudioPlayer) Play(s Sound) {
	if !l.enabled {
		return
	}

	filePath, ok := l.soundMap[s]
	if !ok || filePath == "" {
		return
	}

	// Non-blocking asynchronous playback
	go func() {
		cmd := exec.Command(l.playerBin, filePath)
		_ = cmd.Run()
	}()
}

func (l *LinuxAudioPlayer) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.tempDir != "" {
		_ = os.RemoveAll(l.tempDir)
		l.tempDir = ""
	}
	l.enabled = false
}

func (s Sound) String() string {
	switch s {
	case SoundKeypress:
		return "Keypress"
	case SoundSuccess:
		return "Success"
	case SoundFail:
		return "Fail"
	case SoundLaunch:
		return "Launch"
	default:
		return fmt.Sprintf("Sound(%d)", s)
	}
}
