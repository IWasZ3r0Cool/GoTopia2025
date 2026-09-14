import type { Building } from '../types/game'
import { BUILDINGS } from '../types/game'

interface TileProps {
  x: number
  y: number
  building?: Building
  onSelect: (x: number, y: number) => void
  isSelected: boolean
  interactive?: boolean
}

export function Tile({ x, y, building, onSelect, isSelected, interactive = true }: TileProps) {
  const definition = building ? BUILDINGS.find((item) => item.type === building.type) : undefined
  const label = building ? `${definition?.name ?? building.type} at ${x}, ${y}` : `Empty tile at ${x}, ${y}`

  return (
    <button
      type="button"
      className={`tile${building ? ` tile-${building.type.toLowerCase()}` : ''}${isSelected ? ' is-selected' : ''}`}
      onClick={() => onSelect(x, y)}
      aria-label={label}
      aria-pressed={isSelected}
      disabled={!interactive}
      title={label}
    >
      {definition?.icon && <span aria-hidden="true">{definition.icon}</span>}
    </button>
  )
}
