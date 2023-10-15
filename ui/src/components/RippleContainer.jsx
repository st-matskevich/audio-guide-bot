import { isTouchDevice } from "../api/utils";
import "./RippleContainer.css";
import { useState, useEffect, useCallback } from "react";

function RippleContainer(props) {
    const { className, children, onClick } = props;
    const [rippleArray, setRippleArray] = useState([]);
    const [cleanupTimer, setCleanupTimer] = useState(null);
    const rippleDuration = 500;

    const addRipple = (event) => {
        const container = event.currentTarget.getBoundingClientRect();
        const size = Math.max(container.width, container.height);

        const x = event.pageX - container.x - size / 2;
        const y = event.pageY - container.y - size / 2;
        const ripple = { x, y, size };

        setRippleArray([...rippleArray, ripple]);
    };

    const onPointerDown = (event) => {
        clearTimeout(cleanupTimer);
        addRipple(event);
    };

    const onPointerUp = useCallback(() => {
        if(rippleArray.length > 0) {
            setRippleArray((array) => array.map((ripple) => ({ ...ripple, inactive: true })));
            setCleanupTimer(setTimeout(() => { setRippleArray([]); }, rippleDuration * 4));
        }
    }, [rippleArray]);

    useEffect(() => {
        if(!isTouchDevice())
        {
            window.addEventListener("mouseup", onPointerUp);
            return () => {
                window.removeEventListener("mouseup", onPointerUp);
            };
        }
    }, [onPointerUp]);

    return (
        <div className={`ripple-container ${className}`} onClick={onClick} onPointerDown={onPointerDown} onTouchEnd={onPointerUp} onTouchCancel={onPointerUp}>
            {rippleArray.length > 0 &&
                rippleArray.map((ripple, index) => {
                    return (
                        <span
                            key={"span" + index}
                            className={ripple.inactive ? "inactive" : ""}
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
}

export default RippleContainer;
