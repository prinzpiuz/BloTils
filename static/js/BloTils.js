// BloTils: https://www.blotils.com
// This file is released under the ISC license: https://opensource.org/licenses/ISC

(function () {
    'use strict';
    if (window.blotils && window.blotils.initialized) return;

    window.blotils = window.blotils || {};
    window.blotils.initialized = true;

    // Configuration
    const CONFIG = {
        countUrl: 'api/v1/count_like',
        debounceDelay: 300,
        animationDuration: 1000,
        floatingHeartsCount: 5,
    };

    // State
    const state = {
        isLiked: false,
        count: 0,
        isLoading: false,
        hasError: false,
    };

    // DOM Elements (initialized later)
    let elements = {
        button: null,
        countDisplay: null,
        heartIcon: null,
    };

    // Get configuration from script tag
    function getScriptConfig() {
        const scriptTag = document.querySelector('script[data-blotils_url]');
        const btnScript = document.querySelector(
            'script[data-blotils_like_btn]',
        );

        return {
            baseUrl: scriptTag?.dataset.blotils_url || 'https://blotils.com',
            buttonId: btnScript?.dataset.blotils_like_btn || 'blotils_like_btn',
            countId: btnScript?.dataset.blotils_count || 'blotils_count',
        };
    }

    const config = getScriptConfig();
    window.blotils.base_url = config.baseUrl;

    // URL Builder
    function buildUrl(uri, params = {}) {
        const url = new URL(uri, window.blotils.base_url);
        Object.entries(params).forEach(([key, value]) => {
            url.searchParams.append(key, value);
        });
        return url;
    }

    // Debounce utility
    function debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    // Initialize DOM elements
    function initElements() {
        elements.button = document.getElementById(config.buttonId);

        if (!elements.button) {
            console.warn(
                'BloTils: Like button not found with id:',
                config.buttonId,
            );
            return false;
        }

        // Find or create count display
        elements.countDisplay =
            document.getElementById(config.countId) ||
            elements.button.parentElement?.querySelector(
                '.blotils-count, .love-count, [data-blotils-count]',
            );

        // Find heart icon if exists
        elements.heartIcon = elements.button.querySelector(
            '.heart-icon, svg, .blotils-icon',
        );

        return true;
    }

    // Update UI based on state
    function updateUI() {
        if (!elements.button) return;

        // Update button state
        elements.button.classList.toggle('liked', state.isLiked);
        elements.button.classList.toggle('loading', state.isLoading);
        elements.button.classList.toggle('error', state.hasError);
        elements.button.disabled = state.isLoading;

        // Update count display
        if (elements.countDisplay) {
            elements.countDisplay.textContent = state.count;

            // Trigger count animation
            elements.countDisplay.classList.remove('bump');
            void elements.countDisplay.offsetWidth; // Force reflow
            elements.countDisplay.classList.add('bump');
        }

        // Update aria attributes for accessibility
        elements.button.setAttribute('aria-pressed', state.isLiked);
        elements.button.setAttribute(
            'aria-label',
            state.isLiked
                ? `Liked. ${state.count} likes`
                : `Like this. ${state.count} likes`,
        );
    }

    // Create floating hearts animation
    function createFloatingHearts() {
        if (!elements.button) return;

        const hearts = ['❤️', '💕', '💗', '💖', '💓'];

        for (let i = 0; i < CONFIG.floatingHeartsCount; i++) {
            setTimeout(() => {
                const heart = document.createElement('span');
                heart.className = 'blotils-floating-heart';
                heart.textContent =
                    hearts[Math.floor(Math.random() * hearts.length)];
                heart.style.cssText = `
                    position: absolute;
                    pointer-events: none;
                    font-size: 20px;
                    left: ${15 + Math.random() * 25}px;
                    top: ${5 + Math.random() * 15}px;
                    animation: blotilsFloatUp ${CONFIG.animationDuration}ms ease-out forwards;
                `;

                elements.button.appendChild(heart);

                setTimeout(() => heart.remove(), CONFIG.animationDuration);
            }, i * 80);
        }
    }

    // Trigger heart pop animation
    function triggerHeartPop() {
        if (!elements.heartIcon) return;

        elements.heartIcon.style.animation = 'none';
        void elements.heartIcon.offsetWidth; // Force reflow
        elements.heartIcon.style.animation = 'blotilsHeartPop 0.4s ease';
    }

    // API: Get current like count
    window.blotils.get_likes = async function () {
        state.isLoading = true;
        state.hasError = false;
        updateUI();

        try {
            const response = await fetch(
                buildUrl(CONFIG.countUrl, {
                    page: window.location.pathname,
                }),
            );

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();

            state.count = data.count || 0;
            state.isLiked = false; // Will be determined by server in future
            state.hasError = false;

            return data;
        } catch (error) {
            console.error('BloTils: Error fetching likes:', error);
            state.hasError = true;
            throw error;
        } finally {
            state.isLoading = false;
            updateUI();
        }
    };

    // API: Submit a like
    window.blotils.count_like = async function () {
        // Prevent duplicate requests
        if (state.isLoading) return;

        state.isLoading = true;
        state.hasError = false;
        updateUI();

        try {
            const response = await fetch(buildUrl(CONFIG.countUrl), {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    page: window.location.pathname,
                }),
            });

            const data = await response.json();

            if (data.success) {
                state.isLiked = true;
                state.count = data.count || state.count + 1;

                // Trigger animations
                triggerHeartPop();
                createFloatingHearts();
            } else {
                // Handle "already liked" or other non-success responses
                if (data.message === 'Clap Already Counted') {
                    state.isLiked = true;
                    state.count = data.count || state.count;
                }
            }

            return data;
        } catch (error) {
            console.error('BloTils: Error submitting like:', error);
            state.hasError = true;
            throw error;
        } finally {
            state.isLoading = false;
            updateUI();
        }
    };

    // Debounced click handler
    const handleClick = debounce(async () => {
        await window.blotils.count_like();
    }, CONFIG.debounceDelay);

    // Inject required CSS animations
    function injectStyles() {
        if (document.getElementById('blotils-styles')) return;

        const styles = document.createElement('style');
        styles.id = 'blotils-styles';
        styles.textContent = `
            @keyframes blotilsHeartPop {
                0% { transform: scale(1); }
                25% { transform: scale(1.3); }
                50% { transform: scale(0.9); }
                75% { transform: scale(1.1); }
                100% { transform: scale(1); }
            }

            @keyframes blotilsFloatUp {
                0% {
                    opacity: 1;
                    transform: translateY(0) scale(1);
                }
                100% {
                    opacity: 0;
                    transform: translateY(-80px) scale(0.5);
                }
            }

            @keyframes blotilsCountBump {
                0% { transform: scale(1); }
                50% { transform: scale(1.3); color: #e53935; }
                100% { transform: scale(1); }
            }

            .blotils-count.bump,
            .love-count.bump,
            [data-blotils-count].bump {
                animation: blotilsCountBump 0.3s ease;
            }

            [id*="blotils"].loading {
                opacity: 0.7;
                pointer-events: none;
            }

            [id*="blotils"].error {
                animation: blotilsShake 0.4s ease;
            }

            @keyframes blotilsShake {
                0%, 100% { transform: translateX(0); }
                25% { transform: translateX(-5px); }
                75% { transform: translateX(5px); }
            }
        `;
        document.head.appendChild(styles);
    }

    // Initialize when DOM is ready
    function init() {
        injectStyles();

        if (!initElements()) {
            // Retry once after a short delay (for async loaded content)
            setTimeout(() => {
                if (initElements()) {
                    setupEventListeners();
                    window.blotils.get_likes();
                }
            }, 100);
            return;
        }

        setupEventListeners();
        window.blotils.get_likes();
    }

    function setupEventListeners() {
        if (!elements.button) return;

        // Remove any existing listeners (in case of re-init)
        elements.button.removeEventListener('click', handleClick);
        elements.button.addEventListener('click', handleClick);

        // Store reference for external access
        window.blotils.like_button = elements.button;
    }

    // Start initialization
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // Expose for external use
    window.blotils.refresh = () => window.blotils.get_likes();
    window.blotils.getState = () => ({ ...state });
})();
