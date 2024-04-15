import React from "react";

// DEVELOPER NOTE:
// To get a full list of devices please see:
// https://developers.google.com/fit/datatypes
const GOOGLE_FIT_DATAPOINT_TYPES = {
  "com.google.activity.segment": "Activity",
  "com.google.calories.bmr": "Basal metabolic rate",
  "com.google.calories.expended": "Calories burned",
  "com.google.cycling.pedaling.cadence": "Cycling pedaling cadence",
  "com.google.cycling.pedaling.cumulative": "Cycling pedaling cumulative",
  "com.google.heart_minutes": "Heart Points",
  "com.google.active_minutes": "Move Minutes",
  "com.google.power.sample": "Power",
  "com.google.step_count.cadence": "Step count cadence",
  "com.google.step_count.delta": "Step count delta",
  "com.google.activity.exercise": "Workout",
  "com.google.cycling.wheel_revolution.rpm": "Cycling wheel revolution",
  "com.google.cycling.wheel_revolution.cumulative":
    "Cycling wheel revolution cumulative",
  "com.google.distance.delta": "Distance delta",
  "com.google.location.sample": "Location sample",
  "com.google.speed": "Speed",
  "com.google.hydration": "Hydration",
  "com.google.nutrition": "Nutrition",
  "com.google.blood_glucose": "Blood glucose",
  "com.google.blood_pressure": "Blood pressure",
  "com.google.body.fat.percentage": "Body fat percentage",
  "com.google.body.temperature": "Body temperature",
  "com.google.cervical_mucus": "Cervical mucus",
  "com.google.cervical_position": "Cervical position",
  "com.google.heart_rate.bpm": "Heart Rate",
  "com.google.height": "Height",
  "com.google.menstruation": "Menstruation",
  "com.google.ovulation_test": "Ovulation test",
  "com.google.oxygen_saturation": "Oxygen saturation",
  "com.google.sleep.segment": "Sleep",
  "com.google.vaginal_spotting": "Vaginal spotting",
  "com.google.weight": "Weight",
};
const GoogleFitDataPointTypeToString = ({ datapoint }) => {
  try {
    return GOOGLE_FIT_DATAPOINT_TYPES[datapoint.dataTypeName];
  } catch (e) {
    return "Unknown type: " + datapoint.dataTypeName;
  } finally {
    // Do nothing...
  }
};

export default GoogleFitDataPointTypeToString;
