package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	clock "github.com/nzoschke/midiclock"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"golang.org/x/xerrors"
)

func main() {
	if err := mainE(); err != nil {
		log.Fatalf("ERROR: %+v\n", err)
	}
}

func mainE() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	out, err := midi.FindOutPort("IAC Driver Bus 1")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	send, err := midi.SendTo(out)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	send(midi.Start())

	bpm := 160.0
	tickNum := -1
	tickTime := time.Now()
	ticker := time.NewTicker(clock.Dur(bpm))

	for {
		select {
		case s := <-sig:
			slog.Info("main", "signal", s)
			return nil
		case t := <-ticker.C:
			if err := send(midi.TimingClock()); err != nil {
				return xerrors.Errorf(": %w", err)
			}

			dur := t.Sub(tickTime)
			tickTime = t
			tickNum++
			if tickNum%24 == 0 {
				slog.Info("main", "send", "clock", "bpm", fmt.Sprintf("%.1f", bpm), "beat", tickNum/24, "dur", dur)
			}
		}
	}
}
