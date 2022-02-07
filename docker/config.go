package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"io/ioutil"
	"os"
	"strings"
)

// getEnvVars Create a lists of formatted env vars usable to create a container
func (d *Config) getEnvVars() []string {
	envVars := make([]string, 0, len(d.Env))
	for k, v := range d.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	return envVars
}

// createEntrypointFile Create  an executable temporary file that contains the entrypoint code to run the container.
// WARNING this file might contain sensitive values. It should be deleted when the container is destroyed.
func (d *Config) createEntrypointFile() (*os.File, error) {
	tmpFile, err := ioutil.TempFile("/tmp", "wf-entrypoint-")
	if err != nil {
		return nil, fmt.Errorf("failed to create entrypoint file %v", err)
	}

	// Add each command into the file
	for _, command := range d.Commands {
		_, err := tmpFile.Write([]byte(command))
		if err != nil {
			return nil, fmt.Errorf("failed to write into entrypoint file %v", err)
		}
		// Add line breaks between each command
		_, err = tmpFile.Write([]byte("\n"))
		if err != nil {
			return nil, fmt.Errorf("failed to write into entrypoint file %v", err)
		}
	}

	err = tmpFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to set permission on entrypoint file %v", err)
	}

	// Set the script as executable
	err = os.Chmod(tmpFile.Name(), 0740)
	if err != nil {
		return nil, fmt.Errorf("failed to set permission on entrypoint file %v", err)
	}

	d.entrypointFile = tmpFile

	return tmpFile, nil
}

// clearEntrypointFile Delete the entrypoint file
func (d *Config) clearEntrypointFile() error {
	err := os.RemoveAll(d.entrypointFile.Name())
	d.entrypointFile = nil
	return err
}

// getDockerVolumes construct the docker volumes list that can be used to create a container.
func (d *Config) getDockerVolumes() []mount.Mount {
	volumesMounts := make([]mount.Mount, 0, len(d.Volumes)+1)

	// Mount configParser volumes
	for _, v := range d.Volumes {
		volume := mount.Mount{
			Type:     mount.TypeVolume,
			Source:   v.Label,
			Target:   v.ContainerMapping,
			ReadOnly: v.ReadOnly,
		}
		volumesMounts = append(volumesMounts, volume)
	}

	// Mount entrypoint
	volume := mount.Mount{
		Type:     mount.TypeBind,
		Source:   d.entrypointFile.Name(),
		Target:   "/entrypoint/entrypoint.sh",
		ReadOnly: true,
	}
	return append(volumesMounts, volume)
}

// getTimeout returns the maximum time that a container can live
func (d *Config) getTimeout() *int {
	timeout := 10 * 60
	return &timeout
}

// getWorkingDirectory Return the working directory of the container. Uses /app as default
func (d *Config) getWorkingDirectory() string {
	if d.WorkingDir == "" {
		return "/app"
	}
	return d.WorkingDir
}

func (d *Config) getEntrypoint() []string {
	return strings.Split(d.Entrypoint, " ")
}
