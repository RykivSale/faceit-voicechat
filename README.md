<div align="center">

# ⬇️ [**DOWNLOAD HERE**](https://github.com/RykivSale/faceit-voicechat/releases/latest) ⬇️

</div>

[![GitHub Repo stars](https://img.shields.io/github/stars/RykivSale/faceit-voicechat?style=social)](https://github.com/RykivSale/faceit-voicechat/stargazers)
[![Latest release](https://img.shields.io/github/v/release/RykivSale/faceit-voicechat)](https://github.com/RykivSale/faceit-voicechat/releases/latest)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)

# faceit-voicechat (enhanced fork)

This script generates commands to listen only to your team or the opponents' team in the cs2 demo.

This is a fork of [boris-on/faceit-voicechat](https://github.com/boris-on/faceit-voicechat) that turns
the original one-shot console printout into an actual usable tool: a local web page, one bind line
instead of three, saved settings, CS2 folder auto-detect, and automatic demo handling.

**⭐ If this saves you time, please star the repo — it's the only payment I'm asking for.**

## Why use this fork instead of the original?

| | Original | This fork |
|---|---|---|
| Bind output | 3 separate command blocks you must copy one by one | **1 single line**, paste once and you're done |
| Keybinds | Fixed, not configurable | **F5 / F6 / F7 by default, fully rebindable** in-app |
| Settings | None | **Persistent settings** (`%AppData%`), survive restarts |
| Demo file | You move it manually to your game folder | **Auto-copied** to your configured game folder |
| `playdemo` command | Not provided | **Printed automatically** with the correct filename |
| Interface | Prints once, then exits | **Local web page** (or `--cli` menu): copy buttons, hints, CS2 auto-detect |

## How to use

1. Right-click a `.dem` file
2. Choose **"Open with..." → select `faceit-voicechat.exe`**
3. A local page opens in your browser (`http://127.0.0.1:8765`)
4. Copy the bind, paste it into the CS2 console (`~`)

You can also launch the exe with no file and drop a `.dem` / `.dem.zst` onto the page.

Leave the program window open while you use the page. Press **Quit** on the page or Ctrl+C in the console to stop.

Download the latest build: **[Releases → faceit-voicechat.exe](https://github.com/RykivSale/faceit-voicechat/releases/latest)**

## Local page

- **Bind** — one ready-to-paste line, with a copy button. Default keys: **F5** CT, **F6** T, **F7** all.
- **playdemo** — printed with the correct filename, also copyable.
- **Find CS2** — scans Steam libraries and common install paths (`...\Counter-Strike Global Offensive\game\csgo`). If a folder is saved, the demo is copied there automatically.
- **Keybinds** — change F5/F6/F7 to anything CS2 accepts.
- Language toggle: Russian / English.

Settings are saved to `%AppData%\faceit-voicechat\config.json` and persist between runs.

## Console menu (`--cli`)

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

## Changes in this fork

- Default UI is a local browser page with copy buttons, hints, and drag-and-drop demos.
- **Find CS2** tries Steam libraries and common `game\csgo` paths instead of typing the folder by hand.
- Bind commands are now emitted as a single `bind "KEY" "..."; bind "KEY" "..."` line instead of
  three separate command blocks.
- Added saved settings and custom keybinds (default `F5`, `F6`, `F7`).
- When a game folder is set, the opened demo is copied there automatically and a matching
  `playdemo <name>` command is printed.
- Old console menu is still available with `--cli`.

## License

This project is licensed under the GNU General Public License v3.0.

Copyright (C) 2026 boris-on

## Attribution

If you use this project in videos, streams, tutorials, articles, or public posts,
please credit the author:

Author: https://github.com/boris-on
