import { useState, useEffect } from 'react'
import { useAppSelector, useAppDispatch } from '../store/hooks'
import { connectToGame, buildBuilding } from '../store/middleware'
import { IslandGrid } from './IslandGrid'
import { ControlPanel } from './ControlPanel'
import type { BuildingType } from '../types/game'

export const GameOverlay = () => {
    const dispatch = useAppDispatch();
    const { data: gameState } = useAppSelector(state => state.game);
    const { playerId, isConnected } = useAppSelector(state => state.user);

    // Local selection state
    const [selectedTile, setSelectedTile] = useState<{ x: number, y: number } | null>(null);

    useEffect(() => {
        // Auto connect on mount
        if (!isConnected) {
            dispatch(connectToGame("Player_" + Math.floor(Math.random() * 1000)));
        }
    }, [dispatch, isConnected]);

    if (!isConnected || !gameState) {
        return (
            <div className="min-h-screen bg-gray-900 flex items-center justify-center text-white flex-col gap-4">
                <div className="animate-spin text-4xl">🌀</div>
                <h1 className="text-2xl font-light tracking-widest uppercase">Connecting to GoTopia...</h1>
            </div>
        );
    }

    const islandKeys = Object.keys(gameState.islands);

    return (
        <div className="min-h-screen bg-slate-900 overflow-auto pb-32 relative">
            {/* Header Info */}
            <div className="absolute top-4 left-0 right-0 text-center pointer-events-none z-10">
                <div className="inline-block bg-black/50 backdrop-blur-md px-6 py-2 rounded-full border border-white/10 text-white font-mono shadow-lg">
                    Turn: <span className="text-green-400 font-bold">{gameState.turn}</span>
                    <span className="mx-4 text-gray-500">|</span>
                    Player: <span className="text-blue-400 font-bold">{playerId}</span>
                </div>
            </div>

            {/* Map Area */}
            <div className="grid grid-cols-2 gap-8 p-12 min-w-fit mx-auto justify-items-center items-center h-full">
                {islandKeys.map(key => (
                    <div key={key} className={key === playerId ? "ring-4 ring-yellow-400 rounded-xl" : "opacity-80 hover:opacity-100 transition-opacity"}>
                        <IslandGrid
                            island={gameState.islands[key]}
                            onTileClick={(x, y) => {
                                // Only allow selecting MY island tiles
                                if (key === playerId) {
                                    setSelectedTile({ x, y });
                                }
                            }}
                            selectedTile={selectedTile}
                        />
                    </div>
                ))}
            </div>

            <ControlPanel
                player={playerId && gameState.players[playerId] ? gameState.players[playerId] : null}
                selectedTile={selectedTile}
                onBuild={(type: BuildingType) => {
                    if (selectedTile) {
                        dispatch(buildBuilding(type, selectedTile.x, selectedTile.y));
                    }
                }}
            />
        </div>
    )
}
