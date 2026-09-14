import { configureStore } from '@reduxjs/toolkit'
import gameReducer from './gameSlice'
import userReducer from './userSlice'
import { websocketMiddleware } from './middleware'

export const store = configureStore({
  reducer: {
    game: gameReducer,
    user: userReducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(websocketMiddleware),
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
