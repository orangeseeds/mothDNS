package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/orangeseeds/mothDNS/pkg/server"
	"github.com/spf13/cobra"
)

var (
	port       string
	rootServer string
	logPath    string

	rootCmd = &cobra.Command{
		Use:   "mothDNS",
		Short: "mothDNS is a simple recursive DNS server",
		Long:  `A basic recursive DNS server implementation that follows the conventions given by RFC 1035.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			var writers []io.Writer
			writers = append(writers, os.Stdout)
			if logPath != "" {
				fs, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					return err
				}
				writers = append(writers, fs)
			}
			outputs := io.MultiWriter(writers...)
			log.SetOutput(outputs)

			udpServer := server.UPDServer{
				RootServer: rootServer,
				Handler:    server.HandleConnection,
			}
			udpServer.Serve(port)

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", ":1053", "Port number to use to serve on.")
	rootCmd.PersistentFlags().StringVarP(&rootServer, "root", "r", "198.41.0.4", "Root server to connect to.")
	rootCmd.PersistentFlags().StringVar(&logPath, "log", "", "Path of the log file, if not set default log is stdout only")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
