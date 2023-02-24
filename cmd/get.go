package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	Property string
}

var (
	getOpt GetOptions
)

var getCmd = &cobra.Command{
	Use:   "get [resource]",
	Short: "Get resource by key",
	Long:  `Get a single resource by key`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		m := manifest

		log.Debug().Msg("validating resource type...")
		rt, err := m.ValidateResourceType(args[0])
		if err != nil {
			return err
		}

		resourceKey := args[1]

		log.Debug().Msgf("finding %s with key (%s)...", rt.Names.Singular, resourceKey)

		from := m.ResourcesOfType(rt)
		for _, p := range from {
			if p.Key == resourceKey {
				log.Debug().Msgf("generating output for %s...", p.Key)

				if getOpt.Property == "" {
					b, _ := json.Marshal(p)
					fmt.Println(string(b))
					return nil
				} else {
					if val, ok := p.Properties[getOpt.Property]; ok {
						fmt.Println(val)
						return nil
					}
				}
			}
		}

		os.Exit(1)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&getOpt.Property, "select", "s", "", "selects value of a property (eg. build)")
	getCmd.SetUsageFunc(customUsageFunc)
}
