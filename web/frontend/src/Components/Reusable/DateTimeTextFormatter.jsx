function DateTimeTextFormatter({ value, timeZone = "America/Toronto" }) {
    if (value === "0001-01-01T00:00:00Z") {
        return "-";
    }

    try {
        // Create a JavaScript Date object from the input string
        const date = new Date(value);

        // Use the provided time zone
        const options = {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: 'numeric',
            minute: '2-digit',
            hour12: true,
            timeZone: timeZone,
        };

        // Format the date and time as "MM/DD/YYYY h:mm AM/PM" in the specified time zone
        const formattedDateTime = new Intl.DateTimeFormat('en-US', options).format(date);

        // // For debugging purposes only.
        // console.log("DateTimeTextFormatter | Input:", value);
        // console.log("DateTimeTextFormatter | Output:", formattedDateTime);

        return formattedDateTime;
    } catch (err) {
        return "Invalid ISO Date/Time";
    }
}


export default DateTimeTextFormatter;
