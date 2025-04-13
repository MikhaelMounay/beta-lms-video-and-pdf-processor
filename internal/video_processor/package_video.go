package videoprocessor

import (
	"fmt"
	"os"
	"os/exec"
)

func PackageVideoFile(vFile *VideoFile, vodKeyId string, vodKey string, iv string) error {
	packagerPath, err := GetPackagerExecutable()
	defer os.Remove(packagerPath)

	if err != nil {
		fmt.Println("Error getting Shaka Packager:", err)
		os.Exit(1)
	}

	// Construct the packager command
	cmd := exec.Command(packagerPath,
		fmt.Sprintf("in=%s,stream=video,output=%s,input_format=mp4,output_format=mp4,drm_label=HD", vFile.GetPath(), vFile.GetNewPath()+"_v.mp4"),
		fmt.Sprintf("in=%s,stream=audio,output=%s,input_format=mp4,output_format=mp4,drm_label=AUDIO", vFile.GetPath(), vFile.GetNewPath()+"_a.mp4"),
		"--enable_raw_key_encryption",
		"--keys", fmt.Sprintf("key_id=%s:key=%s", vodKeyId, vodKey),
		"--iv", iv,
		"--protection_scheme", "cenc",
		"--clear_lead", "0",
		fmt.Sprintf("--mpd_output=%s", vFile.GetNewPath()+"_m.mpd"),
	)

	// Uncomment the following lines to see the output of the packager
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// fmt.Print("\n")

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running packager:", err)
		os.Exit(1)
	}

	return nil
}
