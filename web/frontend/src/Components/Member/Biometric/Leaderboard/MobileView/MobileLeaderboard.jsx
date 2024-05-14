import React, { useState } from 'react';

function MobileLeaderboard({
    listRank,
}) {

  return (
    <div>
        {/* Top 3 Entries */}
      <div className="columns ">
        {listRank && listRank.results && listRank.results.slice(0, 3).map((entry, index) => (
          <div
            key={index}
            className={`column has-background-link box`}
          >
            <LeaderboardEntry entry={entry} topThree={true} />
          </div>
        ))}
      </div>

        {/* Entries 4 and below */}
        <div className="has-background-primary columns is-multiline">
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
    <div className={`leaderboard-entry ${topThree ? 'has-text-white' : 'has-text-dark'}`}>
      
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
        <p className="score">{Math.round(entry.value)}</p>
        <p className="username">@{entry.userFirstName.toLowerCase()}</p>
      </div>
    </div>
  );
}

export default MobileLeaderboard;