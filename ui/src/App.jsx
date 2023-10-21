import "./App.css";
import "./components/CommonStyles.css";
import { useEffect, useState } from "react";
import { addTokenListener, removeTokenListener } from "./api/auth";
import ObjectViewerComponent from "./components/ObjectViewerComponent";
import { isTelegramAPISupported } from "./api/telegram";
import { i18n } from "./api/i18n";
import { ButtonComponent } from "./components/ButtonComponents";

function App() {
    const isSupported = isTelegramAPISupported();
    const onScanQRClicked = () => {
        window.Telegram.WebApp.showScanQrPopup({});
    };

    const onCloseClicked = () => {
        window.Telegram.WebApp.close();
    };

    const [scannedObject, setScannedObject] = useState(null);
    const [tokenState, setTokenState] = useState({
        loaded: false,
        token: null
    });

    useEffect(() => {
        const onTokenChanged = (event) => {
            setTokenState({
                loaded: true,
                token: event.detail
            });
        };

        addTokenListener(onTokenChanged);
        return () => {
            removeTokenListener(onTokenChanged);
        };
    }, []);

    useEffect(() => {
        const onBackClicked = () => {
            window.Telegram.WebApp.BackButton.isVisible = false;
            setScannedObject(null);
        };

        window.Telegram.WebApp.BackButton.onClick(onBackClicked);
        return () => {
            window.Telegram.WebApp.BackButton.offClick(onBackClicked);
        };
    }, []);

    useEffect(() => {
        const QR_EVENT = "qrTextReceived";
        const onQRTextReceived = (event) => {
            window.Telegram.WebApp.closeScanQrPopup();
            window.Telegram.WebApp.BackButton.isVisible = true;
            setScannedObject(event.data);
        };

        window.Telegram.WebApp.onEvent(QR_EVENT, onQRTextReceived);
        return () => {
            window.Telegram.WebApp.offEvent(QR_EVENT, onQRTextReceived);
        };
    }, []);

    const getUI = () => {
        if (!isSupported) {
            return (
                <div className="scanner-wrapper">
                    <span>{i18n.t("MESSAGE_NOT_SUPPORTED_LINE_1")}</span>
                    <span>{i18n.t("MESSAGE_NOT_SUPPORTED_LINE_2")}</span>
                    <ButtonComponent onClick={onCloseClicked}>{i18n.t("BUTTON_CLOSE_APP")}</ButtonComponent>
                </div>
            );
        }
        else if (!tokenState.loaded) {
            return (<div className="preloader" />);
        } else if (tokenState.token == null) {
            return (
                <div className="scanner-wrapper">
                    <span>{i18n.t("MESSAGE_NO_TICKET_LINE_1")}</span>
                    <span>{i18n.t("MESSAGE_NO_TICKET_LINE_2")}</span>
                    <ButtonComponent onClick={onCloseClicked}>{i18n.t("BUTTON_CLOSE_APP")}</ButtonComponent>
                </div>
            );
        } else if (scannedObject == null) {
            return (
                <div className="scanner-wrapper">
                    <span>{i18n.t("MESSAGE_WELCOME_LINE_1")}</span>
                    <span>{i18n.t("MESSAGE_WELCOME_LINE_2")}</span>
                    <ButtonComponent onClick={onScanQRClicked}>{i18n.t("BUTTON_SCAN_QR")}</ButtonComponent>
                </div>
            );
        } else {
            return <ObjectViewerComponent accessToken={tokenState.token} objectCode={scannedObject} />;
        }
    };

    return (
        <div className="App">
            {getUI()}
        </div>
    );
}

export default App;
