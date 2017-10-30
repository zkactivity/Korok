package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx/bk"

	"fmt"
	"unsafe"
)

/// MeshComp and MeshTable

type Mesh struct {
	// vertex data <x,y,u,v>
	vertex []PosTexColorVertex
	index  []uint16

	// res handle
	TextureId uint16
	IndexId   uint16
	VertexId  uint16
}

type MeshComp Mesh

func (*Mesh) Type() int32{
	return 0
}

func (m *Mesh) Setup() {
	mem_v := bk.Memory{unsafe.Pointer(&m.vertex[0]), uint32(len(m.vertex)) * 20 }
	if id, _:= bk.R.AllocVertexBuffer(mem_v, 20); id != bk.InvalidId {
		m.VertexId = id
	}

	mem_i := bk.Memory{unsafe.Pointer(&m.index[0]), uint32(len(m.index)) * 2}
	if id, _:= bk.R.AllocIndexBuffer(mem_i); id != bk.InvalidId {
		m.IndexId = id
	}
}

func (m*Mesh) SRT(pos mgl32.Vec2, rot float32, scale mgl32.Vec2) {
	for i := 0; i < 4; i++ {
		v := &m.vertex[i]
		v.X += pos[0]
		v.Y += pos[1]
	}
}

func (m*Mesh) SetVertex(v []PosTexColorVertex) {
	m.vertex = v
}

func (m*Mesh) SetIndex(v []uint16) {
	m.index = v
}


func (m *Mesh) Update() {
	if ok, ib := bk.R.IndexBuffer(m.IndexId); ok {
		ib.Update(0, uint32(len(m.index)) * 4, unsafe.Pointer(&m.index[0]), false)
	}

	if ok, vb := bk.R.VertexBuffer(m.VertexId); ok {
		vb.Update(0, uint32(len(m.vertex)) * 4, unsafe.Pointer(&m.vertex[0]), false)
	}
}

func (m*Mesh) Delete() {
	if ok, ib := bk.R.IndexBuffer(m.IndexId); ok {
		ib.Destroy()
	}
	if ok, vb := bk.R.VertexBuffer(m.VertexId); ok {
		vb.Destroy()
	}
	if ok, tex := bk.R.Texture(m.TextureId); ok {
		tex.Destroy()
	}
}

// new mesh from Texture
func NewQuadMesh(texId uint16) *Mesh{
	m := new(Mesh)

	m.TextureId = texId
	_, tex := bk.R.Texture(texId)

	fmt.Println("w:", tex.Width, " h:", tex.Height)

	tex.Width = 50
	tex.Height = 50

	m.vertex = []PosTexColorVertex{
		// Pos      	 // Tex
		{0.0, tex.Height, 0.0, 1.0, 0},
		{tex.Width, 0.0 , 1.0, 0.0, 0},
		{0.0, 0.0  	   , 0.0, 0.0, 0},

		{0.0, tex.Height, 0.0, 1.0, 0},
		{tex.Width, tex.Height, 1.0, 1.0, 0},
		{tex.Width, 0.0 , 1.0, 0.0, 0},
	}
	return m
}

// new mesh from SubTexture
func NewQuadMeshSubTex(texId uint16, tex *bk.SubTex) *Mesh {
	m := new(Mesh)

	m.TextureId = texId
	if tex.Texture2D == nil {
		_, tex.Texture2D = bk.R.Texture(texId)
	}

	h, w := tex.Height, tex.Width
	m.vertex = []PosTexColorVertex{
		// pos 			 // tex
		{0, h, tex.Min[0]/tex.Width, tex.Max[1]/tex.Height, 0},
		{w, 0, tex.Max[0]/tex.Width, tex.Min[1]/tex.Height, 0},
		{0, 0, tex.Min[0]/tex.Width, tex.Min[1]/tex.Height, 0},

		{0, h, tex.Min[0]/tex.Width, tex.Max[1]/tex.Height, 0},
		{w, h, tex.Max[0]/tex.Width, tex.Max[1]/tex.Height, 0},
		{w, 0, tex.Max[0]/tex.Width, tex.Min[1]/tex.Height, 0},
	}
	return m
}

func NewIndexedMesh(texId uint16, tex *bk.Texture2D) *Mesh {
	m := new(Mesh)

	m.TextureId = texId
	if tex == nil {
		_, tex = bk.R.Texture(texId)
	}

	h, w := tex.Height, tex.Width
	m.vertex = []PosTexColorVertex{
		{0,  h,  0.0, 1.0, 0},
		{w,  0,  1.0, 0.0, 0},
		{0,  0,  0.0, 0.0, 0},
		{w,  h,  1.0, 1.0, 0},
	}
	m.index = []uint16{
		0, 1, 2,
		0, 3, 1,
	}
	return m
}

// Configure VAO/VBO TODO
// 如果每个Sprite都创建一个VBO还是挺浪费的，
// 但是如果不创建新的VBO，那么怎么处理纹理坐标呢？
// 2D 场景中会出现大量模型和纹理相同的物体，仅仅
// 位置不同，比如满屏的子弹
// 或许可以通过工厂来构建mesh，这样自动把重复的mesh丢弃
// mesh 数量最终 <= 精灵的数量
var vertices = []float32{
	// Pos      // Tex
	0.0, 1.0, 0.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 0.0,

	0.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
}


type MeshTable struct {
	_comps []MeshComp
	_index uint32
	_map  map[int]int
}

func (mt *MeshTable) NewComp() (mc *MeshComp) {
	mc = &mt._comps[mt._index]
	mt._index ++
	return
}

func (mt *MeshTable) Comp() (mc *MeshComp) {
	mc = &mt._comps[mt._index]
	return
}

