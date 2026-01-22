import { createSlice, type PayloadAction } from '@reduxjs/toolkit'

interface UserState {
    playerId: string | null;
    isConnected: boolean;
    error: string | null;
}

const initialState: UserState = {
    playerId: null,
    isConnected: false,
    error: null,
}

export const userSlice = createSlice({
    name: 'user',
    initialState,
    reducers: {
        setPlayerId: (state, action: PayloadAction<string>) => {
            state.playerId = action.payload;
        },
        setConnected: (state, action: PayloadAction<boolean>) => {
            state.isConnected = action.payload;
        },
        setError: (state, action: PayloadAction<string | null>) => {
            state.error = action.payload;
        },
    },
})

export const { setPlayerId, setConnected, setError } = userSlice.actions
export default userSlice.reducer
