package solarlog_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/torfjor/solarlog"
)

type doerFunc func(*http.Request) (*http.Response, error)

func (d doerFunc) Do(r *http.Request) (*http.Response, error) {
	return d(r)
}

func stubDayValuesDoer(path string) doerFunc {
	return func(r *http.Request) (*http.Response, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       f,
		}, nil
	}
}

func TestClient_CurrentDayValues(t *testing.T) {
	client := solarlog.NewClient("", "", 0, solarlog.WithDoer(stubDayValuesDoer("./testdata/current_day_values.json")))
	if _, err := client.CurrentDayValues(context.Background()); err != nil {
		t.Fatalf("CurrentDayValues() err=%v", err)
	}
}

func TestClient_DayValues(t *testing.T) {
	to := time.Date(2021, time.September, 20, 0, 0, 0, 0, time.UTC)
	from := to.AddDate(-1, 0, 0)

	client := solarlog.NewClient("", "", 0, solarlog.WithDoer(stubDayValuesDoer("./testdata/current_day_values.json")))
	if _, err := client.DayValues(context.Background(), from, to); err != nil {
		t.Fatalf("DayValues() err=%v", err)
	}
}
