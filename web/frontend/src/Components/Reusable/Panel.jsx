import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

const PanelItem = ({ 
    imageSrc, 
    alt, 
    title, 
    subtitle, 
    description, 
    ctaLink, 
    ctaText, 
    imagePosition, 
    backgroundColor 
}) => {
  const isImageLeft = imagePosition === 'left';
  const panelBackgroundColor = `has-background-${backgroundColor || 'black'}`;

  return (
    <div className={`my-3 column is-12 ${panelBackgroundColor}`}>
      <div className="columns is-vcentered is-vcentered-mobile">
        {isImageLeft ? (
          <>
            <div className="column">
              <figure className="image">
                <img src={imageSrc} alt={alt} />
              </figure>
            </div>
            <div className="column">
              <div className="content has-text-centered">
                <h2 className="has-text-white">{title}</h2>
                <h3 className="is-size-5 has-text-primary">{subtitle}</h3>
                <p className="has-text-white">{description}</p>
                <div>
                  <Link to={ctaLink} className="button is-primary">
                    {ctaText}
                  </Link>
                </div>
              </div>
            </div>
          </>
        ) : (
          <>
            <div className="column">
              <div className="content has-text-centered">
                <h2 className="has-text-white">{title}</h2>
                <h3 className="is-size-5 has-text-primary">{subtitle}</h3>
                <p className="has-text-white">{description}</p>
                <div>
                  <Link to="/membership/monthly-membership" className="button is-primary">
                    {ctaText}
                  </Link>
                </div>
              </div>
            </div>
            <div className="column">
              <figure className="image">
                <img src={imageSrc} alt={alt} />
              </figure>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

// PropTypes for the PanelItem component
PanelItem.propTypes = {
  imageSrc: PropTypes.string.isRequired,
  alt: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  subtitle: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  ctaLink: PropTypes.string.isRequired,
  ctaText: PropTypes.string.isRequired,
  imagePosition: PropTypes.oneOf(['left', 'right']), // 'left' or 'right'
  backgroundColor: PropTypes.string, // Custom background color class
};

export default PanelItem;
