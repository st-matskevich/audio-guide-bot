import "./MarqueeComponent.css";
import { useRef, useState, useEffect } from "react";

function MarqueeComponent(props) {
    const { className, string } = props;

    const containerRef = useRef();
    const [isOverflow, setIsOverflow] = useState(false);

    useEffect(() => {
        const current = containerRef.current;
        if (current) {
            const hasOverflow = current.scrollWidth > current.clientWidth;
            setIsOverflow(hasOverflow);
        }
    }, [containerRef]);

    return (
        <div className={`marquee-container ${className}`} ref={containerRef}>
            <div className={`marquee-content ${isOverflow ? "" : "inactive"}`}>
                {string}
            </div>
        </div>
    );
}

export default MarqueeComponent;