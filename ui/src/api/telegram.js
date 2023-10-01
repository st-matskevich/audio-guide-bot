const MIN_TELEGRAM_BOT_API = "6.9";

const isVersionGreaterOrEqual = (versionA, versionB) => {
    const vAParts = versionA.split('.').map(Number);
    const vBParts = versionB.split('.').map(Number);

    for (let i = 0; i < Math.min(vAParts.length, vBParts.length); i++) {
        const partA = vAParts[i];
        const partB = vBParts[i];

        if (partA < partB) {
            return false;
        } else if (partA > partB) {
            return true;
        }
    }

    return true;
}

export const isTelegramAPISupported = () => {
    return isVersionGreaterOrEqual(window.Telegram.WebApp.version, MIN_TELEGRAM_BOT_API);
}

export const getCloudValue = (key) => {
    return new Promise((resolve, reject) => {
        window.Telegram.WebApp.CloudStorage.getItem(key, (err, value) => {
            if (err != null) {
                return reject(err);
            } else {
                return resolve(value);
            }
        })
    })
};

export const setCloudValue = (key, value) => {
    return new Promise((resolve, reject) => {
        window.Telegram.WebApp.CloudStorage.setItem(key, value, (err, result) => {
            if (err != null) {
                return reject(err);
            } else {
                return resolve(result);
            }
        })
    })
};