# Roadmap

English | [中文](ROADMAP-zh.md)

What's planned next for ARK Server Commander. There are no fixed dates; items are grouped into near-term, mid-term and long-term by priority, and the order may change based on feedback. See [Releases](https://github.com/TBro1998/ArkServerCommander/releases) for what changed in each version.

Want a feature, or think something should come sooner? Please open an [issue](https://github.com/TBro1998/ArkServerCommander/issues).

## Recently Completed

- 🤖 MCP endpoint: agents such as Claude Code and Cursor can manage servers directly, with full-control or read-only tokens
- 🌍 One-click publishing to the [arkfgf.com](https://arkfgf.com) server list, with real online player counts reported automatically
- 🔗 P2P rooms: players can join through [ArkServerLink](ArkServerLink/) even without a public IP
- 🧩 Plugin catalog: install and update ASE (ArkApi) and ASA (AsaApi) plugins from ark-server-api.com in one click, with a tree editor for config.json
- 🕹️ Choice of game version (official release / preaquatica)
- ⬆️ One-click update of the commander itself from the settings page

## Near-term: Verify in Real Environments

These features are implemented but have only been tested against mocks or on local Docker Desktop. They need a full run on a real Linux server:

- [ ] Server list publishing, arkfgf.com account linking and P2P rooms (including multiple servers in one room) against the live website
- [ ] Native Linux server image, and switching a server between Linux and Proton
- [ ] Whether preaquatica works with the Linux server and with ArkApi
- [ ] The full one-click commander update flow on Linux (compose-created containers, registry mirrors, rollback)
- [ ] Whether the game finishes saving within the 30 seconds `docker stop` allows
- [ ] The player portal end to end on Linux (image pull, port mapping, ArkShop config from a real server)

Plus a leftover security item:

- [ ] Invalidate existing login tokens immediately after a password change, and add "sign out of all devices"

## Mid-term

- [ ] 👤 Multiple users: admin / user roles; add, remove and reset passwords from the settings page
- [ ] 🔐 Per-server access: regular users only see and manage the servers assigned to them
- [ ] 📜 Audit log: who did what to which server, and when
- [ ] 🔑 MCP tokens that can be limited to specific servers
- [ ] 🔍 In-app Steam Workshop mod search (requires your own Steam Web API key)
- [ ] 💽 Lower disk usage: each server currently keeps its own ~20 GB copy of the server files; explore sharing them between servers (approach TBD)

## Long-term

- [ ] ☸️ Multi-host management: one commander managing servers on several machines
- [ ] 👥 Player-facing interface: the player portal ([ArkPlayerPortal](ArkPlayerPortal/)) now has announcements, server status and the ArkShop store; next, let ArkServerLink open it directly
