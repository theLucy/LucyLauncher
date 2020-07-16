package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/inkyblackness/imgui-go/v2"
)

// GLFWClientAPI identifies the render system that shall be initialized.
type GLFWClientAPI string

// GLFWClientAPI constants
const (
	MonitorWidth       int = 1920
	MonitorHeight      int = 1080
	MonitorRefreshRate int = 60
	AppWidth           int = 550
	AppHeight          int = 380
)

// GLFW implements a platform based on github.com/go-gl/glfw (v3.2).
type GLFW struct {
	imguiIO imgui.IO

	window *glfw.Window

	time             float64
	mouseJustPressed [3]bool
	mouseCursors     [imgui.MouseCursorCount]*glfw.Cursor
}

// NewGLFW attempts to initialize a GLFW context.
func NewGLFW(imguiIO imgui.IO) (*GLFW, error) {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %v", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(AppWidth, AppHeight, "BusturasLauncher", nil, nil)
	//window.SetAttrib(glfw.Decorated, 0)
	//window.SetMonitor(glfw.GetPrimaryMonitor(), 0, 0, 1920, 1080, 60)
	window.SetMonitor(nil, ((MonitorWidth - AppWidth) / 2), ((MonitorHeight - AppHeight) / 2), AppWidth, AppHeight, MonitorRefreshRate)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create window: %v", err)
	}
	setIcon(window)
	window.MakeContextCurrent()
	glfw.SwapInterval(2)

	platform := &GLFW{
		imguiIO: imguiIO,
		window:  window,
	}
	platform.setKeyMapping()
	platform.setMouseCursors()
	platform.installCallbacks()
	platform.imguiIO.SetBackendFlags(imgui.BackendFlagHasMouseCursors | imgui.BackendFlagsRendererHasVtxOffset)

	return platform, nil
}

// Dispose cleans up the resources.
func (platform *GLFW) Dispose() {
	for cursor_n := 0; cursor_n < imgui.MouseCursorCount; cursor_n++ {
		platform.mouseCursors[cursor_n].Destroy()
		platform.mouseCursors[cursor_n] = nil
	}

	platform.window.Destroy()
	glfw.Terminate()
}

// ShouldStop returns true if the window is to be closed.
func (platform *GLFW) ShouldStop() bool {
	return platform.window.ShouldClose()
}

// ProcessEvents handles all pending window events.
func (platform *GLFW) ProcessEvents() {
	glfw.PollEvents()
}

// DisplaySize returns the dimension of the display.
func (platform *GLFW) DisplaySize() [2]float32 {
	w, h := platform.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

// FramebufferSize returns the dimension of the framebuffer.
func (platform *GLFW) FramebufferSize() [2]float32 {
	w, h := platform.window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

// NewFrame marks the begin of a render pass. It forwards all current state to imgui IO.
func (platform *GLFW) NewFrame() {
	// Setup display size (every frame to accommodate for window resizing)
	displaySize := platform.DisplaySize()
	platform.imguiIO.SetDisplaySize(imgui.Vec2{X: displaySize[0], Y: displaySize[1]})

	// Setup time step
	currentTime := glfw.GetTime()
	if platform.time > 0 {
		platform.imguiIO.SetDeltaTime(float32(currentTime - platform.time))
	}
	platform.time = currentTime

	// Setup inputs
	if platform.window.GetAttrib(glfw.Focused) != 0 {
		x, y := platform.window.GetCursorPos()
		platform.imguiIO.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
	} else {
		platform.imguiIO.SetMousePosition(imgui.Vec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}

	for i := 0; i < len(platform.mouseJustPressed); i++ {
		down := platform.mouseJustPressed[i] || (platform.window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		platform.imguiIO.SetMouseButtonDown(i, down)
		platform.mouseJustPressed[i] = false
	}

	imgui_cursor := imgui.MouseCursor()
	if imgui_cursor == imgui.MouseCursorNone {
		// Hide OS mouse cursor if imgui is drawing it or if it wants no cursor
		platform.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	} else {
		// Show OS mouse cursor
		// FIXME-PLATFORM: Unfocused windows seems to fail changing the mouse cursor with GLFW 3.2, but 3.3 works here.
		platform.window.SetCursor(platform.mouseCursors[imgui_cursor])
		platform.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}

}

// PostRender performs a buffer swap.
func (platform *GLFW) PostRender() {
	platform.window.SwapBuffers()
}

func (platform *GLFW) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	platform.imguiIO.KeyMap(imgui.KeyTab, int(glfw.KeyTab))
	platform.imguiIO.KeyMap(imgui.KeyLeftArrow, int(glfw.KeyLeft))
	platform.imguiIO.KeyMap(imgui.KeyRightArrow, int(glfw.KeyRight))
	platform.imguiIO.KeyMap(imgui.KeyUpArrow, int(glfw.KeyUp))
	platform.imguiIO.KeyMap(imgui.KeyDownArrow, int(glfw.KeyDown))
	platform.imguiIO.KeyMap(imgui.KeyPageUp, int(glfw.KeyPageUp))
	platform.imguiIO.KeyMap(imgui.KeyPageDown, int(glfw.KeyPageDown))
	platform.imguiIO.KeyMap(imgui.KeyHome, int(glfw.KeyHome))
	platform.imguiIO.KeyMap(imgui.KeyEnd, int(glfw.KeyEnd))
	platform.imguiIO.KeyMap(imgui.KeyInsert, int(glfw.KeyInsert))
	platform.imguiIO.KeyMap(imgui.KeyDelete, int(glfw.KeyDelete))
	platform.imguiIO.KeyMap(imgui.KeyBackspace, int(glfw.KeyBackspace))
	platform.imguiIO.KeyMap(imgui.KeySpace, int(glfw.KeySpace))
	platform.imguiIO.KeyMap(imgui.KeyEnter, int(glfw.KeyEnter))
	platform.imguiIO.KeyMap(imgui.KeyEscape, int(glfw.KeyEscape))
	platform.imguiIO.KeyMap(imgui.KeyA, int(glfw.KeyA))
	platform.imguiIO.KeyMap(imgui.KeyC, int(glfw.KeyC))
	platform.imguiIO.KeyMap(imgui.KeyV, int(glfw.KeyV))
	platform.imguiIO.KeyMap(imgui.KeyX, int(glfw.KeyX))
	platform.imguiIO.KeyMap(imgui.KeyY, int(glfw.KeyY))
	platform.imguiIO.KeyMap(imgui.KeyZ, int(glfw.KeyZ))
}

func (platform *GLFW) setMouseCursors() {
	platform.mouseCursors[imgui.MouseCursorArrow] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	platform.mouseCursors[imgui.MouseCursorTextInput] = glfw.CreateStandardCursor(glfw.IBeamCursor)
	platform.mouseCursors[imgui.MouseCursorResizeAll] = glfw.CreateStandardCursor(glfw.ArrowCursor) // FIXME: GLFW doesn't have this.
	platform.mouseCursors[imgui.MouseCursorResizeNS] = glfw.CreateStandardCursor(glfw.VResizeCursor)
	platform.mouseCursors[imgui.MouseCursorResizeEW] = glfw.CreateStandardCursor(glfw.HResizeCursor)
	platform.mouseCursors[imgui.MouseCursorResizeNESW] = glfw.CreateStandardCursor(glfw.ArrowCursor) // FIXME: GLFW doesn't have this.
	platform.mouseCursors[imgui.MouseCursorResizeNWSE] = glfw.CreateStandardCursor(glfw.ArrowCursor) // FIXME: GLFW doesn't have this.
	platform.mouseCursors[imgui.MouseCursorHand] = glfw.CreateStandardCursor(glfw.HandCursor)
}

func setIcon(window *glfw.Window) {
	window_icon_file, err := Assets.Open("icon.png")
	if err != nil {
		fmt.Printf("Window icon failas nerastas")
		return
	}
	defer window_icon_file.Close()

	window_icon, err := png.Decode(window_icon_file)
	if err != nil {
		fmt.Printf("Nepavyko dekoduoti window icon failo")
		return
	}

	window.SetIcon([]image.Image{window_icon})
}

func (platform *GLFW) installCallbacks() {
	platform.window.SetMouseButtonCallback(platform.mouseButtonChange)
	platform.window.SetScrollCallback(platform.mouseScrollChange)
	platform.window.SetKeyCallback(platform.keyChange)
	platform.window.SetCharCallback(platform.charChange)
}

var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: 0,
	glfw.MouseButton2: 1,
	glfw.MouseButton3: 2,
	glfw.MouseButton4: 3,
	glfw.MouseButton5: 4,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	0: glfw.MouseButton1,
	1: glfw.MouseButton2,
	2: glfw.MouseButton3,
	3: glfw.MouseButton4,
	4: glfw.MouseButton5,
}

func (platform *GLFW) mouseButtonChange(window *glfw.Window, rawButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonIndex, known := glfwButtonIndexByID[rawButton]

	if known && (action == glfw.Press) {
		platform.mouseJustPressed[buttonIndex] = true
	}
}

func (platform *GLFW) mouseScrollChange(window *glfw.Window, x, y float64) {
	platform.imguiIO.AddMouseWheelDelta(float32(x), float32(y))
}

func (platform *GLFW) keyChange(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		platform.imguiIO.KeyPress(int(key))
	}
	if action == glfw.Release {
		platform.imguiIO.KeyRelease(int(key))
	}

	// Modifiers are not reliable across systems
	platform.imguiIO.KeyCtrl(int(glfw.KeyLeftControl), int(glfw.KeyRightControl))
	platform.imguiIO.KeyShift(int(glfw.KeyLeftShift), int(glfw.KeyRightShift))
	platform.imguiIO.KeyAlt(int(glfw.KeyLeftAlt), int(glfw.KeyRightAlt))
	platform.imguiIO.KeySuper(int(glfw.KeyLeftSuper), int(glfw.KeyRightSuper))
}

func (platform *GLFW) charChange(window *glfw.Window, char rune) {
	platform.imguiIO.AddInputCharacters(string(char))
}

// ClipboardText returns the current clipboard text, if available.
func (platform *GLFW) ClipboardText() (string, error) {
	return platform.window.GetClipboardString(), nil
}

// SetClipboardText sets the text as the current clipboard text.
func (platform *GLFW) SetClipboardText(text string) {
	platform.window.SetClipboardString(text)
}
