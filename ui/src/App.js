import './App.css';
import { useEffect, useState } from "react"
import { addTokenListener, removeTokenListener } from './api/auth';
import ObjectViewerComponent from './ObjectViewerComponent'
import { isTelegramAPISupported } from './api/telegram';

function App() {
  const isSupported = isTelegramAPISupported();
  const onScanQRClicked = () => {
    window.Telegram.WebApp.showScanQrPopup({});
  };

  const onCloseClicked = () => {
    window.Telegram.WebApp.close();
  }

  const [scannedObject, setScannedObject] = useState(null)
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

  useEffect(() => {
    const onBackClicked = () => {
      window.Telegram.WebApp.BackButton.isVisible = false;
      setScannedObject(null);
    }

    window.Telegram.WebApp.BackButton.onClick(onBackClicked);
    return () => {
      window.Telegram.WebApp.BackButton.offClick(onBackClicked);
    }
  }, [])

  useEffect(() => {
    const QR_EVENT = "qrTextReceived";
    const onQRTextReceived = (event) => {
      window.Telegram.WebApp.closeScanQrPopup();
      window.Telegram.WebApp.BackButton.isVisible = true;
      setScannedObject(event.data);
    }

    window.Telegram.WebApp.onEvent(QR_EVENT, onQRTextReceived);
    return () => {
      window.Telegram.WebApp.offEvent(QR_EVENT, onQRTextReceived);
    }
  }, [])

  const getUI = () => {
    if (!isSupported) {
      return (
        <div className="scanner-wrapper">
          <span>Your Telegram version is not supported.</span>
          <span>Please update to the latest one.</span>
          <div className="button" onClick={onCloseClicked}>Close app</div>
        </div>
      )
    }
    else if (!tokenState.loaded) {
      return (<div className="preloader" />)
    } else if (tokenState.token == null) {
      return (
        <div className="scanner-wrapper">
          <span>It seems you haven't purchased a ticket yet.</span>
          <span>To start our tour, please go back to the bot and buy a ticket.</span>
          <div className="button" onClick={onCloseClicked}>Close app</div>
        </div>
      )
    } else if (scannedObject == null) {
      return (
        <div className="scanner-wrapper">
          <span>Welcome to the tour!</span>
          <span>Scan QR codes to start listening.</span>
          <div className="button" onClick={onScanQRClicked}>Scan QR</div>
        </div>
      )
    } else {
      return <ObjectViewerComponent AccessToken={tokenState.token} ObjectCode={scannedObject} />
    }
  }

  return (
    <div className="App">
      {getUI()}
    </div>
  );
}

export default App;
