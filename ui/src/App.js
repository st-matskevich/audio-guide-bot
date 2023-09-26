import './App.css';
import { useEffect, useState } from "react"
import jwt_decode from "jwt-decode";
import { getCloudValue, setCloudValue } from './api/telegram';
import { exchangeTicketForToken } from './api/guide';

function App() {
  const onScanQRClicked = () => {
    // Example of interaction with telegram-web-app.js script to show a native alert
    window.Telegram.WebApp.showScanQrPopup({}, () => { return true; });
  };

  const onCloseClicked = () => {
    window.Telegram.WebApp.close();
  }

  const [tokenState, setTokenState] = useState({
    loaded: false,
    token: null
  })

  useEffect(() => {
    // - load token from CloudStorage
    // - verify token
    // - if token is null or expired AND ticket is not null
    // -- exchange ticket for token
    // -- save token to CloudStorage
    const CLOUD_TOKEN_KEY = "AUTH_TOKEN"

    // refreshes token if needed
    const refreshToken = (jwt) => {
      return new Promise((resolve, reject) => {
        const isValidToken = jwt?.claims?.exp > Date.now() / 1000;
        if (isValidToken) {
          return resolve(jwt.token);
        }

        const queryParameters = new URLSearchParams(window.location.search);
        const ticket = queryParameters.get("ticket");
        if (ticket == null) {
          return resolve(null);
        }

        exchangeTicketForToken(ticket).then((response) => {
          const token = response.data.data.token;
          setCloudValue(CLOUD_TOKEN_KEY, token).then(() => {
            return resolve(token)
          })
        }).catch((err) => {
          return reject(err)
        })
      })
    }

    // decodes JWT if valid
    const decodeJWT = (value) => {
      try {
        return { token: value, claims: jwt_decode(value) }
      } catch (err) {
        return null;
      }
    }

    getCloudValue(CLOUD_TOKEN_KEY).then((value) => {
      return decodeJWT(value)
    }).then((jwt) => {
      return refreshToken(jwt)
    }).then((token) => {
      setTokenState({
        loaded: true,
        token: token
      })
    })

  }, [])

  const getUI = () => {
    if (!tokenState.loaded) {
      return (<div className="preloader" />)
    } else if (tokenState.token == null) {
      return (
        <div className="no-ticket-wrapper">
          <span>It seems you haven't purchased a ticket yet.</span>
          <span>To start our tour, please go back to the bot and buy a ticket.</span>
          <div className="button" onClick={onCloseClicked}>Close app</div>
        </div>
      )
    } else {
      return (
        <div className="scanner-wrapper">
          <span>Welcome to the tour!</span>
          <span>Scan QR codes to start listening.</span>
          <div className="button" onClick={onScanQRClicked}>Scan QR</div>
        </div>
      )
    }
  }

  return (
    <div className="App">
      {getUI()}
    </div>
  );
}

export default App;
