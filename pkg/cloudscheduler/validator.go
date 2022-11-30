package cloudscheduler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/waggle-sensor/edge-scheduler/pkg/datatype"
	"github.com/waggle-sensor/edge-scheduler/pkg/logger"
)

type JobValidator struct {
	dataPath string
	Plugins  map[string]*datatype.PluginManifest
	Nodes    map[string]*datatype.NodeManifest
}

func NewJobValidator(dataPath string) *JobValidator {
	return &JobValidator{
		dataPath: dataPath,
		Plugins:  make(map[string]*datatype.PluginManifest),
		Nodes:    make(map[string]*datatype.NodeManifest),
	}
}

func (jv *JobValidator) GetNodeManifest(nodeName string) *datatype.NodeManifest {
	if n, exist := jv.Nodes[nodeName]; exist {
		return n
	} else {
		return nil
	}
}

func (jv *JobValidator) GetPluginManifest(pluginImage string) *datatype.PluginManifest {
	if p, exist := jv.Plugins[pluginImage]; exist {
		return p
	} else {
		return nil
	}
}

func (jv *JobValidator) LoadDatabase() error {
	nodeFiles, err := ioutil.ReadDir(path.Join(jv.dataPath, "nodes"))
	if err != nil {
		return err
	}
	for _, nodeFile := range nodeFiles {
		nodeFilePath := path.Join(jv.dataPath, "nodes", nodeFile.Name())
		raw, err := os.ReadFile(nodeFilePath)
		if err != nil {
			logger.Debug.Printf("Failed to read %s:%s", nodeFilePath, err.Error())
			continue
		}
		var n datatype.NodeManifest
		err = json.Unmarshal(raw, &n)
		if err != nil {
			logger.Debug.Printf("Failed to parse %s:%s", nodeFilePath, err.Error())
			continue
		}
		jv.Nodes[n.Name] = &n
	}
	pluginFiles, err := ioutil.ReadDir(path.Join(jv.dataPath, "plugins"))
	if err != nil {
		return err
	}
	for _, pluginFile := range pluginFiles {
		pluginFilePath := path.Join(jv.dataPath, "plugins", pluginFile.Name())
		raw, err := os.ReadFile(pluginFilePath)
		if err != nil {
			logger.Debug.Printf("Failed to read %s:%s", pluginFilePath, err.Error())
			continue
		}
		var p datatype.PluginManifest
		err = json.Unmarshal(raw, &p)
		if err != nil {
			logger.Debug.Printf("Failed to parse %s:%s", pluginFilePath, err.Error())
			continue
		}
		jv.Plugins[p.Image] = &p
	}
	return nil
}

// GetNodeNamesByTags returns a list of node names matched with given tags
func (jv *JobValidator) GetNodeNamesByTags(tags []string) (nodesFound []string) {
	if len(tags) == 0 {
		return
	}
	for _, node := range jv.Nodes {
		if node.MatchTags(tags, true) {
			nodesFound = append(nodesFound, node.Name)
		}
	}
	return
}
