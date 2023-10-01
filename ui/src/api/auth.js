import jwt_decode from "jwt-decode";
import { isTelegramAPISupported, getCloudValue, setCloudValue } from './telegram';
import { exchangeTicketForToken } from './guide';

const TICKET_URL_PARAM = "ticket"
const CLOUD_TOKEN_KEY = "AUTH_TOKEN"
const TOKEN_CHANGED_EVENT = "TOKEN_CHANGED"

const dispatcher = new EventTarget()
let state = {
    token: null,
    loaded: false,
}

// refreshes token if needed
// - if token is (null or expired) AND ticket is not null
// -- exchange ticket for token
// -- save token to CloudStorage
const refreshToken = (jwt) => {
    return new Promise((resolve, reject) => {
        const isValidToken = jwt?.claims?.exp > Date.now() / 1000;
        if (isValidToken) {
            return resolve(jwt.token);
        }

        const queryParameters = new URLSearchParams(window.location.search);
        const ticket = queryParameters.get(TICKET_URL_PARAM);
        if (ticket == null) {
            return resolve(null);
        }

        exchangeTicketForToken(ticket).then((response) => {
            const token = response.data.data.token;
            return setCloudValue(CLOUD_TOKEN_KEY, token).then(() => {
                return resolve(token);
            })
        }).catch((err) => {
            // 403 = ticket is already activated and cannot be exchanged
            if (err?.response?.status === 403) {
                resolve(null);
            } else {
                reject(err);
            }
        })
    })
}

const decodeJWT = (value) => {
    try {
        return { token: value, claims: jwt_decode(value) }
    } catch (err) {
        return null;
    }
}

const getToken = () => {
    state.loaded = false;
    getCloudValue(CLOUD_TOKEN_KEY).then((value) => {
        return decodeJWT(value)
    }).then((jwt) => {
        return refreshToken(jwt)
    }).then((token) => {
        state.token = token;
        state.loaded = true;
        const event = new CustomEvent(TOKEN_CHANGED_EVENT, { detail: token });
        dispatcher.dispatchEvent(event);
    })
}

if(isTelegramAPISupported()) {
    getToken();
}

export const addTokenListener = (callback) => {
    if (state.loaded) {
        callback(state.token);
    }
    dispatcher.addEventListener(TOKEN_CHANGED_EVENT, callback)
}

export const removeTokenListener = (callback) => {
    dispatcher.removeEventListener(TOKEN_CHANGED_EVENT, callback)
}