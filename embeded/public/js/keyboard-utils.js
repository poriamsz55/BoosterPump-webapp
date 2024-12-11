export const handleEscKey = (callback) => {
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape' || event.key === 'Esc' || event.keyCode === 27) {
            callback();
        }
    });
};