/** g007_transform.go
Transforms the rectangle.*/

package main

import (
  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/img"
  "github.com/go-gl/gl/v3.3-core/gl"
  "github.com/go-gl/mathgl/mgl32"
  "fmt"
  "io/ioutil"
  "strings"
  "math"
)

const (
  sW = 720
  sH = 720
)

var (
  shaderProg uint32
  vao uint32
  vbo uint32
  ebo uint32
  texture0 uint32
  texture1 uint32
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

func loadTexture(fname string) {
  var im *sdl.Surface
  var mode uint32
  var err error

  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
  gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

  im, err = img.Load(fname)
  if err != nil {
    fmt.Println("Could not load " + fname + ".")
    panic(err)
  }

  if (im.Format.BytesPerPixel == 4) {
    mode = gl.RGBA
  } else {
    mode = gl.RGB
  }
  gl.TexImage2D(gl.TEXTURE_2D, 0, int32(mode), im.W, im.H, 0, mode,
    gl.UNSIGNED_BYTE, gl.Ptr(im.Pixels()))

  im.Free()
}

func initTextures() {
  gl.GenTextures(1, &texture0)
  gl.BindTexture(gl.TEXTURE_2D, texture0)
  loadTexture("textures/texture0.jpg")

  gl.GenTextures(1, &texture1)
  gl.BindTexture(gl.TEXTURE_2D, texture1)
  loadTexture("textures/texture1.jpg")
}

func init1() {
  var err error
  var textureLocation int32
  err = sdl.Init(sdl.INIT_EVERYTHING)
  if err != nil { panic(err) }

  win, err = sdl.CreateWindow("g007 Transform", 200, 100, sW, sH, sdl.WINDOW_OPENGL)
  if err != nil { panic(err) }

  win.GLCreateContext()

  gl.Init()
  version := gl.GoStr(gl.GetString(gl.VERSION))
  fmt.Println("OpenGL version", version)

  initShaderProgram(&shaderProg)

  /*Vertices*/
  v := []float32{
  /*   X     Y  Z    Color RGB   tX tY - texture coordinates*/
    -0.5, -0.5, 0,   0, 0.5, 0,   0, 1, /*tY flipped because Y points up*/
     0.5, -0.5, 0,   0, 0.5, 0,   1, 1,
     0.5,  0.5, 0,   0,   0, 1,   1, 0,
    -0.5,  0.5, 0,   0, 0.5, 1,   0, 0,
  }

  /*Indices*/
  ind := []uint32{
    0, 1, 3, /*first triangle*/
    1, 2, 3, /*second triangle*/
  }

  gl.GenVertexArrays(1, &vao)
  gl.GenBuffers(1, &vbo)
  gl.GenBuffers(1, &ebo)

  gl.BindVertexArray(vao)

  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.BufferData(gl.ARRAY_BUFFER, len(v) * 4, gl.Ptr(v), gl.STATIC_DRAW)

  gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
  gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(ind) * 4, gl.Ptr(ind),
    gl.STATIC_DRAW)

  gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8 * 4, nil)
  gl.EnableVertexAttribArray(0) /*Position XYZ*/

  /*in 3 * 4, 3 is the number of floats, 4 is the number of bytes in float32*/
  gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8 * 4, gl.PtrOffset(3 * 4))
  gl.EnableVertexAttribArray(1) /*Color RGB*/

  gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8 * 4, gl.PtrOffset(6 * 4))
  gl.EnableVertexAttribArray(2) /*Texture coordinates*/

  initTextures()

  gl.UseProgram(shaderProg)

  textureLocation = gl.GetUniformLocation(shaderProg, gl.Str("texture0\x00"))
  gl.Uniform1i(textureLocation, 0)

  textureLocation = gl.GetUniformLocation(shaderProg, gl.Str("texture1\x00"))
  gl.Uniform1i(textureLocation, 1)
}

func draw() {
  var t float32

  var globalColor float32
  var globalColorLocation int32

  var transform mgl32.Mat4
  var transformLocation int32

  var rotate mgl32.Mat4

  t = float32((math.Sin(float64(sdl.GetTicks()) / 1000) + 1) / 2)
  globalColor = t

  transform = mgl32.Diag4(mgl32.Vec4{1.4, 1.4, 1.4, 1})
  rotate = mgl32.HomogRotate3D((t - 0.5) * 1.5 * math.Pi, mgl32.Vec3{0, 0, 1})
  transform = transform.Mul4(rotate)

  gl.ClearColor(0, 0, 0, 0)
  gl.Clear(gl.COLOR_BUFFER_BIT)

  gl.ActiveTexture(gl.TEXTURE0)
  gl.BindTexture(gl.TEXTURE_2D, texture0)

  gl.ActiveTexture(gl.TEXTURE1)
  gl.BindTexture(gl.TEXTURE_2D, texture1)

  globalColorLocation = gl.GetUniformLocation(shaderProg, gl.Str("globalColor\x00"))
  gl.Uniform1f(globalColorLocation, globalColor)

  transformLocation = gl.GetUniformLocation(shaderProg, gl.Str("transform\x00"))
  gl.UniformMatrix4fv(transformLocation, 1, false, &transform[0])

  gl.BindVertexArray(vao)
  gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
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
