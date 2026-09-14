import { configureStore } from '@reduxjs/toolkit'
import { describe, expect, it } from 'vitest'
import gameReducer from './gameSlice'
import { buildBuilding, connectToGame, createWebsocketMiddleware, disconnectFromGame } from './middleware'
import userReducer from './userSlice'

class FakeSocket {
  readyState: number = WebSocket.CONNECTING
  sent: string[] = []
  closed = false
  onopen: ((event: Event) => void) | null = null
  onmessage: ((event: MessageEvent) => void) | null = null
  onerror: ((event: Event) => void) | null = null
  onclose: ((event: CloseEvent) => void) | null = null

  send(message: string) { this.sent.push(message) }
  close() { this.closed = true }
  open() {
    this.readyState = WebSocket.OPEN
    this.onopen?.(new Event('open'))
  }
  receive(value: unknown) {
    this.onmessage?.(new MessageEvent('message', { data: JSON.stringify(value) }))
  }
}

function setup() {
  const sockets: FakeSocket[] = []
  const middleware = createWebsocketMiddleware(() => {
    const socket = new FakeSocket()
    sockets.push(socket)
    return socket as unknown as WebSocket
  }, () => 'ws://test/ws')
  const store = configureStore({
    reducer: { game: gameReducer, user: userReducer },
    middleware: (defaults) => defaults().concat(middleware),
  })
  return { store, sockets }
}

describe('websocket middleware', () => {
  it('joins, receives identity and state, builds, and disconnects', () => {
    const { store, sockets } = setup()
    store.dispatch(connectToGame('Ada'))
    expect(store.getState().user.status).toBe('connecting')

    const socket = sockets[0]
    socket.open()
    expect(store.getState().user.status).toBe('connected')
    expect(JSON.parse(socket.sent[0])).toEqual({ type: 'JOIN_GAME', payload: { name: 'Ada' } })

    socket.receive({ type: 'WELCOME', payload: { playerId: 'player-1' } })
    socket.receive({ type: 'GAME_STATE', payload: { mapWidth: 200, mapHeight: 200, islands: {}, players: {}, turn: 1 } })
    expect(store.getState().user.playerId).toBe('player-1')
    expect(store.getState().game.data?.turn).toBe(1)

    store.dispatch(buildBuilding({ buildingType: 'FARM', x: 10, y: 10 }))
    expect(JSON.parse(socket.sent[1])).toEqual({ type: 'BUILD', payload: { buildingType: 'FARM', x: 10, y: 10 } })

    store.dispatch(disconnectFromGame())
    expect(socket.closed).toBe(true)
    expect(store.getState().user.status).toBe('idle')
    expect(store.getState().game.data).toBeNull()
  })

  it('surfaces malformed server messages', () => {
    const { store, sockets } = setup()
    store.dispatch(connectToGame('Lin'))
    sockets[0].open()
    sockets[0].onmessage?.(new MessageEvent('message', { data: 'not-json' }))
    expect(store.getState().user.error).toBe('The server sent an unreadable response.')
  })
})
