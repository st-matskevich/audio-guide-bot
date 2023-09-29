import './ObjectViewerComponent.css';
import ReactSlider from "react-slider";
import { ReactComponent as PlayIcon } from './assets/play.svg'
import { ReactComponent as PauseIcon } from './assets/pause.svg'
import { ReactComponent as QRIcon } from './assets/qr-code.svg'
import { useRef, useState, useEffect } from "react"
import { getObjectAudioURL, getObjectCoverURL } from './api/guide';

function ObjectViewerComponent(props) {
    const objectCode = props.ObjectCode;
    const accessToken = props.AccessToken;
    const audioURL = getObjectAudioURL(accessToken, objectCode);
    const coverURL = getObjectCoverURL(accessToken, objectCode);

    const onScanQRClicked = () => {
        window.Telegram.WebApp.showScanQrPopup({});
    };

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
        const updateProgress = () => {
            const progress = ref.currentTime / ref.duration;
            setAudioProgress(progress);
        };

        ref.addEventListener('timeupdate', updateProgress);
        return () => {
            ref.removeEventListener('timeupdate', updateProgress);
        };
    }, []);

    const onSeekAudio = (value) => {
        console.log(audioRef.current.duration + " " + value)
        audioRef.current.currentTime = audioRef.current.duration * value;
        setAudioProgress(value);
    }

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

    return (
        <div className="object-viewer-wrapper">
            <div className="image-viewer">
                <img src={coverURL} alt="cover" />
            </div>
            <div className="object-title">Cat King</div>
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

export default ObjectViewerComponent;