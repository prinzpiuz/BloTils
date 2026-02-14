// Toast notification system
(function () {
    'use strict';

    const TOAST_DURATION = 5000; // 5 seconds
    const ICONS = {
        success: '✓',
        error: '✕',
        warning: '⚠',
        info: 'ℹ',
    };

    /**
     * Show a toast notification
     * @param {string} type - 'success', 'error', 'warning', 'info'
     * @param {string} message - The message to display
     * @param {number} duration - Auto-dismiss time in ms (0 = no auto-dismiss)
     */
    window.showToast = function (type, message, duration = TOAST_DURATION) {
        const container = document.getElementById('toast-container');
        const template = document.getElementById('toast-template');

        if (!container || !template) {
            console.error('Toast container or template not found');
            return;
        }

        // Clone the template
        const toast = template.content.cloneNode(true).querySelector('.toast');

        // Set type and content
        toast.classList.add(type);
        toast.querySelector('.toast-icon').textContent =
            ICONS[type] || ICONS.info;
        toast.querySelector('.toast-message').textContent = message;

        // Close button handler
        const closeBtn = toast.querySelector('.toast-close');
        closeBtn.addEventListener('click', function () {
            dismissToast(toast);
        });

        // Add to container
        container.appendChild(toast);

        // Auto-dismiss
        if (duration > 0) {
            setTimeout(function () {
                dismissToast(toast);
            }, duration);
        }

        return toast;
    };

    /**
     * Dismiss a toast with animation
     * @param {HTMLElement} toast - The toast element to dismiss
     */
    function dismissToast(toast) {
        if (!toast || toast.classList.contains('hiding')) return;

        toast.classList.add('hiding');

        // Remove after animation completes
        toast.addEventListener('animationend', function () {
            toast.remove();
        });
    }

    // Convenience functions
    window.showSuccess = function (message, duration) {
        return showToast('success', message, duration);
    };

    window.showError = function (message, duration) {
        return showToast('error', message, duration);
    };

    window.showWarning = function (message, duration) {
        return showToast('warning', message, duration);
    };

    window.showInfo = function (message, duration) {
        return showToast('info', message, duration);
    };
})();
