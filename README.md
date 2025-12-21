# Oreon Defense

security daemon + tray app for linux (specifically KDE Plasma). wraps ClamAV for scanning and nftables for firewall stuff. the goal is something that just works and doesn't get in your way.

## what it does

- antivirus scanning via ClamAV (quick scan, full scan)
- firewall management via nftables
- system tray icon that shows protection status
- desktop notifications when stuff happens
- all the usual settings you'd expect

## architecture

two binaries:
- `defensed` - the daemon, runs as root, does the actual work
- `defense-ui` - tray app + dashboard, runs as your user

they talk over a unix socket (`/run/oreon/defense.sock`) using a simple JSON protocol.

## tech stack

| thing | what we're using |
|-------|------------------|
| language | Go 1.21+ |
| gui | [therecipe/qt](https://github.com/therecipe/qt) (Qt6) |
| antivirus | ClamAV via [go-clamd](https://github.com/dutchcoders/go-clamd) |
| firewall | [google/nftables](https://github.com/google/nftables) |
| service | systemd |
| config | TOML |
| local db | SQLite |

## building

```bash
make build
```

binaries end up in `bin/`.

## project structure

```
cmd/defensed/       daemon entry point
cmd/defense-ui/     tray/gui entry point
internal/daemon/    daemon internals (state machine, etc)
pkg/config/         config loading/saving
pkg/ipc/            IPC protocol definitions
```

## for collaborators

if you're working on this with me:

- **firewall stuff**: check `pkg/ipc/protocol.go` for the firewall commands (`CmdFirewallStatus`, `CmdFirewallEnable`, etc)
- **tray/ui stuff**: check `internal/daemon/state.go` for the state machine and `pkg/ipc/protocol.go` for response types
- config lives at `/etc/oreon/defense.toml` (system) or `~/.config/oreon/defense.toml` (user override)
- need something to work on? just ask me

## status

early days, actively being built.

## license

GPL-3.0 - see [LICENSE](LICENSE)
