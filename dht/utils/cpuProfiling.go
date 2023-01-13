package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/testground/sdk-go/runtime"
)

func runPs(file *os.File) error {
	pid := os.Getpid()

	cmd := exec.Command(
		"ps",
		"-p",
		fmt.Sprintf("%d", pid),
		"-o",
		"pid,tid,psr,pcpu",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	strs := strings.Split(string(out), "\n")

	_, err = file.Write([]byte(strs[1] + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func runPsRoutine(file *os.File, runenv *runtime.RunEnv) {
	time.Sleep(time.Second)
	timer := time.NewTicker(time.Second)
	for {
		select {
		case <-timer.C:
			err := runPs(file)
			if err != nil {
				runenv.RecordMessage("runPsRoutine: %s", err)
			}
			return
		}
	}

}

func CollectCpuUsage(runenv *runtime.RunEnv) error {
	cpuprofile := "" // TODO: add flag

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}

		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not start CPU profile: %w", err)
		}

		defer pprof.StopCPUProfile()
	}

	// TODO: add flag
	psFile, err := os.Create("psfile.out")
	if err != nil {
		return err
	}

	defer psFile.Close()

	go runPsRoutine(psFile, runenv)

	return nil
}
