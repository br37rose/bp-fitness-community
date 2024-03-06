import React from 'react';
import {
    Chart as ChartJS,
    ArcElement,
    Tooltip,
    Legend
} from 'chart.js';
import { Doughnut } from 'react-chartjs-2';

// Register the necessary components
ChartJS.register(
    ArcElement,
    Tooltip,
    Legend
);

const DoughnutChart = () => {
    // Sample data
    const data = {
        labels: ['Present', 'Medical Leave', 'Unpaid Leave', 'Absent'],
        datasets: [{
            data: [70, 20, 5, 5],
            backgroundColor: ['#62BF92', '#7159EF', '#E4C261', '#E67270'],
            borderColor: 'rgba(54, 162, 235, 1)',
            borderWidth: 0,
        }],
    };

    // Chart options
    const options = {
        responsive: true,
        maintainAspectRatio: false, // Add this to control the height independently of the width
        animation: {
            duration: 1000, // Duration of animations in milliseconds
            easing: 'easeOutCubic',
            animateScale: true,
            animateRotate: true
        },
        plugins: {
            legend: {
                position: 'bottom',
                labels: {
                    usePointStyle: true,
                    boxWidth: 6
                }
            },
            title: {
                display: true,
                text: 'Advanced Doughnut Chart'
            }
        },
        cutout: '32%', // This is now a string, representing the cutout percentage
    };

    return (
        <div className='box' style={{ height: '300px' }}>
            <Doughnut data={data} options={options} />
        </div>
    );
};

export default DoughnutChart;
