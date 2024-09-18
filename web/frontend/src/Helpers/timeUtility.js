// A function to calculate the difference in seconds between two dates
const differenceInSeconds = (date1, date2) => {
    if (!(date1 instanceof Date && date2 instanceof Date)) {
        throw new Error("Invalid date objects provided.");
    }
    return Math.floor((date2 - date1) / 1000);
};

// Function to format the time difference
const formatTimeDifference = (diffInSeconds) => {
    // If the difference is negative, it means the date is in the future
    if (diffInSeconds < 0) {
        diffInSeconds = Math.abs(diffInSeconds);
        // Format future dates accordingly
        return formatFutureTimeDifference(diffInSeconds);
    }

    if (diffInSeconds < 60) {
        return `${diffInSeconds} seconds ago`;
    } else if (diffInSeconds < 3600) {
        return `${Math.floor(diffInSeconds / 60)} mins ago`;
    } else if (diffInSeconds < 86400) {
        return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    } else {
        return `${Math.floor(diffInSeconds / 86400)} days ago`;
    }
};

// Function to handle future dates
const formatFutureTimeDifference = (diffInSeconds) => {
    if (diffInSeconds < 60) {
        return `in ${diffInSeconds} seconds`;
    } else if (diffInSeconds < 3600) {
        return `in ${Math.floor(diffInSeconds / 60)} mins`;
    } else if (diffInSeconds < 86400) {
        return `in ${Math.floor(diffInSeconds / 3600)} hours`;
    } else {
        return `in ${Math.floor(diffInSeconds / 86400)} days`;
    }
};

// Main function to calculate the 'time ago'
const timeAgo = (createdTimestamp) => {
    try {
        const createdDate = new Date(createdTimestamp);
        if (isNaN(createdDate.getTime())) {
            throw new Error("Invalid timestamp provided.");
        }
        const currentDate = new Date();
        const diffInSeconds = differenceInSeconds(createdDate, currentDate);
        return formatTimeDifference(diffInSeconds);
    } catch (error) {
        console.error(error.message);
        return "Invalid date";
    }
};
function formatDateStringWithTimezone(dateString) {
    try {
        const date = new Date(dateString);
        if (isNaN(date.getTime())) {
            throw new Error("Invalid date");
        }

        const options = {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            second: 'numeric',
            timeZoneName: 'short'
        };
        const formattedDate = new Intl.DateTimeFormat('en-US', options).format(date);

        return formattedDate;
    } catch (error) {
        console.error('Error parsing date:', error.message);
        return null;
    }
}

function formatTimeDifferenceSince(start, end) {
    try {
        const startDate = new Date(start);
        const endDate = new Date(end);
        if (isNaN(startDate.getTime()) || isNaN(endDate.getTime())) {
            throw new Error("Invalid date");
        }

        const differenceInSeconds = Math.floor((endDate - startDate) / 1000);
        const differenceInMinutes = differenceInSeconds / 60;
        const differenceInHours = differenceInMinutes / 60;
        const differenceInDays = differenceInHours / 24;
        const differenceInMonths = differenceInDays / 30; // Approximate
        const differenceInYears = differenceInDays / 365; // Approximate

        if (differenceInSeconds < 60) {
            return `since a few seconds ago`;
        } else if (differenceInSeconds < 3600) {
            return `since ${Math.floor(differenceInMinutes)} minute(s) ago`;
        } else if (differenceInSeconds < 86400) {
            return `since ${Math.floor(differenceInHours)} hour(s) ago`;
        } else if (differenceInDays < 30) {
            return `since ${Math.floor(differenceInDays)} day(s) ago`;
        } else if (differenceInMonths < 12) {
            return `since ${Math.floor(differenceInMonths)} month(s) ago`;
        } else {
            return `since ${Math.floor(differenceInYears)} year(s) ago`;
        }
    } catch (error) {
        console.error('Error parsing dates:', error.message);
        return null;
    }
}

export { timeAgo, formatDateStringWithTimezone, formatTimeDifferenceSince };
// Auto-generated comment for change 10
