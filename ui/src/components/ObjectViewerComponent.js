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

function ObjectViewerComponent(props) {
    const { accessToken, objectCode } = props;

    const [objectData, setObjectData] = useState({ loaded: false, data: null, error: null })
    useEffect(() => {
        setObjectData({ loaded: false, data: null, error: null });
        getObjectData(accessToken, objectCode).then((response) => {
            setObjectData({ loaded: true, data: response.data.data, error: null });
        }).catch((error) => {
            setObjectData({ loaded: true, data: null, error: error.response.data.data });
        })
    }, [objectCode, accessToken])

    const audioURL = objectData.loaded ? getObjectAudioURL(accessToken, objectCode) : null;
    const coverURL = objectData.loaded ? getObjectCoverURL(accessToken, objectCode) : null;

    const audioRef = useRef(new Audio());
    useEffect(() => {
        const ref = audioRef.current;
        ref.src = audioURL;
        setAudioPlaying(false);
        setAudioProgress(0);

        return () => {
            ref.pause();
        };
    }, [audioURL]);

    const [audioPlaying, setAudioPlaying] = useState(false);
    const toggleAudioPlay = () => {
        if (audioPlaying) {
            audioRef.current.pause();
        } else {
            audioRef.current.play();
        }
        setAudioPlaying(!audioPlaying);
    }

    const [audioProgress, setAudioProgress] = useState(0);
    useEffect(() => {
        const ref = audioRef.current;
        const onTimeUpdate = () => {
            if (audioRef.current.duration > 0) {
                const progress = ref.currentTime / ref.duration;
                setAudioProgress(progress);
            }
        };

        const onAudioEnded = () => {
            setAudioPlaying(false);
        }

        ref.addEventListener('timeupdate', onTimeUpdate);
        ref.addEventListener('ended', onAudioEnded);
        return () => {
            ref.removeEventListener('timeupdate', onTimeUpdate);
            ref.removeEventListener('ended', onAudioEnded);
        };
    }, []);

    const onSeekAudio = (value) => {
        audioRef.current.currentTime = audioRef.current.duration * value;
        setAudioProgress(value);
    }

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
                    <div className="image-viewer">
                        <ImageComponent src={coverURL} alt="cover" />
                    </div>
                    <MarqueeComponent className="object-title" string={objectData.data.title} />
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
                        <RippleContainer className="icon-button" onClick={toggleAudioPlay}>
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