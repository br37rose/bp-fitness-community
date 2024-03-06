import React from "react";

// Function to format time (e.g., 09:00 AM)
function formatTime(time) {
  let [hour, minute] = time.split(":");
  let period = hour >= 12 ? "PM" : "AM";
  hour = hour % 12 || 12;
  return `${hour}:${minute} ${period}`;
}

const OpeningHourText = ({
  item
}) => {
  let [day, timeRange] = item.split(" ");
  let [startTime, endTime] = timeRange.split("-");

  // Format the day
  let formattedDay = "";
  switch (day) {
    case "Mo":
      formattedDay = "Monday";
      break;
    case "Tu":
      formattedDay = "Tuesday";
      break;
    case "We":
      formattedDay = "Wednesday";
      break;
    case "Th":
      formattedDay = "Thursday";
      break;
    case "Fr":
      formattedDay = "Friday";
      break;
    case "Sa":
      formattedDay = "Saturday";
      break;
    case "Su":
      formattedDay = "Sunday";
      break;
    default:
      formattedDay = day;
  }

  // Format the start time
  let formattedStartTime = formatTime(startTime);

  // Format the end time
  let formattedEndTime = formatTime(endTime);

  return `${formattedDay} ${formattedStartTime} - ${formattedEndTime}`;
};

export default OpeningHourText;
