package demoscenes

import "github/chschmidt99/pt/pkg/pt"

type DemoScene struct {
	Name    string
	Scene   *pt.Scene
	Cameras []*pt.Camera
}
