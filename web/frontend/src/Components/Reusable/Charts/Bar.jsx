import React from 'react';
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    Title,
    Tooltip,
    Legend,
    registerables
} from 'chart.js';
import { Bar } from 'react-chartjs-2';

// Register all controllers, elements, and scales with Chart.js
ChartJS.register(...registerables);

const BarChart = ({ data }) => {

    // Chart options
    const options = {
        responsive: true,
        maintainAspectRatio: false, // Add this to control the height independently of the width
        animation: {
            duration: 1000, // Duration of animations in milliseconds
            easing: 'easeOutCubic',
        },
        plugins: {
            legend: {
                position: 'top',
                labels: {
                    usePointStyle: true,
                    boxWidth: 6
                }
            },

            zoom: {
                pan: {
                    enabled: true,
                    mode: 'xy', // or 'x', 'y'
                },
                zoom: {
                    wheel: {
                        enabled: true,
                    },
                    pinch: {
                        enabled: true,
                    },
                    mode: 'xy',
                },
            },
            title: {
                display: true,
                text: data && data.text
            },
        },
        scales: {
            x: {
                display: true,
                grid: {
                    drawOnChartArea: false,
                },
                ticks: {
                    display: true,
                },
            },
            y: {
                display: true,
                grid: {
                    drawOnChartArea: true,
                    borderDash: [8, 4],
                },
                ticks: {
                    display: true,
                },
            },
        },
        elements: {
            bar: {
                tension: 0.5
            }
        }
    };

    return <div className='box' style={{ height: '300px', width: 'auto' }}>
        <Bar data={data} options={options} />
    </div>
};

export default BarChart;
