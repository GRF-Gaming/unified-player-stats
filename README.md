# Unified Player Stats

The UPS project is developed with community engagement in mind – to keep track of player statistics across various games that are hosted by a community. The goal in mind is to bridge the gap between in-game and out-of-game interplayer interactions. 

1. Common API exposed, allowing multiple games to provide updates using mods.  
2. Persistent player statistics and analytics (such as kills, total playtime, etc.) that are viewable to players, to give them this feeling of accomplishment when they are out of the game.
3. Players can (ideally) talk about their stats with other players/community members in friendly competition that is enabled through periodic leaderboards.
4. Game server statistics on progress and game engine performance.
5. These data points would also provide community owners with valuable insights on player activity, engagement and overall player trends.

## Features 

**Backend**
- [ ] Individual Player Statistics
  - [x] Tracks the number of kills by each player (unique identifier).
  - [ ] Graphical representation of kill count over time.
  - [ ] Tracking of player name changes.
- [ ] Game Server Statistics
  - [ ] Online players
  - [ ] Server state
    - [ ] Total number of captures
    - [ ] Server performance (server FPS)
    - [ ] Number of entities on the server
- [ ] Discord integration
  - [ ] Kill feed updates
  - [ ] Event feed updates
  - [ ] Leaderboard updates
  - [ ] Player commands

**Further expansion**
- Dedicated documentation on API and internals.
- Web UI for admins and players to view statistics.
- Reporter for server/player statistics for integration into Prometheus/Grafana.
- Player verification across various games to allow players to view overall statistics.
