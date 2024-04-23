import DataDisplayRowText from "./DataDisplayRowText";
import DraggableItem from "./dragable";
import { useState } from "react";

function ExerciseDisplay({
  exercises,
  label,
  helpText,
  isdraggable,
  wrapperclass,
  onAdd,
  showindex = true,
  showDescription = true,
}) {
  const [showMore, setShowMore] = useState(false);

  const toggleDescription = () => {
    setShowMore(!showMore);
  };
  const trimDescription = (description) => {
    const maxLength = 100;
    return description.length > maxLength
      ? description.slice(0, maxLength) + "..."
      : description;
  };

  const exerciseItemJSX = (exercise, index) => (
    <div key={index} className="mb-3 exercise-item">
      <h2 className="mb-3 has-text-weight-bold is-size-6">
        {showindex && `${index + 1}. `}
        {exercise.name}
      </h2>
      {exercise.instructions && <p>{exercise.instructions}</p>}
      {showDescription && (
        <div>
          <p className="mb-3 exercise-description">
            {showMore
              ? exercise.description
              : trimDescription(exercise.description)}
            {exercise.description.length > 100 && (
              <button
                className="is-small transparent-button is-family-sans-serif has-text-weight-bold"
                onClick={toggleDescription}
              >
                {showMore ? "Show Less" : "More"}
              </button>
            )}
          </p>
        </div>
      )}
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
