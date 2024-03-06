import React from 'react';

const UserDetail = ({ avatarObjectUrl, name }) => {
    return (
        <div className="card">
            <div className="card-content">
                <div className="columns is-vcentered is-centered is-flex-direction-column is-justify-content-center is-align-items-center">
                    <div className="column">
                        <figure className="image is-128x128" style={{ margin: 'auto' }}>
                            <img className="is-rounded" src={avatarObjectUrl} alt={name} style={{ objectFit: 'cover', height: '128px', width: '128px' }} />
                        </figure>
                    </div>
                    <div className="column">
                        <div className="has-text-centered">
                            <p className="title is-6">{name}</p>
                            <p className="subtitle is-italic has-text-link is-6">Software Engineer</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

    );
};

export default UserDetail;
