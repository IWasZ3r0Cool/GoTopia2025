import type { Island } from '../types/game'
import { Tile } from './Tile'

interface IslandGridProps {
    island: Island;
    onTileClick: (x: number, y: number) => void;
    selectedTile: { x: number, y: number } | null;
}

export const IslandGrid = ({ island, onTileClick, selectedTile }: IslandGridProps) => {
    const tiles = [];

    for (let dy = 0; dy < island.height; dy++) {
        for (let dx = 0; dx < island.width; dx++) {
            const gx = island.x + dx;
            const gy = island.y + dy;
            const key = `${gx},${gy}`;
            const building = island.buildings[key];

            const isSelected = selectedTile?.x === gx && selectedTile?.y === gy;

            tiles.push(
                <Tile
                    key={key}
                    x={gx}
                    y={gy}
                    building={building}
                    onBuild={onTileClick}
                    isSelected={isSelected}
                />
            );
        }
    }

    return (
        <div className="bg-blue-300 p-4 rounded-xl shadow-2xl border-4 border-blue-400/50 backdrop-blur-sm">
            <h2 className="text-center font-bold text-blue-900 mb-2 drop-shadow-sm uppercase tracking-wider">
                Island {island.ownerId}
            </h2>
            <div
                className="grid gap-0 bg-emerald-900/10 p-1 rounded-lg"
                style={{
                    gridTemplateColumns: `repeat(${island.width}, minmax(0, 1fr))`
                }}
                data-testid={`island-grid-${island.ownerId}`}
            >
                {tiles}
            </div>
        </div>
    )
}
