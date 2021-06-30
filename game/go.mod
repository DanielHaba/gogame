module github.com/dhindustries/game

go 1.16

require (
	github.com/dhindustries/graal v0.0.0
	github.com/dhindustries/graal/lib/font v0.0.0-00010101000000-000000000000
	github.com/dhindustries/graal/lib/glfw v0.0.0
	github.com/dhindustries/graal/lib/opengl v0.0.0
	github.com/go-gl/mathgl v1.0.0
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e
)

replace (
	github.com/dhindustries/graal => ../graal
	github.com/dhindustries/graal/lib/font => ../graal/lib/font
	github.com/dhindustries/graal/lib/glfw => ../graal/lib/glfw
	github.com/dhindustries/graal/lib/opengl => ../graal/lib/opengl
)
