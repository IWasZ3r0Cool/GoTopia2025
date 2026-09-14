export const BUILDINGS = [
  { type: 'HOUSE', name: 'House', cost: 150, icon: '⌂', description: '+50 population capacity' },
  { type: 'FARM', name: 'Farm', cost: 100, icon: '≋', description: '+5 gold each turn' },
  { type: 'FACTORY', name: 'Factory', cost: 200, icon: '▥', description: '+20 gold each turn' },
  { type: 'FORT', name: 'Fort', cost: 200, icon: '♜', description: 'Protects nearby waters' },
  { type: 'SCHOOL', name: 'School', cost: 300, icon: '◇', description: 'Supports island wellbeing' },
  { type: 'HOSPITAL', name: 'Hospital', cost: 300, icon: '✚', description: 'Protects island health' },
] as const

export type BuildingType = (typeof BUILDINGS)[number]['type']

export interface Building {
  type: BuildingType
  health: number
}

export interface Island {
  ownerId: string
  buildings: Record<string, Building>
  x: number
  y: number
  width: number
  height: number
}

export interface Player {
  id: string
  name: string
  gold: number
  population: number
  populationCapacity: number
  mood: number
}

export interface GameState {
  mapWidth: number
  mapHeight: number
  islands: Record<string, Island>
  players: Record<string, Player>
  turn: number
}

export type ServerMessage =
  | { type: 'WELCOME'; payload: { playerId: string } }
  | { type: 'GAME_STATE'; payload: GameState }
  | { type: 'ERROR'; payload: { code: string; message: string } }

export type ConnectionStatus = 'idle' | 'connecting' | 'connected'
