import "./RippleContainer.css"
import React, { useState, useEffect } from "react";

const useDebouncedCleanUp = (rippleCount, duration, onClean) => {
    useEffect(() => {
        let bounce = null;
        if (rippleCount > 0) {
            clearTimeout(bounce);

            bounce = setTimeout(() => {
                onClean();
                clearTimeout(bounce);
            }, duration * 4);
        }

        return () => clearTimeout(bounce);
    }, [rippleCount, duration, onClean]);
};

function RippleContainer(props) {
    const { className, children, onClick } = props;
    const [rippleArray, setRippleArray] = useState([]);
    const rippleDuration = 1000;

    useDebouncedCleanUp(rippleArray.length, rippleDuration, () => {
        setRippleArray([]);
    });

    const addRipple = (event) => {
        const container = event.currentTarget.getBoundingClientRect();
        const size = Math.max(container.width, container.height);

        const x = event.pageX - container.x - size / 2;
        const y = event.pageY - container.y - size / 2;
        const ripple = { x, y, size };

        setRippleArray([...rippleArray, ripple]);
    };

    const handleClick = (event) => {
        addRipple(event);
        if (onClick) {
            onClick();
        }
    }

    return (
        <div className={`ripple-container ${className}`} onClick={handleClick}>
            {rippleArray.length > 0 &&
                rippleArray.map((ripple, index) => {
                    return (
                        <span
                            key={"span" + index}
                            style={{
                                top: ripple.y,
                                left: ripple.x,
                                width: ripple.size,
                                height: ripple.size,
                                "--duration": `${rippleDuration}ms`
                            }}
                        />
                    );
                })}
            {children}
        </div>
    );
};

export default RippleContainer;
