package physarum

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Video struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
	// stdout      io.ReadCloser
	// stdout_scan *bufio.Scanner
	// stderr      io.ReadCloser
	// stderr_scan *bufio.Scanner
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func NewVideo() *Video {
	// fmt.Println("Making new Video struct")
	v := &Video{}
	v.StartVideo()
	return v
}

func (v *Video) StartVideo() {
	// v.cmd = exec.Command("ffmpeg", "-h")
	v.cmd = exec.Command(
		"ffmpeg.exe",
		"-y",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-s", "2048x1024",
		"-r", "60",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-profile:v", "high",
		"-preset", "veryslow",
		"-tune", "film",
		"-bf", "2",
		"-rc-lookahead", "2",
		"-g", "30",
		"-crf", "20",
		"-pix_fmt", "yuv420p",
		// "-movflags", "+faststart",
		"out.mp4",
	)

	// Set up pipe to send data to FFMPEG
	stdin, err := v.cmd.StdinPipe()
	check(err)
	v.stdin = stdin

	// Redirect both stdout and stderr from the process to the console
	v.cmd.Stderr = os.Stderr
	v.cmd.Stdout = os.Stdout

	// Start the process
	check(v.cmd.Start())
}

func (v *Video) SaveVideoFfmpeg(videoFameChann <-chan []uint8, videoDoneChann chan<- bool) {
	for {
		frame, more := <-videoFameChann
		if more {
			// fmt.Println("frame")
			_, err := v.stdin.Write(frame)
			check(err)
		} else {
			fmt.Println("Video Frame Channel Closed, all done!")
			// Close the pipe, this should end the process?
			check(v.stdin.Close())

			v.cmd.Wait()           // Wait for the process to complete
			videoDoneChann <- true // Send sync signal
			return
		}
	}
}
