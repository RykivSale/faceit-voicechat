<div align="center">

**[🇷🇺 RUSSIAN README HERE](RU_README.md)**

<img src="docs/logo.png" width="128" alt="faceit-voicechat">

# THE EASIEST WAY<br>TO LOAD FACEIT DEMOS<br>WITH VOICE CHAT!

**Drop a `.dem` → copy the bind → paste it into the CS2 console.**<br>
Hear only CT, only T, or everyone — three keys.

<br>

<a href="https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe">
  <img src="https://img.shields.io/badge/%E2%AC%87%20DOWNLOAD-Windows%20.exe-e00000?style=for-the-badge&logo=windows&logoColor=white" alt="Download Windows .exe" height="52">
</a>
&nbsp;
<a href="https://github.com/RykivSale/faceit-voicechat/releases/latest">
  <img src="https://img.shields.io/github/v/release/RykivSale/faceit-voicechat?style=for-the-badge&label=LATEST&color=111111" alt="Latest release" height="52">
</a>

<br>

[![GitHub stars](https://img.shields.io/github/stars/RykivSale/faceit-voicechat?style=social)](https://github.com/RykivSale/faceit-voicechat/stargazers)
[![Downloads](https://img.shields.io/github/downloads/RykivSale/faceit-voicechat/total?label=downloads&logo=github)](https://github.com/RykivSale/faceit-voicechat/releases)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![☕ coffee](https://img.shields.io/badge/%E2%98%95%20coffee-dalink-FFDD00)](https://dalink.to/rykisosali)
[![💰 crypto](https://img.shields.io/badge/%F0%9F%92%B0%20crypto-ETH%2FUSDT-26a17b)](#buy-me-a-coffee)

<br>

**[⬇️  DOWNLOAD faceit-voicechat.exe](https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe)**
&nbsp;·&nbsp;
[all releases](https://github.com/RykivSale/faceit-voicechat/releases/latest)

<br>

<img src="docs/screenshot-drop.png" width="920" alt="Drop a .dem on the page">

<p><sub>Drop a demo — get a bind. No console wizardry, no dragging files by hand.</sub></p>

</div>

## Why this exists

You download a FACEIT demo, open it in CS2 — no voice, or everyone at once, and you can't tell who to listen to.

This app turns a `.dem` into **one console line**:

- **F5** — CT only
- **F6** — T only
- **F7** — all voices

Keys are rebindable. The demo is copied into your CS2 folder for you. `playdemo` is one copy-button away.

<img src="docs/screenshot-games.png" width="920" alt="Past games — Launch button">

Past matches from your CS2 folder — map icon, date, a red **LAUNCH** button. Score stays behind a spoiler so you don't get spoiled.

<div align="center">

<a href="https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe">
  <img src="https://img.shields.io/badge/%E2%AC%87%20DOWNLOAD-faceit--voicechat.exe-e00000?style=for-the-badge&logo=windows&logoColor=white" alt="Download .exe" height="48">
</a>

</div>

## How to use (30 seconds)

1. **[Download the `.exe`](https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe)** → double-click.
2. **Don't close** the console window — that's the server. Your browser opens `http://127.0.0.1:8765`.
3. First run: **Find CS2** (or pick `...\game\csgo` yourself).
4. Drop a `.dem` / `.dem.zst` on the page *or* right-click the demo → Open with → this exe.
5. In CS2 hit `~` → paste the **bind** → paste **playdemo**.

Done. F5 / F6 / F7 switch whose voices you hear on the recording.

Don't want to drop a file every time? **Past games** → **Launch**.

## Features

|  |  |
|---|---|
| 🖱️ Drag and drop | `.dem` and compressed `.dem.zst` |
| 📋 One bind | Not three blocks — **one line**, paste and go |
| 📁 Auto-copy | Demo lands in `game\csgo` (skipped if it's already there) |
| ▶️ playdemo | Correct filename, copy button |
| 🎮 Past games | Demos from the CS2 folder + map icon + launch |
| 🔍 Find CS2 | Scans Steam libraries for you |
| ⌨️ Rebind | F5 / F6 / F7 or any keys CS2 accepts |
| 🇷🇺 / 🇬🇧 | Russian and English UI |
| 🔄 Updates | Checks GitHub Releases on startup |

Settings live in `%AppData%\faceit-voicechat\config.json` and survive restarts.

<div align="center">

### missed the button? here it is again

<a href="https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe">
  <img src="https://img.shields.io/badge/%E2%AC%87%20DOWNLOAD-Windows%20.exe-e00000?style=for-the-badge&logo=windows&logoColor=white" alt="Download Windows" height="52">
</a>

**[⬇️  faceit-voicechat.exe](https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe)**

</div>

## This fork vs original

Fork of [boris-on/faceit-voicechat](https://github.com/boris-on/faceit-voicechat): same `tv_listen_voice_indices` idea, actually usable UI.

| | Original | This fork |
|---|---|---|
| Bind | 3 separate blocks | **1 line** |
| Keys | Fixed | **F5 / F6 / F7**, rebindable |
| Settings | None | Saved in `%AppData%` |
| Demo file | You move it | **Auto-copied** into `game\csgo` |
| `playdemo` | No | Printed + copy button |
| UI | Prints and exits | **Local page** (drag-drop, past games) |
| Updates | None | Checks GitHub Releases |

**⭐ If this saves you a tilt-queue — star the repo.** Or ☕ [coffee](https://dalink.to/rykisosali) / 💰 [crypto](#buy-me-a-coffee).

<a id="buy-me-a-coffee"></a>

## ☕💰 Buy me a coffee

The app is free. A star already slaps. If it really helped — you can throw coffee money.

<div align="center">

<a href="https://dalink.to/rykisosali">
  <img src="https://img.shields.io/badge/%E2%98%95%20Buy%20me%20a%20coffee-dalink-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=black" alt="☕ Buy me a coffee" height="48">
</a>
&nbsp;
<a href="https://etherscan.io/address/0x6DE7D78Bdd175B3b35343a2B9473444D3f6AeA16">
  <img src="https://img.shields.io/badge/%F0%9F%92%B0%20Donate%20crypto-ETH%20%2F%20USDT-26a17b?style=for-the-badge&logo=ethereum&logoColor=white" alt="💰 Donate crypto" height="48">
</a>

<br>

☕ **[dalink.to/rykisosali](https://dalink.to/rykisosali)**
&nbsp;·&nbsp;
💰 **[0x6DE7…eA16](https://etherscan.io/address/0x6DE7D78Bdd175B3b35343a2B9473444D3f6AeA16)**

</div>

Or send crypto directly. **EVM** (ETH, USDT ERC-20, same address on Polygon / BSC / Arbitrum):

```
0x6DE7D78Bdd175B3b35343a2B9473444D3f6AeA16
```

Don't send TRC-20 / non-EVM networks — different address format, the donate won't arrive.

## `--cli` (old-school)

```
faceit-voicechat.exe --cli
```

```
Enter  — print bind
S      — settings
U      — check updates
Q      — quit
```

Bind looks like this:

```
bind "F5" "tv_listen_voice_indices <ct>; tv_listen_voice_indices_h <ct>"; bind "F6" "tv_listen_voice_indices <t>; tv_listen_voice_indices_h <t>"; bind "F7" "tv_listen_voice_indices -1; tv_listen_voice_indices_h -1"
```

If a game folder is set, the demo is copied there and `playdemo <file>` is printed too.

## Build from source

Windows exe (what we ship):

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.appVersion=1.2.3"
```

macOS / Linux: `go build .` then run the binary. Same local page.

## License

GNU General Public License v3.0 — see [LICENSE](LICENSE).

Copyright (C) 2026 boris-on

## Attribution

If you use this in videos, streams, tutorials, or posts, credit the original author:

https://github.com/boris-on

<div align="center">

<br>

<a href="https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe">
  <img src="https://img.shields.io/badge/%E2%AC%87%20DOWNLOAD-right%20now-e00000?style=for-the-badge&logo=windows&logoColor=white" alt="Download right now" height="48">
</a>

<br><br>

**skill doesn't download. demo voices do.**

</div>
