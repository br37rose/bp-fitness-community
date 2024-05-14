import React from 'react';

const UserItem = ({ user }) => {
  const { username, score, avatar } = user;
  return (
    <div className="user-item">
      <figure className="avatar">
        <img src={user.avatar} alt={`${user.username}'s avatar`} className="user-avatar" />
        {user.position === 1 && <img src="static/crown.png" alt="Crown" className="crown" />}
      </figure>
      <div className="user-info">
        <p className="username">{user.username}</p>
        <p className="score">{user.score}</p>
      </div>
    </div>
  );
};

export default UserItem;
