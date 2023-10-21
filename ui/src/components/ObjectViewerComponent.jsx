import "./ObjectViewerComponent.css";
import "./CommonStyles.css";
import PlayIcon from "../assets/play.svg?react";
import PauseIcon from "../assets/pause.svg?react";
import QRIcon from "../assets/qr-code.svg?react";
import { useRef, useState, useEffect } from "react";
import { getObjectAudioURL, getObjectCoverURL, getObjectData } from "../api/guide";
import { i18n } from "../api/i18n";
import SliderComponent from "./SliderComponent";
import ImageComponent from "./ImageComponent";
import CarouselComponent from "./CarouselContainer";
import MarqueeComponent from "./MarqueeComponent";
import { ButtonComponent, IconButtonComponent } from "./ButtonComponents";

function ObjectViewerComponent(props) {
    const { accessToken, objectCode } = props;

    const [objectData, setObjectData] = useState({ loaded: false, data: null, error: null });
    useEffect(() => {
        setObjectData({ loaded: false, data: null, error: null });
        getObjectData(accessToken, objectCode).then((response) => {
            const object = response.data.data;
            object.covers.sort((a, b) => a.index - b.index);
            setObjectData({ loaded: true, code: objectCode, data: object, error: null });
        }).catch((error) => {
            setObjectData({ loaded: true, code: null, data: null, error: error.response.data.data });
        });
    }, [objectCode, accessToken]);

    const audioRef = useRef();
    const audioURL = objectData.loaded ? getObjectAudioURL(accessToken, objectData.code) : null;
    const [audioPlaying, setAudioPlaying] = useState(false);
    const [audioProgress, setAudioProgress] = useState(0);

    useEffect(() => {
        setAudioPlaying(false);
        setAudioProgress(0);
    }, [audioURL]);

    const onToggleAudioPlay = () => {
        const current = audioRef.current;
        if (current) {
            if (audioPlaying) {
                current.pause();
            } else {
                current.play();
            }
        }
    };

    const onAudioTimeUpdate = () => {
        const current = audioRef.current;
        if (current) {
            if (current.duration > 0) {
                const progress = current.currentTime / current.duration;
                setAudioProgress(progress);
            }
        }
    };

    const onAudioPlay = () => {
        setAudioPlaying(true);
    };

    const onAudioPause = () => {
        setAudioPlaying(false);
    };

    const onSeekAudio = (value) => {
        const current = audioRef.current;
        if (current) {
            current.currentTime = current.duration * value;
            setAudioProgress(value);
        }
    };

    const onScanQRClicked = () => {
        window.Telegram.WebApp.showScanQrPopup({});
    };

    const getPlayIcon = () => {
        if (audioPlaying) {
            return (
                <PauseIcon alt={i18n.t("ALT_PAUSE_GUIDE")} />
            );
        } else {
            return (
                <PlayIcon alt={i18n.t("ALT_PLAY_GUIDE")} />
            );
        }
    };

    const getUI = () => {
        if (!objectData.loaded) {
            return (<div className="preloader" />);
        } else if (objectData.error != null) {
            return (
                <div className="object-viewer-wrapper">
                    <span>{i18n.t("ERROR_OBJECT_LOAD_FAILED")}</span>
                    <span>{objectData.error}</span>
                    <ButtonComponent onClick={onScanQRClicked}>{i18n.t("BUTTON_SCAN_QR")}</ButtonComponent>
                </div>
            );
        } else {
            return (
                <div className="object-viewer-wrapper">
                    <CarouselComponent className="image-viewer" alt={{ left: i18n.t("ALT_NAVIGATE_LEFT"), right: i18n.t("ALT_NAVIGATE_RIGHT") }}>
                        {objectData.data.covers.map((cover) => {
                            return (
                                <ImageComponent key={cover.index} src={getObjectCoverURL(accessToken, objectData.code, cover.index)} alt={i18n.t("ALT_OBJECT_COVER")} />
                            );
                        })}
                    </CarouselComponent>
                    <MarqueeComponent className="object-title" string={objectData.data.title} />
                    <audio ref={audioRef} src={audioURL} onTimeUpdate={onAudioTimeUpdate} onPlay={onAudioPlay} onPause={onAudioPause} />
                    <SliderComponent className="audio-range" value={audioProgress} min={0} max={1} step={0.01} onChange={onSeekAudio} />
                    <div className="controls-bar">
                        <IconButtonComponent className="icon-button" onClick={onScanQRClicked}>
                            <QRIcon alt={i18n.t("ALT_SCAN_QR")} />
                        </IconButtonComponent>
                        <IconButtonComponent className="icon-button" onClick={onToggleAudioPlay}>
                            {getPlayIcon()}
                        </IconButtonComponent>
                    </div>
                </div>
            );
        }
    };

    return getUI();
}

export default ObjectViewerComponent;