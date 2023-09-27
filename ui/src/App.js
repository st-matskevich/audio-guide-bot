import './App.css';
import { useEffect, useState } from "react"
import { addTokenListener, removeTokenListener } from './api/auth';

function App() {
  const onScanQRClicked = () => {
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
    const onTokenChanged = (event) => {
      setTokenState({
        loaded: true,
        token: event.detail
      })
    };

    addTokenListener(onTokenChanged);
    return () => {
      removeTokenListener(onTokenChanged);
    }
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
