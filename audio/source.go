package audio

import (
	"bytes"
	"io"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type AudioType int

var ctx = audio.NewContext(48000)

const (
	OGG = iota
	WAV
	MP3
)

type Object interface {
	Position() (float64, float64)
}

type AudioSource struct {
	source      Object
	controllers map[string]*controller
}

type stream struct {
	io.ReadSeeker
	pan float64
	buf []byte
}

type controller struct {
	stream *stream
	player *audio.Player
}

func newStream(src io.ReadSeeker) *stream {
	return &stream{
		ReadSeeker: src,
	}
}

func (as *AudioSource) Pan(key string) float64 {
	return as.controllers[key].stream.pan
}

func (as *AudioSource) SetPan(key string, pan float64) {
	as.controllers[key].stream.pan = math.Min(math.Max(-1, pan), 1)
}

func (s *stream) Read(p []byte) (int, error) {
	var bufN int
	if len(s.buf) > 0 {
		bufN = copy(p, s.buf)
		s.buf = s.buf[bufN:]
	}

	readN, err := s.ReadSeeker.Read(p[bufN:])
	if err != nil && err != io.EOF {
		return 0, err
	}

	totalN := bufN + readN
	extra := totalN - totalN/8*8
	s.buf = append(s.buf, p[totalN-extra:totalN]...)
	alignedN := totalN - extra

	ls := float32(math.Min(s.pan*-1+1, 1))
	rs := float32(math.Min(s.pan+1, 1))

	for i := 0; i < alignedN; i += 8 {
		lc := math.Float32frombits(uint32(p[i])|(uint32(p[i+1])<<8)|(uint32(p[i+2])<<16)|(uint32(p[i+3])<<24)) * ls
		rc := math.Float32frombits(uint32(p[i+4])|(uint32(p[i+5])<<8)|(uint32(p[i+6])<<16)|(uint32(p[i+7])<<24)) * rs

		lcBits := math.Float32bits(lc)
		rcBits := math.Float32bits(rc)
		p[i] = byte(lcBits)
		p[i+1] = byte(lcBits >> 8)
		p[i+2] = byte(lcBits >> 16)
		p[i+3] = byte(lcBits >> 24)
		p[i+4] = byte(rcBits)
		p[i+5] = byte(rcBits >> 8)
		p[i+6] = byte(rcBits >> 16)
		p[i+7] = byte(rcBits >> 24)
	}

	return alignedN, err
}

func NewSource(obj Object) *AudioSource {
	return &AudioSource{
		source:      obj,
		controllers: make(map[string]*controller, 1),
	}
}
func (as *AudioSource) Play(key string) {
	as.controllers[key].player.Play()
}

func (as *AudioSource) Playing(key string) bool {
	return as.controllers[key].player.IsPlaying()
}
func (as *AudioSource) AddController(key string, data []byte, loop bool) error {
	s, err := vorbis.DecodeF32(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	stream := newStream(audio.NewInfiniteLoop(s, s.Length()))
	stream.pan = 0
	player, err := ctx.NewPlayerF32(stream)

	as.controllers[key] = &controller{
		stream: stream,
		player: player,
	}
	if err != nil {
		return err
	}
	return nil
}
