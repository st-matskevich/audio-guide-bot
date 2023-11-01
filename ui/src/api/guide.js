import axios from "axios";
import { getTelegramLanguage } from "./telegram";

const URL_BASE = window.REACT_APP_ENV.REACT_APP_BOT_API_URL;

export const exchangeTicketForToken = (ticket) => {
    return axios.post(`${URL_BASE}/tickets/${ticket}/token`, {});
};

export const getObjectData = (accessToken, objectCode) => {
    return axios.get(`${URL_BASE}/objects/${objectCode}?language=${getTelegramLanguage()}`, {
        headers: {
            "Authorization": accessToken
        }
    });
};

export const getObjectCoverURL = (accessToken, objectCode, coverIndex) => {
    return `${URL_BASE}/objects/${objectCode}/covers/${coverIndex}?access-token=${accessToken}&language=${getTelegramLanguage()}`;
};

export const getObjectAudioURL = (accessToken, objectCode) => {
    return `${URL_BASE}/objects/${objectCode}/audio?access-token=${accessToken}&language=${getTelegramLanguage()}`;
};