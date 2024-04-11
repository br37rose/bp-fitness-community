import { faP } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import DraggableItem from "./dragable";

function ExerciseDisplay({
  exercises,
  label,
  helpText,
  isdraggable,
  wrapperclass,
  onAdd,
  showindex = true,
}) {
  const exerciseItemJSX = (exercise, index) => (
    <div key={index} className="mb-4 exercise-item">
      <h2 className="mb-3 has-text-weight-medium is-size-6">
        {showindex && `${index + 1}. `}
        {exercise.name}
      </h2>
      <video className="exercise-video" poster={exercise.thumbnailUrl} controls>
        <source src={exercise.videoUrl} type="video/mp4" />
        Your browser does not support the video tag.
      </video>
    </div>
  );

  return (
    <div className={wrapperclass}>
      {label && (
        <label className="is-size-4 has-text-weight-bold is-family-secondary grey-1 ">
          {label}
        </label>
      )}
      <div className="mt-4 is-family-secondary exercise-container">
        {exercises.map((exercise, index) =>
          isdraggable ? (
            <DraggableItem
              className={"exercise-item"}
              key={index}
              id={exercise.id}
              content={exerciseItemJSX(exercise, index)}
              onAdd={() => onAdd(exercise)}
            />
          ) : (
            exerciseItemJSX(exercise, index)
          )
        )}
      </div>
      {helpText && <p className="help">{helpText}</p>}
    </div>
  );
}

export default ExerciseDisplay;
