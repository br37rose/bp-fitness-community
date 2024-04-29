import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTrash,
  faArrowLeft,
  faGauge,
  faPencil,
  faVideo,
  faEye,
  faDumbbell,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
  EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER,
  EXERCISE_VIDEO_TYPE_YOUTUBE,
  EXERCISE_VIDEO_TYPE_VIMEO,
} from "../../../Constants/App";
import { deleteWorkoutAPI, getWorkoutDetailAPI } from "../../../API/workout";
import Vimeo from "@u-wave/react-vimeo";
import YouTubeVideo from "../../Reusable/YoutubePlayer";

function MemberWorkoutDetail() {
  ////
  //// URL Parameters.
  ////

  const { id } = useParams();

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [datum, setDatum] = useState({});
  const [selectedWorkoutForDeletion, setselectedWorkoutForDeletion] =
    useState(null);

  ////
  //// Event handling.
  ////

  const onDeleteConfirmButtonClick = () => {
    deleteWorkoutAPI(id, ondeleteSuccess, ondeleteError, onDeleteDone);
    setselectedWorkoutForDeletion(null);
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onVideoCollectionDetailSuccess(response) {
    setDatum(response);
  }

  function onVideoCollectionDetailError(apiErr) {
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onVideoCollectionDetailDone() {
    setFetching(false);
  }

  // --- Delete --- //

  function ondeleteSuccess(response) {
    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("workout deleted");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect back to the members page.
    setForceURL("/workouts");
  }

  function ondeleteError(apiErr) {
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onDeleteDone() {
    setFetching(false);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getWorkoutDetailAPI(
        id,
        onVideoCollectionDetailSuccess,
        onVideoCollectionDetailError,
        onVideoCollectionDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, [id]);
  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/workouts" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Workouts
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/workouts" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to workouts
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modal */}
          <nav>
            {/* Delete modal */}
            <div
              class={`modal ${
                selectedWorkoutForDeletion !== null ? "is-active" : ""
              }`}
            >
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">Are you sure?</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={(e, ses) => setselectedWorkoutForDeletion(null)}
                  ></button>
                </header>
                <section class="modal-card-body">
                  You are about to delete this workout and all the data
                  associated with it. This action is cannot be undone. Are you
                  sure you would like to continue?
                </section>
                <footer class="modal-card-foot">
                  <button
                    class="button is-success"
                    onClick={onDeleteConfirmButtonClick}
                  >
                    Confirm
                  </button>
                  <button
                    class="button"
                    onClick={(e, ses) => setselectedWorkoutForDeletion(null)}
                  >
                    Cancel
                  </button>
                </footer>
              </div>
            </div>
          </nav>

          {/* Page */}
          <nav class="box">
            {/* Title + Options */}
            {datum && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faVideo} />
                    &nbsp;Workouts
                  </p>
                </div>
                <div class="column has-text-right">
                  <Link
                    to={`/workouts/${id}/update`}
                    class="button is-warning is-small is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                  &nbsp;
                  <Link
                    onClick={(e, s) => {
                      setselectedWorkoutForDeletion(datum);
                    }}
                    class="button is-danger is-small is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrash} />
                    &nbsp;Delete
                  </Link>
                </div>
              </div>
            )}

            {/* Tab Navigation */}
            <div class="tabs is-medium is-size-7-mobile">
              <ul>
                <li class="is-active">
                  <Link>
                    <strong>Detail</strong>
                  </Link>
                </li>
              </ul>
            </div>

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                <div className="columns">
                  <div className="column ">
                    <h2 className="title is-4 mb-3">{datum?.name}</h2>
                    {datum?.workoutExercises?.map((exercise, index) => (
                      <div className="box" key={index}>
                        <h3 className="title is-5 mb-3">
                          Exercise {index + 1}:{" "}
                          {exercise.isRest ? "REST" : exercise.exerciseName}
                        </h3>
                        <div className="columns">
                          <div className="column is-three-quarters">
                            <p>
                              <span className="label">Description: </span>
                              {exercise.excercise.description}
                            </p>
                            <p>
                              <span className="label">Reps: </span>
                              {exercise.sets}
                            </p>
                            <p>
                              <span className="label">Rest Period: </span>
                              {exercise.restPeriodInSecs}
                            </p>
                            {/* Add other details regarding the exercise */}
                          </div>
                          <div className="column">
                            {/* Add video here */}
                            <div className="video-container">
                              {/* Render the video element */}
                              {(() => {
                                switch (exercise.excercise.videoType) {
                                  case EXERCISE_VIDEO_TYPE_SIMPLE_STORAGE_SERVER:
                                    return (
                                      <>
                                        <video
                                          style={{
                                            width: "100%",
                                            height: "100%",
                                          }}
                                          controls
                                          poster={exercise.excercise.thumbnailObjectUrl ||
                                          exercise.excercise.thumbnailUrl}
                                        >
                                          <source
                                            src={
                                              exercise.excercise.videoObjectUrl
                                            }
                                            type="video/mp4"
                                          />
                                        </video>
                                      </>
                                    );
                                  case EXERCISE_VIDEO_TYPE_YOUTUBE:
                                    return (
                                      <>
                                        <YouTubeVideo
                                          videoId={exercise.excercise.videoUrl}
                                        />
                                      </>
                                    );
                                  case EXERCISE_VIDEO_TYPE_VIMEO:
                                    return (
                                      <div className="vimeo-container">
                                        <Vimeo
                                        className="vimeo-wrapper"
                                          video={`${exercise.excercise.videoUrl}`}
                                          autoplay
                                        />
                                      </div>
                                    );
                                  default:
                                    return null;
                                }
                              })()}
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
                <FormErrorBox errors={errors} />
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default MemberWorkoutDetail;
