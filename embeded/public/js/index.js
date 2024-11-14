// Create animated background
function createBackground() {
    const background = document.getElementById('background');
    for (let i = 0; i < 50; i++) {
        const span = document.createElement('span');
        span.style.width = Math.random() * 30 + 'px';
        span.style.height = span.style.width;
        span.style.left = Math.random() * 100 + '%';
        span.style.animationDelay = Math.random() * 5 + 's';
        background.appendChild(span);
    }
}

// Handle splash screen and main content transition
window.onload = () => {
    createBackground();
    
    const splashScreen = document.getElementById('splashScreen');
    const mainContent = document.getElementById('mainContent');
    
    // Show splash screen for 2.5 seconds
    setTimeout(() => {
        splashScreen.style.opacity = '0';
        mainContent.style.display = 'block';
        
        // Small delay to ensure smooth transition
        setTimeout(() => {
            mainContent.style.opacity = '1';
            splashScreen.style.display = 'none';
        }, 500);
    }, 2500);
};

// Add click effect to buttons
document.querySelectorAll('.glass-button').forEach(button => {
    button.addEventListener('click', function(e) {
        let ripple = document.createElement('span');
        ripple.style.position = 'absolute';
        ripple.style.background = '#fff';
        ripple.style.transform = 'translate(-50%, -50%)';
        ripple.style.pointerEvents = 'none';
        ripple.style.borderRadius = '50%';
        ripple.style.animation = 'ripple 0.6s linear';
        this.appendChild(ripple);
        setTimeout(() => ripple.remove(), 600);
    });
});

document.getElementById("btnProjectPage").addEventListener("click", function() {
    window.location.href = "/projects";
});

document.getElementById("btnDevicePage").addEventListener("click", function() {
    window.location.href = "/devices.html";
});

document.getElementById("btnPartPage").addEventListener("click", function() {
    window.location.href = "./parts.html";
});