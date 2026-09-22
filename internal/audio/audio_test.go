package audio_test

import (
	"embed"
	"testing"

	"stratagem-zero/internal/audio"
)

func TestAudioPlayer_NoCrash(t *testing.T) {
	var emptyFS embed.FS
	player := audio.NewLinuxAudioPlayer(true, emptyFS)
	player.Play(audio.SoundKeypress)
	player.Play(audio.SoundSuccess)
	player.Play(audio.SoundFail)
	player.Play(audio.SoundLaunch)
	player.Close()
}
