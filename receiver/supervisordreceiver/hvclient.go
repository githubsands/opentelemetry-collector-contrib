package supervisordreceiver

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

var superVisordUnixReq = []byte("testtesttest")

// See supervisord API documents: http://supervisord.org/api.html
type superVisordResponse struct {
	Stats []*Statistic `xml:""`
}

type svClient struct {
	unixHTTP     *http.Client
	svUnixSocket string
	logger       *zap.Logger
	buffer       []byte
}

func newSVClient(cfg *Config, logger *zap.Logger) (*svClient, error) {
	if logger == nil {
		logger.Panic("hypervisord client requires a logger")
	}
	ok := strings.HasPrefix(cfg.SvUnixSocket, "unix")
	if !ok {
		return nil, errors.New("HvUnixSocket host should have unix prefixj")
	}
	svClient := new(svClient)
	svClient.buffer = make([]byte, 64)
	svClient.unixHTTP = new(http.Client)
	svClient.unixHTTP.Transport = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", svClient.svUnixSocket)
		}}
	return svClient, nil
}

func (hvc *svClient) getStats(stats []*Statistic) []*Statistic {
	buf := hvc.do()
	resp := superVisordResponse{}
	xml.Unmarshal(buf, resp)
	return resp.Stats
}

func (hvc *svClient) do() []byte {
	resp, err := hvc.unixHTTP.Do(&http.Request{Method: "GET", Body: io.NopCloser(bytes.NewReader(superVisordUnixReq))})
	if err != nil {
		hvc.logger.Error(err.Error())
	}
	if resp.StatusCode != 200 {
		hvc.logger.Error(err.Error())
	}
	defer resp.Body.Close()
	buffer := bytes.NewBuffer(hvc.buffer)
	io.Copy(buffer, resp.Body)
	return buffer.Bytes()
}
