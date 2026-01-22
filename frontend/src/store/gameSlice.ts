import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { GameState } from '../types/game'

interface GameSliceState {
    data: GameState | null;
    lastUpdated: number;
}

const initialState: GameSliceState = {
    data: null,
    lastUpdated: 0,
}

export const gameSlice = createSlice({
    name: 'game',
    initialState,
    reducers: {
        updateGameState: (state, action: PayloadAction<GameState>) => {
            state.data = action.payload;
            state.lastUpdated = Date.now();
        },
    },
})

export const { updateGameState } = gameSlice.actions
export default gameSlice.reducer
