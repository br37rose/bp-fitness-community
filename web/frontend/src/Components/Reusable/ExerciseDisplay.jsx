import { faP } from "@fortawesome/free-solid-svg-icons";
import React from "react";

function ExerciseDisplay({ exercises, label, helpText }) {
  console.log("excercises", exercises);
  return (
    <div>
      {label && (
        <label className="is-size-4 has-text-weight-bold is-family-secondary grey-1 ">
          {label}
        </label>
      )}
      <div className="mt-4 is-family-secondary exercise-container">
        {exercises.map((exercise, index) => (
          <div key={index} className="mb-4 exercise-item">
            <h2 className="mb-3  has-text-weight-medium is-size-6">{`${
              index + 1
            }. ${exercise.name}`}</h2>
            <video
              className="exercise-video"
              poster={exercise.thumbnailUrl}
              controls
            >
              <source src={exercise.videoUrl} type="video/mp4" />
              Your browser does not support the video tag.
            </video>
          </div>
        ))}
      </div>
      {helpText && <p className="help">{helpText}</p>}
    </div>
  );
}

export default ExerciseDisplay;
