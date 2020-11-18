package datatype

import "fmt"

// ScienceGoal structs local goals and success criteria
type ScienceGoal struct {
	ID         string     `yaml:"id"`
	Name       string     `yaml:"name,omitempty"`
	SubGoals   []*SubGoal `yaml:"subgoals,omitempty"`
	Conditions []string   `yaml:"conditions,omitempty"`
}

// GetMySubGoal returns the subgoal assigned to node
func (g *ScienceGoal) GetMySubGoal(nodeName string) *SubGoal {
	for _, subGoal := range g.SubGoals {
		if subGoal.Node.Name == nodeName {
			return subGoal
		}
	}
	return nil
}

// SubGoal structs node-specific goal along with conditions and rules
type SubGoal struct {
	Node       *Node     `yaml:"node,omitempty"`
	Plugins    []*Plugin `yaml:"plugins,omitempty"`
	Rules      []string  `yaml:"rules,omitempty"`
	Statements []string  `yaml:"statements,omitempty"`
}

// UpdatePluginContext updates plugin's context event within the subgoal
// It returns an error if it fails to update context status of the plugin
func (sg *SubGoal) UpdatePluginContext(contextEvent EventPluginContext) error {
	for _, plugin := range sg.Plugins {
		if plugin.Name == contextEvent.PluginName {
			plugin.ContextStatus = contextEvent.Status
			return nil
		}
	}
	return fmt.Errorf("Failed to update context (%s) of plugin %s", contextEvent.Status, contextEvent.PluginName)
}

// GetSchedulablePlugins returns a list of plugins that are schedulable.
// A plugin is schedulable when its ContextStatus is Runnable and
// SchedulingStatus is not Running
func (sg *SubGoal) GetSchedulablePlugins() (schedulable []*Plugin) {
	for _, plugin := range sg.Plugins {
		if plugin.ContextStatus == Runnable &&
			plugin.SchedulingStatus != Running {
			schedulable = append(schedulable, plugin)
		}
	}
	return
}

// type Goal struct {
// 	APIVersion string `yaml:"apiVersion"`
// 	Header     struct {
// 		GoalId      string   `yaml:"goalId"`
// 		GoalName    string   `yaml:"goalName"`
// 		Priority    int      `yaml:"priority"`
// 		TargetNodes []string `yaml:"targetNodes"`
// 		UserId      string   `yaml:"userId"`
// 	}
// 	Body struct {
// 		AppConfig    []PluginConfig `yaml:"appConfig"`
// 		Rules        []string       `yaml:"rules"`
// 		SensorConfig struct {
// 			Plugins []Plugin `yaml:"plugins"`
// 		}
// 	}
// }
