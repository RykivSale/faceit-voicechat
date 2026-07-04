This script generates commands to listen only to your team or the opponents' team in the cs2 demo.

> This is a fork of [boris-on/faceit-voicechat](https://github.com/boris-on/faceit-voicechat) with a
> few workflow additions (see [Changes in this fork](#changes-in-this-fork) below). All credit for the
> original demo-parsing logic goes to the original author.

## How to use

1. Right-click a `.dem` file
2. Choose **"Open with..." → select `faceit-voicechat.exe`**
3. Press **Enter** to get the bind command, or **S** for settings

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
