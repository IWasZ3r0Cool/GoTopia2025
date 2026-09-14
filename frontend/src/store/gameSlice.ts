import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { GameState } from '../types/game'

interface GameSliceState {
  data: GameState | null
}

const initialState: GameSliceState = { data: null }

export const gameSlice = createSlice({
  name: 'game',
  initialState,
  reducers: {
    updateGameState: (state, action: PayloadAction<GameState>) => {
      state.data = action.payload
    },
    clearGameState: (state) => {
      state.data = null
    },
  },
})

export const { updateGameState, clearGameState } = gameSlice.actions
export default gameSlice.reducer
