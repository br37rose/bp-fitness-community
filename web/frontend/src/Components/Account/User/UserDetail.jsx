import React from 'react';
import { USER_ROLE_MAP } from '../../../Constants/FieldOptions';

const UserDetail = ({ avatarObjectUrl, name, role }) => {
    const userRole = USER_ROLE_MAP[role];
    const defaultAvatarUrl = 'static/default_user.jpg';

    return (
        <div className="card">
            <div className="card-content">
                <div className="columns is-vcentered is-centered is-flex-direction-column is-justify-content-center is-align-items-center">
                    <div className="column">
                        <figure className="image is-128x128" style={{ margin: 'auto' }}>
                            <img className="is-rounded" src={avatarObjectUrl || defaultAvatarUrl} alt={name} style={{ objectFit: 'cover', height: '128px', width: '128px' }} />
                        </figure>
                    </div>
                    <div className="column">
                        <div className="has-text-centered">
                            <p className="title is-6">{name}</p>
                            <p className="subtitle is-italic has-text-link is-6">{userRole}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default UserDetail;
