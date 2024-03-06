import React from 'react';
import { timeAgo } from '../../../../Helpers/timeUtility';

const Row = ({
    place,
    imageUrl,
    name,
    dailyAvg,
    sevenDayAvg, // New prop for 7 days average
    sevenDayAvgStart,
    sevenDayAvgEnd
}) => {
    function getDateFormatFromISOString(isoString) {
        let currentDate = new Date(isoString);
        if (isNaN(currentDate.getTime())) {
            throw new Error("Invalid ISO date string.");
        }

        return currentDate;
    }
    return (
        <tr>
            <th class=" is-vcentered"><h2 class="is-size-6 has-text-weight-semibold has-text-grey">{place}</h2></th>
            <th className="is-vcentered">
                <div className="media" style={{ alignItems: 'center' }}>  {/* Align items vertically */}
                    <figure className="media-left image" style={{ width: '64px', height: '64px', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                        {imageUrl ? (
                            <img className="is-rounded" src={imageUrl} alt={name} style={{ borderRadius: '50%', width: '64px', height: '64px', objectFit: 'cover' }} />
                        ) : (
                            <img className="is-rounded" src='/static/default_user.jpg' alt={name} style={{ borderRadius: '50%', width: '64px', height: '64px', objectFit: 'cover' }} />
                        )}
                    </figure>
                    <div className="media-content" style={{ display: 'flex', alignItems: 'center' }}>  {/* This ensures vertical center alignment */}
                        <h5 className="is-size-6 has-text-weight-semibold mb-0 is-size-6-mobile" style={{ margin: '0' }}>{name}</h5>
                    </div>
                </div>
            </th>

            <td class="is-vcentered"><h4 class="is-size-6 mr-4 has-text-weight-semibold has-text-grey is-size-6-mobile">{Math.round(dailyAvg)}</h4></td>
            <td class="is-vcentered">
                {sevenDayAvg && (
                    <h5 class="is-size-4 has-text-centered has has-text-weight-semibold mb-0 is-size-5-mobile">{Math.round(sevenDayAvg)}</h5>
                )}
                <p class="is-size-7 has-text-weight-normal has-text-centered">{timeAgo(getDateFormatFromISOString(sevenDayAvgStart))}</p>
            </td>
        </tr>
    );
};

export default React.memo(Row);
