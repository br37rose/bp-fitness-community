import React, { useEffect } from "react";
import DragSortableList from "./dragsortableList";
import FormSelectField from "./FormSelectField";
import FormInputField from "./FormInputField";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMultiply } from "@fortawesome/free-solid-svg-icons";
import { RESTPERIOD_OPTIONS } from "../../Constants/FieldOptions";

const DragSortListForSelectedWorkouts = ({
  selectedWorkouts,
  onRemove,
  onSortChange,
  handleInputChange,
}) => {
  useEffect(() => {
    return () => {};
  }, [selectedWorkouts]);

  const renderWorkoutContent = (exercise) => {
    return (
      <>
        <div className="is-flex is-align-content-center" style={{}}>
          <div>
            {exercise.thumbnailUrl ||
            (exercise.excercise && exercise.excercise?.thumbnailUrl) ? (
              <div className="is-flex is-align-items-center">
                <img
                  src={
                    exercise.thumbnailUrl ||
                    (exercise.excercise && exercise.excercise?.thumbnailUrl)
                  }
                  alt="Exercise Thumbnail"
                  style={{
                    width: "50px",
                    height: "50px",
                    marginRight: "10px",
                    borderRadius: "5px",
                  }}
                />
                <span>{exercise.name}</span>
              </div>
            ) : exercise.isRest ? (
              "REST"
            ) : (
              exercise.name
            )}
            <div className="columns">
              {!exercise.isRest && (
                <div className="column">
                  <FormInputField
                    label={"Reps"}
                    maxWidth={"80px"}
                    type="number"
                    value={exercise.reps}
                    onChange={(e) => handleInputChange(e, exercise.id, "reps")}
                    placeholder="Reps"
                  />
                </div>
              )}
              <div className="column">
                <FormSelectField
                  label={"Rest "}
                  maxWidth={"80px"}
                  options={RESTPERIOD_OPTIONS}
                  selectedValue={exercise.restPeriod}
                  onChange={(e) =>
                    handleInputChange(e, exercise.id, "restPeriod")
                  }
                  placeholder="Rest"
                />
              </div>
              {!exercise.isRest && (
                <>
                  <div className="column">
                    <FormSelectField
                      label={"TargetType"}
                      placeholder={"Target type"}
                      options={[
                        { value: 0, label: "select" },
                        { value: 1, label: "time" },
                        { value: 2, label: "text" },
                      ]}
                      onChange={(e) =>
                        handleInputChange(e, exercise.id, "targetType")
                      }
                    />
                  </div>
                  {exercise.targetType == 1 && (
                    <div className="column">
                      <div className="is-flex is-align-items-center">
                        <FormInputField
                          label={"TargetTime"}
                          maxWidth={"80px"}
                          type="number"
                          value={exercise.targetTime}
                          onChange={(e) =>
                            handleInputChange(e, exercise.id, "targetTime")
                          }
                          placeholder="Time"
                        />
                        <span className="label is-small is-inline has-text-grey ml-1">
                          Sec
                        </span>
                      </div>
                    </div>
                  )}

                  <div className="column">
                    <FormInputField
                      label={"Target"}
                      maxWidth={"150px"}
                      type="text"
                      value={exercise.targetText}
                      onChange={(e) =>
                        handleInputChange(e, exercise.id, "targetText")
                      }
                      placeholder="Target Text"
                    />
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
        <div className="column"></div>
        <span className="ml-2 mr-2" onClick={() => onRemove(exercise)}>
          <FontAwesomeIcon color="darkred" icon={faMultiply} />
        </span>
      </>
    );
  };

  return (
    <div>
      <DragSortableList
        items={selectedWorkouts.map((exercise) => ({
          ...exercise,
          content: renderWorkoutContent(exercise),
        }))}
        onSortChange={onSortChange}
      />
    </div>
  );
};

export default DragSortListForSelectedWorkouts;
