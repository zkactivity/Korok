package assets

import (
	"korok.io/korok/gfx/bk"
	"korok.io/korok/gfx"

	"log"
	"os"
	"fmt"
	"image"
	"errors"
	"io/ioutil"
	"encoding/json"
)

type TextureManager struct {
	repo map[string]RefCount
	names map[string]uint32
}

func NewTextureManager() *TextureManager {
	return &TextureManager{
		make(map[string]RefCount),
		make(map[string]uint32),
		}
}

// Load loads a single Texture file.
func (tm *TextureManager) Load(file string) {
	var rid, cnt uint16
	if v, ok := tm.repo[file]; ok {
		cnt = v.cnt
	} else {
		// create bk.Texture2D
		id, err := tm.loadTexture(file)
		if err != nil {
			log.Println(err)
		}
		rid = id
	}
	tm.repo[file] = RefCount{rid, cnt}
}

// LoadAtlas loads the atlas with a description file.
// The SubTexture can be found by SubTexture's name.
func (tm *TextureManager) LoadAtlas(img, desc string) {
	id, data := tm.loadAtlas(img, desc)
	size := len(data.texes)

	// new atlas
	at := gfx.R.NewAtlas(id, size, img)

	// fill
	for _, item := range data.texes {
		at.AddItem(float32(item.X), float32(item.Y), float32(item.W), float32(item.H), item.Name)
	}
}

// LoadAtlasIndexed loads the atlas with specified with/height/num.
func (tm *TextureManager) LoadAtlasIndexed(file string, width, height float32, row, col int) {
	id, err := tm.loadTexture(file)
	if err != nil {
		log.Println(err)
	}
	size := row * col

	// new atlas
	at := gfx.R.NewAtlas(id, size, file)

	// fill
	for i := 0; i < row; i ++ {
		for j := 0; j < col; j ++ {
			at.AddItem(float32(j)*width, float32(i)*height, width, height, "")
		}
	}
}

// Get returns the low-level Texture.
func (tm *TextureManager) Get(file string) gfx.Sprite {
	rid := tm.repo[file]
	return gfx.NewTex(rid.rid)
}

// Get returns the low-level Texture.
func (tm *TextureManager) GetRaw(file string) (uint16, *bk.Texture2D)  {
	if v, ok := tm.repo[file]; ok {
		if ok, tex := bk.R.Texture(v.rid); ok {
			return v.rid, tex
		}
	}
	return bk.InvalidId, nil
}


// Atlas returns the Atlas.
func (tm *TextureManager) Atlas(file string) *gfx.Atlas {
	return gfx.R.Atlas(file)
}

// Unload delete raw Texture and any related SubTextures.
func (tm *TextureManager) Unload(file string) {
	if v, ok := tm.repo[file]; ok {
		if v.cnt > 1 {
			tm.repo[file] = RefCount{v.rid, v.cnt -1}
		} else {
			delete(tm.repo, file)
			bk.R.Free(v.rid)
		}
	}
}

func (tm *TextureManager) loadTexture(file string)(uint16, error)  {
	log.Println("load file:" + file)
	// 1. load file
	imgFile, err := os.Open(file)
	if err != nil {
		return bk.InvalidId, fmt.Errorf("texture %q not found: %v", file, err)
	}
	// 2. decode image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return bk.InvalidId, err
	}
	// 3. create raw texture
	if id, _ := bk.R.AllocTexture(img); id != bk.InvalidId {
		return id, nil
	}
	return bk.InvalidId, errors.New("fail to load texture")
}

// 加载纹理图集
func (tm *TextureManager) loadAtlas(img, desc string)(id uint16, at *atlas) {
	id, err := tm.loadTexture(img)
	if err != nil {
		// todo
	}
	file, err := os.Open(desc)
	defer file.Close()

	if err != nil {
		// todo
	}
	d, err := ioutil.ReadAll(file)
	if err != nil {
		// todo
	}
	at = &atlas{}
	json.Unmarshal(d, at)
	return
}

type atlas struct {
	texes []struct{
		Name string
		X, Y, W, H int32
	}
}