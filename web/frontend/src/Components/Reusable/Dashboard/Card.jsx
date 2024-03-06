// Card.js
import React from 'react';

const Card = ({ title, value }) => {
    return (
        <div className="card">
            <div className="card-content">
                <p className="title is-4">{title}</p>
                <p className="subtitle is-6">{value}</p>
            </div>
        </div>
    );
};

export default Card;
