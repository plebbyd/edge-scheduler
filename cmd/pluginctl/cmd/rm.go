package cmd

import (
	"fmt"

	"github.com/sagecontinuum/ses/pkg/logger"
	"github.com/sagecontinuum/ses/pkg/pluginctl"
	"github.com/spf13/cobra"
)

func init() {
	cmdRm.Flags().BoolVarP(&followLog, "follow", "f", false, "Specified if logs should be streamed")
	rootCmd.AddCommand(cmdRm)
}

var cmdRm = &cobra.Command{
	Use:              "rm APP_NAME",
	Short:            "Remove plugin",
	TraverseChildren: true,
	Args:             cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug.Printf("kubeconfig: %s", kubeconfig)
		name = args[0]
		logger.Debug.Printf("args: %v", args)
		pluginCtl, err := pluginctl.NewPluginCtl(kubeconfig)
		if err != nil {
			logger.Error.Println(err.Error())
		}
		err = pluginCtl.Terminate(name)
		if err != nil {
			logger.Error.Printf("%s", err.Error())
		} else {
			fmt.Printf("Terminated the plugin %s successfully\n", name)
		}
	},
}