import type { Island } from '../types/game'
import { Tile } from './Tile'

interface IslandGridProps {
  island: Island
  ownerName: string
  onTileClick?: (x: number, y: number) => void
  selectedTile?: { x: number; y: number } | null
  compact?: boolean
}

export function IslandGrid({ island, ownerName, onTileClick, selectedTile = null, compact = false }: IslandGridProps) {
  const tiles = []
  for (let dy = 0; dy < island.height; dy += 1) {
    for (let dx = 0; dx < island.width; dx += 1) {
      const x = island.x + dx
      const y = island.y + dy
      const key = `${x},${y}`
      tiles.push(
        <Tile
          key={key}
          x={x}
          y={y}
          building={island.buildings[key]}
          onSelect={onTileClick ?? (() => undefined)}
          isSelected={selectedTile?.x === x && selectedTile?.y === y}
          interactive={Boolean(onTileClick)}
        />,
      )
    }
  }

  return (
    <section className={`island${compact ? ' island-compact' : ''}`} aria-label={`${ownerName}'s island`}>
      <div className="island-title-row">
        <div>
          <p className="eyebrow">{compact ? 'Across the water' : 'Your island'}</p>
          <h2>{ownerName}</h2>
        </div>
        <span className="building-count">{Object.keys(island.buildings).length} structures</span>
      </div>
      <div
        className="shoreline"
        role={compact ? undefined : 'region'}
        aria-label={compact ? undefined : 'Scrollable island map'}
        tabIndex={compact ? undefined : 0}
      >
        <div
          className="island-grid"
          style={{
            gridTemplateColumns: `repeat(${island.width}, minmax(0, 1fr))`,
            minWidth: compact ? undefined : `${island.width * 36}px`,
          }}
          data-testid={`island-grid-${island.ownerId}`}
        >
          {tiles}
        </div>
      </div>
    </section>
  )
}
