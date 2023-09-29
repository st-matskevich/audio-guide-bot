import './ObjectViewerComponent.css';
import ReactSlider from "react-slider";
import { ReactComponent as PlayIcon } from './assets/play.svg'
import { ReactComponent as PauseIcon } from './assets/pause.svg'
import { ReactComponent as QRIcon } from './assets/qr-code.svg'
import { useRef, useState, useEffect } from "react"
import { getObjectAudioURL, getObjectCoverURL, getObjectData } from './api/guide';

function ObjectViewerComponent(props) {
    const objectCode = props.ObjectCode;
    const accessToken = props.AccessToken;

    const [objectData, setObjectData] = useState({ loaded: false, data: null })
    useEffect(() => {
        setObjectData({ loaded: false, data: null });
        getObjectData(accessToken, objectCode).then((response) => {
            setObjectData({ loaded: true, data: response.data.data });
        })
    }, [objectCode, accessToken])

    const audioURL = objectData.loaded ? getObjectAudioURL(accessToken, objectCode) : null;
    const coverURL = objectData.loaded ? getObjectCoverURL(accessToken, objectCode) : null;

    const audioRef = useRef(new Audio());
    useEffect(() => {
        const ref = audioRef.current;
        ref.src = audioURL;

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
        } else {
            return (
                <div className="object-viewer-wrapper">
                    <div className="image-viewer">
                        <img src={coverURL} alt="cover" />
                    </div>
                    <div className="object-title">{objectData.data.title}</div>
                    <ReactSlider className="audio-range" value={audioProgress} min={0} max={1} step={0.01} onChange={onSeekAudio} />
                    <div className="controls-bar">
                        <div className="icon-button" onClick={onScanQRClicked}>
                            <QRIcon
                                width="31"
                                height="31"
                                fill="var(--tg-theme-button-text-color)"
                                stroke="var(--tg-theme-button-text-color)"
                                alt="scan qr code"
                            />
                        </div>
                        <div className="icon-button" onClick={toggleAudioPlay}>
                            {getPlayIcon()}
                        </div>
                    </div>
                </div>
            )
        }
    }

    return getUI();
}

export default ObjectViewerComponent;