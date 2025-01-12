package cli

import (
	"github.com/korol8484/gophkeeper/internal/client/bubble"
	"github.com/korol8484/gophkeeper/internal/client/service"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "clx",
		Short:   "use clx",
		Version: "1",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := newConfig()
			if err != nil {
				return err
			}

			return bubble.Run(service.NewClient(cfg.client))
		},
	}

	return rootCmd
}
