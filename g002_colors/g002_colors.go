/** g002_colors.go
Applies a gradient color to the triangle.*/

package main

import (
  "github.com/veandco/go-sdl2/sdl"
  "github.com/go-gl/gl/v3.3-core/gl"
  "fmt"
  "io/ioutil"
  "strings"
)

const (
  sW = 1280
  sH = 720
)

var (
  shaderProg uint32
  vao uint32
  vbo uint32
  win *sdl.Window
)

/**Reads a file, appends byte 0X and stores it in s. Halts on error.*/
func readShader(fname string, s *string) {
  b, err := ioutil.ReadFile(fname)
  if err != nil {
    fmt.Println("Could not read file '" + fname + "'.")
    panic(err)
  }
  *s = string(b) + "\x00"
}

/**Loads shader source from a file and creates a shader of type sType.
   sType should be either gl.VERTEX_SHADER or gl.FRAGMENT_SHADER.
   Halts on error.*/
func loadShader(fname string, shader *uint32, sType uint32) {
  sh := gl.CreateShader(sType)

  var src0 string
  readShader(fname, &src0)

  src, free := gl.Strs(src0)
  gl.ShaderSource(sh, 1, src, nil)
  free()

  gl.CompileShader(sh)

  var status int32
  gl.GetShaderiv(sh, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE {
    var logLen int32
    gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &logLen)
    log := strings.Repeat("\x00", int(logLen + 1))
    gl.GetShaderInfoLog(sh, logLen, nil, gl.Str(log))
    panic("Failed to compile shader '" + fname + "':\n" + log)
  }

  *shader = sh;
}

func initShaderProgram(shaderProg *uint32) {
  var vShader, fShader uint32
  loadShader("shaders/vertex.txt", &vShader, gl.VERTEX_SHADER)
  loadShader("shaders/fragment.txt", &fShader, gl.FRAGMENT_SHADER)

  pr := gl.CreateProgram()
  gl.AttachShader(pr, vShader)
  gl.AttachShader(pr, fShader)
  gl.LinkProgram(pr)

  var status int32
  gl.GetProgramiv(pr, gl.LINK_STATUS, &status)
  if status == gl.FALSE {
    var logLen int32
    gl.GetProgramiv(pr, gl.INFO_LOG_LENGTH, &logLen)
    log := strings.Repeat("\x00", int(logLen + 1))
    gl.GetProgramInfoLog(pr, logLen, nil, gl.Str(log))
    panic("Failed to link shader program\n" + log)
  }

  gl.DeleteShader(vShader)
  gl.DeleteShader(fShader)

  *shaderProg = pr
}

func init1() {
  var err error
  err = sdl.Init(sdl.INIT_EVERYTHING)
  if err != nil { panic(err) }

  win, err = sdl.CreateWindow("g002 Colors", 200, 100, sW, sH, sdl.WINDOW_OPENGL)
  if err != nil { panic(err) }

  win.GLCreateContext()

  gl.Init()
  version := gl.GoStr(gl.GetString(gl.VERSION))
  fmt.Println("OpenGL version", version)

  initShaderProgram(&shaderProg)

  v := []float32{
  /*   X     Y  Z    Color RGB */
    -0.5, -0.5, 0,   1,   0, 0,
     0.5, -0.5, 0,   1, 0.5, 0,
     0,    0.5, 0,   0, 0.5, 1,
  }

  gl.GenBuffers(1, &vbo)
  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

  gl.GenVertexArrays(1, &vao)
  gl.BindVertexArray(vao)

  gl.BufferData(gl.ARRAY_BUFFER, len(v) * 4, gl.Ptr(v), gl.STATIC_DRAW)

  gl.EnableVertexAttribArray(0) /*Position XYZ*/
  gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6 * 4, nil)

  gl.EnableVertexAttribArray(1) /*Color RGB*/
  /*in 3 * 4, 3 is the number of floats, 4 is the number of bytes in float32*/
  gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6 * 4, gl.PtrOffset(3 * 4))

  gl.BindVertexArray(0)
}

func draw() {
  gl.ClearColor(0, 0, 0, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)
  gl.UseProgram(shaderProg)
  gl.BindVertexArray(vao)
  gl.DrawArrays(gl.TRIANGLES, 0, 3)
  win.GLSwap()
}

func handleKey(e *sdl.KeyboardEvent, done *bool) {
  if e.Keysym.Sym == sdl.K_ESCAPE {
    *done = true
  }
}

func run() {
  var event sdl.Event
  var done bool

  done = false
  for !done {
    event = sdl.PollEvent()
    for event != nil {
      switch e := event.(type) {
        case *sdl.QuitEvent: done = true
        case *sdl.KeyboardEvent: handleKey(e, &done)
      }
      event = sdl.PollEvent()
    }
    draw()
  }
}

func done() {
  win.Destroy()
  sdl.Quit()
}

func main() {
  init1()
  run()
  done()
}
