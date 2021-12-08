package cmd

import (
	"fmt"

	"github.com/sagecontinuum/ses/pkg/logger"
	"github.com/sagecontinuum/ses/pkg/pluginctl"
	"github.com/sagecontinuum/ses/pkg/runplugin"
	"github.com/spf13/cobra"
)

func init() {
	cmdDeploy.Flags().StringVarP(&name, "name", "n", "", "Specify plugin name")
	cmdDeploy.Flags().StringVar(&node, "node", "", "run plugin on node")
	cmdDeploy.Flags().StringVarP(&job, "job", "j", "sage", "Specify job name")
	cmdDeploy.Flags().StringVar(&selectorStr, "selector", "", "Specify where plugin can run")
	cmdDeploy.Flags().BoolVarP(&privileged, "privileged", "p", false, "Deploy as privileged plugin")
	rootCmd.AddCommand(cmdDeploy)
}

var cmdDeploy = &cobra.Command{
	Use:              "deploy [FLAGS] PLUGIN_IMAGE [-- PLUGIN ARGUMENTS]",
	Short:            "Deploy a plugin",
	TraverseChildren: true,
	Args:             cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug.Printf("kubeconfig: %s", kubeconfig)
		logger.Debug.Printf("name: %s", name)
		logger.Debug.Printf("job: %s", job)
		logger.Debug.Printf("selector: %s", selectorStr)
		logger.Debug.Printf("args: %v", args)
		pluginCtl, err := pluginctl.NewPluginCtl(kubeconfig)
		if err != nil {
			logger.Error.Println(err.Error())
		}
		selector, err := pluginctl.ParseSelector(selectorStr)
		if err != nil {
			logger.Error.Println(err.Error())
		}
		spec := &runplugin.Spec{
			Privileged: privileged,
			Node:       node,
			Image:      args[0],
			Args:       args[1:],
			Job:        job,
			Name:       name,
			Selector:   selector,
		}
		if err = pluginCtl.Deploy(spec); err != nil {
			logger.Error.Println(err.Error())
		}

		fmt.Printf("Launched the plugin %s successfully\n", name)
		fmt.Printf("You may check the log: pluginctl log %s\n", name)
		fmt.Printf("To terminate the job: pluginctl rm %s\n", name)
	},
}