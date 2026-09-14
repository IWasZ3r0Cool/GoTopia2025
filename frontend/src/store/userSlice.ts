import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { ConnectionStatus } from '../types/game'

interface UserState {
  playerId: string | null
  status: ConnectionStatus
  error: string | null
}

const initialState: UserState = { playerId: null, status: 'idle', error: null }

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    connectionStarted: (state) => {
      state.status = 'connecting'
      state.playerId = null
      state.error = null
    },
    connectionOpened: (state) => {
      state.status = 'connected'
    },
    connectionClosed: (state) => {
      state.status = 'idle'
      state.playerId = null
    },
    setPlayerId: (state, action: PayloadAction<string>) => {
      state.playerId = action.payload
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload
    },
  },
})

export const { connectionStarted, connectionOpened, connectionClosed, setPlayerId, setError } = userSlice.actions
export default userSlice.reducer
