package engine

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const screenVertexShaderSource = `
#version 330 core
layout (location = 0) in vec4 vertex;

out vec2 TexCoords;

uniform mat4 projection;

void main()
{
   TexCoords = vertex.zw;
   gl_Position = projection * vec4(vertex.xy, 0.0, 1.0);
}
` + "\x00"

const screenFragmentShaderSource = `
#version 330 core
in vec2 TexCoords;
out vec4 color;

uniform usampler2D screen;
uniform sampler1D palette;

void main()
{
   uint paletteIndex = texture(screen, TexCoords).r;
   color = texelFetch(palette, int(paletteIndex), 0);
}
` + "\x00"

var window *glfw.Window
var screenTexture uint32
var paletteTexture uint32
var screenQuadVAO uint32
var shaderProgram uint32

func InitGraphicsGLFW() {
	initGLFWWindow()
	initGLScreen()
}

func ShutdownGraphicsGLFW() {
	glfw.Terminate()
}

func initGLFWWindow() {
	err := glfw.Init()
	if err != nil {
		FatalError("InitGraphics: Could not initialize GLFW:", "err", err)
	}
	slog.Info(fmt.Sprintf("InitGraphics: GLFW version %s", glfw.GetVersionString()))

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	window, err = glfw.CreateWindow(3*ScreenWidth, 3*ScreenWidth*3/4, "go-doom", nil, nil)
	if err != nil {
		FatalError("InitGraphics: Could not create window", "err", err)
	}

	window.SetSizeLimits(ScreenWidth, ScreenWidth*3/4, glfw.DontCare, glfw.DontCare)
	window.SetAspectRatio(4, 3)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.SetKeyCallback(glfwKeyCallback)
	window.SetMouseButtonCallback(glfwMouseButtonCallback)
	window.SetCursorPosCallback(glfwCursorPosCallback)

	window.MakeContextCurrent()
}

func initGLScreen() {
	if err := gl.Init(); err != nil {
		FatalError("InitGraphics: Could not initialize OpenGL", "err", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	slog.Info(fmt.Sprintf("InitGraphics: OpenGL version %s", version))

	vertexShader, err := compileShader(screenVertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		FatalError("InitGraphics: Failed to compile vertex shader", "err", err)
	}

	fragmentShader, err := compileShader(screenFragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		FatalError("InitGraphics: Failed to compile fragment shader", "err", err)
	}

	shaderProgram, err = newShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		FatalError("InitGraphics: Failed to create shaders program", "err", err)
	}

	// Set up screen quad with ortho projection for rendering screen texture to window
	var projectionMatrix = []float32{
		2.0, 0.0, 0.0, -1.0,
		0.0, -2.0, 0.0, 1.0,
		0.0, 0.0, 0.5, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}

	var vertices = []float32{
		// pos      // tex
		0.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
	}
	var vbo uint32

	gl.GenVertexArrays(1, &screenQuadVAO)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindVertexArray(screenQuadVAO)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	// Set up screen and palette textures
	gl.GenTextures(1, &screenTexture)
	gl.BindTexture(gl.TEXTURE_2D, screenTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.GenTextures(1, &paletteTexture)
	gl.BindTexture(gl.TEXTURE_1D, paletteTexture)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// Shader uses screen as texture 0 and palette as texture 1
	gl.UseProgram(shaderProgram)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram, gl.Str("screen\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(shaderProgram, gl.Str("palette\x00")), 1)
	gl.UniformMatrix4fv(gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00")), 1, true, &projectionMatrix[0])
}

func newShaderProgram(vertexShader, fragmentShader uint32) (uint32, error) {
	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func FinishUpdate(screen *[ScreenWidth * ScreenHeight]uint8) {
	gl.UseProgram(shaderProgram)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, screenTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, ScreenWidth, ScreenHeight, 0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(&screen[0]))
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_1D, paletteTexture)

	gl.BindVertexArray(screenQuadVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)

	window.SwapBuffers()
}

func SetPaletteGLFW(palette *[3 * 256]uint8) {
	gl.BindTexture(gl.TEXTURE_1D, paletteTexture)
	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA, 256, 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(&palette[0]))
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
}

func WindowShouldClose() bool {
	return window.ShouldClose()
}

// PollInputEvents TODO: move to GLFW input code
func PollInputEvents() {
	glfw.PollEvents()
}
