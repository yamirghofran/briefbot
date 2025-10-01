package services

import (
	"github.com/stretchr/testify/mock"
)

type MockSpeechService struct {
	mock.Mock
}

func (m *MockSpeechService) TextToSpeech(prompt string, voice VoiceEnum, speed float64, dialogueInfo *DialogueInfo) (*File, error) {
	args := m.Called(prompt, voice, speed, dialogueInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*File), args.Error(1)
}

func (m *MockSpeechService) DownloadAudio(url string) ([]byte, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockSpeechService) SaveAudio(data []byte, filename string) error {
	args := m.Called(data, filename)
	return args.Error(0)
}
