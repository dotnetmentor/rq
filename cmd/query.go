package cmd

import (
	"fmt"
	"sort"

	"github.com/dotnetmentor/rq/internal/pkg/output"
	"github.com/dotnetmentor/rq/internal/pkg/query"
	"github.com/dotnetmentor/rq/internal/pkg/schema"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type QueryOptions struct {
	Parameters            []string
	DisableStrictMatching bool
	Sort                  bool
	Output                string
}

var (
	queryOpt QueryOptions
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query resources",
	Long:  `Query resources`,
	Args:  cobra.ExactArgs(1),
	PostRun: func(cmd *cobra.Command, args []string) {
		resetQueryOpt()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug().Msg("reading manifest...")
		m, err := schema.NewManifest(opt.FilePath)
		if err != nil {
			return err
		}

		log.Debug().Msg("validating resource type...")
		rt, ok := m.ResourceType(args[0])
		if !ok {
			return fmt.Errorf("unknown resource type %s, valid resource types: %s", args[0], m.ResourceTypeNames())
		}

		conditions, err := query.ParseArgs(queryOpt.Parameters)
		if err != nil {
			return err
		}

		log.Debug().Msgf("querying %s based on conditions (%v) (strict-mode=%v)...", rt.Names.Plural, conditions, !queryOpt.DisableStrictMatching)
		from := m.ResourcesOfType(rt)
		results := from.Where(func(r schema.Resource) bool {
			if !queryOpt.DisableStrictMatching && len(r.Conditions) == 0 && len(conditions) > 0 {
				return false
			}

			cs := make([]query.Condition, 0)
			for _, c := range r.Conditions {
				cs = append(cs, query.Condition{
					Name:         "condition",
					MatchMissing: queryOpt.DisableStrictMatching,
					MatchRegExp:  true,
					Left:         conditions,
					Right:        c,
				})
			}

			match := query.OneOf(cs...)
			log.Debug().Str("type", rt.Names.Plural).Str("match", fmt.Sprintf("%v", match)).Msg(r.Key)
			return match
		})

		// Prepare list of keys
		keys := results.SelectMany(func(i schema.Resource) string {
			return i.Key
		})

		if queryOpt.Sort {
			sort.Strings(keys)
		}

		log.Debug().Msgf("writing output for %d of %d %s...", len(keys), len(from), rt.Names.Plural)
		output.WriteKeys(keys, queryOpt.Output)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringSliceVarP(&queryOpt.Parameters, "parameter", "p", []string{}, "parameter (eg. c:environment=prod)")
	queryCmd.Flags().BoolVarP(&queryOpt.DisableStrictMatching, "disable-strict-matching", "d", false, "disable strict matching (matches conditions where parameter is missing)")
	queryCmd.Flags().BoolVarP(&queryOpt.Sort, "sort", "s", false, "sort output (lexicographically)")
	queryCmd.Flags().StringVarP(&queryOpt.Output, "out", "o", output.Newline, "output type (eg. json/xargs)")
}

func resetQueryOpt() {
	queryOpt.DisableStrictMatching = false
	queryOpt.Parameters = []string{}
}
