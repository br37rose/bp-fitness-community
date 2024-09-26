package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/spf13/cobra"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"

	cli_config "github.com/bci-innovation-labs/bp8fitnesscommunity-cli/config"
)

func init() {
	rootCmd.AddCommand(proxyCmd)
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts ngrok server to generate a public free address which https forwards to your local running developer server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BP8 Fitness Community Proxy Running using Ngrok")
		cfg := cli_config.New()
		ngrokAuthToken := cfg.Ngrok.AuthToken
		if err := run(context.Background(), ngrokAuthToken); err != nil {
			log.Fatal(err)
		}
	},
}

func run(ctx context.Context, ngrokAuthToken string) error {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(ngrokAuthToken),
	)
	if err != nil {
		return err
	}

	log.Println("App URL", listener.URL())
	return http.Serve(listener, http.HandlerFunc(proxyHandler))
}

func ngrokListener(ctx context.Context) (net.Listener, error) {
	// Special thanks: https://ngrok.com/docs/http/?cty=go-sdk#rewrite-host-header
	return ngrok.Listen(ctx,
		config.HTTPEndpoint(
			config.WithRequestHeader("host", "localhost"),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
}
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Create a request to the local API server
	backendURL := "http://localhost:8000" + r.URL.Path + "?" + r.URL.RawQuery
	log.Printf("Forwarding request to %s", backendURL)

	// Perform the request to the local API server
	client := http.Client{}
	backendResp, err := client.Get(backendURL)
	if err != nil {
		http.Error(w, "Failed to perform backend request", http.StatusInternalServerError)
		return
	}
	defer backendResp.Body.Close()

	// Copy the response from the local API server to the client
	w.WriteHeader(backendResp.StatusCode)
	w.Header().Set("Content-Type", backendResp.Header.Get("Content-Type"))
	io.Copy(w, backendResp.Body)
}

// Auto-generated comment for change 6
