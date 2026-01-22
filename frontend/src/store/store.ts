import { configureStore } from '@reduxjs/toolkit'
import gameReducer from './gameSlice'
import userReducer from './userSlice'
import { websocketMiddleware } from './middleware'

export const store = configureStore({
    reducer: {
        game: gameReducer,
        user: userReducer,
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(websocketMiddleware),
})

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch
