package main

import (
	"fmt"
	_ "go_launcher_app/icons"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/inkyblackness/imgui-go/v2"
)

var (
	build_type string

	exePath, _ = os.Executable()
	exeDir     = filepath.Dir(exePath)
	errorLog   *log.Logger
	queryMap   = make(map[string]string)
	logText    string

	year1    int32
	year2    int32
	month1   int32
	month2   int32
	day1     int32
	day2     int32
	today    time.Time
	tomorrow time.Time
)

/* IMGUI */
func (glfw GLFW) Text() (string, error) {
	return glfw.ClipboardText()
}

func (glfw GLFW) SetText(text string) {
	glfw.SetClipboardText(text)
}

func StyleColorsLight() {
	/*Lerp := func(a imgui.Vec4, b imgui.Vec4, t float32) imgui.Vec4 {
		return imgui.Vec4{
			a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t, a.Z + (b.Z-a.Z)*t, a.W + (b.W-a.W)*t,
		}
	}*/

	style := imgui.CurrentStyle()

	style.SetColor(imgui.StyleColorText, imgui.Vec4{0.95, 0.96, 0.98, 1.00})
	style.SetColor(imgui.StyleColorTextDisabled, imgui.Vec4{0.36, 0.42, 0.47, 1.00})
	style.SetColor(imgui.StyleColorWindowBg, imgui.Vec4{0.11, 0.15, 0.17, 1.00})
	style.SetColor(imgui.StyleColorChildBg, imgui.Vec4{0.15, 0.18, 0.22, 1.00})
	style.SetColor(imgui.StyleColorPopupBg, imgui.Vec4{0.08, 0.08, 0.08, 0.94})
	style.SetColor(imgui.StyleColorBorder, imgui.Vec4{0.08, 0.10, 0.12, 1.00})
	style.SetColor(imgui.StyleColorBorderShadow, imgui.Vec4{0.00, 0.00, 0.00, 0.00})
	style.SetColor(imgui.StyleColorFrameBg, imgui.Vec4{0.20, 0.25, 0.29, 1.00})
	style.SetColor(imgui.StyleColorFrameBgHovered, imgui.Vec4{0.12, 0.20, 0.28, 1.00})
	style.SetColor(imgui.StyleColorFrameBgActive, imgui.Vec4{0.09, 0.12, 0.14, 1.00})
	style.SetColor(imgui.StyleColorTitleBg, imgui.Vec4{0.09, 0.12, 0.14, 0.65})
	style.SetColor(imgui.StyleColorTitleBgActive, imgui.Vec4{0.08, 0.10, 0.12, 1.00})
	style.SetColor(imgui.StyleColorTitleBgCollapsed, imgui.Vec4{0.00, 0.00, 0.00, 0.51})
	style.SetColor(imgui.StyleColorMenuBarBg, imgui.Vec4{0.15, 0.18, 0.22, 1.00})
	style.SetColor(imgui.StyleColorScrollbarBg, imgui.Vec4{0.02, 0.02, 0.02, 0.39})
	style.SetColor(imgui.StyleColorScrollbarGrab, imgui.Vec4{0.20, 0.25, 0.29, 1.00})
	style.SetColor(imgui.StyleColorScrollbarGrabHovered, imgui.Vec4{0.18, 0.22, 0.25, 1.00})
	style.SetColor(imgui.StyleColorScrollbarGrabActive, imgui.Vec4{0.09, 0.21, 0.31, 1.00})
	style.SetColor(imgui.StyleColorCheckMark, imgui.Vec4{0.28, 0.56, 1.00, 1.00})
	style.SetColor(imgui.StyleColorSliderGrab, imgui.Vec4{0.28, 0.56, 1.00, 1.00})
	style.SetColor(imgui.StyleColorSliderGrabActive, imgui.Vec4{0.37, 0.61, 1.00, 1.00})
	style.SetColor(imgui.StyleColorButton, imgui.Vec4{0.33, 0.60, 0.33, 1.00})
	style.SetColor(imgui.StyleColorButtonHovered, imgui.Vec4{0.33, 0.75, 0.33, 1.00})
	style.SetColor(imgui.StyleColorButtonActive, imgui.Vec4{0.33, 0.90, 0.33, 1.00})
	style.SetColor(imgui.StyleColorHeader, imgui.Vec4{0.20, 0.25, 0.29, 0.55})
	style.SetColor(imgui.StyleColorHeaderHovered, imgui.Vec4{0.26, 0.59, 0.98, 0.80})
	style.SetColor(imgui.StyleColorHeaderActive, imgui.Vec4{0.26, 0.59, 0.98, 1.00})
	style.SetColor(imgui.StyleColorSeparator, imgui.Vec4{0.20, 0.25, 0.29, 1.00})
	style.SetColor(imgui.StyleColorSeparatorHovered, imgui.Vec4{0.10, 0.40, 0.75, 0.78})
	style.SetColor(imgui.StyleColorSeparatorActive, imgui.Vec4{0.10, 0.40, 0.75, 1.00})
	style.SetColor(imgui.StyleColorResizeGrip, imgui.Vec4{0.26, 0.59, 0.98, 0.25})
	style.SetColor(imgui.StyleColorResizeGripHovered, imgui.Vec4{0.26, 0.59, 0.98, 0.67})
	style.SetColor(imgui.StyleColorResizeGripActive, imgui.Vec4{0.26, 0.59, 0.98, 0.95})
	style.SetColor(imgui.StyleColorTab, imgui.Vec4{0.11, 0.15, 0.17, 1.00})
	style.SetColor(imgui.StyleColorTabHovered, imgui.Vec4{0.26, 0.59, 0.98, 0.80})
	style.SetColor(imgui.StyleColorTabActive, imgui.Vec4{0.20, 0.25, 0.29, 1.00})
	style.SetColor(imgui.StyleColorTabUnfocused, imgui.Vec4{0.11, 0.15, 0.17, 1.00})
	style.SetColor(imgui.StyleColorTabUnfocusedActive, imgui.Vec4{0.11, 0.15, 0.17, 1.00})
	style.SetColor(imgui.StyleColorPlotLines, imgui.Vec4{0.61, 0.61, 0.61, 1.00})
	style.SetColor(imgui.StyleColorPlotLinesHovered, imgui.Vec4{1.00, 0.43, 0.35, 1.00})
	style.SetColor(imgui.StyleColorPlotHistogram, imgui.Vec4{0.90, 0.70, 0.00, 1.00})
	style.SetColor(imgui.StyleColorPlotHistogramHovered, imgui.Vec4{1.00, 0.60, 0.00, 1.00})
	style.SetColor(imgui.StyleColorTextSelectedBg, imgui.Vec4{0.26, 0.59, 0.98, 0.35})
	style.SetColor(imgui.StyleColorDragDropTarget, imgui.Vec4{1.00, 1.00, 0.00, 0.90})
	style.SetColor(imgui.StyleColorNavHighlight, imgui.Vec4{0.26, 0.59, 0.98, 1.00})
	style.SetColor(imgui.StyleColorNavWindowingHighlight, imgui.Vec4{1.00, 1.00, 1.00, 0.70})
	style.SetColor(imgui.StyleColorNavWindowingDarkening, imgui.Vec4{0.80, 0.80, 0.80, 0.20})
	style.SetColor(imgui.StyleColorModalWindowDarkening, imgui.Vec4{0.80, 0.80, 0.80, 0.35})

}

func Run(glfw GLFW, ogl OpenGL) {
	imgui.CurrentIO().SetClipboard(glfw)
	clearColor := [4]float32{0.0, 0.0, 0.0, 1.0}

	for !glfw.ShouldStop() {
		glfw.ProcessEvents()

		glfw.NewFrame()
		imgui.NewFrame()

		imgui.SetNextWindowPos(imgui.Vec2{})
		imgui.SetNextWindowSize(imgui.Vec2{X: glfw.DisplaySize()[0], Y: glfw.DisplaySize()[1]})

		imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 0.0)
		imgui.PushStyleVarFloat(imgui.StyleVarFrameRounding, 4.0)
		imgui.PushStyleVarFloat(imgui.StyleVarGrabRounding, 4.0)

		imgui.BeginV("BusturasLauncher", nil, imgui.WindowFlagsNoTitleBar|
			imgui.WindowFlagsNoResize|
			imgui.WindowFlagsNoMove|
			imgui.WindowFlagsNoCollapse|
			imgui.WindowFlagsNoBringToFrontOnFocus)

		imgui.PopStyleVarV(3)
		imgui.PushItemWidth(120.0)
		imgui.Text("Programos")
		imgui.Button("KurasKM")
		imgui.SameLine()
		imgui.Text("Up to date")
		imgui.Button("Tautvydo Korteles")
		imgui.SameLine()
		imgui.Text("Up to date")

		imgui.End()

		// This call only creates the draw data list. Actual rendering to framebuffer is done below.
		imgui.Render()
		ogl.PreRender(clearColor)
		ogl.Render(glfw.DisplaySize(), glfw.FramebufferSize(), imgui.RenderedDrawData())
		glfw.PostRender()
	}
}

/**
 * ANSI to UTF8
 */

var ansiToUtf = make(map[byte][2]byte)

func init() {
	ansiToUtf[0xE0] = [2]byte{0xC4, 0x85}
	ansiToUtf[0xE8] = [2]byte{0xC4, 0x8D}
	ansiToUtf[0xE6] = [2]byte{0xC4, 0x99}
	ansiToUtf[0xEB] = [2]byte{0xC4, 0x97}
	ansiToUtf[0xE1] = [2]byte{0xC4, 0xAF}
	ansiToUtf[0xF0] = [2]byte{0xC5, 0xA1}
	ansiToUtf[0xF8] = [2]byte{0xC5, 0xB3}
	ansiToUtf[0xFB] = [2]byte{0xC5, 0xAB}
	ansiToUtf[0xFE] = [2]byte{0xC5, 0xBE}
	ansiToUtf[0xC0] = [2]byte{0xC4, 0x84}
	ansiToUtf[0xC8] = [2]byte{0xC4, 0x8C}
	ansiToUtf[0xC6] = [2]byte{0xC4, 0x98}
	ansiToUtf[0xCB] = [2]byte{0xC4, 0x96}
	ansiToUtf[0xC1] = [2]byte{0xC4, 0xAE}
	ansiToUtf[0xD0] = [2]byte{0xC5, 0xA0}
	ansiToUtf[0xD8] = [2]byte{0xC5, 0xB2}
	ansiToUtf[0xDB] = [2]byte{0xC5, 0xAA}
	ansiToUtf[0xDE] = [2]byte{0xC5, 0xBD}
}

func utf(ansiString string) string {
	ansiBytesLen := len([]byte(ansiString))

	result := make([]byte, 0, ansiBytesLen*2)

	for _, byte := range []byte(ansiString) {
		if bytes, is_utf := ansiToUtf[byte]; is_utf {
			result = append(result, bytes[0], bytes[1])
		} else {
			result = append(result, byte)
		}
	}

	return string(result)
}

/**
 * Entry point
 */
func main() {

	/* Logging */
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("[INFO ] ")
	logFile, err := os.OpenFile(filepath.Join(exeDir, "go_launcher_app.LOG"),
		os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Unable open/create log file!")
	} else {
		log.Println("Log file was opened/created.")
	}
	w := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(w)
	errorLog = log.New(w, "[ERROR] ", log.LstdFlags|log.Lshortfile)

	/* ImGUI init */
	imguiContext := imgui.CreateContext(nil)
	defer imguiContext.Destroy()
	imguiIO := imgui.CurrentIO()
	imguiIO.SetIniFilename("")
	StyleColorsLight()

	glfw, err := NewGLFW(imguiIO)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer glfw.Dispose()

	ogl, err := NewOpenGL(imguiIO)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer ogl.Dispose()
	/* ImGUI init end */

	Run(*glfw, *ogl)
}
