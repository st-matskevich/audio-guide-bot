import axios from 'axios';

const URL_BASE = process.env.REACT_APP_BOT_API_URL;

export const exchangeTicketForToken = (ticket) => {
    return axios.post(`${URL_BASE}/tickets/${ticket}/token`, {});
};

export const getObjectCoverURL = (accessToken, objectCode) => {
    return `${URL_BASE}/objects/${objectCode}/cover?access-token=${accessToken}`
}

export const getObjectAudioURL = (accessToken, objectCode) => {
    return `${URL_BASE}/objects/${objectCode}/audio?access-token=${accessToken}`
}