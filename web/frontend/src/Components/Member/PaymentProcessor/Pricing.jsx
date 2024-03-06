import { faArrowLeft, faCheck, faCheckCircle } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React from 'react';
import { Link } from 'react-router-dom';

const PricingOption = ({ title, url, subtitle, price, period, features, recommended }) => (
    <div className={`column`}>
        <div className={`card has-shadow`}>
            {recommended && (
                <div className="recommended-badge is-uppercase has-text-centered has-background-primary has-text-white has-text-weight-semibold py-2">Recommended</div>
            )}
            {url && title && <div className="card-image">
                <div>
                    <figure className="image is-16by9">
                        <img src={url} alt={title} />
                    </figure>
                </div>
            </div>}
            <div className="card-content">
                {title && <h3 className="title is-4 has-text-centered is-primary mb-4">{title}</h3>}
                {subtitle && <p className="subtitle has-text-centered is-6 mb-5">{subtitle}</p>}
                {price && period && <div className="pricing has-text-centered my-4">
                    <span className="price is-size-2 has-text-weight-bold">${price}</span>{period}
                </div>}
                {features && <ul className="features-list subtitle has-text-centered">
                    {features.map((feature, index) => (
                        <li key={index} className="has-text-dark">
                            <span className="icon is-large has-text-primary mr-2">
                                <FontAwesomeIcon icon={faCheck} />
                            </span>
                            {feature}
                        </li>
                    ))}
                </ul>}
                <div className="buttons is-centered mt-5">
                    <Link
                        to={`/subscriptions`}
                        className="button is-primary is-medium"
                        type="button"
                    >
                        <span>Buy Now</span>
                    </Link>
                </div>
            </div>
        </div>
    </div>
);

const PricingTable = () => {
    const pricingOptions = [
        {
            title: 'Basic Plan',
            subtitle: 'For beginners',
            price: 19,
            period: '/month',
            features: ['Access to basic features', 'Limited support'],
            url: 'https://images.pexels.com/photos/841130/pexels-photo-841130.jpeg?cs=srgb&dl=pexels-victor-freitas-841130.jpg&fm=jpg'
        },
        {
            title: 'Pro Plan',
            subtitle: 'For advanced users',
            price: 39,
            period: '/month',
            features: ['Access to all features', 'Premium support'],
            recommended: true,
            url: 'https://images.pexels.com/photos/841130/pexels-photo-841130.jpeg?cs=srgb&dl=pexels-victor-freitas-841130.jpg&fm=jpg'
        },
        {
            title: 'Enterprise Plan',
            subtitle: 'For businesses',
            price: 99,
            period: '/month',
            features: ['Custom features', 'Premium support', 'Priority access'],
            url: 'https://images.pexels.com/photos/841130/pexels-photo-841130.jpeg?cs=srgb&dl=pexels-victor-freitas-841130.jpg&fm=jpg'
        },
    ];

    return (
        <div className="columns is-multiline">
            {pricingOptions.map((option, index) => (
                <PricingOption key={index} {...option} />
            ))}
        </div>
    );
};

export default PricingTable;
