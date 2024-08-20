package proxy

import (
	"context"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"tailscale.com/tsnet"
)

type Config struct {
	Hostname  string
	Ephemeral bool
	DataDir   string
	Listeners []Listener
	Funnel    bool
	Debug     bool
}

func (c *Config) Run(ctx context.Context) error {
	server := c.getServer()

	status, err := server.Up(ctx)
	if err != nil {
		log.Error().Err(err).Str("hostname", c.Hostname).Str("data-dir", c.DataDir).Msg("Could not start server")
		return err
	}

	defer server.Close()

	log.Info().Str("hostname", status.Self.DNSName[:len(status.Self.DNSName)-1]).Msg("Server started successfully")

	for _, listener := range c.Listeners {
		err := listener.Start(ctx, &server, c.Funnel)
		if err != nil {
			log.Error().Err(err).Msg("Could not start listener")
			return err
		}
	}

	<-ctx.Done()

	return nil
}

func (c *Config) resolveDataDir() string {
	dir := c.DataDir
	if dir == "" {
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			log.Warn().Err(err).Msg("Could not get user config dir, falling back on using the current directory.")
			cfgDir = "./"
		}

		dir = path.Join(cfgDir, "tailservice", c.Hostname)
	}

	return dir
}

func (c *Config) getServer() tsnet.Server {
	return tsnet.Server{
		Hostname:  c.Hostname,
		Ephemeral: c.Ephemeral,
		Dir:       c.resolveDataDir(),
		Logf: func(f string, args ...any) {
			if !isDebugLog(f, args...) {
				log.Info().Msgf(f, args...)
			} else if c.Debug {
				log.Debug().Msgf(f, args...)
			}
		},
	}
}
