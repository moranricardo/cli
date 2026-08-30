package cmd

import (
"runtime/debug"
"strings"
)

var version string

func Version() string {
if version != "" && version != "0.0.0-dev" {
return version
}

if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
return info.Main.Version
}

if info, ok := debug.ReadBuildInfo(); ok {
var rev, dirty string
for _, s := range info.Settings {
if s.Key == "vcs.revision" {
rev = s.Value
if len(rev) > 7 {
rev = rev[:7]
}
}
if s.Key == "vcs.modified" && s.Value == "true" {
dirty = "-dirty"
}
}
if rev != "" {
if version == "" {
version = "0.0.0-dev"
}
if strings.HasPrefix(version, "0.0.0-dev") {
return version + "-" + rev + dirty
}
return version + "-" + rev + dirty
}
}

if version == "" {
version = "0.0.0-dev"
}
return version
}
