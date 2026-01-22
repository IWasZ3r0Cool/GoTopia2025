import type { Middleware, UnknownAction } from 'redux'
import { updateGameState } from './gameSlice'
import { setConnected, setError } from './userSlice'
import type { WSMessage } from '../types/game'

// Check for custom action types
interface CustomAction extends UnknownAction {
    payload?: any;
}

export const websocketMiddleware: Middleware = store => {
    let socket: WebSocket | null = null;

    return next => (action: unknown) => {
        const typedAction = action as CustomAction;
        // Defines middleware actions
        if (typedAction.type === 'WS_CONNECT') {
            if (socket !== null) {
                socket.close();
            }

            // Connect to localhost:8080 (or config)
            socket = new WebSocket('ws://localhost:8080/ws');

            socket.onopen = () => {
                store.dispatch(setConnected(true));
                store.dispatch(setError(null));

                // Join Game automatically for now with random name or prompted
                const name = typedAction.payload || "Player_" + Math.floor(Math.random() * 1000);
                socket?.send(JSON.stringify({
                    type: "JOIN_GAME",
                    payload: { name }
                }));
            };

            socket.onmessage = (event) => {
                try {
                    const msg: WSMessage = JSON.parse(event.data);

                    switch (msg.type) {
                        case "GAME_STATE":
                            store.dispatch(updateGameState(msg.payload));
                            break;
                        case "ERROR":
                            store.dispatch(setError(msg.payload.message));
                            break;
                        default:
                            console.log("Unknown WS Message:", msg);
                    }
                } catch (e) {
                    console.error("WS Parse Error", e);
                }
            };

            socket.onclose = () => {
                store.dispatch(setConnected(false));
            };

            socket.onerror = () => {
                store.dispatch(setError("WebSocket connection failed"));
            };
        } else if (typedAction.type === 'WS_SEND') {
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify(typedAction.payload));
            }
        }
        // Actions to trigger specific Game Actions
        else if (typedAction.type === 'GAME_BUILD') {
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({
                    type: "BUILD",
                    payload: typedAction.payload
                }));
            }
        }

        return next(action);
    };
};

// Action creators for the middleware
export const connectToGame = (playerName?: string) => ({ type: 'WS_CONNECT', payload: playerName });
export const sendWSMessage = (msg: any) => ({ type: 'WS_SEND', payload: msg });
export const buildBuilding = (type: string, x: number, y: number) => ({
    type: 'GAME_BUILD',
    payload: { buildingType: type, x, y }
});
