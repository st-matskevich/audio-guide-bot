import "./CarouselContainer.css";
import { Children, useRef, useState } from "react";
import LeftIcon from "../assets/left.svg?react";
import RightIcon from "../assets/right.svg?react";

function CarouselContainer(props) {
    const { className, children, alt } = props;
    const [activePage, setActivePage] = useState(0);
    const containerRef = useRef();

    const pageCount = Children.count(children);
    const canNavigateRight = activePage < pageCount - 1;
    const canNavigateLeft = activePage > 0;

    const onScroll = (e) => {
        const pageWidth = e.target.scrollWidth / pageCount;
        const page = Math.round(e.target.scrollLeft / pageWidth);
        setActivePage(page);
    };

    const onNavigateRight = () => {
        const current = containerRef.current;
        if (current) {
            const pageWidth = current.scrollWidth / pageCount;
            current.scrollTo({ left: pageWidth * (activePage + 1), behavior: "smooth" });
        }
    };

    const onNavigateLeft = () => {
        const current = containerRef.current;
        if (current) {
            const pageWidth = current.scrollWidth / pageCount;
            current.scrollTo({ left: pageWidth * (activePage - 1), behavior: "smooth" });
        }
    };

    const getPaginator = () => {
        if (pageCount > 1) {
            return (
                <div className='carousel-paginator'>
                    {Children.map(children, (_, index) =>
                        <div className={`carousel-page-icon ${index === activePage ? "active" : ""}`} />
                    )}
                </div>
            );
        }

        return null;
    };

    return (
        <div className={`carousel-wrapper ${className}`}>
            <div className="carousel-container" onScroll={onScroll} ref={containerRef}>
                {Children.map(children, (child) => (
                    <div className="carousel-page">
                        {child}
                    </div>
                ))}
            </div>
            {getPaginator()}
            {canNavigateLeft && (
                <button className='carousel-navigator left' onClick={onNavigateLeft}>
                    <LeftIcon
                        width="20"
                        height="20"
                        alt={alt.left}
                    />
                </button>
            )}
            {canNavigateRight && (
                <button className='carousel-navigator right' onClick={onNavigateRight}>
                    <RightIcon
                        width="20"
                        height="20"
                        alt={alt.right}
                    />
                </button>
            )}
        </div>
    );
}

export default CarouselContainer;