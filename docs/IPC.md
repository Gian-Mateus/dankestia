# IPC Commands Reference

Dankestia provides comprehensive IPC (Inter-Process Communication) functionality that allows external control of the shell through command-line commands. All IPC commands follow the format:

```bash
dankestia ipc call <target> <function> [parameters...]
```

## Target: `audio`

Audio system control and information.

### Functions

**`setvolume <percentage>`**
- Set output volume to specific percentage (0-100)
- Returns: Confirmation message

**`increment <step>`**
- Increase output volume by step amount
- Parameters: `step` - Volume increase amount (default: 5)
- Returns: Confirmation message

**`decrement <step>`**
- Decrease output volume by step amount
- Parameters: `step` - Volume decrease amount (default: 5)
- Returns: Confirmation message

**`mute`**
- Toggle output device mute state
- Returns: Current mute status

**`setmic <percentage>`**
- Set input (microphone) volume to specific percentage (0-100)
- Returns: Confirmation message

**`micmute`**
- Toggle input device mute state
- Returns: Current mic mute status

**`status`**
- Get current audio status for both input and output devices
- Returns: Volume levels and mute states

### Examples
```bash
dankestia ipc call audio setvolume 50
dankestia ipc call audio increment 10
dankestia ipc call audio mute
```

## Target: `brightness`

Display brightness control for internal and external displays.

### Functions

**`set <percentage> [device]`**
- Set brightness to specific percentage (1-100)
- Parameters:
  - `percentage` - Brightness level (1-100)
  - `device` - Optional device name (empty string for default)
- Returns: Confirmation with device info

**`increment <step> [device]`**
- Increase brightness by step amount
- Parameters:
  - `step` - Brightness increase amount
  - `device` - Optional device name (empty string for default)
- Returns: Confirmation with new brightness level

**`decrement <step> [device]`**
- Decrease brightness by step amount
- Parameters:
  - `step` - Brightness decrease amount
  - `device` - Optional device name (empty string for default)
- Returns: Confirmation with new brightness level

**`status`**
- Get current brightness status
- Returns: Current device and brightness level

**`list`**
- List all available brightness devices
- Returns: Device names and classes

### Examples
```bash
dankestia ipc call brightness set 80
dankestia ipc call brightness increment 10 ""
dankestia ipc call brightness decrement 5 "intel_backlight"
```

## Target: `night`

Night mode (gamma/color temperature) control.

### Functions

**`toggle`**
- Toggle night mode on/off
- Returns: Current night mode state

**`enable`**
- Enable night mode
- Returns: Confirmation message

**`disable`**
- Disable night mode
- Returns: Confirmation message

**`status`**
- Get current night mode status
- Returns: Night mode enabled/disabled state

**`temperature [value]`**
- Get or set night mode color temperature
- Parameters:
  - `value` - Optional temperature in Kelvin (2500-6000, steps of 500)
- Returns: Current or newly set temperature

**`automation [mode]`**
- Get or set night mode automation mode
- Parameters:
  - `mode` - Optional automation mode: "manual", "time", or "location"
- Returns: Current or newly set automation mode

**`schedule <start> <end>`**
- Set time-based automation schedule
- Parameters:
  - `start` - Start time in HH:MM format (e.g., "20:00")
  - `end` - End time in HH:MM format (e.g., "06:00")
- Returns: Confirmation of schedule update

**`location <latitude> <longitude>`**
- Set manual coordinates for location-based automation
- Parameters:
  - `latitude` - Latitude coordinate (e.g., 40.7128)
  - `longitude` - Longitude coordinate (e.g., -74.0060)
- Returns: Confirmation of coordinates update

### Examples
```bash
dankestia ipc call night toggle
dankestia ipc call night temperature 4000
dankestia ipc call night automation time
dankestia ipc call night schedule 20:00 06:00
dankestia ipc call night location 40.7128 -74.0060
```

## Target: `mpris`

Media player control via MPRIS interface.

### Functions

**`list`**
- List all available media players
- Returns: Player names

**`play`**
- Start playback on active player
- Returns: Nothing

**`pause`**
- Pause playback on active player
- Returns: Nothing

**`playPause`**
- Toggle play/pause state on active player
- Returns: Nothing

**`previous`**
- Skip to previous track
- Returns: Nothing

**`next`**
- Skip to next track
- Returns: Nothing

**`stop`**
- Stop playback on active player
- Returns: Nothing

### Examples
```bash
dankestia ipc call mpris playPause
dankestia ipc call mpris next
```

## Target: `lock`

Screen lock control and status.

### Functions

**`lock`**
- Lock the screen immediately
- Returns: Nothing

**`demo`**
- Show lock screen in demo mode (doesn't actually lock)
- Returns: Nothing

**`isLocked`**
- Check if screen is currently locked
- Returns: Boolean lock state

### Examples
```bash
dankestia ipc call lock lock
dankestia ipc call lock isLocked
```

## Target: `sessions`

Logind session enumeration and seat-local session switching. Wraps `loginctl list-sessions` and `loginctl activate`. Only switches between sessions that are *already running* on the current seat — creating a fresh login as another user requires a multi-session greeter setup (greetd-flexiserver / GDM / LightDM) and is out of scope.

### Functions

**`list`**
- Print every session DANKESTIA knows about as tab-separated columns: `sessionId\tusername\tseat\ttty\ttype\tcurrent-marker`
- Returns: Multi-line string. The current session is marked with `*current*`.

**`refresh`**
- Re-enumerate sessions in the background (the picker also refreshes itself on open)
- Returns: `"ok"`

**`open`**
- Refresh and open the Switch User picker on the focused screen
- Returns: `"ok"`

**`activate <sessionId>`**
- Activate a session by its numeric logind ID (the `Id=` field from `loginctl show-session`). Performs a VT switch
- Parameters: `sessionId` - Numeric session ID
- Returns: `"ok"` on dispatch, `"ERROR: missing session id"` if blank
- Note: Failures from `loginctl activate` surface through the `switchFailed` QML signal and a Log warning — the IPC call returns success once the spawn is queued, not after activation completes

**`switchTo <target>`**
- Switch to another session by username *or* session ID. The first non-current session matching the username wins; if there's no match, the call fails through the same logging path as `activate`
- Parameters: `target` - Username (e.g. `testuser2`) or numeric session ID
- Returns: `"ok"` on dispatch, `"ERROR: missing target (username or session id)"` if blank

### Examples
```bash
# Inspect what's switchable
dankestia ipc call sessions list

# Open the picker (useful for a keybind)
dankestia ipc call sessions open

# Jump straight to another logged-in user without the picker
dankestia ipc call sessions switchTo testuser2

# Or by session ID, when the user has multiple sessions
dankestia ipc call sessions activate 4
```

The dedicated `dankestia switch-user [target]` CLI command wraps the same behavior with a friendlier error path (it prints the switchable list when no target matches).

## Target: `inhibit`

Idle inhibitor control to prevent automatic sleep/lock.

### Functions

**`toggle`**
- Toggle idle inhibit state
- Returns: Current inhibit state message

**`enable`**
- Enable idle inhibit (prevent sleep/lock)
- Returns: Confirmation message

**`disable`**
- Disable idle inhibit (allow sleep/lock)
- Returns: Confirmation message

### Examples
```bash
dankestia ipc call inhibit toggle
dankestia ipc call inhibit enable
```

## Target: `powerprofile`

Power profile control via `power-profiles-daemon`. Changes stay in sync with DANKESTIA UI and trigger the power profile OSD when enabled.

Requires `power-profiles-daemon` to be installed and running. Works on all compositors.

### Functions

**`open`**
- Show the power profile picker modal
- Returns: Success confirmation or error if daemon unavailable

**`close`**
- Close the power profile picker modal
- Returns: Success confirmation

**`toggle`**
- Toggle power profile picker modal visibility
- Returns: Success confirmation or error if daemon unavailable

**`list`**
- List available profile slugs, one per line
- Returns: `power-saver`, `balanced`, and `performance` when supported

**`status`**
- Get the currently active profile slug
- Returns: `power-saver`, `balanced`, `performance`, or error if daemon unavailable

**`set <profile>`**
- Set the active power profile
- Parameters: Profile slug or alias — `power-saver` (`powersaver`, `saver`, `0`), `balanced` (`1`), `performance` (`2`)
- Returns: Success confirmation or error if profile unknown, unsupported, or write failed

**`cycle`**
- Cycle to the next available profile in order: power-saver → balanced → performance → power-saver
- Returns: Success confirmation or error if daemon unavailable or write failed

### Examples
```bash
dankestia ipc call powerprofile status
dankestia ipc call powerprofile list
dankestia ipc call powerprofile cycle
dankestia ipc call powerprofile set balanced
dankestia ipc call powerprofile set performance
dankestia ipc call powerprofile toggle
```

## Target: `wallpaper`

Wallpaper management and retrieval with support for per-monitor configurations.

### Legacy Functions (Global Wallpaper Mode)

**`get`**
- Get current wallpaper path
- Returns: Full path to current wallpaper file, or error if per-monitor mode is enabled

**`set <path>`**
- Set wallpaper to specified path
- Parameters: `path` - Absolute or relative path to image file
- Returns: Confirmation message or error if per-monitor mode is enabled

**`clear`**
- Clear all wallpapers and disable per-monitor mode
- Returns: Success confirmation

**`next`**
- Cycle to next wallpaper in the same directory
- Returns: Success confirmation or error if per-monitor mode is enabled

**`prev`**
- Cycle to previous wallpaper in the same directory
- Returns: Success confirmation or error if per-monitor mode is enabled

### Per-Monitor Functions

**`getFor <screenName>`**
- Get wallpaper path for specific monitor
- Parameters: `screenName` - Monitor name (e.g., "DP-2", "eDP-1")
- Returns: Full path to wallpaper file for the specified monitor

**`setFor <screenName> <path>`**
- Set wallpaper for specific monitor (automatically enables per-monitor mode)
- Parameters:
  - `screenName` - Monitor name (e.g., "DP-2", "eDP-1")
  - `path` - Absolute or relative path to image file
- Returns: Success confirmation with monitor and path info

**`nextFor <screenName>`**
- Cycle to next wallpaper for specific monitor
- Parameters: `screenName` - Monitor name (e.g., "DP-2", "eDP-1")
- Returns: Success confirmation

**`prevFor <screenName>`**
- Cycle to previous wallpaper for specific monitor
- Parameters: `screenName` - Monitor name (e.g., "DP-2", "eDP-1")
- Returns: Success confirmation

### Examples

**Global wallpaper mode:**
```bash
dankestia ipc call wallpaper get
dankestia ipc call wallpaper set /path/to/image.jpg
dankestia ipc call wallpaper next
dankestia ipc call wallpaper clear
```

**Per-monitor wallpaper mode:**
```bash
# Set different wallpapers for each monitor
dankestia ipc call wallpaper setFor DP-2 /path/to/image1.jpg
dankestia ipc call wallpaper setFor eDP-1 /path/to/image2.jpg

# Get wallpaper for specific monitor
dankestia ipc call wallpaper getFor DP-2

# Cycle wallpapers for specific monitor
dankestia ipc call wallpaper nextFor eDP-1
dankestia ipc call wallpaper prevFor DP-2

# Clear all wallpapers and return to global mode
dankestia ipc call wallpaper clear
```

**Error handling:**
When per-monitor mode is enabled, legacy functions will return helpful error messages:
```bash
dankestia ipc call wallpaper get
# Returns: "ERROR: Per-monitor mode enabled. Use getFor(screenName) instead."

dankestia ipc call wallpaper set /path/to/image.jpg
# Returns: "ERROR: Per-monitor mode enabled. Use setFor(screenName, path) instead."
```

## Target: `profile`

User profile image management.

### Functions

**`getImage`**
- Get current profile image path
- Returns: Full path to profile image or empty string if not set

**`setImage <path>`**
- Set profile image to specified path
- Parameters: `path` - Absolute or relative path to image file
- Returns: Success message with path or error message

**`clearImage`**
- Clear the profile image
- Returns: Success confirmation message

### Examples
```bash
dankestia ipc call profile getImage
dankestia ipc call profile setImage /path/to/avatar.png
dankestia ipc call profile clearImage
```

## Target: `theme`

Theme mode control (light/dark mode switching).

### Functions

**`toggle`**
- Toggle between light and dark themes
- Returns: Current theme mode ("light" or "dark")

**`light`**
- Switch to light theme mode
- Returns: "light"

**`dark`**
- Switch to dark theme mode
- Returns: "dark"

**`getMode`**
- Returns current mode
- Returns: "dark" or "light"

### Examples
```bash
dankestia ipc call theme toggle
dankestia ipc call theme dark
```

## Target: `bar`

Top bar visibility control.

### Functions

**`reveal`**
- Show the top bar
- Returns: Success confirmation

**`hide`**
- Hide the top bar
- Returns: Success confirmation

**`toggle`**
- Toggle top bar visibility
- Returns: Success confirmation with current state

**`toggleReveal`**
- Toggle the runtime reveal/tuck state for an autohidden bar
- Returns: Success confirmation with current reveal state

**`status`**
- Get current top bar visibility status
- Returns: "visible" or "hidden"

### Examples
```bash
dankestia ipc call bar toggle
dankestia ipc call bar toggleReveal index 0
dankestia ipc call bar hide
dankestia ipc call bar status
```

## Target: `systemupdater`

System updater widget control and background update checks.

### Functions

**`toggle`**
- Toggle the system updater popout open/closed

**`open`**
- Open the system updater popout

**`close`**
- Close the system updater popout

**`updatestatus`**
- Trigger a background update check
- Returns: Success confirmation

### Examples
```bash
dankestia ipc call systemupdater toggle
dankestia ipc call systemupdater open
dankestia ipc call systemupdater close
dankestia ipc call systemupdater updatestatus
```

## Target: `defaultApp`

Launch applications configured in Settings > Default Apps.

### Functions

**`browser`**
- Launch the configured default web browser
- Returns: Launch request confirmation

**`fileManager`**
- Launch the configured default file manager
- Returns: Launch request confirmation

**`textEditor`**
- Launch the configured default text editor
- Returns: Launch request confirmation

**`pdfReader`**
- Launch the configured default PDF reader
- Returns: Launch request confirmation

**`imageViewer`**
- Launch the configured default image viewer
- Returns: Launch request confirmation

**`videoPlayer`**
- Launch the configured default video player
- Returns: Launch request confirmation

**`musicPlayer`**
- Launch the configured default music player
- Returns: Launch request confirmation

**`mail`**
- Launch the configured default mail client
- Returns: Launch request confirmation

**`calendar`**
- Launch the configured default calendar application
- Returns: Launch request confirmation

### Examples
```bash
dankestia ipc call defaultApp browser
dankestia ipc call defaultApp fileManager
```

## Modal Controls

These targets control various modal windows and overlays.

### Target: `spotlight`
Application launcher modal control.

**Functions:**
- `open` - Show the spotlight launcher
- `close` - Hide the spotlight launcher
- `toggle` - Toggle spotlight launcher visibility
- `openQuery <query>` - Show the spotlight launcher with pre-filled search query
  - Parameters: `query` - Search text to pre-fill in the search box
  - Returns: Success confirmation
- `toggleQuery <query>` - Toggle spotlight launcher with pre-filled search query
  - Parameters: `query` - Search text to pre-fill in the search box (only used when opening)
  - Returns: Success confirmation

### Target: `clipboard`
Clipboard history modal control.

**Functions:**
- `open` - Show clipboard history
- `close` - Hide clipboard history
- `toggle` - Toggle clipboard history visibility

### Target: `notifications`
Notification center modal control.

**Functions:**
- `open` - Show notification center
- `close` - Hide notification center
- `toggle` - Toggle notification center visibility

### Target: `settings`
Settings modal control.

**Functions:**
- `open` - Show settings modal
- `close` - Hide settings modal
- `toggle` - Toggle settings modal visibility

### Target: `processlist`
System process list and performance modal control.

**Functions:**
- `open` - Show process list modal
- `close` - Hide process list modal
- `toggle` - Toggle process list modal visibility

### Target: `powermenu`
Power menu modal control for system power actions.

**Functions:**
- `open` - Show power menu modal
- `close` - Hide power menu modal
- `toggle` - Toggle power menu modal visibility

### Target: `powerprofile`
Power profile picker modal and profile control via `power-profiles-daemon`.

**Functions:**
- `open` - Show power profile picker modal
- `close` - Hide power profile picker modal
- `toggle` - Toggle power profile picker modal visibility
- `list` - List available profile slugs
- `status` - Get current profile slug
- `set <profile>` - Set profile by slug or alias (`power-saver`, `balanced`, `performance`)
- `cycle` - Cycle to the next available profile

### Target: `control-center`
Control Center popout containing network, bluetooth, audio, power, and other quick settings.

**Functions:**
- `open` - Show the control center
- `close` - Hide the control center
- `toggle` - Toggle control center visibility

**Examples**
```bash
dankestia ipc call control-center toggle
dankestia ipc call control-center open
dankestia ipc call control-center close
```

### Target: `notepad`
Notepad/scratchpad modal control for quick note-taking.

**Functions:**
- `open` - Show notepad modal
- `close` - Hide notepad modal
- `toggle` - Toggle notepad modal visibility
- `expand` - Expand the active notepad width and open it if hidden
- `collapse` - Collapse the active notepad width without changing visibility
- `toggleExpand` - Toggle the active notepad width between collapsed and expanded

### Target: `dash`
Dashboard popup control with tab selection for overview, media, and weather information.

**Functions:**
- `open [tab]` - Show dashboard popup with optional tab selection
  - Parameters: `tab` - Tab to open: "", "overview", "media", or "weather"
  - Returns: Success/failure message
- `close` - Hide dashboard popup
  - Returns: Success/failure message
- `toggle [tab]` - Toggle dashboard popup visibility with optional tab selection
  - Parameters: `tab` - Tab to open when showing: "", "overview", "media", or "weather"
  - Returns: Success/failure message

### Target: `dankdash`
DankDash wallpaper browser control.

**Functions:**
- `wallpaper` - Toggle DankDash popup on focused screen with wallpaper tab selected
  - Returns: Success/failure message

### Target: `file`
File browser controls for selecting wallpapers and profile images.

**Functions:**
- `browse <type>` - Open file browser for specific file type
  - Parameters: `type` - Either "wallpaper" or "profile"
  - `wallpaper` - Opens wallpaper file browser in Pictures directory
  - `profile` - Opens profile image file browser in Pictures directory
  - Both browsers support common image formats (jpg, jpeg, png, bmp, gif, webp)

### Target: `color-picker`
Color picker modal control.

**Functions:**
- `open` - Show color picker modal
- `openColor <color>` - Show color picker modal with a pre-selected color
  - Parameters: `color` - Color string (e.g. "#ff0000", "#3f51b5")
- `close` - Hide color picker modal
- `closeInstant` - Hide color picker modal without animation
- `toggle` - Toggle color picker modal visibility
- `toggleInstant` - Toggle color picker modal visibility without animation on hide

### Target: `hypr`
Hyprland-specific controls including keybinds cheatsheet and workspace overview (Hyprland only).

**Functions:**
- `openBinds` - Show Hyprland keybinds cheatsheet modal
  - Returns: Success/failure message
  - Note: Returns "HYPR_NOT_AVAILABLE" if not running Hyprland
- `closeBinds` - Hide Hyprland keybinds cheatsheet modal
  - Returns: Success/failure message
  - Note: Returns "HYPR_NOT_AVAILABLE" if not running Hyprland
- `toggleBinds` - Toggle Hyprland keybinds cheatsheet modal visibility
  - Returns: Success/failure message
  - Note: Returns "HYPR_NOT_AVAILABLE" if not running Hyprland
- `openOverview` - Show Hyprland workspace overview
  - Returns: "OVERVIEW_OPEN_SUCCESS" or "HYPR_NOT_AVAILABLE"
  - Displays all workspaces across all monitors with live window previews
  - Allows drag-and-drop window movement between workspaces and monitors
- `closeOverview` - Hide Hyprland workspace overview
  - Returns: "OVERVIEW_CLOSE_SUCCESS" or "HYPR_NOT_AVAILABLE"
- `toggleOverview` - Toggle Hyprland workspace overview visibility
  - Returns: "OVERVIEW_OPEN_SUCCESS", "OVERVIEW_CLOSE_SUCCESS", or "HYPR_NOT_AVAILABLE"

**Keybinds Cheatsheet Description:**
Displays an auto-categorized cheatsheet of all Hyprland keybinds parsed from `~/.config/hypr`. Keybinds are organized into three columns:
- **Window / Monitor** - Window and monitor management keybinds (sorted by dispatcher)
- **Workspace** - Workspace switching and management (sorted by dispatcher)
- **Execute** - Application launchers and commands (sorted by keybind)

**Workspace Overview Description:**
Displays a live overview of all workspaces across all monitors with window previews:
- **Multi-monitor support** - Shows workspaces from all connected monitors with monitor name labels
- **Live window previews** - Real-time screen capture of all windows on each workspace
- **Drag-and-drop** - Move windows between workspaces and monitors by dragging
- **Keyboard navigation** - Use Left/Right arrow keys to switch between workspaces on current monitor
- **Visual indicators** - Active workspace highlighted when it contains windows
- **Click to switch** - Click any workspace to switch to it
- **Click outside or press Escape** - Close the overview

### Modal Examples
```bash
# Open application launcher
dankestia ipc call spotlight toggle

# Open spotlight with pre-filled search
dankestia ipc call spotlight openQuery browser
dankestia ipc call spotlight toggleQuery "!"

# Show clipboard history
dankestia ipc call clipboard open

# Toggle notification center
dankestia ipc call notifications toggle

# Show settings
dankestia ipc call settings open

# Show system monitor
dankestia ipc call processlist toggle

# Show power menu
dankestia ipc call powermenu toggle

# Cycle or set power profile (requires power-profiles-daemon)
dankestia ipc call powerprofile cycle
dankestia ipc call powerprofile toggle

# Open notepad
dankestia ipc call notepad toggle

# Open the active notepad expanded
dankestia ipc call notepad expand

# Collapse the active notepad width
dankestia ipc call notepad collapse

# Toggle the active notepad width
dankestia ipc call notepad toggleExpand

# Show dashboard with specific tabs
dankestia ipc call dash open overview
dankestia ipc call dash toggle media
dankestia ipc call dash open weather

# Open wallpaper browser
dankestia ipc call dankdash wallpaper

# Open file browsers
dankestia ipc call file browse wallpaper
dankestia ipc call file browse profile

# Open color picker
dankestia ipc call color-picker toggle

# Show Hyprland keybinds cheatsheet (Hyprland only)
dankestia ipc call hypr toggleBinds
dankestia ipc call hypr openBinds

# Show Hyprland workspace overview (Hyprland only)
dankestia ipc call hypr toggleOverview
dankestia ipc call hypr openOverview
dankestia ipc call hypr closeOverview
```

## Common Usage Patterns

### Keybinding Integration

These IPC commands are designed to be used with window manager keybindings.

**Example niri configuration:**
```kdl
binds {
    Mod+Space { spawn "qs" "-c" "dankestia" "ipc" "call" "spotlight" "toggle"; }
    Mod+V { spawn "qs" "-c" "dankestia" "ipc" "call" "clipboard" "toggle"; }
    Mod+P { spawn "qs" "-c" "dankestia" "ipc" "call" "notepad" "toggle"; }
    Mod+Shift+P { spawn "qs" "-c" "dankestia" "ipc" "call" "notepad" "expand"; }
    Mod+Ctrl+P { spawn "qs" "-c" "dankestia" "ipc" "call" "notepad" "toggleExpand"; }
    Mod+X { spawn "qs" "-c" "dankestia" "ipc" "call" "powermenu" "toggle"; }
    XF86AudioRaiseVolume { spawn "qs" "-c" "dankestia" "ipc" "call" "audio" "increment" "3"; }
    XF86MonBrightnessUp { spawn "qs" "-c" "dankestia" "ipc" "call" "brightness" "increment" "5" ""; }
}
```

**Example Hyprland configuration:**
```conf
bind = SUPER, Space, exec, qs -c dankestia ipc call spotlight toggle
bind = SUPER, V, exec, qs -c dankestia ipc call clipboard toggle
bind = SUPER, P, exec, qs -c dankestia ipc call notepad toggle
bind = SUPER SHIFT, P, exec, qs -c dankestia ipc call notepad expand
bind = SUPER CTRL, P, exec, qs -c dankestia ipc call notepad toggleExpand
bind = SUPER, X, exec, qs -c dankestia ipc call powermenu toggle
bind = SUPER, slash, exec, qs -c dankestia ipc call hypr toggleBinds
bind = SUPER, Tab, exec, qs -c dankestia ipc call hypr toggleOverview
bind = , XF86AudioRaiseVolume, exec, qs -c dankestia ipc call audio increment 3
bind = , XF86MonBrightnessUp, exec, qs -c dankestia ipc call brightness increment 5 ""
```

### Scripting and Automation

IPC commands can be used in scripts for automation:

```bash
#!/bin/bash
# Toggle night mode based on time of day
hour=$(date +%H)
if [ $hour -ge 20 ] || [ $hour -le 6 ]; then
    dankestia ipc call night enable
else
    dankestia ipc call night disable
fi
```

### Status Checking

Many commands provide status information useful for scripts:

```bash
# Check if screen is locked before performing action
if dankestia ipc call lock isLocked | grep -q "false"; then
    # Perform action only if unlocked
    dankestia ipc call notifications open
fi
```

## Return Values

Most IPC functions return string messages indicating:
- Success confirmation with current values
- Error messages if operation fails
- Status information for query functions
- Empty/void return for simple action functions

Functions that return void (like media controls) execute the action but don't provide feedback. Check the application state through other means if needed.
