import './ObjectViewerComponent.css';
import './CommonStyles.css';
import MarqueeComponent from './MarqueeComponent';
import { ReactComponent as PlayIcon } from '../assets/play.svg';
import { ReactComponent as PauseIcon } from '../assets/pause.svg';
import { ReactComponent as QRIcon } from '../assets/qr-code.svg';
import { useRef, useState, useEffect } from "react";
import { getObjectAudioURL, getObjectCoverURL, getObjectData } from '../api/guide';
import SliderComponent from './SliderComponent';
import RippleContainer from './RippleContainer';
import ImageComponent from './ImageComponent';
import CarouselComponent from './CarouselContainer';

function ObjectViewerComponent(props) {
    const { accessToken, objectCode } = props;

    const [objectData, setObjectData] = useState({ loaded: false, data: null, error: null })
    useEffect(() => {
        setObjectData({ loaded: false, data: null, error: null });
        getObjectData(accessToken, objectCode).then((response) => {
            const object = response.data.data;
            object.covers.sort((a, b) => a.index - b.index);
            setObjectData({ loaded: true, data: object, error: null });
        }).catch((error) => {
            setObjectData({ loaded: true, data: null, error: error.response.data.data });
        })
    }, [objectCode, accessToken]);

    const audioRef = useRef();
    const audioURL = objectData.loaded ? getObjectAudioURL(accessToken, objectCode) : null;
    const [audioPlaying, setAudioPlaying] = useState(false);
    const [audioProgress, setAudioProgress] = useState(0);

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

    const onAudioLoadStarted = () => {
        setAudioPlaying(false);
        setAudioProgress(0);
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
                <PauseIcon
                    width="31"
                    height="31"
                    fill="var(--tg-theme-button-text-color)"
                    stroke="var(--tg-theme-button-text-color)"
                    alt="pause guide"
                />
            )
        } else {
            return (
                <PlayIcon
                    width="31"
                    height="31"
                    fill="var(--tg-theme-button-text-color)"
                    stroke="var(--tg-theme-button-text-color)"
                    alt="play guide"
                />
            )
        }
    }

    const getUI = () => {
        if (!objectData.loaded) {
            return (<div className="preloader" />)
        } else if (objectData.error != null) {
            return (
                <div className="object-viewer-wrapper">
                    <span>{"An error occurred while loading object: "}</span>
                    <span>{objectData.error}</span>
                    <RippleContainer className="button" onClick={onScanQRClicked}>Scan QR</RippleContainer>
                </div>
            )
        } else {
            return (
                <div className="object-viewer-wrapper">
                    <CarouselComponent className="image-viewer">
                        {objectData.data.covers.map((cover) => {
                            return (
                                <ImageComponent key={cover.index} src={getObjectCoverURL(accessToken, objectCode, cover.index)} alt="cover" />
                            );
                        })}
                    </CarouselComponent>
                    <MarqueeComponent className="object-title" string={objectData.data.title} />
                    <audio ref={audioRef} src={audioURL} onTimeUpdate={onAudioTimeUpdate} onLoadStart={onAudioLoadStarted} onPlay={onAudioPlay} onPause={onAudioPause} />
                    <SliderComponent className="audio-range" value={audioProgress} min={0} max={1} step={0.01} onChange={onSeekAudio} />
                    <div className="controls-bar">
                        <RippleContainer className="icon-button" onClick={onScanQRClicked}>
                            <QRIcon
                                width="31"
                                height="31"
                                fill="var(--tg-theme-button-text-color)"
                                stroke="var(--tg-theme-button-text-color)"
                                alt="scan qr code"
                            />
                        </RippleContainer>
                        <RippleContainer className="icon-button" onClick={onToggleAudioPlay}>
                            {getPlayIcon()}
                        </RippleContainer>
                    </div>
                </div>
            )
        }
    }

    return getUI();
}

export default ObjectViewerComponent;