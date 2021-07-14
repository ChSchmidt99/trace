package pt

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
	prims := s.root.collectPrimitives(IdentityMatrix())
	return NewBVH(prims)
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

// TODO: Make Multi Thread
// TODO: Check if zero alloc with accumulator is more efficient
func (n *SceneNode) collectPrimitives(t Matrix4) []Primitive {
	t = n.transformation.MultiplyMatrix(t)
	out := make([]Primitive, 0)
	if n.mesh != nil {
		out = append(out, n.mesh.transformed(t)...)
	}
	for _, child := range n.children {
		out = append(out, child.collectPrimitives(t)...)
	}
	return out
}

type Mesh struct {
	transformation Matrix4
	geometry       geometry
	//material       Material
}

func NewSphereMesh(center Vector3, radius float64) *Mesh {
	geo := geometry{newSphere(center, radius)}
	return newMesh(geo, nil)
}

func newMesh(geometry geometry, material Material) *Mesh {
	return &Mesh{
		transformation: IdentityMatrix(),
		geometry:       geometry,
		//material:       material,
	}
}

func (m Mesh) transformed(t Matrix4) []Primitive {
	return m.geometry.transformed(m.transformation.MultiplyMatrix(t))
}

type geometry []Primitive

func (g geometry) transformed(t Matrix4) []Primitive {
	prims := make([]Primitive, len(g))
	for i, prim := range g {
		prims[i] = prim.transformed(t)
	}
	return prims
}
