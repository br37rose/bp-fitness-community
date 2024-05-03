import React from 'react';
import DoughnutChart from '../../Reusable/Charts/Doughnut';

const ActivityStream = () => {
    return (
        <div className="card">
            <div className="card-content">
                <p className="title is-6">Recent Activities</p>
                <DoughnutChart text="BMI Global Stats" />
            </div>
        </div>
    );
};

export default ActivityStream;
