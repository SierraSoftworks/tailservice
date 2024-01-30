package proxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	humane "github.com/sierrasoftworks/humane-errors-go"
	"tailscale.com/tsnet"
)

func (l *Listener) listenHttp(ctx context.Context, srv *tsnet.Server, listener net.Listener) {
	defer listener.Close()

	httpClient := http.DefaultClient

	// Detect if we should use the tailscale HTTP client for these requests
	uri, err := url.Parse(l.Target)
	if err != nil {
		if isTailscaleHost(uri.Hostname()) {
			httpClient = srv.HTTPClient()
		}
	}

	http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.Clone(ctx)
		resp, err := httpClient.Do(req)

		if err != nil {
			herr := humane.Wrap(
				err,
				"Could not forward the request to the target service.",
				"Make sure that the target service is still running and reachable.",
				"Make sure that your client is still connected to the Tailnet and re-authenticate if necessary.",
				"Make sure that you've specified the correct target service in your listener configuration.",
			)

			log.Error().Err(err).Msg(herr.Display())
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}

		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			herr := humane.Wrap(
				err,
				"Could not forward the response from the target service.",
				"Make sure that the target service is still running and reachable.",
			)

			log.Error().Err(err).Msg(herr.Display())
		} else {
			log.Debug().Str("method", r.Method).Str("path", r.URL.Path).Int("status", resp.StatusCode).Msg("Forwarded request")
		}

	}))
}
