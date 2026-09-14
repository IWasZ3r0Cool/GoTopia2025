import { createAction, type Middleware } from '@reduxjs/toolkit'
import { clearGameState, updateGameState } from './gameSlice'
import { connectionClosed, connectionOpened, connectionStarted, setError, setPlayerId } from './userSlice'
import type { BuildingType, ServerMessage } from '../types/game'

export const connectToGame = createAction<string>('websocket/connect')
export const disconnectFromGame = createAction('websocket/disconnect')
export const buildBuilding = createAction<{ buildingType: BuildingType; x: number; y: number }>('game/build')

type SocketFactory = (url: string) => WebSocket

export function websocketURL(): string {
  const configured = import.meta.env.VITE_WS_URL as string | undefined
  if (configured) return configured
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/ws`
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

function isServerMessage(value: unknown): value is ServerMessage {
  if (!isRecord(value) || !isRecord(value.payload)) return false
  switch (value.type) {
    case 'WELCOME':
      return typeof value.payload.playerId === 'string'
    case 'ERROR':
      return typeof value.payload.code === 'string' && typeof value.payload.message === 'string'
    case 'GAME_STATE':
      return (
        typeof value.payload.mapWidth === 'number'
        && typeof value.payload.mapHeight === 'number'
        && typeof value.payload.turn === 'number'
        && isRecord(value.payload.players)
        && isRecord(value.payload.islands)
      )
    default:
      return false
  }
}

export const createWebsocketMiddleware = (
  createSocket: SocketFactory = (url) => new WebSocket(url),
  getURL: () => string = websocketURL,
): Middleware => (store) => {
  let socket: WebSocket | null = null

  return (next) => (action) => {
    const result = next(action)

    if (connectToGame.match(action)) {
      socket?.close(1000, 'Starting a new session')
      store.dispatch(clearGameState())
      store.dispatch(connectionStarted())

      const nextSocket = createSocket(getURL())
      socket = nextSocket

      nextSocket.onopen = () => {
        if (socket !== nextSocket) return
        store.dispatch(connectionOpened())
        nextSocket.send(JSON.stringify({ type: 'JOIN_GAME', payload: { name: action.payload } }))
      }

      nextSocket.onmessage = (event) => {
        if (socket !== nextSocket || typeof event.data !== 'string') return
        try {
          const message: unknown = JSON.parse(event.data)
          if (!isServerMessage(message)) throw new Error('Invalid server message')
          switch (message.type) {
            case 'WELCOME':
              store.dispatch(setPlayerId(message.payload.playerId))
              break
            case 'GAME_STATE':
              store.dispatch(updateGameState(message.payload))
              break
            case 'ERROR':
              store.dispatch(setError(message.payload.message))
              break
          }
        } catch {
          store.dispatch(setError('The server sent an unreadable response.'))
        }
      }

      nextSocket.onerror = () => {
        if (socket === nextSocket) store.dispatch(setError('Could not connect to the game server.'))
      }

      nextSocket.onclose = () => {
        if (socket !== nextSocket) return
        socket = null
        store.dispatch(connectionClosed())
        store.dispatch(clearGameState())
      }
    }

    if (disconnectFromGame.match(action)) {
      const activeSocket = socket
      socket = null
      activeSocket?.close(1000, 'Player left the game')
      store.dispatch(connectionClosed())
      store.dispatch(clearGameState())
    }

    if (buildBuilding.match(action)) {
      if (!socket || socket.readyState !== WebSocket.OPEN) {
        store.dispatch(setError('The game connection is not ready.'))
      } else {
        socket.send(JSON.stringify({ type: 'BUILD', payload: action.payload }))
      }
    }

    return result
  }
}

export const websocketMiddleware = createWebsocketMiddleware()
