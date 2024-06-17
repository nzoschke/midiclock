package main

import (
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

	tick := 0
	ticker := time.NewTicker(clock.Dur(165.0))

	send(midi.Start())

	for {
		select {
		case s := <-sig:
			slog.Info("main", "signal", s)
			return nil
		case <-ticker.C:
			if tick%24 == 0 {
				slog.Info("main", "send", "clock")
			}
			tick++

			if err := send(midi.TimingClock()); err != nil {
				return xerrors.Errorf(": %w", err)
			}
		}
	}
}
