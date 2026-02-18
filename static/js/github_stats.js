// github-stats.js - Fetch and display GitHub repository stats
(function () {
    'use strict';

    const REPO = 'prinzpiuz/BloTils';
    const CACHE_KEY = 'blotils_github_stats';
    const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

    // Check cache first
    function getCachedStats() {
        try {
            const cached = localStorage.getItem(CACHE_KEY);
            if (cached) {
                const { data, timestamp } = JSON.parse(cached);
                if (Date.now() - timestamp < CACHE_DURATION) {
                    return data;
                }
            }
        } catch (e) {
            // Ignore cache errors
        }
        return null;
    }

    // Save to cache
    function setCachedStats(data) {
        try {
            localStorage.setItem(
                CACHE_KEY,
                JSON.stringify({
                    data,
                    timestamp: Date.now(),
                }),
            );
        } catch (e) {
            // Ignore cache errors
        }
    }

    // Format number (1234 -> 1.2k)
    function formatNumber(num) {
        if (num >= 1000) {
            return (num / 1000).toFixed(1).replace(/\.0$/, '') + 'k';
        }
        return num.toString();
    }

    // Update DOM with stats
    function updateStats(stats) {
        const starsEl = document.getElementById('github-stars');
        const forksEl = document.getElementById('github-forks');
        const versionEl = document.getElementById('github-version');

        if (starsEl) starsEl.textContent = formatNumber(stats.stars);
        if (forksEl) forksEl.textContent = formatNumber(stats.forks);
        if (versionEl) versionEl.textContent = stats.version;

        // Show the stats container
        const container = document.getElementById('github-stats');
        if (container) container.classList.add('loaded');
    }

    // Fetch stats from GitHub API
    async function fetchGitHubStats() {
        // Check cache first
        const cached = getCachedStats();
        if (cached) {
            updateStats(cached);
            return;
        }

        try {
            // Fetch repo info (stars, forks)
            const repoResponse = await fetch(
                `https://api.github.com/repos/${REPO}`,
            );
            if (!repoResponse.ok) throw new Error('Failed to fetch repo');
            const repoData = await repoResponse.json();

            // Fetch latest release (version)
            let version = 'v1.0.0';
            try {
                const releaseResponse = await fetch(
                    `https://api.github.com/repos/${REPO}/releases/latest`,
                );
                if (releaseResponse.ok) {
                    const releaseData = await releaseResponse.json();
                    version = releaseData.tag_name || version;
                }
            } catch (e) {
                // Use default version if no releases
            }

            const stats = {
                stars: repoData.stargazers_count || 0,
                forks: repoData.forks_count || 0,
                version: version,
            };

            // Cache and update
            setCachedStats(stats);
            updateStats(stats);
        } catch (error) {
            console.error('Failed to fetch GitHub stats:', error);
            // Show fallback
            const container = document.getElementById('github-stats');
            if (container) container.classList.add('loaded');
        }
    }

    // Initialize
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', fetchGitHubStats);
    } else {
        fetchGitHubStats();
    }
})();
