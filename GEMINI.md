# Project: Gotopia2025

This document outlines the plan for a modern, real-time, 4-player remake of the classic Intellivision strategy game, "Utopia". The game will be a web-based application where each player manages an island, aiming for the highest score by fostering population growth, economic prosperity, and citizen happiness.

## Project Overview
The backend will be developed in Golang, serving as the authoritative game server. It will use WebSockets to manage and synchronize game state across all players in real-time. The frontend will be a responsive single-page application built with React and Redux, emphasizing a reusable component architecture for a clean and maintainable user interface.

## Core Gameplay Concepts
Objective: Achieve the highest score by the end of the game (e.g., after 20 rounds). Score is a composite of population, treasury, and island happiness. Game length is a configurable parameter at the start of the game. When the game “expires”, the players can vote for a 5 turn extension.

### Game Structure
- **Map**: A 200x200 global grid containing 4 separate islands surrounded by sea.
- **Islands**: Each player starts with one island of approximately 25x25 tiles. Islands are equidistant with at least 2 tiles of sea separation to ensure fair naval access from all sides.
- **Time System**: Hybrid Real-Time/Turn-Based.
    - **Actions**: Players can build and move units in real-time during a "Turn".
    - **Economy**: Revenue, population growth, and consumption occur only at the **End of Turn** (e.g., every 60 seconds). 

4-Player Real-Time: All four players build and manage their islands simultaneously. Actions and their consequences are broadcast to all players instantly.

### Island Management: 
Each player is the ruler of their own island. They can construct various buildings to influence the island's development:

- Housing: Increases population capacity.
- Farms and Factories: Generate income.
- Schools and Hospitals: Increase happiness and mitigate disasters.
- Forts: Allow players to build patrol boats to sink pirates or fishing boats, potentially impacting other players.
- Resorts: Allow players to generate income and employ their population.

### Resource Management: 
Players manage a single primary resource: Gold. Gold is generated at the **End of Turn** from economic buildings and is used to construct new buildings instantly.

### Buildings & Economics
- **Housing**: Cost 150. Significance: High population capacity.
- **Farms**: Cost 100. Generates Food/Income at end of turn.
- **Factory**: Cost 200. Generates high Income at end of turn.
- **Fort**: Cost 200. Range 5. Enables Patrol Boats.
- **School/Hospital**: higher costs, boost happiness/health.

Events & Disasters: The game server will periodically trigger events that affect one or all players.

Natural Disasters: Hurricanes, volcanoes, tidal waves, disease.

Player-Driven Events: A player's fishing boat can be sunk by a rival's patrol boat. Pirates may attack islands with low military presence.

### Technology Stack
| Area | Technology | Rationale |
|:-----|:-----------|:----------|
| Backend | Golang | Excellent for concurrency and networking. Lightweight, performant, and ideal for a real-time server.|
| Real-time | WebSockets (e.g., gorilla/websocket| Provides persistent, bidirectional communication for instant state synchronization with low overhead.|
| Frontend | React | Component-based architecture is perfect for building a complex, interactive game UI.|
| State Mgmt | Redux Toolkit | Provides predictable, centralized state management, essential for syncing complex game state from the server.|
| Styling | Tailwind CSS or Styled-Components | Facilitates the creation of reusable, scoped styles that align with the component-based architecture.|

### Architecture
- Backend (Golang)
The server is the single source of truth for all game state.

- Game Lobby: Manages game creation and allows players to join a game room. A unique room ID will be generated for each 4-player session.

- WebSocket Hub: A concurrent-safe structure that manages active WebSocket connections. It will handle player registration, de-registration, and broadcasting messages to all players in a specific game room.

#### Game Engine:

Contains the primary game loop. The loop handles the "Round Timer".

On each tick (real-time):
- Processes immediate player actions (Build, Move Boat).
- Updates boat positions/combat.

At **End of Turn**:
- Calculates income generation for each player (Factories/Farms).
- Updates population growth and food consumption.
- Checks victory conditions or round limits.
- Broadcasts the "Turn Complete" summary.

On each tick, the engine:



- Checks for and triggers random events.

- Updates the game round/timer.

- Broadcasts the updated GameState to all clients in the room.

#### State Model:
 Go structs will define the game's state, including Game, Player, Island, and Building.

#### API:

/ws?roomID=<id>: HTTP endpoint to upgrade the connection to a WebSocket.

Client actions (e.g., building a house) are sent as structured JSON messages over the WebSocket.

Frontend (React/Redux)
The UI is a dumb client that primarily renders the state received from the server.

Redux Store:

gameSlice: Holds the entire synchronized game state from the server. The WebSocket handler will dispatch actions to this slice to keep it up-to-date.

userSlice: Stores local user information, like their player ID and the current game room ID.

WebSocket Middleware: A Redux middleware will handle the WebSocket connection. It will listen for incoming messages from the server and dispatch corresponding actions to the Redux store. It will also send messages to the server when a user performs an action.

#### Reusable Components:

<IslandGrid>: A container component that renders the 2D grid of a player's island.

<Tile>: A small component representing a single square on the grid. It can display terrain or a <Building> component.

<Building>: A component to render a specific building type, with its own icon and state (e.g., health).

<ControlPanel>: The main UI for player actions, showing build options, resources, and player stats.

<Scoreboard>: A shared component that displays the current scores and status for all 4 players in the game.

<Notification>: A component to display alerts for game events and disasters.

### Data Flow Example: 
Player Builds a Factory
UI Action: A player clicks a <Tile> and selects "Factory" from the <ControlPanel>.

Dispatch: An onClick handler dispatches a Redux action: game/buildRequest({ building: 'FACTORY', position: {x: 5, y: 10} }).

Middleware: The WebSocket middleware intercepts this action and sends a JSON message to the Go server: {"action": "BUILD", "payload": {"type": "FACTORY", "x": 5, "y": 10}}.

Backend Logic: The Go server receives the message. It validates the request (Does the player have enough gold? Is the tile empty?).

State Update (Server): If the action is valid, the server updates the main GameState object in memory, adding the factory to the player's island grid and deducting the cost.

Broadcast: The server serializes the entire new GameState object and broadcasts it to all four connected clients via WebSockets.

Middleware (Client): Each client's WebSocket middleware receives the new state and dispatches a Redux action: game/stateSync(newState).

UI Update: The gameSlice reducer updates the store. React components subscribed to the store re-render automatically, and all four players now see the new factory on the board.