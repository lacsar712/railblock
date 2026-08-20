package app_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lacsar712/railblock/internal/app"
	"github.com/lacsar712/railblock/internal/config"
	"github.com/lacsar712/railblock/internal/ingest"
)

func testState(t *testing.T) *app.State {
	t.Helper()
	cfg := config.Defaults()
	st := app.NewState(cfg)
	return st
}

func TestHealthz(t *testing.T) {
	st := testState(t)
	srv := httptest.NewServer(app.NewRouter(st).Handler(st))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", res.StatusCode)
	}
}

func TestIngestAndQuery(t *testing.T) {
	st := testState(t)
	h := app.NewRouter(st).Handler(st)
	ts := httptest.NewServer(h)
	defer ts.Close()

	encBody := `{"block_id":12,"occupied":1,"seq":100}`
	encRes, err := http.Post(ts.URL+"/v1/frames/encode", "application/json", bytes.NewBufferString(encBody))
	if err != nil {
		t.Fatal(err)
	}
	defer encRes.Body.Close()
	var encData struct {
		Hex string `json:"hex"`
	}
	_ = json.NewDecoder(encRes.Body).Decode(&encData)

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/frames", bytes.NewBuffer(mustHex(encData.Hex)))
	req.Header.Set("Content-Type", ingest.MimeOctets)
	req.Header.Set("X-Source-ID", "station-a")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status=%d body=%s", resp.StatusCode, body)
	}
	resp.Body.Close()

	blockRes, err := http.Get(ts.URL + "/v1/blocks/12")
	if err != nil {
		t.Fatal(err)
	}
	defer blockRes.Body.Close()
	if blockRes.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", blockRes.StatusCode)
	}
}

func TestClearanceAPI(t *testing.T) {
	st := testState(t)
	ts := httptest.NewServer(app.NewRouter(st).Handler(st))
	defer ts.Close()

	body := `{"route_id":"t","blocks":[1,2]}`
	res, err := http.Post(ts.URL+"/v1/clearance/check", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", res.StatusCode)
	}
}

func TestIndexPage(t *testing.T) {
	st := testState(t)
	ts := httptest.NewServer(app.NewRouter(st).Handler(st))
	defer ts.Close()
	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", res.StatusCode)
	}
}

func mustHex(h string) []byte {
	out := make([]byte, len(h)/2)
	for i := 0; i < len(h); i += 2 {
		var v byte
		for _, c := range []byte{h[i], h[i+1]} {
			v <<= 4
			switch {
			case c >= '0' && c <= '9':
				v |= c - '0'
			case c >= 'a' && c <= 'f':
				v |= c - 'a' + 10
			case c >= 'A' && c <= 'F':
				v |= c - 'A' + 10
			}
		}
		out[i/2] = v
	}
	return out
}
