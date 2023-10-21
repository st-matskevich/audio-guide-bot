import "./ButtonComponents.css";
import RippleContainer from "./RippleContainer";

function ButtonComponent(props) {
    const { children, onClick } = props;

    return (
        <button className="button-wrapper" onClick={onClick}>
            <RippleContainer className="button-content">
                {children}
            </RippleContainer>
        </button>
    );
}

function IconButtonComponent(props) {
    const { children, onClick } = props;

    return (
        <button className="icon button-wrapper" onClick={onClick}>
            <RippleContainer className="icon button-content">
                {children}
            </RippleContainer>
        </button>
    );
}

export { ButtonComponent, IconButtonComponent };