package main

import (
	"fmt"
	"math"
	"time"

	"github.com/dhindustries/graal"

	"github.com/dhindustries/graal/component"
	"github.com/dhindustries/graal/lib/font"
	"github.com/dhindustries/graal/lib/glfw"
	"github.com/dhindustries/graal/lib/opengl"

	"github.com/go-gl/mathgl/mgl64"
)

const vertexShader = `
#version 330
layout(location=0) in vec3 vPosition;
layout(location=1) in vec3 vNormal;
layout(location=2) in vec2 vTexCoords;
layout(location=3) in vec4 vColor;

out vec3 fPosition;
out vec3 fNormal;
out vec2 fTexCoords;
out vec4 fColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    fPosition = vPosition;
    fNormal = vNormal;
    fTexCoords = vTexCoords;
    fColor = vColor;
    gl_Position = projection * view * model * vec4(vPosition, 1);
}
`

const uiVertexShader = `
#version 330
layout(location=0) in vec3 vPosition;
layout(location=1) in vec3 vNormal;
layout(location=2) in vec2 vTexCoords;
layout(location=3) in vec4 vColor;

out vec3 fPosition;
out vec3 fNormal;
out vec2 fTexCoords;
out vec4 fColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    fPosition = vPosition;
    fNormal = vNormal;
    fTexCoords = vTexCoords;
    fColor = vColor;
    gl_Position = projection * model * vec4(vPosition, 1);
}
`

const fragmentShader = `
#version 330
out vec4 oColor;

in vec3 fPosition;
in vec3 fNormal;
in vec2 fTexCoords;
in vec4 fColor;

uniform sampler2D texture1;
uniform vec4 color;

void main() {
    oColor = fColor * texture(texture1, fTexCoords);
}
`

type Model struct {
	component.Meshed
	component.Textured
	component.Transformed
}

type Label struct {
	graal.Text
	component.Transformed
}

func (label *Label) Render(api graal.Api, transform mgl64.Mat4) {
	api.Render(label.Text, transform)
}

type App struct {
	graal.Bindable
	renderable []interface{}
	model      *Model
	label      *Label
	time       float64

	sceneProgram graal.Program
	uiProgram    graal.Program
	sceneCamera  graal.Camera
	uiCamera     graal.Camera
}

func (app *App) Setup() error {
	app.renderable = []interface{}{}
	graal.SetClearColor(graal.Color{1, 0, 0, 1})
	return nil
}

func (app *App) Load() error {
	if err := app.loadProgram(); err != nil {
		return err
	}
	if err := app.loadScene(); err != nil {
		return err
	}
	return nil
}

func (app *App) Update(dt time.Duration) {
	if graal.IsKeyDown(graal.KeyEscape) {
		graal.Close()
	}

	app.time += dt.Seconds()
	r := math.Sin(app.time)
	g := math.Cos(app.time)
	b := 1 - (r * g)
	graal.SetClearColor(graal.Color{r, g, b, 1})
	rot := app.model.Rotation()
	rot[2] += 0.3 * dt.Seconds()
	app.model.SetRotation(rot)

	x, y := graal.MousePosition()
	w, h := graal.WindowSize()
	rx := (float64(x)/float64(w))*16 - 8
	ry := (float64(y)/float64(h))*12 - 6
	app.label.SetPosition(mgl64.Vec3{
		rx,
		ry,
		0,
	})
	app.label.SetMaxWidth(8 - rx)

	//pos := app.label.Position()
	//pos[1] -= 0.25 * dt.Seconds()
	//app.label.SetPosition(pos)
}

func (app *App) Render() {
	graal.UseProgram(app.sceneProgram)
	graal.UseCamera(app.sceneCamera)
	for _, obj := range app.renderable {
		graal.Render(obj, mgl64.Ident4())
	}
	//graal.UseProgram(app.uiProgram)
	//graal.UseCamera(app.uiCamera)
	graal.Render(app.label, mgl64.Ident4())
}

func (app *App) loadProgram() error {
	var vert graal.VertexShader
	var frag graal.FragmentShader
	var prog graal.Program
	var err error

	vert, err = graal.NewVertexShader(vertexShader)
	if err != nil {
		return err
	}

	frag, err = graal.NewFragmentShader(fragmentShader)
	if err != nil {
		return err
	}

	prog, err = graal.NewProgram(vert, frag)
	if err != nil {
		return err
	}
	app.Bind(prog)
	app.sceneProgram = prog

	vert, err = graal.NewVertexShader(uiVertexShader)
	if err != nil {
		return err
	}
	prog, err = graal.NewProgram(vert, frag)
	if err != nil {
		return err
	}
	app.Bind(prog)
	app.uiProgram = prog

	return nil
}

func (app *App) loadScene() error {
	var cam graal.OrthoCamera
	var err error
	if cam, err = graal.NewOrthoCamera(); err != nil {
		return err
	}
	cam.SetNear(0)
	cam.SetFar(1000)
	cam.SetPosition(mgl64.Vec3{0, 0, 15})
	cam.SetLookAt(mgl64.Vec3{0, 0, 0})
	cam.SetUp(mgl64.Vec3{0, 1, 0})
	cam.SetViewport(-8, -6, 8, 6)
	app.Bind(cam)
	app.sceneCamera = cam

	if cam, err = graal.NewOrthoCamera(); err != nil {
		return err
	}
	cam.SetNear(0)
	cam.SetFar(1000)
	cam.SetPosition(mgl64.Vec3{0, 0, 15})
	cam.SetLookAt(mgl64.Vec3{0, 0, 0})
	cam.SetUp(mgl64.Vec3{0, 1, 0})
	cam.SetViewport(0, 0, 800, 600)
	app.Bind(cam)

	app.model = &Model{}
	app.Bind(app.model)

	var mesh graal.Mesh
	if mesh, err = graal.NewSimpleQuad(-1, -1, 1, 1); err != nil {
		return err
	}
	app.model.SetMesh(mesh)

	var texture graal.Texture
	if texture, err = graal.LoadTexture("hero.png"); err != nil {
		return err
	}
	app.model.SetTexture(texture)

	font, err := graal.LoadFont("C:\\Windows\\Fonts\\Arial.ttf")
	if err != nil {
		return err
	}

	text, err := graal.NewText(font)
	if err != nil {
		return err
	}
	app.label = &Label{Text: text}
	app.Bind(app.label)

	//app.label.SetPosition(mgl64.Vec3{8, -6, 0})
	//app.label.SetRotation(mgl64.Vec3{0, math.Pi, 0})
	//app.label.SetScale(mgl64.Vec3{0.5, 0.5, 0.5})

	app.renderable = append(app.renderable, app.model)

	//go func() {
	//	text := "witaj podróżniku!\nzanim wyruszysz\nw drogę musisz\nzebrać drużynę"
	//	index := 0
	//	time.Sleep(1 * time.Second)
	//	ticker := time.NewTicker(500 * time.Millisecond)
	//	for range ticker.C {
	//		if index > len(text) {
	//			break
	//		}
	//		app.label.SetString(text[0:index])
	//		index++
	//	}
	//	ticker.Stop()
	//}()

	go func() {
		ch, _ := graal.KeyboardInput()
		text := ""
		for char := range ch {
			if char == '\x08' {
				if len(text) > 0 {
					text = text[:len(text)-1]
				}
			} else {
				text += string(char)
			}
			app.label.SetString(text)
		}
	}()

	return nil
}

func main() {
	graal.UseLibrary(&glfw.Library{})
	graal.UseLibrary(&opengl.Library{})
	graal.UseLibrary(&font.Library{
		Ranges: font.Ranges{
			{rune(32), rune(127)},
			{rune(243), rune(381)},
		},
	})
	if err := graal.Run(&App{}); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
