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