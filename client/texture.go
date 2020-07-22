package main

import (
	"bytes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/inkyblackness/imgui-go/v2"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

type Texture struct {
	ID          imgui.TextureID
	Pavadinimas string
	Width       float32
	Height      float32
}

// VEIKIA TIK SU OPENGL2
func NewTextureFromFile(filename string) (Texture, error) {
	img_file, err := os.Open(filename)
	if err != nil {
		return Texture{}, err
	}
	defer img_file.Close()

	// Decode detects the type of image as long as its image/<type> is imported
	img, _, err := image.Decode(img_file)
	if err != nil {
		log.Printf("Nepalaikomas nuotraukos dekodavimas: %s", err)
		return Texture{}, err
	}
	return NewTexture(img, filename), nil
}

func NewTextureFromBytes(b []byte, pavadinimas string) (Texture, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Printf("Nepalaikomas nuotraukos dekodavimas: %s", err)
		return Texture{}, err
	}
	return NewTexture(img, pavadinimas), nil
}

// png užkoduoti baitai, nepriklausomai ar šaltinis buvo png/jpg
func BytesFromFilename(filename string) []byte {
	img_file, err := os.Open(filename)
	if err != nil {
		return []byte{}
	}
	defer img_file.Close()

	img, _, err := image.Decode(img_file)
	if err != nil {
		log.Printf("Nepalaikomas nuotraukos dekodavimas: %s", err)
		return []byte{}
	}

	return BytesFromImage(img)
}

func BytesFromImage(img image.Image) []byte {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		log.Printf("Nepavyko užkoduoti nuotraukos į baitus: %s", err)
		return []byte{}
	}

	return buffer.Bytes()
}

func NewTexture(img image.Image, pavadinimas string) Texture {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)
	rgba_width := rgba.Rect.Dx()
	rgba_height := rgba.Rect.Dy()

	var last_texture_handle int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &last_texture_handle)
	// Create a OpenGL texture identifier
	var texture_handle uint32
	gl.GenTextures(1, &texture_handle)
	gl.BindTexture(gl.TEXTURE_2D, texture_handle)
	// Setup filtering parameters for display
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP_SGIS, gl.TRUE)
	// Upload pixels into texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba_width), int32(rgba_height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	// restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(last_texture_handle))

	return Texture{ID: imgui.TextureID(texture_handle), Pavadinimas: pavadinimas, Width: float32(rgba_width), Height: float32(rgba_height)}
}

func DeleteTextures(textures []Texture) {
	for _, x := range textures {
		id := uint32(x.ID)
		gl.DeleteTextures(1, &id)
		x.ID = 0
	}
}
