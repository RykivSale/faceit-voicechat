[![GitHub Repo stars](https://img.shields.io/github/stars/RykivSale/faceit-voicechat?style=social)](https://github.com/RykivSale/faceit-voicechat/stargazers)
[![Latest release](https://img.shields.io/github/v/release/RykivSale/faceit-voicechat)](https://github.com/RykivSale/faceit-voicechat/releases/latest)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)

# faceit-voicechat (enhanced fork)

This script generates commands to listen only to your team or the opponents' team in the cs2 demo.

This is a fork of [boris-on/faceit-voicechat](https://github.com/boris-on/faceit-voicechat) that turns
the original one-shot console printout into an actual usable tool: one bind line instead of three,
a real menu, saved settings, and automatic demo handling.

**⭐ If this saves you time, please star the repo — it's the only payment I'm asking for.**

## Why use this fork instead of the original?

| | Original | This fork |
|---|---|---|
| Bind output | 3 separate command blocks you must copy one by one | **1 single line**, paste once and you're done |
| Keybinds | Fixed, not configurable | **F5 / F6 / F7 by default, fully rebindable** in-app |
| Settings | None | **Persistent settings** (`%AppData%`), survive restarts |
| Demo file | You move it manually to your game folder | **Auto-copied** to your configured game folder |
| `playdemo` command | Not provided | **Printed automatically** with the correct filename |
| Interface | Prints once, then exits | **Interactive menu**: get bind / settings / quit, reusable without relaunching |

## How to use

1. Right-click a `.dem` file
2. Choose **"Open with..." → select `faceit-voicechat.exe`**
3. Press **Enter** to get the bind command, or **S** for settings

Download the latest build: **[Releases → faceit-voicechat.exe](https://github.com/RykivSale/faceit-voicechat/releases/latest)**

## Menu

```
Press Enter - get bind
Press S - settings
Press Q - quit
```

- **Enter** prints a single ready-to-paste bind line:
  ```
  bind "F5" "tv_listen_voice_indices <ct_mask>; tv_listen_voice_indices_h <ct_mask>"; bind "F6" "tv_listen_voice_indices <t_mask>; tv_listen_voice_indices_h <t_mask>"; bind "F7" "tv_listen_voice_indices -1; tv_listen_voice_indices_h -1"
  ```
  If a game folder is configured (see Settings), the demo is also copied there and a
  `playdemo <filename>` command is printed.
- **S** opens settings:
  1. Set game folder — the folder demos get copied to (e.g. `...\Counter-Strike Global Offensive\game\csgo`)
  2. Change keybinds — replace F5/F6/F7 with any keys you want
  3. Back
- **Q** quits.

Settings are saved to `%AppData%\faceit-voicechat\config.json` and persist between runs.

## Changes in this fork

- Bind commands are now emitted as a single `bind "KEY" "..."; bind "KEY" "..."` line instead of
  three separate command blocks.
- Added an interactive menu (get bind / settings / quit) instead of a one-shot print.
- Added a settings screen to configure a game folder and custom keybinds (default `F5`, `F6`, `F7`).
- When a game folder is set, the opened demo is copied there automatically and a matching
  `playdemo <name>` command is printed.

## License

This project is licensed under the GNU General Public License v3.0.

Copyright (C) 2026 boris-on

## Attribution

If you use this project in videos, streams, tutorials, articles, or public posts,
please credit the author:

Author: https://github.com/boris-on
