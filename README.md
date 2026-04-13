# Spindle — Spotify Controller for Kindle

> Turn your Kindle into a Spotify remote.

Kindles have a built-in browser but no Spotify client. Spindle solves that — it's a web application optimized for e-ink screens that lets you browse and control Spotify directly from your Kindle, using any active device on your account as the audio output.

---

## The Problem

The Kindle browser is functional but limited — no app installs, no Spotify, and a touchscreen that works best with large tap targets. The official Spotify Web Player is too heavy and not designed for e-ink. I wanted something minimal, fast, and actually usable on the device.

---

## How It Works

Authentication happens on your phone to avoid typing on the Kindle keyboard. The Kindle displays a QR code — you scan it, authorize on Spotify, and the Kindle page detects the login automatically via polling and redirects you in.

```
Kindle (browser) ──► Spindle Server ──► Spotify Web API
       ▲                                        │
       └────────────── playback state ──────────┘
```

From there, you can browse your Playlists, Albums, and followed Artists in a paginated grid, tap to play, and control playback — all from the Kindle.

---

## Technical Highlights

**Backend — Go**
The server is written in Go using [chi](https://github.com/go-chi/chi) for routing and [GORM](https://gorm.io) with MySQL for persistence. The architecture is layered into handlers, services, repositories, and entities — keeping HTTP concerns separate from business logic and data access.

**Authentication & Security**
- Spotify OAuth 2.0 flow with a stateful pairing system — a short-lived token ties the Kindle session to the phone authorization, preventing replay attacks
- Access and refresh tokens are encrypted at rest using **AES-256-GCM** before being stored in the database
- Session cookies are signed and encrypted with [gorilla/securecookie](https://github.com/gorilla/securecookie), with `HttpOnly`, `Secure`, and a 30-day expiry
- Token refresh is handled transparently — when the access token expires, it's refreshed, re-encrypted, and persisted without any user interaction

**Frontend**
No frameworks. Vanilla HTML, CSS, and JavaScript served via Go's template engine. The UI is intentionally minimal — high contrast, large tap targets, no animations or transitions (e-ink doesn't render them). The layout uses CSS Grid with fixed row heights to keep the card grid stable across all content tabs.

**Login UX**
The QR code login was the most interesting UX challenge. The Kindle polls the server every 2 seconds after displaying the QR code. The phone scans it, authorizes on Spotify, the server marks the pairing as authenticated, and the next poll from the Kindle gets an `ok` response and redirects automatically — no manual action needed on the Kindle side.

---

## Stack

| | |
|---|---|
| Language | Go 1.24 |
| Router | chi |
| ORM | GORM |
| Database | MySQL |
| Crypto | AES-256-GCM (stdlib) |
| Session | gorilla/securecookie |
| Config | Viper |
| QR Code | go-qrcode |
| Frontend | Vanilla HTML/CSS/JS |
| Hosting | Discloud |

---

## Project Structure

```
.
├── cmd/server/             # Entrypoint
├── configs/                # Viper config loader
├── infra/
│   ├── database/           # GORM repository and connection
│   └── webserver/
│       └── handlers/       # HTTP handlers
├── internal/
│   ├── dao/                # Spotify API response types
│   ├── entity/             # Domain entities
│   ├── services/           # Spotify API calls
│   └── utils/              # Auth, cookies, crypto, pairing, templates
└── website/
    ├── assets/             # CSS and JS
    └── templates/          # Go HTML templates
```

---

## License

MIT © Jacksom Guilherme
