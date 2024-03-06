import React from 'react';
import Row from './Row';

const LeaderBoardTable = ({
    data,
    headers
}) => {
    const rows = data.map((item, index) => {

        return (
            <Row
                key={index}
                place={item.place}
                imageUrl={item.userAvatarObjectUrl}
                name={item.userFirstName.concat(" ", item.userLastName)}
                dailyAvg={item.value}
                sevenDayAvg={item.weeklyAvg} // Use value from second API response
                sevenDayAvgStart={item.weeklyAvgStart}
                sevenDayAvgEnd={item.weeklyAvgEnd}

            />
        );
    });

    return (
        <section className="table-main">
            <div className="container">
                <div className="columns">
                    <div className="column">
                        <div className="table-container">
                            <table className="table is-fullwidth is-scrollable">
                                <thead>
                                    <tr>
                                        {headers && headers.length > 0 && headers.map((header, index) => (
                                            <th key={index} className={header.className}>
                                                {header.component ? (
                                                    header.component
                                                ) : (
                                                    <h2 className="is-size-6 is-size-6-mobile has-text-weight-semibold has-text-grey">
                                                        {header.title}
                                                    </h2>
                                                )}
                                            </th>
                                        ))}
                                    </tr>
                                </thead>
                                <tbody>
                                    {rows}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </section>

    );
};

export default LeaderBoardTable;
