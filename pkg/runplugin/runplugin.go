package runplugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//         env:
//         - name: WAGGLE_PLUGIN_NAME
//           value: "{{.Name}}:{{.Version}}"
//         - name: WAGGLE_PLUGIN_VERSION
//           value: "{{.Version}}"
//         - name: WAGGLE_PLUGIN_USERNAME
//           value: "{{.Username}}"
//         - name: WAGGLE_PLUGIN_PASSWORD
//           value: "{{.Password}}"
//         - name: WAGGLE_PLUGIN_HOST
//           value: "rabbitmq-server"
//         - name: WAGGLE_PLUGIN_PORT
//           value: "5672"
//         envFrom:
//           - configMapRef:
//               name: waggle-config
//         resources:
//           limits:
//             cpu: 200m
//             memory: 20Mi
//           requests:
//             cpu: 100m
//             memory: 10Mi
//         volumeMounts:
//           - name: uploads
//             mountPath: /run/waggle/uploads
//           - name: waggle-data-config
//             mountPath: /run/waggle/data-config.json
//             subPath: data-config.json
//       volumes:
//       - name: uploads
//         hostPath:
//           path: /media/plugin-data/uploads/{{.Name}}/{{.Version}}
//           type: DirectoryOrCreate
//       - name: waggle-data-config
//         configMap:
//           name: waggle-data-config
// `))

type pluginConfig struct {
	Image    string
	Name     string
	Version  string
	Username string
	Password string
	Args     []string
}

var hostPathDirectoryOrCreate = apiv1.HostPathDirectoryOrCreate

// RunPlugin prepares to run a plugin image
func RunPlugin(image string, args ...string) error {
	base := path.Base(image)

	// split name:version from image string
	parts := strings.Split(base, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid plugin name %q", image)
	}

	config := &pluginConfig{
		Image:    image,
		Name:     parts[0],
		Version:  parts[1],
		Username: strings.ReplaceAll(base, ":", "-"),
		Password: "averysecurepassword",
		Args:     args,
	}

	// cmd := `
	// 	while ! rabbitmqctl -q authenticate_user ${plugin_username} ${plugin_password}; do
	// 	  echo "adding user ${plugin_username} to rabbitmq"
	// 	  rabbitmqctl -q add_user ${plugin_username} ${plugin_password} || \
	// 	  rabbitmqctl -q change_password ${plugin_username} ${plugin_password}
	// 	done

	// 	rabbitmqctl set_permissions ${plugin_username} ".*" ".*" ".*"
	// `

	/*
		kubectl exec --stdin service/rabbitmq-server -- sh -s <<EOF
		while ! rabbitmqctl -q authenticate_user ${plugin_username} ${plugin_password}; do
		  echo "adding user ${plugin_username} to rabbitmq"
		  rabbitmqctl -q add_user ${plugin_username} ${plugin_password} || \
		  rabbitmqctl -q change_password ${plugin_username} ${plugin_password}
		done

		rabbitmqctl set_permissions ${plugin_username} ".*" ".*" ".*"
		EOF
	*/

	res := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "app",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": config.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  config.Name,
						"role": "plugin",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  config.Name,
							Image: config.Image,
							Args:  config.Args,
							Env: []apiv1.EnvVar{
								{
									Name:  "WAGGLE_PLUGIN_NAME",
									Value: config.Name + ":" + config.Version,
								},
								{
									Name:  "WAGGLE_PLUGIN_VERSION",
									Value: config.Version,
								},
								{
									Name:  "WAGGLE_PLUGIN_USERNAME",
									Value: config.Username,
								},
								{
									Name:  "WAGGLE_PLUGIN_PASSWORD",
									Value: config.Password,
								},
								{
									Name:  "WAGGLE_PLUGIN_HOST",
									Value: "rabbitmq-server",
								},
								{
									Name:  "WAGGLE_PLUGIN_PORT",
									Value: "5672",
								},
							},
							EnvFrom: []apiv1.EnvFromSource{
								{
									ConfigMapRef: &apiv1.ConfigMapEnvSource{
										LocalObjectReference: apiv1.LocalObjectReference{
											Name: "waggle-config",
										},
									},
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "uploads",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: path.Join("/media/plugin-data/uploads", config.Name, config.Version),
									Type: &hostPathDirectoryOrCreate,
								},
							},
						},
						{
							Name: "waggle-data-config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "waggle-data-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	json.NewEncoder(os.Stdout).Encode(res)

	return nil
}
