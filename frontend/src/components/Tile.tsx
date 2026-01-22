import type { Building, BuildingType } from '../types/game'

interface TileProps {
    x: number;
    y: number;
    building?: Building;
    onBuild: (x: number, y: number) => void;
    isSelected: boolean;
}

const getBuildingColor = (type?: BuildingType) => {
    switch (type) {
        case "HOUSE": return "bg-blue-500 shadow-blue-500/50";
        case "FACTORY": return "bg-amber-600 shadow-amber-600/50";
        case "FARM": return "bg-green-500 shadow-green-500/50";
        case "FORT": return "bg-stone-700 shadow-stone-700/50";
        case "SCHOOL": return "bg-purple-500 shadow-purple-500/50";
        case "HOSPITAL": return "bg-red-500 shadow-red-500/50";
        default: return "bg-emerald-200 hover:bg-emerald-300"; // Grass
    }
}

const getBuildingIcon = (type?: BuildingType) => {
    switch (type) {
        case "HOUSE": return "🏠";
        case "FACTORY": return "🏭";
        case "FARM": return "🌾";
        case "FORT": return "🏰";
        case "SCHOOL": return "🏫";
        case "HOSPITAL": return "🏥";
        default: return "";
    }
}

export const Tile = ({ x, y, building, onBuild, isSelected }: TileProps) => {
    const baseClass = "w-8 h-8 m-0.5 rounded cursor-pointer transition-all duration-200 flex items-center justify-center text-sm shadow-md";
    const colorClass = getBuildingColor(building?.type);
    const borderClass = isSelected ? "ring-2 ring-white scale-110 z-10" : "";

    return (
        <div
            className={`${baseClass} ${colorClass} ${borderClass}`}
            onClick={() => onBuild(x, y)}
            title={`(${x},${y}) ${building?.type || "Empty"}`}
            role="button"
            aria-label={`Tile ${x},${y}`}
        >
            {getBuildingIcon(building?.type)}
        </div>
    )
}
