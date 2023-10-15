const fs = require("fs");
const variables = process.env;

const config = {};
for (const key in variables) {
    if (key.startsWith("REACT_APP")) {
        config[key] = variables[key];
    }
}

const script = `window.REACT_APP_ENV = ${JSON.stringify(config, null, 2)};`;
fs.writeFileSync(__dirname + "/public/env-config.js", script);