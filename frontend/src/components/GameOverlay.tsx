import { useMemo, useState } from 'react'
import { useAppDispatch, useAppSelector } from '../store/hooks'
import { buildBuilding, connectToGame, disconnectFromGame } from '../store/middleware'
import { setError } from '../store/userSlice'
import { ControlPanel } from './ControlPanel'
import { IslandGrid } from './IslandGrid'
import { JoinScreen } from './JoinScreen'
import { Scoreboard } from './Scoreboard'
import type { BuildingType } from '../types/game'

export function GameOverlay() {
  const dispatch = useAppDispatch()
  const gameState = useAppSelector((state) => state.game.data)
  const { playerId, status, error } = useAppSelector((state) => state.user)
  const [selection, setSelection] = useState<{ playerId: string; x: number; y: number } | null>(null)
  const selectedTile = selection?.playerId === playerId ? selection : null

  const players = useMemo(() => Object.values(gameState?.players ?? {}), [gameState])

  if (status === 'idle' || (status === 'connecting' && !gameState)) {
    return (
      <JoinScreen
        status={status}
        error={error}
        onJoin={(name) => dispatch(connectToGame(name))}
      />
    )
  }

  const player = playerId ? gameState?.players[playerId] : undefined
  const island = playerId ? gameState?.islands[playerId] : undefined
  if (!gameState || !playerId || !player || !island) {
    return (
      <main className="loading-screen" aria-live="polite">
        <div className="loading-mark" aria-hidden="true">≈</div>
        {error ? (
          <>
            <p className="eyebrow">Unable to join</p>
            <h1>{error}</h1>
          </>
        ) : (
          <>
            <p className="eyebrow">Charting the waters</p>
            <h1>Preparing your island…</h1>
          </>
        )}
        <button type="button" className="text-button" onClick={() => dispatch(disconnectFromGame())}>
          {error ? 'Back' : 'Cancel'}
        </button>
      </main>
    )
  }

  const selectedKey = selectedTile ? `${selectedTile.x},${selectedTile.y}` : ''
  const selectedBuilding = island.buildings[selectedKey]
  const opponents = players.filter((candidate) => candidate.id !== playerId)

  function build(type: BuildingType) {
    if (!selectedTile || selectedBuilding) return
    dispatch(setError(null))
    dispatch(buildBuilding({ buildingType: type, x: selectedTile.x, y: selectedTile.y }))
  }

  return (
    <div className="game-shell">
      <header className="game-header">
        <a className="wordmark" href="/" aria-label="GoTopia home"><span>G</span> GoTopia</a>
        <div className="turn-indicator"><span>Turn</span><strong>{gameState.turn}</strong></div>
        <div className="header-actions">
          <span className="connection-pill"><i /> Live · {players.length}/4 rulers</span>
          <button type="button" className="text-button" onClick={() => dispatch(disconnectFromGame())}>Leave game</button>
        </div>
      </header>

      {error && (
        <div className="toast" role="alert">
          <span>{error}</span>
          <button type="button" aria-label="Dismiss error" onClick={() => dispatch(setError(null))}>×</button>
        </div>
      )}

      <main className="game-content">
        <div className="game-board-column">
          <IslandGrid
            island={island}
            ownerName={player.name}
            onTileClick={(x, y) => setSelection({ playerId, x, y })}
            selectedTile={selectedTile}
          />

          {opponents.length > 0 && (
            <section className="opponent-section" aria-labelledby="opponents-title">
              <div className="panel-heading">
                <p className="eyebrow">World view</p>
                <h2 id="opponents-title">Neighboring islands</h2>
              </div>
              <div className="opponent-grid">
                {opponents.map((opponent) => (
                  <IslandGrid
                    key={opponent.id}
                    island={gameState.islands[opponent.id]}
                    ownerName={opponent.name}
                    compact
                  />
                ))}
              </div>
            </section>
          )}
        </div>

        <div className="sidebar-column">
          <ControlPanel
            player={player}
            selectedTile={selectedTile}
            selectedBuilding={selectedBuilding}
            onBuild={build}
          />
          <Scoreboard players={players} currentPlayerId={playerId} />
        </div>
      </main>
    </div>
  )
}
