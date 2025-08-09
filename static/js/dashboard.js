// Global variables for managing charts and data
window.agentData = window.agentData || {};
window.charts = window.charts || {};

// Initialize dashboard updates
function initDashboard() {
    // Start polling for updates every second
    setInterval(updateDashboard, 1000);
    
    // Initial load via HTMX for the agent cards
    setTimeout(() => {
        updateDashboard();
    }, 100);
}

// Update dashboard with new data
async function updateDashboard() {
    try {
        const response = await fetch('/agents/data');
        const agents = await response.json();
        
        // Debug logging
        if (agents.length > 0) {
            console.log('Agent data structure:', agents[0]);
            console.log('Metrics structure:', agents[0].Metrics);
        }
        
        // Get current agent hostnames
        const currentAgents = agents.map(agent => agent.Hostname);
        
        // Clean up disconnected agents
        cleanupOldCharts(currentAgents);
        
        // Remove loading message on first successful update
        const agentsContainer = document.getElementById('agents');
        const loadingMsg = agentsContainer.querySelector('.loading');
        if (loadingMsg && agents.length > 0) {
            loadingMsg.remove();
        }
        
        // Update each agent
        agents.forEach(agent => {
            updateAgentCard(agent);
            createOrUpdateChart(agent.Hostname, {
                cpu: agent.Metrics.cpu_usage,
                memory: agent.Metrics.memory_usage,
                temp: agent.Metrics.cpu_temp
            });
        });
        
        // Handle no agents case
        if (agents.length === 0 && !agentsContainer.querySelector('.no-agents')) {
            agentsContainer.innerHTML = '<div class="no-agents">No agents currently online.</div>';
        }
        
    } catch (error) {
        console.error('Failed to update dashboard:', error);
    }
}

// Update individual agent card data without destroying the chart
function updateAgentCard(agent) {
    const existingCard = document.querySelector(`[data-agent-id="${agent.Hostname}"]`);
    
    if (!existingCard) {
        // Agent doesn't exist, create new card
        createAgentCard(agent);
        return;
    }
    
    // Update existing card metrics
    const metrics = [
        { selector: '.cpu .metric-value', value: `${agent.Metrics.cpu_usage}%` },
        { selector: '.memory .metric-value', value: `${agent.Metrics.memory_usage}%` },
        { selector: '.disk .metric-value', value: `${agent.Metrics.disk_usage}%` },
        { selector: '.network .metric-value', value: `↑${agent.Metrics.network_upload} ↓${agent.Metrics.network_download}` },
        { selector: '.temp .metric-value', value: `${agent.Metrics.cpu_temp}°C` },
        { selector: '.uptime .metric-value', value: agent.FormattedUptime },
        { selector: '.last-seen', value: `Last seen ${agent.LastSeenAgo}` }
    ];
    
    metrics.forEach(metric => {
        const element = existingCard.querySelector(metric.selector);
        if (element) {
            element.textContent = metric.value;
        }
    });
}

// Create a new agent card
function createAgentCard(agent) {
    const agentsContainer = document.getElementById('agents');
    
    // Remove "no agents" or "loading" message if it exists
    const noAgentsMsg = agentsContainer.querySelector('.no-agents');
    const loadingMsg = agentsContainer.querySelector('.loading');
    if (noAgentsMsg) {
        noAgentsMsg.remove();
    }
    if (loadingMsg) {
        loadingMsg.remove();
    }
    
    const cardHTML = `
        <div class="agent-card" data-agent-id="${agent.Hostname}">
            <div class="agent-header">
                <div class="agent-title">
                    <span class="status-indicator"></span>
                    ${agent.Hostname}
                </div>
                <div class="agent-os">${agent.OS}</div>
            </div>
            
            <div class="metrics-grid">
                <div class="metric-item cpu">
                    <div class="metric-icon">⚡</div>
                    <div class="metric-content">
                        <div class="metric-label">CPU</div>
                        <div class="metric-value">${agent.Metrics.cpu_usage}%</div>
                    </div>
                </div>
                
                <div class="metric-item memory">
                    <div class="metric-icon">▣</div>
                    <div class="metric-content">
                        <div class="metric-label">Memory</div>
                        <div class="metric-value">${agent.Metrics.memory_usage}%</div>
                    </div>
                </div>
                
                <div class="metric-item disk">
                    <div class="metric-icon">◉</div>
                    <div class="metric-content">
                        <div class="metric-label">Disk</div>
                        <div class="metric-value">${agent.Metrics.disk_usage}%</div>
                    </div>
                </div>
                
                <div class="metric-item network">
                    <div class="metric-icon">⟲</div>
                    <div class="metric-content">
                        <div class="metric-label">Network</div>
                        <div class="metric-value">↑${agent.Metrics.network_upload} ↓${agent.Metrics.network_download}</div>
                    </div>
                </div>
                
                <div class="metric-item temp">
                    <div class="metric-icon">◐</div>
                    <div class="metric-content">
                        <div class="metric-label">Temp</div>
                        <div class="metric-value">${agent.Metrics.cpu_temp}°C</div>
                    </div>
                </div>
                
                <div class="metric-item uptime">
                    <div class="metric-icon">○</div>
                    <div class="metric-content">
                        <div class="metric-label">Uptime</div>
                        <div class="metric-value">${agent.FormattedUptime}</div>
                    </div>
                </div>
            </div>
            
            <div class="chart-container">
                <canvas id="chart-${agent.Hostname}"></canvas>
            </div>
            
            <div class="last-seen">
                Last seen ${agent.LastSeenAgo}
            </div>
        </div>
    `;
    
    agentsContainer.insertAdjacentHTML('beforeend', cardHTML);
}

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
            temp: []
        };
    }
    
    const data = window.agentData[hostname];
    const now = new Date();
    
    // Add current data point
    data.timestamps.push(now);
    data.cpu.push(metrics.cpu);
    data.memory.push(metrics.memory);
    data.temp.push(metrics.temp);
    
    // Keep only last 15 data points for performance
    const maxPoints = 15;
    if (data.timestamps.length > maxPoints) {
        data.timestamps.shift();
        data.cpu.shift();
        data.memory.shift();
        data.temp.shift();
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
                    label: 'Temp',
                    data: data.temp,
                    borderColor: '#e53e3e',
                    backgroundColor: 'rgba(229, 62, 62, 0.1)',
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
                        maxTicksLimit: 4
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
    chart.data.datasets[2].data = [...data.temp];
    
    chart.update('none'); // Update without animation to prevent flicker
}

// Clean up old charts when agents disconnect
function cleanupOldCharts(currentAgents) {
    Object.keys(window.charts).forEach(hostname => {
        if (!currentAgents.includes(hostname)) {
            // Remove the card from DOM
            const card = document.querySelector(`[data-agent-id="${hostname}"]`);
            if (card) {
                card.remove();
            }
            
            // Clean up chart resources
            window.charts[hostname].destroy();
            delete window.charts[hostname];
            delete window.agentData[hostname];
        }
    });
}

// Start dashboard when page loads
document.addEventListener('DOMContentLoaded', initDashboard);