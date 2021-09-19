package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/torfjor/solarlog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		userFlag     = flags.String("user", "", "Solarlog username")
		pwFlag       = flags.String("password", "", "Solarlog password")
		solarlogFlag = flags.Int("solarlog", 0, "Solarlog ID")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	client := solarlog.NewClient(*userFlag, *pwFlag, *solarlogFlag)

	w := csv.NewWriter(os.Stdout)

	values, err := client.CurrentDayValues(ctx)
	if err != nil {
		return err
	}

	w.Write([]string{"date", "type", "channel", "description", "unit", "value"})
	for _, v := range values {
		for date, va := range v.Values {
			c, ok := v.Channel(va.Channel())
			_ = ok
			w.Write([]string{date, v.Device.Type, va.Channel(), c.Description, c.Unit, va.Value()})
		}
	}

	w.Flush()
	return w.Error()
}
