import "./SliderComponent.css";

function SliderComponent(props) {
    const {className, value, min, max, step, onChange} = props;


    const handleChange = (e) => {
        const newValue = parseFloat(e.target.value, 10);
        if (onChange) {
            onChange(newValue);
        }
    };

    const progress = ((value - min) / (max - min)) * 100;

    return (
        <input
            className={`range-input ${className}`}
            style={{ "--progress": `${progress}%` }}
            type="range"
            min={min}
            max={max}
            step={step}
            value={value}
            onChange={handleChange}
        />
    );
}

export default SliderComponent;