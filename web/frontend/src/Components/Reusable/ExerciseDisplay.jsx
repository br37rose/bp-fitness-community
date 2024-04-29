import YouTubeVideo from "./YoutubePlayer";
import DraggableItem from "./dragable";
import { useState } from "react";
import Vimeo from "@u-wave/react-vimeo";
import {
  EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
  EXERCISE_VIDEO_TYPE_VIMEO,
  EXERCISE_VIDEO_TYPE_YOUTUBE,
} from "../../Constants/App";

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

      {(exercise.videoUrl || exercise.videoObjectUrl) &&
        (() => {
          switch (exercise.videoType) {
            case EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER:
              return (
                <>
                  <video style={{ width: "100%", height: "100%" }} controls>
                    <source src={exercise.videoObjectUrl} type="video/mp4" />
                  </video>
                </>
              );
            case EXERCISE_VIDEO_TYPE_YOUTUBE:
              return (
                <>
                  <YouTubeVideo
                    width={"100%"}
                    height={"auto"}
                    videoId={exercise.videoUrl}
                    minHeight={"50vh"}
                  />
                </>
              );
            case EXERCISE_VIDEO_TYPE_VIMEO:
              return (
                <div className="vimeo-container">
                  <Vimeo
                    className="vimeo-wrapper"
                    video={`${exercise.videoUrl}`}
                  />
                </div>
              );
            default:
              return null;
          }
        })()}
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
