import React from "react";

const GoogleFitDataPointValueToString = ({ datapoint }) => {
  try {
    // DEVELOPER NOTE:
    // To get a full list of devices please see:
    // https://developers.google.com/fit/datatypes
    switch (datapoint.dataTypeName) {
      case "com.google.activity.segment": // https://developers.google.com/fit/datatypes/activity#activity
        if (datapoint.activitySegment.durationInMinutes === 1) {
          return (
            datapoint.activitySegment.activityTypeLabel +
            " for " +
            datapoint.activitySegment.durationInMinutes +
            " minute"
          );
        } else {
          return (
            datapoint.activitySegment.activityTypeLabel +
            " for " +
            datapoint.activitySegment.durationInMinutes +
            " minutes"
          );
        }

      case "com.google.calories.bmr": // https://developers.google.com/fit/datatypes/activity#basal_metabolic_rate_bmr
        return "-";
      // return datapoint.basalMetabolicRate.calories + " kcal per day";

      case "com.google.calories.expended": // https://developers.google.com/fit/datatypes/activity#calories_burned
        // return "-";
        return datapoint.caloriesBurned.calories + " kcal";

      case "com.google.cycling.pedaling.cadence": // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cadence
        return datapoint.cyclingPedalingCadence.rpm + " rpm";

      case "com.google.cycling.pedaling.cumulative": // https://developers.google.com/fit/datatypes/activity#cycling_pedaling_cumulative
        return datapoint.cyclingPedalingCumulative.revolutions + " count";

      case "com.google.heart_minutes": // https://developers.google.com/fit/datatypes/activity#heart_points
        return datapoint.heartPoints.intensity + " HP";

      case "com.google.active_minutes": // https://developers.google.com/fit/datatypes/activity#move_minutes
        return datapoint.moveMinutes.duration + " milliseconds";

      case "com.google.power.sample": // https://developers.google.com/fit/datatypes/activity#power
        return datapoint.power.watts + " watts";

      case "com.google.step_count.cadence": // https://developers.google.com/fit/datatypes/activity#step_count_cadence
        return datapoint.stepCountCadence.rpm + " steps/minute";

      case "com.google.step_count.delta": // https://developers.google.com/fit/datatypes/activity#step_count_delta
        return datapoint.stepCountDelta.steps + " count";

      case "com.google.activity.exercise": // https://developers.google.com/fit/datatypes/activity#workout
        return "deprecated";

      case "com.google.cycling.wheel_revolution.rpm": // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_rpm
        return datapoint.cyclingWheelRevolutionRpm.rpm + " rpm";

      case "com.google.cycling.wheel_revolution.cumulative": // https://developers.google.com/fit/datatypes/location#cycling_wheel_revolution_cumulative
        return (
          datapoint.cyclingWheelRevolutionCumulative.revolutions +
          " revolutions"
        );

      case "com.google.distance.delta": // https://developers.google.com/fit/datatypes/location#distance_delta
        return datapoint.distanceDelta.distance + " meters";

      case "com.google.location.sample": // https://developers.google.com/fit/datatypes/location#location_sample
        return (
          " latitude:" +
          datapoint.locationSample.latitude +
          " degrees, longitude" +
          datapoint.locationSample.longitude +
          " degrees, altitude:" +
          datapoint.locationSample.altitude +
          " meters, accuracy:" +
          datapoint.locationSample.accuracy +
          " meters"
        );

      case "com.google.speed": // https://developers.google.com/fit/datatypes/location#speed
        return datapoint.speed.speed + " meters per second";

      case "com.google.hydration": // https://developers.google.com/fit/datatypes/nutrition#hydration
        return datapoint.hydration.volume + " liters";

      case "com.google.nutrition": // https://developers.google.com/fit/datatypes/nutrition#nutrition
        return datapoint.nutrition.name + ""; //TODO

      case "com.google.blood_glucose": // https://developers.google.com/fit/datatypes/health#blood_glucose
        let out1 = "";
        out1 +=
          "level:" + datapoint.bloodGlucose.bloodGlucoseLevel + " mmol/L, ";
        switch (datapoint.bloodGlucose.temporalRelationToMeal) {
          case 1:
            out1 += "Reading wasn't taken before or after a meal";
            break;
          case 2:
            out1 += "Reading was taken during a fasting period";
            break;
          case 3:
            out1 += "Reading was taken before a meal";
            break;
          case 4:
            out1 += "Reading was taken after a meal";
            break;
          default:
          // Do nothing...
        }
        out1 += ", ";
        switch (datapoint.bloodGlucose.mealType) {
          case 1:
            out1 += "Unknown";
            break;
          case 2:
            out1 += "Breakfast";
            break;
          case 3:
            out1 += "Lunch";
            break;
          case 4:
            out1 += "Dinner";
            break;
          case 5:
            out1 += "Snack";
            break;
          default:
          // Do nothing.
        }
        out1 += ", ";
        switch (datapoint.bloodGlucose.temporalRelationToSleep) {
          case 1:
            out1 += "User was fully awake";
            break;
          case 2:
            out1 += "Before the user fell asleep";
            break;
          case 3:
            out1 += "After the user woke up";
            break;
          case 4:
            out1 += "While the user was still sleeping";
            break;
          default:
        }
        out1 += ", ";
        switch (datapoint.bloodGlucose.specimenSource) {
          case 1:
            out1 += "Interstitial fluid";
            break;
          case 2:
            out1 += "Capillary blood";
            break;
          case 3:
            out1 += "Plasma";
            break;
          case 4:
            out1 += "Serum";
            break;
          case 5:
            out1 += "Tears";
            break;
          case 6:
            out1 += "Whole blood";
            break;
          default:
        }
        return out1;

      case "com.google.blood_pressure": // https://developers.google.com/fit/datatypes/health#blood_pressure
        let out2 = "";
        out2 += "systolic:" + datapoint.bloodPressure.systolic + " mmHg, ";
        out2 += "diastolic:" + datapoint.bloodPressure.diastolic + " mmHg, ";
        switch (datapoint.bloodPressure.bodyPosition) {
          case 1:
            out2 += "Standing up";
            break;
          case 2:
            out2 += "Sitting down";
            break;
          case 3:
            out2 += "Lying down";
            break;
          case 4:
            out2 += "Reclining";
            break;
          default:
          // Do nothing...
        }
        out2 += ", ";
        switch (datapoint.bloodPressure.measurementLocation) {
          case 1:
            out2 += "Left wrist";
            break;
          case 2:
            out2 += "Right wrist";
            break;
          case 3:
            out2 += "Left upper arm";
            break;
          case 4:
            out2 += "Right upper arm";
            break;
          default:
          // Do nothing.
        }
        return out2;

      case "com.google.body.fat.percentage": // https://developers.google.com/fit/datatypes/health#body_fat_percentage
        return datapoint.bodyFatPercentage.percentage + " %";
      case "com.google.body.temperature": // https://developers.google.com/fit/datatypes/health#body_temperature
        let out3 = datapoint.bodyTemperature.bodyTemperature + " celsius, ";
        switch (datapoint.bodyTemperature.measurementLocation) {
          case 1:
            out1 += "Armpit";
            break;
          case 2:
            out1 += "Finger";
            break;
          case 3:
            out1 += "Forehead";
            break;
          case 4:
            out1 += "Mouth (oral)";
            break;
          case 5:
            out1 += "Rectum";
            break;
          case 6:
            out1 += "Temporal artery";
            break;
          case 7:
            out1 += "Toe";
            break;
          case 8:
            out1 += "Ear (tympanic)";
            break;
          case 9:
            out1 += "Wrist";
            break;
          case 10:
            out1 += "Vagina";
            break;
          default:
            // Do nothing...
            break;
        }
        return out3;
      case "com.google.cervical_mucus": // https://developers.google.com/fit/datatypes/health#cervical_mucus
        return "not supported";
      case "com.google.cervical_position": // https://developers.google.com/fit/datatypes/health#cervical_position
        return "not supported";
      case "com.google.heart_rate.bpm": // https://developers.google.com/fit/datatypes/health#heart_rate
        return datapoint.heartRateBpm.bpm + " bpm";
      case "com.google.height": // https://developers.google.com/fit/datatypes/health#height
        return datapoint.height.height + " meters";
      case "com.google.menstruation": // https://developers.google.com/fit/datatypes/health#menstruation
        return "not supported";
      case "com.google.ovulation_test": // https://developers.google.com/fit/datatypes/health#ovulation_test
        return "not supported";
      case "com.google.oxygen_saturation": // https://developers.google.com/fit/datatypes/health#oxygen_saturation
        let out4 = datapoint.oxygenSaturation.oxygenSaturation + " %,";
        out4 +=
          datapoint.oxygenSaturation.supplementalOxygenFlowRate + " L/min,";
        if (datapoint.oxygenSaturation.oxygenTherapyAdministrationMode === 1) {
          out4 += "administered by nasal canula";
        }
        if (datapoint.oxygenSaturation.oxygenSaturationSystem === 1) {
          out4 += "measured in peripheral capillaries";
        }
        if (
          datapoint.oxygenSaturation.oxygenSaturationMeasurementMethod === 1
        ) {
          out4 += "measured by pulse oximetry";
        }
        return out4;
      case "com.google.sleep.segment": // https://developers.google.com/fit/datatypes/health#sleep
        switch (datapoint.sleep.sleepSegmentType) {
          case 0:
            return "Unspecified or unknown if user is sleeping.";
          case 1:
            return "Awake; user is awake.";
          case 2:
            return "Sleeping; generic or non-granular sleep description.";
          case 3:
            return "Out of bed; user gets out of bed in the middle of a sleep session.";
          case 4:
            return "Light sleep; user is in a light sleep cycle.";
          case 5:
            return "Deep sleep; user is in a deep sleep cycle.";
          case 6:
            return "REM sleep; user is in a REM sleep cyle.";
          default:
            // Do nothing.
            return "-";
        }
      case "com.google.vaginal_spotting": // https://developers.google.com/fit/datatypes/health
        return "not supported";
      case "com.google.weight": // https://developers.google.com/fit/datatypes/health#weight
        return datapoint.weight.weight + " kgs";
      default:
        return "not implemented yet";
    }
  } catch (e) {
    console.log(
      "GoogleFitDataPointValueToString: Error: datapoint.dataTypeName:",
      datapoint.dataTypeName,
    );
    console.log(
      "GoogleFitDataPointValueToString: Error: datapoint:",
      datapoint,
    );
    console.log("GoogleFitDataPointValueToString: Error:", e);
    return "ERROR";
  }
};

export default GoogleFitDataPointValueToString;
