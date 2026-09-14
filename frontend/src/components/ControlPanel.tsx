import { BUILDINGS, type Building, type BuildingType, type Player } from '../types/game'

interface ControlPanelProps {
  player: Player
  selectedTile: { x: number; y: number } | null
  selectedBuilding?: Building
  onBuild: (type: BuildingType) => void
}

export function ControlPanel({ player, selectedTile, selectedBuilding, onBuild }: ControlPanelProps) {
  return (
    <aside className="control-panel" aria-labelledby="build-title">
      <div className="stats-grid">
        <div><span>Treasury</span><strong data-testid="stat-gold">{player.gold.toLocaleString()}</strong><small>gold</small></div>
        <div><span>Population</span><strong>{player.population.toLocaleString()} / {player.populationCapacity.toLocaleString()}</strong><small>citizens / capacity</small></div>
        <div><span>Public mood</span><strong>{player.mood}%</strong><small>content</small></div>
      </div>

      <div className="panel-heading build-heading">
        <div>
          <p className="eyebrow">Development office</p>
          <h2 id="build-title">Build</h2>
        </div>
        <span className="selection-label">
          {selectedTile ? `${selectedTile.x}, ${selectedTile.y}` : 'Select a tile'}
        </span>
      </div>

      {selectedBuilding && (
        <p className="occupied-notice">That tile contains a {selectedBuilding.type.toLowerCase()}.</p>
      )}

      <div className="building-menu">
        {BUILDINGS.map((building) => {
          const cannotAfford = player.gold < building.cost
          const disabled = !selectedTile || Boolean(selectedBuilding) || cannotAfford
          return (
            <button
              type="button"
              key={building.type}
              onClick={() => onBuild(building.type)}
              disabled={disabled}
              data-testid={`build-btn-${building.type}`}
              title={cannotAfford ? `Requires ${building.cost} gold` : building.description}
            >
              <span className="building-icon" aria-hidden="true">{building.icon}</span>
              <span className="building-details">
                <strong>{building.name}</strong>
                <small>{building.description}</small>
              </span>
              <span className="building-cost">{building.cost}</span>
            </button>
          )
        })}
      </div>
    </aside>
  )
}
