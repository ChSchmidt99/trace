package pt

import (
	"runtime"
)

type Scene struct {
	root *SceneNode
}

func NewScene() *Scene {
	return &Scene{
		root: NewSceneNode(nil),
	}
}

func (s *Scene) Add(node *SceneNode) {
	s.root.Add(node)
}

func (s *Scene) Compile() BVH {
	prims := s.root.collectTracables(IdentityMatrix())
	builder := NewDefaultBuilder(prims)
	return builder.Build()
}

func (s *Scene) CompileLBVH() BVH {
	prims := s.root.collectTracables(IdentityMatrix())
	return DefaultLBVH(prims)
}

func (s *Scene) CompilePHR(alpha, delta float64, branchingFactor int) BVH {
	prims := s.root.collectTracables(IdentityMatrix())
	builder := NewPHRBuilder(prims, alpha, delta, branchingFactor, runtime.GOMAXPROCS(0))
	return builder.Build()
}

func (s *Scene) UntransformedTracables() []tracable {
	return s.root.collectTracablesRaw()
}

type SceneNode struct {
	transformation Matrix4
	children       []*SceneNode
	mesh           *Mesh
}

func NewSceneNode(mesh *Mesh) *SceneNode {
	return &SceneNode{
		transformation: IdentityMatrix(),
		mesh:           mesh,
	}
}

func (n *SceneNode) Add(node *SceneNode) {
	n.children = append(n.children, node)
}

func (n *SceneNode) ScaleUniform(factor float64) {
	n.transform(scaleUniform(factor))
}

func (n *SceneNode) Scale(x, y, z float64) {
	n.transform(scale(x, y, z))
}

func (n *SceneNode) Translate(x, y, z float64) {
	n.transform(translate(x, y, z))
}

func (n *SceneNode) Rotate(dir Vector3, angle float64) {
	n.transform(rotate(dir, angle))
}

func (n *SceneNode) ResetTransformation(dir Vector3, angle float64) {
	n.transformation = IdentityMatrix()
}

func (n *SceneNode) transform(t Matrix4) {
	n.transformation = n.transformation.MultiplyMatrix(t)
}

// Returns all Tracables without transforming
func (n *SceneNode) collectTracablesRaw() []tracable {
	out := make([]tracable, 0)
	if n.mesh != nil {
		out = append(out, n.mesh.raw()...)
	}
	for _, child := range n.children {
		out = append(out, child.collectTracablesRaw()...)
	}
	return out
}

// TODO: Make Multi Thread
// TODO: Check if zero alloc with accumulator is more efficient
// Returns all Tracables transformed by t
func (n *SceneNode) collectTracables(t Matrix4) []tracable {
	t = n.transformation.MultiplyMatrix(t)
	out := make([]tracable, 0)
	if n.mesh != nil {
		out = append(out, n.mesh.Transformed(t)...)
	}
	for _, child := range n.children {
		out = append(out, child.collectTracables(t)...)
	}
	return out
}

type Geometry []primitive

// TODO: Does it make sense to have seperate transformation for mesh?
type Mesh struct {
	transformation Matrix4
	geometry       Geometry
	material       Material
}

func NewMesh(geometry Geometry, mat Material) *Mesh {
	return &Mesh{
		transformation: IdentityMatrix(),
		geometry:       geometry,
		material:       mat,
	}
}

func (m Mesh) raw() []tracable {
	tracables := make([]tracable, len(m.geometry))
	for i, prim := range m.geometry {
		tracables[i] = tracable{
			prim: prim,
			mat:  m.material,
		}
	}
	return tracables
}

func (m Mesh) Transformed(t Matrix4) []tracable {
	tracables := make([]tracable, len(m.geometry))
	t = m.transformation.MultiplyMatrix(t)
	for i, prim := range m.geometry {
		tracables[i] = tracable{
			prim: prim.transformed(t),
			mat:  m.material,
		}
	}
	return tracables
}
