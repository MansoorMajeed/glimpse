// Global variables for managing charts and data
window.agentData = window.agentData || {};
window.charts = window.charts || {};

function createOrUpdateChart(hostname, metrics) {
    const chartId = 'chart-' + hostname;
    const canvas = document.getElementById(chartId);
    
    if (!canvas) return;
    
    // Initialize agent data if it doesn't exist
    if (!window.agentData[hostname]) {
        window.agentData[hostname] = {
            timestamps: [],
            cpu: [],
            memory: [],
            disk: []
        };
    }
    
    const data = window.agentData[hostname];
    const now = new Date();
    
    // Add current data point
    data.timestamps.push(now);
    data.cpu.push(metrics.cpu);
    data.memory.push(metrics.memory);
    data.disk.push(metrics.disk);
    
    // Keep only last 15 data points for performance
    const maxPoints = 15;
    if (data.timestamps.length > maxPoints) {
        data.timestamps.shift();
        data.cpu.shift();
        data.memory.shift();
        data.disk.shift();
    }
    
    // Create chart if it doesn't exist, otherwise update data
    if (!window.charts[hostname]) {
        createChart(hostname, canvas, data);
    } else {
        updateChart(hostname, data);
    }
}

function createChart(hostname, canvas, data) {
    const ctx = canvas.getContext('2d');
    
    window.charts[hostname] = new Chart(ctx, {
        type: 'line',
        data: {
            labels: data.timestamps.map(() => ''),
            datasets: [
                {
                    label: 'CPU',
                    data: data.cpu,
                    borderColor: '#f56565',
                    backgroundColor: 'rgba(245, 101, 101, 0.1)',
                    borderWidth: 1.5,
                    fill: false,
                    tension: 0.3,
                    pointRadius: 0,
                    pointHoverRadius: 3
                },
                {
                    label: 'Memory',
                    data: data.memory,
                    borderColor: '#4299e1',
                    backgroundColor: 'rgba(66, 153, 225, 0.1)',
                    borderWidth: 1.5,
                    fill: false,
                    tension: 0.3,
                    pointRadius: 0,
                    pointHoverRadius: 3
                },
                {
                    label: 'Disk',
                    data: data.disk,
                    borderColor: '#48bb78',
                    backgroundColor: 'rgba(72, 187, 120, 0.1)',
                    borderWidth: 1.5,
                    fill: false,
                    tension: 0.3,
                    pointRadius: 0,
                    pointHoverRadius: 3
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: {
                duration: 0 // Disable animations to prevent flickering
            },
            plugins: {
                legend: {
                    display: true,
                    position: 'top',
                    labels: {
                        usePointStyle: true,
                        boxWidth: 4,
                        font: {
                            size: 9
                        }
                    }
                }
            },
            scales: {
                x: {
                    display: false
                },
                y: {
                    beginAtZero: true,
                    max: 100,
                    display: true,
                    grid: {
                        color: 'rgba(0,0,0,0.08)',
                        lineWidth: 0.5
                    },
                    ticks: {
                        display: true,
                        font: {
                            size: 8
                        },
                        maxTicksLimit: 4,
                        callback: function(value) {
                            return value + '%';
                        }
                    }
                }
            },
            interaction: {
                intersect: false,
                mode: 'index'
            }
        }
    });
}

function updateChart(hostname, data) {
    const chart = window.charts[hostname];
    if (!chart) return;
    
    chart.data.labels = data.timestamps.map(() => '');
    chart.data.datasets[0].data = [...data.cpu];
    chart.data.datasets[1].data = [...data.memory];
    chart.data.datasets[2].data = [...data.disk];
    
    chart.update('none'); // Update without animation to prevent flicker
}

// Clean up old charts when agents disconnect
function cleanupOldCharts(currentAgents) {
    Object.keys(window.charts).forEach(hostname => {
        if (!currentAgents.includes(hostname)) {
            window.charts[hostname].destroy();
            delete window.charts[hostname];
            delete window.agentData[hostname];
        }
    });
}