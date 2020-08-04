package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	pubcfg "github.com/buildpacks/pack/config"

	"github.com/buildpacks/pack"
	"github.com/buildpacks/pack/internal/config"
	"github.com/buildpacks/pack/internal/style"
	"github.com/buildpacks/pack/logging"
)

func Rebase(logger logging.Logger, cfg config.Config, client PackClient) *cobra.Command {
	var opts pack.RebaseOptions
	var noPull bool
	var policy string

	cmd := &cobra.Command{
		Use:   "rebase <image-name>",
		Args:  cobra.ExactArgs(1),
		Short: "Rebase app image with latest run image",
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			opts.RepoName = args[0]
			opts.AdditionalMirrors = getMirrors(cfg)

			if cmd.Flags().Changed("no-pull") {
				logger.Warn("Flag --no-pull has been deprecated, please use `--pull-policy never` instead")

				if cmd.Flags().Changed("pull-policy") {
					logger.Warn("Flag --no-pull ignored in favor of --pull-policy")
				} else if noPull {
					policy = "never"
				}
			}

			var err error
			opts.PullPolicy, err = pubcfg.ParsePullPolicy(policy)
			if err != nil {
				return errors.Wrapf(err, "parse pull policy %s", policy)
			}

			if err := client.Rebase(cmd.Context(), opts); err != nil {
				return err
			}
			logger.Infof("Successfully rebased image %s", style.Symbol(opts.RepoName))
			return nil
		}),
	}

	cmd.Flags().BoolVar(&opts.Publish, "publish", false, "Publish to registry")
	cmd.Flags().StringVar(&policy, "pull-policy", "", "pull policy to use")
	cmd.Flags().StringVar(&opts.RunImage, "run-image", "", "Run image to use for rebasing")
	AddHelpFlag(cmd, "rebase")

	cmd.Flags().BoolVar(&noPull, "no-pull", false, "Skip pulling app and run images before use")
	cmd.Flags().MarkHidden("no-pull")

	return cmd
}
