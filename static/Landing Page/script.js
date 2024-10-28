// Add smooth scroll for anchor links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth'
            });
        }
    });
});

// Rotating headlines functionality
const headlines = [
    {
        main: 'Aquarc the tool',
        highlight: 'to make ',
        end: 'highschool easier.'
    },
    {
        main: 'SAT studying made',
        highlight: 'easier ',
        end: 'with aquarc.'
    },
    {
        main: 'Find the best',
        highlight: 'extracurriculars ',

    }
];

let currentIndex = 0;
const container = document.querySelector('.headline-container');
const mainText = document.querySelector('.headline-text');
const highlightText = document.querySelector('.highlight-text');
const endText = document.querySelector('.end-text');

function updateText() {
    const current = headlines[currentIndex];
    mainText.textContent = current.main;
    highlightText.textContent = current.highlight;
    endText.textContent = current.end;
}

function rotateHeadlines() {
    // Add fade out class
    container.classList.add('fade');
    
    // After fade out, update text and fade in
    setTimeout(() => {
        currentIndex = (currentIndex + 1) % headlines.length;
        updateText();
        container.classList.remove('fade');
    }, 500);
}

// Start the rotation
setInterval(rotateHeadlines, 6000);

// Add this to your script.js file
document.addEventListener('DOMContentLoaded', function() {
    // Get all navigation buttons
    const startLearningBtn = document.querySelector('.cta-button');
    const getStartedBtn = document.querySelector('.submit-button');

    // Add click event listeners
    const navigateToSAT = (e) => {
        e.preventDefault();
        window.location.href = '/static/sat/index.html';
    };

    if (startLearningBtn) {
        startLearningBtn.addEventListener('click', navigateToSAT);
    }

    if (getStartedBtn) {
        getStartedBtn.addEventListener('click', navigateToSAT);
    }
});