import { useState } from "react";
import DraggableItem from "./dragable";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faAdd, faMultiply } from "@fortawesome/free-solid-svg-icons";

function WorkoutDisplay({
  workouts,
  label,
  helpText,
  isDraggable,
  wrapperClass,
  onAdd,
  showIndex = false,
  getSelectedWorkouts,
  initialState = [],
}) {
  const [selectedWorkouts, setSelectedWorkouts] = useState(initialState);

  const toggleWorkoutSelection = (workout) => {
    const isSelected = selectedWorkouts.some((w) => w.id === workout.id);
    if (isSelected) {
      setSelectedWorkouts((prevSelected) =>
        prevSelected.filter((w) => w.id !== workout.id)
      );
    } else {
      setSelectedWorkouts((prevSelected) => [...prevSelected, workout]);
    }
    if (getSelectedWorkouts) {
      getSelectedWorkouts(
        selectedWorkouts.some((w) => w.id === workout.id)
          ? selectedWorkouts.filter((w) => w.id !== workout.id)
          : [...selectedWorkouts, workout]
      );
    }
  };

  const exercisePanelJSX = (exercise, index) => (
    <div key={index} className="workout-panel">
      <div className="columns">
        <div className="column">
          <span className="has-text-info-dark has-text-weight-bold">
            {showIndex && `${index + 1}. `}
            {exercise.exerciseName}
          </span>
        </div>
      </div>
      <div className="columns">
        <div className="column">
          <span className="small">Order:</span> {exercise.orderNumber}
        </div>
        <div className="column">
          <span className="">Rest Period: </span>
          {exercise.restPeriodInSecs}
        </div>
        <div className="column">
          <span className="">Reps:</span> {exercise.sets}
        </div>
        <div className="column">
          <span className="">Target: </span>
          {exercise.targetTimeInSecs || exercise.targetText}
        </div>
      </div>
    </div>
  );

  const workoutBoxJSX = (workout, index) => (
    <div
      key={index}
      className={`workout-box mb-4 ${isSelected(workout) && "selected"}`}
    >
      <div className="is-flex is-align-content-center is-justify-content-space-between mx-2 px-2">
        <span className="title is-6">
          <span className="small is-inline is-flex is-align-items-center">
            {index + 1}.{" "}
          </span>
          {workout.name}
        </span>

        <button
          className={`button   is-small ${
            isSelected(workout) ? "selected is-danger" : "is-success"
          }`}
          onClick={() => toggleWorkoutSelection(workout)}
        >
          {selectedWorkouts.some((w) => w.id === workout.id) ? (
            <FontAwesomeIcon icon={faMultiply} />
          ) : (
            <FontAwesomeIcon icon={faAdd} />
          )}
        </button>
      </div>

      <div className="panel p-3">
        {workout.workoutExercises?.map((exercise, index) =>
          isDraggable ? (
            <DraggableItem
              className={"exercise-item"}
              key={index}
              id={exercise.id}
              content={exercisePanelJSX(exercise, index)}
              onAdd={() => onAdd(exercise)}
            />
          ) : (
            exercisePanelJSX(exercise, index)
          )
        )}
      </div>
    </div>
  );

  const isSelected = (workout) =>
    selectedWorkouts.some((w) => w.id === workout.id);

  return (
    <div className={wrapperClass}>
      {label && (
        <label className="is-size-4 has-text-weight-bold is-family-secondary grey-1">
          {label}
        </label>
      )}
      <div className="mt-4 is-family-secondary">
        {workouts.map((workout, index) => workoutBoxJSX(workout, index))}
      </div>
      {helpText && <p className="help">{helpText}</p>}
    </div>
  );
}

export default WorkoutDisplay;
