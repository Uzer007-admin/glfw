# Optional native window visuals

The 3.3, 3.4 and 3.5 Go wrappers expose the same opt-in API. It does not change
existing window/driver interfaces. Call on the main/window-owner thread while
the window is alive; `Destroy` clears the native pointer and is idempotent.

```go
glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True) // before CreateWindow
window, err := glfw.CreateWindow(800, 600, "Visuals", nil, nil)
// handle err, keep the window on its owner thread, initialize/render before Show
state, err := window.QueryWindowVisuals()
result, err := window.ApplyWindowVisuals(glfw.WindowVisualOptions{
    Attributes: glfw.VisualBackdrop | glfw.VisualDarkCaption | glfw.VisualCorners,
    Backdrop: glfw.BackdropMainWindow,
    DarkCaption: true,
    Corners: glfw.CornersRound,
})
// Inspect state, err, result.Accepted and per-attribute HRESULTs.
```

The Windows implementation uses documented `DwmSetWindowAttribute`/
`DwmGetWindowAttribute`. It loads `dwmapi.dll` using an absolute system-directory
path and releases that reference before returning. It does not replace the frame,
caption buttons, HWND, context, or native event procedure.

`QueryWindowVisuals` captures actual OS build, DWM composition, remote session,
high contrast (with availability), framebuffer transparency and actual Win32
pixel-format alpha bits. OS API support is separate from the created surface.
Only the system backdrop is documented as readable; caption/corner/alpha requests
are reported as accepted requests, never as proof of visible compositor output.

`VisualBackdrop` requires build 22621, caption/corners 22000, and
`VisualPremultipliedAlpha` 26100. Alpha enablement additionally requires a verified
transparent WGL pixel format. This bridge does not enable it for an unverified
EGL/NoAPI surface. All pixels must contain premultiplied RGB, with RGB=0 at A=0.
Preflight failures change no attributes; native failures can leave a subset
accepted, reported in `WindowVisualResult`. Unselected attributes stay untouched.

Other platforms return `ErrVisualUnsupported` from the DWM-specific API; the
3.5 Linux compositor API below independently supports transparency and blur.
High contrast/remote sessions/transparency policy remain the application's
decision; an accepted DWM request is not a native-material visual acceptance test.

Sources: [DWM attributes](https://learn.microsoft.com/en-us/windows/win32/api/dwmapi/ne-dwmapi-dwmwindowattribute),
[backdrop roles](https://learn.microsoft.com/en-us/windows/win32/api/dwmapi/ne-dwmapi-dwm_systembackdrop_type),
[public SDK enum](https://github.com/microsoft/win32metadata/blob/main/generation/WinSDK/RecompiledIdlHeaders/um/dwmapi.h).

## Linux compositor backgrounds (3.5)

`QueryBackgroundEffects` distinguishes configured alpha, live composition,
advertised blur and our own request. X11 verifies the native alpha visual and
compositor selection; Wayland GL contexts verify actual EGLConfig alpha.
For NoAPI Wayland windows, callers must verify the alpha of their own buffers.
`SetBackgroundBlur(true)` requests the full surface through KWin's X11 property
or `ext-background-effect-v1`; no desktop/version guessing or screen capture is
performed. An unsupported compositor returns `ErrBackgroundBlur`; an inactive
X11 compositor returns `ErrBackgroundComposition`. An unsuitable format returns
`ErrBackgroundAlpha`. These APIs use the GLFW initialization thread and reject
destroyed windows. Plain transparency needs no blur request: hint
`TransparentFramebuffer` before creation and render premultiplied RGBA.

Wayland uses GLFW's existing display/surface and the normal event dispatch.
Capabilities/global removal update state; surface objects are released before
the owning surface. Blur regions are surface-local, cover subsequent resizes,
and are committed by the next rendered buffer swap. The setter never performs
a surface commit or roundtrip. Before an opaque fallback swap, queue
`SetBackgroundBlur(false)` so effect removal and opaque pixels share a commit.
The compositor controls blur strength and policies. A successful request is
not an optical acceptance test. Old compositors without the standard extension
support transparency but return a reasoned blur error.

Protocol XML and generated client code are retained in the module with their
original license. Regenerate using `wayland-scanner client-header` and
`wayland-scanner private-code`, then refresh the native source hash.

## Surface generations and pixel regions

`SurfaceTracker` is available in all three GLFW versions on every native backend.
It keeps a per-window snapshot without replacing callbacks. After event processing
on the main thread, call `window.UpdateSurface(&tracker)`; it performs the existing
single-transition metrics query. Unchanged snapshots keep their generations.
Geometry, content-scale or explicit `InvalidateContext` changes increment
`ResourceGeneration`. Focus and visibility have separate change flags and do not
invalidate lower-layer material content. `Drawable` is false for hidden, iconified
or zero-sized surfaces. The compositor uses this gate to stop frame work.

`SurfaceMetrics.FramebufferRegion` clips window-coordinate rectangles and rounds
outwards using actual framebuffer/window ratios separately on X and Y.
`FramebufferRegion.BottomLeftY` converts the top-left result to GL crop coordinates.
Do not multiply content scale into this conversion: UI/user/text scale is a
separate application policy. Invalid or non-finite regions return `ok=false`.
Trackers own no context, detect no context loss automatically, and have no
observers, event pumps or scheduling goroutines. Use one tracker per window and
invalidate explicitly on context replacement. A destroyed window returns
`ErrWindowDestroyed` from `UpdateSurface`; a nil tracker returns `ErrSurfaceTracker`.
