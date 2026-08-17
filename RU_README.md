<div align="center">

**[🇬🇧 ENGLISH README HERE](README.md)**

<img src="docs/logo.png" width="128" alt="faceit-voicechat">

# САМЫЙ ЛЕГКИЙ СПОСОБ<br>ГРУЗИТЬ ДЕМКИ FACEIT<br>с VOICE CHAT!

**Кинул `.dem` → скопировал бинд → вставил в консоль CS2.**<br>
Слышишь только CT, только T или всех — на трёх клавишах.

<br>

<a href="https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe">
  <img src="https://img.shields.io/badge/%E2%AC%87%20%20%D0%A1%D0%9A%D0%90%D0%A7%D0%90%D0%A2%D0%AC-Windows%20.exe-e00000?style=for-the-badge&logo=windows&logoColor=white" alt="Скачать Windows .exe" height="52">
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
[![💰 crypto](https://img.shields.io/badge/%F0%9F%92%B0%20crypto-ETH%2FUSDT-26a17b)](#кинь-на-кофе)

<br>

**[⬇️  СКАЧАТЬ faceit-voicechat.exe](https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe)**
&nbsp;·&nbsp;
[все релизы](https://github.com/RykivSale/faceit-voicechat/releases/latest)

<br>

<img src="docs/screenshot-drop.png" width="920" alt="Кинь .dem на страницу — и всё">

<p><sub>Перетащил демку — получил бинд. Без консольной магии и ручного копирования файлов.</sub></p>

</div>

## Зачем это

На FACEIT скачал демку, открыл в CS2 — а войса нет / слышно всех сразу / непонятно, кого слушать.

Эта прога делает из `.dem` **одну строку** для консоли:

- **F5** — только CT
- **F6** — только T
- **F7** — все голоса

Клавиши можно поменять. Демка сама улетает в папку CS2. `playdemo` тоже копируется одной кнопкой.

<img src="docs/screenshot-games.png" width="920" alt="Прошлые игры — кнопка Запустить">

Прошлые катки из папки CS2 — иконка карты, дата, красная **ЗАПУСТИТЬ**. Счёт спрятан под спойлер, чтобы не спойлерить.

<div align="center">

<a href="https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe">
  <img src="https://img.shields.io/badge/%E2%AC%87%20%20%D0%A1%D0%9A%D0%90%D0%A7%D0%90%D0%A2%D0%AC-faceit--voicechat.exe-e00000?style=for-the-badge&logo=windows&logoColor=white" alt="Скачать .exe" height="48">
</a>

</div>

## Как юзать (30 секунд)

1. **[Скачай `.exe`](https://github.com/RykivSale/faceit-voicechat/releases/latest/download/faceit-voicechat.exe)** → двойной клик.
2. Окно консоли **не закрывай** — это сервер. В браузере откроется `http://127.0.0.1:8765`.
3. Первый запуск: **Найти CS2** (или укажи папку `...\game\csgo` сам).
4. Кинь `.dem` / `.dem.zst` на страницу *или* ПКМ по демке → Open with → этот exe.
5. В CS2 жми `~` → вставь **бинд** → вставь **playdemo**.

Всё. F5 / F6 / F7 переключают, чьи войса слышно на записи.

Можно не кидать файл каждый раз: вкладка **Прошлые игры** → **Запустить**.

## Что умеет

|  |  |
|---|---|
| 🖱️ Драг-н-дроп | `.dem` и сжатые `.dem.zst` |
| 📋 Один бинд | Не три блока, а **одна строка** — копипаст и готово |
| 📁 Автокопия | Демка сама кладётся в `game\csgo` (если уже там — не копирует повторно) |
| ▶️ playdemo | Команда с правильным именем файла, кнопка Copy |
| 🎮 Прошлые игры | Демки из папки CS2 + иконка карты + запуск |
| 🔍 Найти CS2 | Сам ищет Steam-библиотеки |
| ⌨️ Ребинд | F5 / F6 / F7 или любые клавиши, которые ест CS2 |
| 🇷🇺 / 🇬🇧 | Русский и английский |
| 🔄 Апдейты | На старте чекает GitHub Releases |

Настройки живут в `%AppData%\faceit-voicechat\config.json` и не слетают.

<div align="center">

### не нашёл кнопку? вот она ещё раз

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

**⭐ If this saves you a tilt-queue — star the repo.** Или ☕ [кофе](https://dalink.to/rykisosali) / 💰 [крипта](#кинь-на-кофе).

<a id="кинь-на-кофе"></a>

## ☕💰 Кинь на кофе

Прога бесплатная. Звезда — уже кайф. Если совсем зашла — можно кинуть на кофе.

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

Или криптой напрямую. **EVM** (ETH, USDT ERC-20, тот же адрес на Polygon / BSC / Arbitrum):

```
0x6DE7D78Bdd175B3b35343a2B9473444D3f6AeA16
```

Не слать TRC-20 / не-EVM сети — другой формат адреса, донат не дойдёт.

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
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.appVersion=1.2.1"
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
  <img src="https://img.shields.io/badge/%E2%AC%87%20%20%D0%A1%D0%9A%D0%90%D0%A7%D0%90%D0%A2%D0%AC-%D0%BF%D1%80%D1%8F%D0%BC%D0%BE%20%D1%81%D0%B5%D0%B9%D1%87%D0%B0%D1%81-e00000?style=for-the-badge&logo=windows&logoColor=white" alt="Скачать прямо сейчас" height="48">
</a>

<br><br>

**скилл не качается. войсы на демке — да.**

</div>
