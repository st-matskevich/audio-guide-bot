import "./SliderComponent.css";
import { useState } from "react";

function SliderComponent(props) {
    const {className, value, min, max, step, onChange} = props;

    const handleChange = (e) => {
        const newValue = parseFloat(e.target.value, 10);
        if (onChange) {
            onChange(newValue);
        }
    };

    const [hover, setHover] = useState(false);
    const onTouchStart = (e) => {
        setHover(true)
    }

    const onTouchEnd = (e) => {
        setHover(false)
    }

    const progress = ((value - min) / (max - min)) * 100;

    return (
        <input
            className={`range-input ${className} ${hover ? "active" : ""}`}
            style={{ "--progress": `${progress}%` }}
            type="range"
            min={min}
            max={max}
            step={step}
            value={value}
            onChange={handleChange}
            onTouchStart={onTouchStart}
            onTouchEnd={onTouchEnd}
            onMouseEnter={onTouchStart}
            onMouseLeave={onTouchEnd}
        />
    );
}

export default SliderComponent;