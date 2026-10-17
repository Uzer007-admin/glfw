# glfw: библиотека проекта Granio

**Исходный репозиторий:** [go-gl/glfw](https://github.com/go-gl/glfw).
**Цель форка:** Окна и ввод для Granio: GLFW 3.5.1, X11/Wayland, прозрачность, системное размытие и диагностика выбранного backend.

Репозиторий развивается самостоятельно. Исходная история, лицензия и авторство
сохранены. Минимальная версия Go — **1.27.1**. Собственный адрес модуля:
`github.com/Uzer007-admin/glfw`
; для GLFW используются вложенные модули `v3.3/glfw`, `v3.4/glfw` и `v3.5/glfw`.

Платформенные ограничения и результаты проверок:
[развитие графического стека](https://github.com/Uzer007-admin/granio/blob/dev/granio/docs/NATIVE_STACK_IMPLEMENTATION.md).
Длительные проверки производительности в этот этап не входят.

---

## Документация исходного проекта и прежних этапов

Адреса установки в сохранённых примерах ниже могут относиться к исходнику.

[![Tests](https://github.com/go-gl/glfw/actions/workflows/tests.yml/badge.svg)](https://github.com/go-gl/glfw/actions/workflows/tests.yml)
[![Static Analysis](https://github.com/go-gl/glfw/actions/workflows/analysis.yml/badge.svg)](https://github.com/go-gl/glfw/actions/workflows/analysis.yml)
[![GoDoc](https://godoc.org/github.com/go-gl/glfw/v3.5/glfw?status.svg)](https://godoc.org/github.com/go-gl/glfw/v3.5/glfw)

# GLFW 3.5.1 for Go (local fork)

The `v3.3` and `v3.4` modules remain available. `v3.5/glfw` wraps native GLFW 3.5.1, commit `d9d6f0f1f967807ffade6598ea9a631ebaf37a56`, and adds `UnlimitedMouseButtons`. This module has no published upstream version; use the local workspace or a `replace` directive.

All three modules add `GetAxesInto`, `GetButtonsInto`, `GetHatsInto`,
`GetGamepadStateInto` and `GetGammaRampInto` for caller-owned reusable buffers.
Returned data is copied from GLFW storage; individual elements no longer require
cgo calls. Monitor/video-mode wrappers share a Go backing allocation.

`SetGammaRamp` validates equal, nonzero channel lengths and passes them through a
C-stack bridge. GLFW copies their data before returning. Six native getters also
return their count/state by value through C-stack bridges, removing hidden Go
output-buffer allocations while preserving callbacks. The main-thread rules still apply.

`Window.GetSurfaceMetrics` returns window and framebuffer dimensions, independent
content scales, focus/visibility/minimization/maximization and framebuffer alpha
availability in one cgo call with no Go output allocations. `PixelScale` computes
framebuffer density without treating it as UI scaling; zero geometry is invalid.
The native queries are sequential, do not pump events and must run on the main
thread. Framebuffer alpha availability does not establish OS backdrop support.

`Window.GetSurfaceGeometry` is a smaller three-query snapshot available in all
three modules. It returns sizes and independent content scales in one cgo call,
without querying focus, visibility or transparency and without Go output
allocations. Prefer collecting it after sizing/scale events and reusing it in
rendering; native geometry queries can still require X11 server round trips.

## Installation

* GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
	* On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
	* On Ubuntu/Debian-like Linux distributions, the default Linux build enables both X11 and Wayland, so you need `libgl1-mesa-dev`, `xorg-dev`, `libwayland-dev`, `libxkbcommon-dev` and `wayland-protocols`.
	* On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel` packages.
	* On FreeBSD, you need the package `pkgconf`. To build for X, you also need the package `xorg`; and to build for Wayland, you need the package `wayland`.
	* On NetBSD, to build for X, you need the X11 sets installed. These are included in all graphical installs, and can be added to the system with `sysinst(8)` on non-graphical systems. Wayland support is incomplete, due to missing wscons support in upstream GLFW. To attempt to build for Wayland, you need to install the `wayland libepoll-shim` packages and set the environment variable `PKG_CONFIG_PATH=/usr/pkg/libdata/pkgconfig`.
	* On OpenBSD, you need the X11 sets. These are installed by default, and can be added from the ramdisk kernel at any time.
	* See [here](http://www.glfw.org/docs/latest/compile.html#compile_deps) for full details.
* This local fork requires Go 1.27.1 or newer and a platform C compiler.

```
go mod edit -require=github.com/go-gl/glfw/v3.5/glfw@v0.0.0
go mod edit -replace=github.com/go-gl/glfw/v3.5/glfw=./deps/glfw/v3.5/glfw
```

### OpenGL ES

If your target system only provides an OpenGL ES implementation (true for some ARM boards), you need to link against that implementation.
You do this by defining the appropriate build tags, e.g.

```
go get -tags=gles2 github.com/go-gl/glfw/v3.5/glfw
```

Supported tags are `gles1`, `gles2`, `gles3` and `vulkan`.
Note that required packages might differ from those listed above; consult your hardware's documentation.

## Usage

```Go
package main

import (
	"runtime"
	"github.com/go-gl/glfw/v3.5/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
```

## Changelog

### GLFW 3.4 Specific Changes

* This section is cumulative for the v3.4 branch and includes changes from multiple commits.
* Added support for dynamic runtime switching between X11 and Wayland on Linux.
* Replaced compile-time Wayland icon checking with `glfwGetPlatform()`.
* Mapped `GLFW_FEATURE_UNAVAILABLE` and `GLFW_FEATURE_UNIMPLEMENTED` to prevent go panics on unsupported platform features.
* Mapped `GLFW_CURSOR_UNAVAILABLE` and `GLFW_PLATFORM_UNAVAILABLE`.
* Added function `GetPlatform`.
* Added function `PlatformSupported`.
* Added function `InitAllocator`.
* Added function `InitVulkanLoader`.
* Added function `GetError`.
* Added function `GetPhysicalDevicePresentationSupport`.
* Added function `Window.GetTitle`.
* Added function `Window.GetCocoaView`.
* Added hint `MousePassthrough`.
* Added hint `PositionX`.
* Added hint `PositionY`.
* Added hint `ContextDebug`.
* Added hint `ContextNoError`.
* Added hint `AnglePlatformType`.
* Added hint `PlatformHint`.
* Added hint `ScaleFramebuffer`.
* Added hint `Win32KeyboardMenu`.
* Added hint `Win32ShowDefault`.
* Added hint `WaylandAppID`.
* Added hint `WaylandLibdecor`.
* Added hint `X11XcbVulkanSurface`.
* Added hint value `AnyPosition`.
* Added hint value `WaylandPreferLibdecor`.
* Added hint value `WaylandDisableLibdecor`.
* Added cursor mode value `CursorCaptured`.
* `GLFW_UNLIMITED_MOUSE_BUTTONS` is not mapped because it is not present in the vendored GLFW C revision used by this package.
* Added multiple new standard cursors like `PointingHandCursor`, `ResizeEWCursor`, etc.
* Added `AnglePlatformType*` hint values for configuring ANGLE rendering backend.
* Added `Platform` type and values (`AnyPlatform`, `PlatformWin32`, `PlatformCocoa`, `PlatformWayland`, `PlatformX11`, `PlatformNull`).

### GLFW 3.3 Specific Changes
- Joystick functions now uses receivers instead of passing the joystick ID as argument.
- Vulkan methods are intentionally not implemented. `Window.Handle` can be used to create a Vulkan surface via the [this](https://github.com/vulkan-go/vulkan) package.

* Renamed `Window.GLFWWindow` to `Window.Handle`
* Added function `Window.SetAttrib`.
* Added function `Window.RequestAttention`.
* Added function `Window.GetContentScale`.
* Added function `Window.GetOpacity`.
* Added function `Window.SetOpacity`.
* Added function `Window.SetMaximizeCallback`.
* Added function `Window.SetContentScaleCallback`.
* Added function `Monitor.GetWorkarea`.
* Added function `Monitor.GetContentScale`.
* Added function `Monitor.SetUserPointer`.
* Added function `Monitor.GetUserPointer`.
* Added function `InitHint`.
* Added function `RawMouseMotionSupported`
* Added function `GetKeyScancode`.
* Added function `WindowHintString`.
* Added function `GetClipboardString`.
* Added function `SetClipboardString`.
* Added function `Joystick.GetHats`.
* Added function `Joystick.IsGamepad`.
* Added function `Joystick.GetGUID`.
* Added function `Joystick.GetGamepadName`.
* Added function `Joystick.GetGamepadState`.
* Added function `Joystick.SetUserPointer`.
* Added function `Joystick.GetUserPointer`.
* Added function `UpdateGamepadMappings`.
* Added function `SetX11SelectionString`.
* Added function `GetX11SelectionString`.
* Added gamepad button IDs.
* Added gamepad axis IDs.
* Added joystick hat state IDs.
* Added ModifierKey `ModCapsLock`.
* Added ModifierKey `ModNumLock`
* Added InputMode `LockKeyMods`.
* Added InputMode `RawMouseMotion`.
* Added hint `Hovered`.
* Added hint `CenterCursor`.
* Added hint `TransparentFramebuffer`.
* Added hint `FocusOnShow`.
* Added hint `ScaleToMonitor`.
* Added hint `JoystickHatButtons`.
* Added hint `CocoaChdirResources`.
* Added hint `CocoaMenubar`.
* Added hint `TransparentFramebuffer`.
* Added hint value `OSMesaContextAPI`.
* Added hint value `CocoaGraphicsSwitching`.
* Added hint value `CocoaRetinaFramebuffer`.
* Added string hint value `CocoaFrameNAME`.
* Added string hint value `X11ClassName`.
* Added string hint value `X11InstanceName`.
* `MonitorEvent` renamed to `PeripheralEvent` for reuse with joystick events.
* `Joystick.GetButtons` Returns `[]Action` instead of `[]byte`.
* `SetMonitorCallback` Returns `MonitorCallback`.
* `Focus` No longer returns an error.
* `Iconify` No longer returns an error.
* `Maximize` No longer returns an error.
* `Restore` No longer returns an error.
* `GetClipboardString` No longer returns an error.


### GLFW 3.2 Specific Changes
- Easy `go get` installation. GLFW source code is now included in-repo and compiled in so you don't have to build GLFW on your own and distribute shared libraries. The revision of GLFW C library used is listed in [GLFW_C_REVISION.txt](https://github.com/go-gl/glfw/blob/master/v3.2/glfw/GLFW_C_REVISION.txt) file.
- The error callback is now set internally. Functions return an error with corresponding code and description (do a type assertion to glfw3.Error for accessing the variables) if the error is recoverable. If not a panic will occur.

* Added function `Window.SetSizeLimits`.
* Added function `Window.SetAspectRatio`.
* Added function `Window.SetMonitor`.
* Added function `Window.Maximize`.
* Added function `Window.SetIcon`.
* Added function `Window.Focus`.
* Added function `GetKeyName`.
* Added function `VulkanSupported`.
* Added function `GetTimerValue`.
* Added function `GetTimerFrequency`.
* Added function `WaitEventsTimeout`.
* Added function `SetJoystickCallback`.
* Added window hint `Maximized`.
* Added hint `NoAPI`.
* Added hint `NativeContextAPI`.
* Added hint `EGLContextAPI`.


### GLFW 3.1 Specific Changes
* Added type `Cursor`.
* Added function `Window.SetDropCallback`.
* Added function `Window.SetCharModsCallback`.
* Added function `PostEmptyEvent`.
* Added function `CreateCursor`.
* Added function `CreateStandardCursor`.
* Added function `Cursor.Destroy`.
* Added function `Window.SetCursor`.
* Added function `Window.GetFrameSize`.
* Added window hint `Floating`.
* Added window hint `AutoIconify`.
* Added window hint `ContextReleaseBehavior`.
* Added window hint `DoubleBuffer`.
* Added hint value `AnyReleaseBehavior`.
* Added hint value `ReleaseBehaviorFlush`.
* Added hint value `ReleaseBehaviorNone`.
* Added hint value `DontCare`.

### API changes
* `Window.Iconify` Returns an error.
* `Window.Restore` Returns an error.
* `Init` Returns an error instead of `bool`.
* `GetJoystickAxes` No longer returns an error.
* `GetJoystickButtons` No longer returns an error.
* `GetJoystickName` No longer returns an error.
* `GetMonitors` No longer returns an error.
* `GetPrimaryMonitor` No longer returns an error.
* `Monitor.GetGammaRamp` No longer returns an error.
* `Monitor.GetVideoMode` No longer returns an error.
* `Monitor.GetVideoModes` No longer returns an error.
* `GetCurrentContext` No longer returns an error.
* `Window.SetCharCallback` Accepts `rune` instead of `uint`.
* Added type `Error`.
* Removed `SetErrorCallback`.
* Removed error code `NotInitialized`.
* Removed error code `NoCurrentContext`.
* Removed error code `InvalidEnum`.
* Removed error code `InvalidValue`.
* Removed error code `OutOfMemory`.
* Removed error code `PlatformError`.
* Removed `KeyBracket`.
* Renamed `Window.SetCharacterCallback` to `Window.SetCharCallback`.
* Renamed `Window.GetCursorPosition` to `GetCursorPos`.
* Renamed `Window.SetCursorPosition` to `SetCursorPos`.
* Renamed `CursorPositionCallback` to `CursorPosCallback`.
* Renamed `Window.SetCursorPositionCallback` to `SetCursorPosCallback`.
* Renamed `VideoMode` to `VidMode`.
* Renamed `Monitor.GetPosition` to `Monitor.GetPos`.
* Renamed `Window.GetPosition` to `Window.GetPos`.
* Renamed `Window.SetPosition` to `Window.SetPos`.
* Renamed `Window.GetAttribute` to `Window.GetAttrib`.
* Renamed `Window.SetPositionCallback` to `Window.SetPosCallback`.
* Renamed `PositionCallback` to `PosCallback`.
* Ranamed `Cursor` to `CursorMode`.
* Renamed `StickyKeys` to `StickyKeysMode`.
* Renamed `StickyMouseButtons` to `StickyMouseButtonsMode`.
* Renamed `ApiUnavailable` to `APIUnavailable`.
* Renamed `ClientApi` to `ClientAPI`.
* Renamed `OpenglForwardCompatible` to `OpenGLForwardCompatible`.
* Renamed `OpenglDebugContext` to `OpenGLDebugContext`.
* Renamed `OpenglProfile` to `OpenGLProfile`.
* Renamed `SrgbCapable` to `SRGBCapable`.
* Renamed `OpenglApi` to `OpenGLAPI`.
* Renamed `OpenglEsApi` to `OpenGLESAPI`.
* Renamed `OpenglAnyProfile` to `OpenGLAnyProfile`.
* Renamed `OpenglCoreProfile` to `OpenGLCoreProfile`.
* Renamed `OpenglCompatProfile` to `OpenGLCompatProfile`.
* Renamed `KeyKp...` to `KeyKP...`.
