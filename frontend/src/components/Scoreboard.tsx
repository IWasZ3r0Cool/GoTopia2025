import type { Player } from '../types/game'

interface ScoreboardProps {
  players: Player[]
  currentPlayerId: string
}

export function Scoreboard({ players, currentPlayerId }: ScoreboardProps) {
  const ranked = [...players].sort((a, b) => b.gold - a.gold)
  return (
    <section className="scoreboard" aria-labelledby="scoreboard-title">
      <div className="panel-heading">
        <p className="eyebrow">The archipelago</p>
        <h2 id="scoreboard-title">Rulers</h2>
      </div>
      <ol>
        {ranked.map((player, index) => (
          <li key={player.id} className={player.id === currentPlayerId ? 'is-you' : ''}>
            <span className="rank">{index + 1}</span>
            <span className="ruler-name">{player.name}{player.id === currentPlayerId && <small>You</small>}</span>
            <strong>{player.gold.toLocaleString()} g</strong>
          </li>
        ))}
      </ol>
    </section>
  )
}
