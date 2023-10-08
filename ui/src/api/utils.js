export const isTouchDevice = () => {
    return 'ontouchstart' in window;
}