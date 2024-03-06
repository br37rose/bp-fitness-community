import React from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faStar } from '@fortawesome/free-solid-svg-icons';



function ScoreRatingFormatter(props) {
    const { value } = props;
    return (
        <>
		{value === 0 && <span>-</span>}
		{value >= 1 && <FontAwesomeIcon className="fas" icon={faStar} />}
		{value >= 2 && <FontAwesomeIcon className="fas" icon={faStar} />}
		{value >= 3 && <FontAwesomeIcon className="fas" icon={faStar} />}
		{value >= 4 && <FontAwesomeIcon className="fas" icon={faStar} />}
		{value === 5 && <FontAwesomeIcon className="fas" icon={faStar} />}
		</>
    );
}

export default ScoreRatingFormatter;
