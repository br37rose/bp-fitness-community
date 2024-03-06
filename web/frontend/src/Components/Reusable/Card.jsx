import React from 'react';
import PropTypes from 'prop-types';

const Card = ({ title, image, content, actions }) => {
  return (
    <div className="card my-3">
      {image && <div className="card-image">{image}</div>}
      <div className="card-content">
        {title && <h1 className="card-title is-size-4 has-text-weight-semibold has-text-info">{title}</h1>}
        {content && <div className="card-text">{content}</div>}
      </div>
      {actions && <div className="card-actions">{actions}</div>}
    </div>
  );
};

Card.propTypes = {
  title: PropTypes.string,
  image: PropTypes.node,
  content: PropTypes.node,
  actions: PropTypes.node,
};

export default Card;
