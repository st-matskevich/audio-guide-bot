import { useEffect, useState } from "react";

function ImageComponent(props) {
    const { src, alt } = props;
    const [state, setState] = useState(null);

    useEffect(() => {
        const image = new Image();
        image.src = src;

        const onImageLoaded = () => {
            setState(src);
        }

        image.addEventListener('load', onImageLoaded);
        return () => {
            image.removeEventListener('load', onImageLoaded);
        }
    }, [src]);

    if (state != null) {
        return (
            <img src={state} alt={alt} draggable="false" />
        )
    }

    return null;
}

export default ImageComponent;