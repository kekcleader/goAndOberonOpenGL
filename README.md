# Go Oberon OpenGL

Go and [Oberon](https://freeoberon.su/en) examples of OpenGL 3D graphics.

Примеры программ, реализующих 3D-графику с помощью OpenGL на языках Go и [Oberon](https://freeoberon.su).

## How To Use
1. Install Go, GCC, OpenGL and SDL2. GCC is required by `go get` commands below.

2. Open terminal (or Command Prompt on Windows) and type:
```
go get -u github.com/go-gl/gl/v3.3-core/gl
go get -u github.com/veandco/go-sdl2/sdl
```
This will install Go bindings for OpenGL and SDL2.

3. Clone this repository:
```
git clone git@github.com:kekcleader/goAndOberonOpenGL.git
```
or download it.

4. Go to subdirectory `g001_simplest` (or another) and type:
```
make
```
It simply compiles and runs the program using `go run`.

![OpenGL-drawn textured rectangle](g004_rectangle/screenshots/01.png)

## Source Files

| File Name | Description |
| --------- | ----------- |
| [g001\_simplest](g001_simplest) | Simplest Go example of 3D graphics using OpenGL + SDL2. Draws a triangle on the screen. |
| [g002\_colors](g002_colors) | Applies a gradient color to the triangle. |
| [g003\_texture](g003_texture) | Applies a simple texture to the triangle. |
| [g004\_rectangle](g004_rectangle) | Textured rectangle. |
| [g005\_animation](g005_animation) | Draws a rectangle with an animated texture. |
| [g006\_twotextures](g006_twotextures) | Rectangle with two textures mixed together. |

See [video of g005](https://youtu.be/ifXDMCMWGI0)
