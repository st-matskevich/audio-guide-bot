# Audio Guide Bot: UI service
UI service is a dockerized React application that is served as Telegram Mini App.

Service container comes in two versions:
- Development image is used to mount source code as docker volume and start React dev-server to provide React fast refresh. Image is built using [dev.dockerfile](./dev.dockerfile).
- Production image is used to serve built React application. It runs nginx daemon inside and allows setting environment variables in runtime. Image is built using [prod.dockerfile](./prod.dockerfile).

## Configuration
Required environment variables:
- `REACT_APP_BOT_API_URL` - URL to Guide API service

## React app structure
Components are used to build UI:
- [App.js](./src/App.js) - entry component, used to handle authentication flow, interact with Telegram's API and show the initial interface
- [ObjectViewerComponent.js](./src/components/ObjectViewerComponent.js) - component that handles object viewing, loads object metadata, and if it exists - shows cover image and loads audio
- [MarqueeComponent.js](./src/components/MarqueeComponent.js) - replacement for an obsolete [\<marquee>](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/marquee) HTML element
- [SliderComponent.js](./src/components/SliderComponent.js) - range input element with custom styling

Helpers are used to decouple and isolate APIs from React:
- [auth.js](./src/api/auth.js) - implements token storage and refresh logic
- [telegram.js](./src/api/telegram.js) - wrapper for Telegram APIs calls
- [guide.js](./src/api/guide.js) - wrapper for Guide APIs calls

## nginx configuration in the production image
[nginx.template](./nginx.template) is used to configure nginx in the production image:
- Allows to set port, where content is served, by `PORT` environment variable
- Sets to serve build React application from the standard folder

## Runtime environment variables
Using runtime environment variables from React is a known problem because it provides only build-time environment variables. UI service implements setting environment variables in runtime.

This is done by creating `env-config.js` in `public` folder and including it in the `index.html` file. `env-config.js` defines `REACT_APP_ENV` field in `window` object with all variables prefixed by `REACT_APP`. 

Environment variables in runtime are implemented by two files:
- [setup-env-config.js](./setup-env-config.js) - is used in `package.json` to allow executing on all OSes and creates an initial version of `env-config.js` file
- [setup-env-config.sh](./setup-env-config.sh) - is used as an additional entrypoint script in the production version of the service container and appends additional variables to `env-config.js` that were specified on container creation