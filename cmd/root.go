package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/sierrasoftworks/tailservice/proxy"
	"github.com/spf13/cobra"
)

var config proxy.Config

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "tailservice",
	Short: "Easily expose services on your Tailscale network.",
	Long: `Tailservice uses the tsnet library to expose a dedicated service on your
	Tailscale network as its own individual node. This allows you to easily access
	it using the corresponding service name (if you have MagicDNS enabled).`,
	Version: Version,
	Example: `tailservice --name my-service --tcp 80:8080 --udp 53:8.8.4.4:53 --tls 443:example.com:80 --tls 8443:https://example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		listeners := []proxy.Listener{}

		ls, err := cmd.Flags().GetStringArray("tcp")
		if err != nil {
			log.Warn().Err(err).Msg("Could not get TCP listeners")
		} else {
			listeners = append(listeners, parseListeners(ls, "tcp", false)...)
		}

		ls, err = cmd.Flags().GetStringArray("tls")
		if err != nil {
			log.Warn().Err(err).Msg("Could not get TLS listeners")
		} else {
			listeners = append(listeners, parseListeners(ls, "tcp", true)...)
		}

		ls, err = cmd.Flags().GetStringArray("udp")
		if err != nil {
			log.Warn().Err(err).Msg("Could not get UDP listeners")
		} else {
			listeners = append(listeners, parseListeners(ls, "udp", false)...)
		}

		config.Listeners = listeners

		if len(config.Listeners) == 0 {
			log.Fatal().Msg("No listeners were specified, specify at least one using --tcp, --tls or --udp")
		}

		err = config.Run(cmd.Context())
		if err != nil {
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&config.Debug, "ts-debug", false, "Enable debug logging of all Tailscale operations.")
	rootCmd.Flags().StringVar(&config.Hostname, "name", "", "Name of the service to expose on your tailnet.")
	rootCmd.Flags().BoolVar(&config.Ephemeral, "ephemeral", false, "Create the service ephemerally (remove it when the app closes).")
	rootCmd.Flags().StringVar(&config.DataDir, "data-dir", "", "Directory to store the service's connection data.")
	rootCmd.Flags().StringArray("tcp", []string{}, "TCP listeners to expose on the service.")
	rootCmd.Flags().StringArray("tls", []string{}, "TLS listeners to expose on the service.")
	rootCmd.Flags().StringArray("udp", []string{}, "UDP listeners to expose on the service.")
}

func parseListeners(listeners []string, proto string, secure bool) []proxy.Listener {
	ls := []proxy.Listener{}

	for _, listener := range listeners {
		l, err := proxy.ParseListener(listener, proto, secure)
		if err != nil {
			log.Error().Err(err).Str("listener", listener).Str("proto", proto).Bool("secure", secure).Msg("Could not parse listener")
			os.Exit(1)
		}

		ls = append(ls, *l)
	}

	return ls
}
