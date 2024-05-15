import React, { useState } from 'react';

function MobileLeaderboard({
    listRank,
}) {

  return (
    <div>
        {/* Top 3 Entries */}
        <div className="columns is-multiline">
      {listRank && listRank.results && listRank.results.slice(0, 3).map((entry, index) => (
        <div 
          key={entry.username} // Assuming 'username' is a unique identifier for each entry
          className={`column is-one-third-desktop is-full-mobile has-background-dark has-text-centered`}
        >
          <div className="card">
            <div className="card-image">
              <figure className="image is-128x128 is-inline-block">
                <img className="is-rounded" src={entry.image} alt={`Profile of ${entry.username}`} />
              </figure>
            </div>
            <div className="card-content">
              <div className="media">
                <div className="media-content">
                  <p className="title is-4">{entry.username}</p>
                  <p className="subtitle is-6">@{entry.username}</p>
                </div>
              </div>
              <div className="content">
                Score: <strong>{entry.score}</strong>
                <br />
                Rank: <strong>{index + 1}</strong>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>


        {/* Entries 4 and below */}
        <div className="has-background-dark columns is-multiline">
        {listRank && listRank.results && listRank.results.slice(3).map((entry) => (
            <div key={entry.place} className="column is-full">
              <LeaderboardEntry entry={entry} />
            </div>
          ))}
        </div>
      </div>
  );
}

function LeaderboardEntry({ entry, topThree }) {
  return (
    <div className={`leaderboard-entry has-text-white`}>
      
      {topThree && entry.place === 1 && <div className="crown-icon">👑</div>}
      {entry.userAvatarObjectUrl
                            ?
                            <figure class="figure-img is-128x128">
                                <img src={entry.userAvatarObjectUrl} />
                            </figure>
                            :
                            <figure class="figure-img is-128x128">
                                <img src="/static/default_user.jpg" />
                            </figure>
                        }
      <div className="entry-details">
        <p className="name">{entry.userFirstName}</p>
        <p className="has-text-white">{Math.round(entry.value)}</p>
        <p className="username">@{entry.userFirstName.toLowerCase()}</p>
      </div>
    </div>
  );
}

export default MobileLeaderboard;