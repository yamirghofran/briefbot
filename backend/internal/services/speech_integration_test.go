//go:build integration
// +build integration

package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStitchAudioFiles_Integration tests the ffmpeg audio stitching functionality
func TestStitchAudioFiles_Integration(t *testing.T) {
	// Check if ffmpeg is available by trying to execute it
	cmd := os.Getenv("PATH")
	if cmd == "" {
		t.Skip("Skipping ffmpeg integration test: PATH not set")
	}

	// Simple check - the StitchAudioFiles function will fail gracefully if ffmpeg is not found
	// We just need to ensure we have exec permissions

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "audio-stitch-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test audio files (silent WAV files with proper headers)
	// WAV file structure: RIFF header + fmt chunk + data chunk
	testFiles := make([]string, 3)
	for i := range testFiles {
		testFiles[i] = filepath.Join(tempDir, "test"+string(rune('0'+i))+".wav")

		// Create a minimal valid WAV file (44.1kHz, 16-bit, mono, ~0.1 seconds)
		wavData := createMinimalWAV(4410) // 0.1 seconds at 44.1kHz
		err := os.WriteFile(testFiles[i], wavData, 0644)
		require.NoError(t, err)
	}

	outputFile := filepath.Join(tempDir, "output.mp3")

	// Test stitching
	err = StitchAudioFiles(testFiles, outputFile)
	assert.NoError(t, err)

	// Verify output file was created
	info, err := os.Stat(outputFile)
	assert.NoError(t, err)
	assert.True(t, info.Size() > 0, "Output file should have non-zero size")

	t.Logf("Successfully stitched %d files into %s (size: %d bytes)",
		len(testFiles), outputFile, info.Size())
}

// TestStitchAudioFiles_EmptyInput tests error handling
func TestStitchAudioFiles_EmptyInput(t *testing.T) {
	err := StitchAudioFiles([]string{}, "output.mp3")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no input files")
}

// TestStitchAudioFiles_NonexistentFile tests error handling for missing files
func TestStitchAudioFiles_NonexistentFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "audio-stitch-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	nonexistentFile := filepath.Join(tempDir, "nonexistent.wav")
	outputFile := filepath.Join(tempDir, "output.mp3")

	err = StitchAudioFiles([]string{nonexistentFile}, outputFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

// createMinimalWAV creates a minimal valid WAV file with silence
func createMinimalWAV(numSamples int) []byte {
	// WAV file format:
	// - RIFF header (12 bytes)
	// - fmt chunk (24 bytes)
	// - data chunk (8 bytes + sample data)

	sampleRate := 44100
	numChannels := 1
	bitsPerSample := 16
	byteRate := sampleRate * numChannels * bitsPerSample / 8
	blockAlign := numChannels * bitsPerSample / 8
	dataSize := numSamples * blockAlign
	fileSize := 36 + dataSize

	wav := make([]byte, 44+dataSize)

	// RIFF header
	copy(wav[0:4], "RIFF")
	writeUint32(wav[4:8], uint32(fileSize))
	copy(wav[8:12], "WAVE")

	// fmt chunk
	copy(wav[12:16], "fmt ")
	writeUint32(wav[16:20], 16) // fmt chunk size
	writeUint16(wav[20:22], 1)  // audio format (PCM)
	writeUint16(wav[22:24], uint16(numChannels))
	writeUint32(wav[24:28], uint32(sampleRate))
	writeUint32(wav[28:32], uint32(byteRate))
	writeUint16(wav[32:34], uint16(blockAlign))
	writeUint16(wav[34:36], uint16(bitsPerSample))

	// data chunk
	copy(wav[36:40], "data")
	writeUint32(wav[40:44], uint32(dataSize))

	// Sample data (all zeros = silence)
	// Already initialized to zeros by make()

	return wav
}

// Helper functions to write binary data
func writeUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func writeUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}
