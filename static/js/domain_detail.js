document.addEventListener('DOMContentLoaded', () => {
    const chartInstances = new Map();

    document.querySelectorAll('.graph-toggle-btn').forEach((btn) => {
        btn.addEventListener('click', function () {
            const uri = this.dataset.uri;
            const domainId = this.dataset.domainId;
            const graphRow = this.closest('tr').nextElementSibling;
            const canvas = graphRow.querySelector('.likes-chart');

            if (graphRow.style.display === 'none') {
                graphRow.style.display = 'table-row';
                this.textContent = 'Hide Graph';
                this.classList.add('active');

                // Only fetch and render if not already done
                if (!chartInstances.has(canvas)) {
                    fetchAndRenderChart(canvas, domainId, uri);
                }
            } else {
                graphRow.style.display = 'none';
                this.textContent = 'View Graph';
                this.classList.remove('active');
            }
        });
    });

    function fetchAndRenderChart(canvas, domainId, uri) {
        const url = `/domain/${domainId}/likes_timeline?uri=${encodeURIComponent(uri)}`;

        fetch(url)
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then((data) => {
                renderChart(canvas, data, uri);
            })
            .catch((error) => {
                console.error('Error fetching timeline:', error);
                canvas.parentElement.innerHTML =
                    '<p class="graph-error">Failed to load graph data.</p>';
            });
    }

    function renderChart(canvas, data, uri) {
        if (!data || data.length === 0) {
            canvas.parentElement.innerHTML =
                '<p class="graph-no-data">No timeline data available for this post.</p>';
            return;
        }

        const labels = data.map((d) => d.date);
        const counts = data.map((d) => d.count);

        // Calculate cumulative likes
        const cumulative = [];
        counts.reduce((acc, val, i) => {
            cumulative[i] = acc + val;
            return cumulative[i];
        }, 0);

        const ctx = canvas.getContext('2d');
        const chart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Likes per day',
                        data: counts,
                        borderColor: '#742ce0',
                        backgroundColor: 'rgba(116, 44, 224, 0.1)',
                        borderWidth: 2,
                        fill: true,
                        tension: 0.3,
                        pointBackgroundColor: '#742ce0',
                        pointBorderColor: '#fff',
                        pointBorderWidth: 2,
                        pointRadius: 4,
                        pointHoverRadius: 6,
                    },
                    {
                        label: 'Cumulative likes',
                        data: cumulative,
                        borderColor: '#2ecc71',
                        backgroundColor: 'rgba(46, 204, 113, 0.05)',
                        borderWidth: 2,
                        borderDash: [5, 5],
                        fill: false,
                        tension: 0.3,
                        pointBackgroundColor: '#2ecc71',
                        pointBorderColor: '#fff',
                        pointBorderWidth: 2,
                        pointRadius: 3,
                        pointHoverRadius: 5,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'top',
                        labels: {
                            usePointStyle: true,
                            padding: 20,
                            font: {
                                family: "'Inter', sans-serif",
                                size: 12,
                            },
                        },
                    },
                    tooltip: {
                        backgroundColor: 'rgba(0, 0, 0, 0.8)',
                        titleFont: {
                            family: "'Inter', sans-serif",
                        },
                        bodyFont: {
                            family: "'Inter', sans-serif",
                        },
                        cornerRadius: 8,
                        padding: 12,
                    },
                },
                scales: {
                    x: {
                        grid: {
                            display: false,
                        },
                        ticks: {
                            font: {
                                family: "'Inter', sans-serif",
                                size: 11,
                            },
                            maxRotation: 45,
                        },
                    },
                    y: {
                        beginAtZero: true,
                        ticks: {
                            stepSize: 1,
                            font: {
                                family: "'Inter', sans-serif",
                                size: 11,
                            },
                        },
                        grid: {
                            color: 'rgba(0, 0, 0, 0.06)',
                        },
                    },
                },
                interaction: {
                    intersect: false,
                    mode: 'index',
                },
            },
        });

        chartInstances.set(canvas, chart);
    }
});
