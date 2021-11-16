package docker

// image get the docker image that will be used by the container
func (d *DockerImageConfig) image() string {
	return d.Image
}

func (d *DockerImageConfig) initialize(_ Container) error {
	// Create entrypoint file
	file, err := d.createEntrypointFile()
	if err == nil {
		d.entrypointFile = file
	}
	return err
}

func (d *DockerImageConfig) destroy() error {
	// Destroy the  entrypoint file
	err := d.clearEntrypointFile()
	return err
}
