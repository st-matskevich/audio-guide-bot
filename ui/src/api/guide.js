import axios from 'axios';

const URL_BASE = process.env.REACT_APP_BOT_API_URL;

export const exchangeTicketForToken = (ticket) => {
    return axios.post(`${URL_BASE}/tickets/${ticket}/token`, {});
};