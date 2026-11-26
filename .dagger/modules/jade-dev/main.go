package jade_dev

import (
	"context"
	"strings"

	"dagger.io/dagger/core"

	"dagger/jade-dev/internal/dagger"
)

type JadeDev struct {
	// +private
	Source *dagger.Directory
}

func New(ws *dagger.Workspace) *JadeDev {
	return &JadeDev{
		Source: ws.Directory("/", dagger.WorkspaceDirectoryOpts{
			Exclude: []string{"**/.git", "**/.dagger", "**/.idea"},
		}),
	}
}

// Integration runs the integration tests on a virtual display.
//
// +check
func (m *JadeDev) Integration(ctx context.Context) error {
	_, err := base().
		With(virtualDisplay).
		WithDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"xvfb-run", "-a", "go", "test", "-tags", "integration", "-run", "^TestIntegration", "./..."}).
		Sync(ctx)
	return err
}

func base() *dagger.Container {
	return core.NewContainer().
		From("golang:1.27.4").
		WithMountedCache("/go/pkg/mod", core.NewCacheVolume("go-mod")).
		WithMountedCache("/root/.cache/go-build", core.NewCacheVolume("go-build"))
}

// virtualDisplay installs Xvfb and the X11/OpenGL libraries Ebitengine needs.
func virtualDisplay(c *dagger.Container) *dagger.Container {
	packages := []string{
		"xvfb", "xauth", "pkg-config", "libgl1-mesa-dev", "libgl1-mesa-dri",
		"libxcursor-dev", "libxi-dev", "libxinerama-dev", "libxrandr-dev", "libxxf86vm-dev", "libasound2-dev",
	}
	install := "apt-get update && apt-get install -y --no-install-recommends " + strings.Join(packages, " ")
	return c.WithExec([]string{"sh", "-c", install})
}
