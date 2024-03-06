import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowUp, faArrowDown } from "@fortawesome/free-solid-svg-icons";
import { formatTimeDifferenceSince } from "../../Helpers/timeUtility";

const AverageAndTimeComparison = ({ lastDay, thisDay, iconState, mode }) => {
    const averageDifference = thisDay.average - lastDay.average;
    const isPositiveDifference = averageDifference > 0;
    const icon = isPositiveDifference ? faArrowUp : faArrowDown;
    const formattedDifference = Math.abs(averageDifference).toFixed(2);

    const startDate = new Date(lastDay.end); // Assuming 'end' of lastDay is the start of thisDay
    const endDate = new Date(thisDay.end);

    let additionalText = '';
    if (mode === 1) {
        // Mode 1: Time difference
        const difference = formatTimeDifferenceSince(startDate, endDate);
        additionalText = difference;
    } else if (mode === 2) {
        // Mode 2: Step difference
        const difference = thisDay.sum - lastDay.sum; // Assuming count is total steps
        if (difference > 0) {
            additionalText = `${difference.toLocaleString()} more than last week`;
        } else if (difference < 0) {
            additionalText = `${Math.abs(difference).toLocaleString()} less than last week`;
        } else {
            additionalText = "the same as last week";
        }
    }

    return (
        <p className="is-size-6 is-size-6-mobile">
            {iconState && (
                <FontAwesomeIcon
                    className={`fas ${isPositiveDifference ? 'has-text-success' : 'has-text-danger'} is-size-6 is-size-6-mobile`}
                    icon={icon}
                />
            )}
            &nbsp; {mode === 1 ? formattedDifference.concat(" ", additionalText) : additionalText}
        </p>
    );
};

export default AverageAndTimeComparison;
