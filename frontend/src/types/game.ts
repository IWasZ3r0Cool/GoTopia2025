export interface Coordinate {
    x: number;
    y: number;
}

export type BuildingType = "HOUSE" | "FACTORY" | "FARM" | "FORT" | "SCHOOL" | "HOSPITAL";

export interface Building {
    type: BuildingType;
    health: number;
}

export interface Island {
    ownerId: string;
    // Key format: "x,y"
    buildings: Record<string, Building>;
    x: number;
    y: number;
    width: number;
    height: number;
}

export interface Player {
    id: string;
    gold: number;
    population: number;
    mood: number;
}

export interface GameState {
    mapWidth: number;
    mapHeight: number;
    islands: Record<string, Island>; // key is playerID
    players: Record<string, Player>; // key is playerID
    turn: number;
    round: number;
}

export interface WSMessage {
    type: string;
    payload?: any;
}
