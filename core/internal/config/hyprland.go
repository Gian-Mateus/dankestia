package config

import _ "embed"

//go:embed embedded/hyprland.lua
var HyprlandLuaConfig string

//go:embed embedded/hypr-colors.lua
var DANKESTIAColorsLuaConfig string

//go:embed embedded/hypr-layout.lua
var DANKESTIALayoutLuaConfig string

//go:embed embedded/hypr-binds.lua
var DANKESTIABindsLuaConfig string

//go:embed embedded/hypr-outputs.lua
var DANKESTIAOutputsLuaConfig string

//go:embed embedded/hypr-cursor.lua
var DANKESTIACursorLuaConfig string

//go:embed embedded/hypr-windowrules.lua
var DANKESTIAWindowRulesLuaConfig string

//go:embed embedded/hypr-binds-user.lua
var DANKESTIABindsUserLuaConfig string
