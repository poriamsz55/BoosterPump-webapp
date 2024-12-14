particlesJS('particles', {
    particles: {
        number: {
            value: 80,
            density: {
                enable: true,
                value_area: 800
            }
        },
        color: {
            value: ["#90cdf4", "#63b3ed", "#4299e1"]  // Different shades of blue
        },
        opacity: {
            value: 0.3,
            random: true
        },
        size: {
            value: 2,
            random: true
        },
        line_linked: {
            enable: true,
            distance: 150,
            color: "#4299e1",
            opacity: 0.2,
            width: 1
        },
        move: {
            enable: true,
            speed: 1.5,
            direction: "none",
            random: true,
            straight: false,
            out_mode: "out",
            bounce: false
        }
    },
    interactivity: {
        detect_on: "canvas",
        events: {
            onhover: {
                enable: true,
                mode: "repulse"
            },
            resize: true
        }
    }
});