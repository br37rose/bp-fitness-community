import React from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import PropTypes from 'prop-types';

/**
 * Breadcrumb component displaying full breadcrumbs on desktop and, if provided,
 * a configurable back link on mobile.
 *
 * @param {Object} breadcrumbItems - Object containing desktop breadcrumb items and optionally mobile back link configuration.
 * @param {Object[]} breadcrumbItems.items - An array of breadcrumb items for desktop view.
 * @param {Object} [breadcrumbItems.mobileBackLinkItems] - Optional configuration for the mobile back link.
 * @returns {React.Component} The rendered Breadcrumb component.
 */
const Breadcrumb = ({ breadcrumbItems }) => {
    const { items, mobileBackLinkItems } = breadcrumbItems;

    // Error handling for items
    if (!items || !Array.isArray(items) || items.length === 0) {
        console.error('Breadcrumb component expects "items" to be a non-empty array.');
        return null;
    }

    return (
        <>
            {/* Desktop Breadcrumbs */}
            <nav className="breadcrumb has-background-light p-4 is-hidden-touch" aria-label="breadcrumbs">
                <ul>
                    {items.map((item, idx) => (
                        <li className={item.isActive ? 'is-active' : ''} key={idx}>
                            {item.link ? (
                                <Link to={item.link} aria-current={item.isActive ? 'page' : undefined}>
                                    <FontAwesomeIcon icon={item.icon} />&nbsp;{item.text}
                                </Link>
                            ) : (
                                <span>
                                    <FontAwesomeIcon icon={item.icon} />&nbsp;{item.text}
                                </span>
                            )}
                        </li>
                    ))}
                </ul>
            </nav>

            {/* Mobile Breadcrumbs - Render only if mobileBackLinkItems is provided */}
            {mobileBackLinkItems && (
                <nav className="breadcrumb has-background-light p-4 is-hidden-desktop" aria-label="breadcrumbs">
                    <ul>
                        <li>
                            <Link to={mobileBackLinkItems.link} aria-current="page">
                                <FontAwesomeIcon icon={mobileBackLinkItems.icon} />&nbsp;{mobileBackLinkItems.text}
                            </Link>
                        </li>
                    </ul>
                </nav>
            )}
        </>
    );
};

Breadcrumb.propTypes = {
    breadcrumbItems: PropTypes.shape({
        items: PropTypes.arrayOf(
            PropTypes.shape({
                text: PropTypes.string.isRequired,
                link: PropTypes.string,
                isActive: PropTypes.bool,
                icon: PropTypes.object.isRequired,
            })
        ).isRequired,
        mobileBackLinkItems: PropTypes.shape({
            link: PropTypes.string,
            text: PropTypes.string,
            icon: PropTypes.object,
        }),
    }).isRequired,
};

export default Breadcrumb;
