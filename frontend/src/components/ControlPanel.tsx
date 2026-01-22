import type { BuildingType, Player } from '../types/game'

interface ControlPanelProps {
    player: Player | null;
    selectedTile: { x: number, y: number } | null;
    onBuild: (type: BuildingType) => void;
}

export const ControlPanel = ({ player, selectedTile, onBuild }: ControlPanelProps) => {
    if (!player) return <div className="p-4 text-white">Loading...</div>;

    const buildings: { type: BuildingType, cost: number, icon: string }[] = [
        { type: "HOUSE", cost: 150, icon: "🏠" },
        { type: "FARM", cost: 100, icon: "🌾" },
        { type: "FACTORY", cost: 200, icon: "🏭" },
        { type: "FORT", cost: 200, icon: "🏰" },
        { type: "SCHOOL", cost: 300, icon: "🏫" },
        { type: "HOSPITAL", cost: 300, icon: "🏥" },
    ];

    return (
        <div className="fixed bottom-0 left-0 right-0 bg-gray-900/95 text-white p-6 border-t border-gray-700 shadow-2xl backdrop-blur-md flex justify-between items-center z-50">
            {/* Stats */}
            <div className="flex gap-8 text-lg font-mono">
                <div className="flex flex-col items-center">
                    <span className="text-gray-400 text-xs uppercase tracking-widest">Gold</span>
                    <span className="text-yellow-400 font-bold text-2xl flex items-center gap-2" data-testid="stat-gold">
                        💰 {player.gold}
                    </span>
                </div>
                <div className="flex flex-col items-center">
                    <span className="text-gray-400 text-xs uppercase tracking-widest">Pop</span>
                    <span className="text-blue-400 font-bold text-2xl flex items-center gap-2">
                        👥 {player.population}
                    </span>
                </div>
                <div className="flex flex-col items-center">
                    <span className="text-gray-400 text-xs uppercase tracking-widest">Mood</span>
                    <span className="text-pink-400 font-bold text-2xl flex items-center gap-2">
                        ❤️ {player.mood}%
                    </span>
                </div>
                <div className="flex flex-col items-center">
                    <span className="text-gray-400 text-xs uppercase tracking-widest">Selection</span>
                    <span className="text-white font-bold text-xl">
                        {selectedTile ? `(${selectedTile.x}, ${selectedTile.y})` : "None"}
                    </span>
                </div>
            </div>

            {/* Build Menu */}
            <div className="flex gap-3">
                {buildings.map((b) => (
                    <button
                        key={b.type}
                        onClick={() => onBuild(b.type)}
                        disabled={!selectedTile || player.gold < b.cost}
                        className={`
                            flex flex-col items-center p-3 rounded-lg border-2 transition-all duration-200
                            ${!selectedTile ? 'opacity-50 cursor-not-allowed border-gray-700 bg-gray-800' :
                                player.gold < b.cost ? 'opacity-50 cursor-not-allowed border-red-900 bg-red-900/20' :
                                    'border-gray-600 bg-gray-800 hover:bg-gray-700 hover:border-gray-400 hover:scale-105 active:scale-95 cursor-pointer'}
                        `}
                        data-testid={`build-btn-${b.type}`}
                    >
                        <span className="text-2xl mb-1">{b.icon}</span>
                        <span className="text-xs font-bold uppercase">{b.type}</span>
                        <span className="text-xs text-yellow-500 font-mono">${b.cost}</span>
                    </button>
                ))}
            </div>
        </div>
    )
}
