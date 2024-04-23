import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTrash,
  faArrowLeft,
  faGauge,
  faEye,
  faTable,
  faCalendarPlus,
  faEdit,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import {
  deleteTrainingProgAPI,
  getTrainingProgDetailAPI,
} from "../../../API/trainingProgram";
import PhasePanel from "./phasepanel";
import FitnessPlanDisplay from "../../Reusable/FitnessPlanDisplay";

function AdminTPDetailView() {
  // URL Parameters
  const { id } = useParams();

  // Global state
  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  // Component states
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [datum, setDatum] = useState({});
  const [selectedWorkoutForDeletion, setSelectedWorkoutForDeletion] =
    useState(null);
  const [selectedPhase, setSelectedPhase] = useState(null);

  // API

  // Detail
  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      setFetching(true);

      getTrainingProgDetailAPI(
        id,
        onDetailSuccess,
        onDetailError,
        onDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, [id]);

  // Delete
  function handleDeleteConfirmButtonClick() {
    deleteTrainingProgAPI(id, ondeleteSuccess, ondeleteError, onDeleteDone);
    setSelectedWorkoutForDeletion(null);
  }

  // Callbacks
  function onDetailSuccess(response) {
    setDatum(response);
    if (response.trainingPhases) {
      const updatedWorkoutForRoutine = {};
      response.trainingPhases.forEach((tp) => {
        if (tp.trainingRoutines && tp.trainingRoutines.length) {
          const phaseWorkout = tp.trainingRoutines.map((routine) => ({
            ...routine.workout, // Embed all fields of workout
            day:
              (routine.trainingDays &&
                routine.trainingDays.length > 0 &&
                routine.trainingDays[0].day) ||
              0,
            week:
              (routine.trainingDays &&
                routine.trainingDays.length > 0 &&
                routine.trainingDays[0].week) ||
              0,
          }));
          updatedWorkoutForRoutine[tp.id] = phaseWorkout;
        } else {
          updatedWorkoutForRoutine[tp.id] = [];
        }
      });
    }
  }

  function onDetailError(apiErr) {
    setErrors(apiErr);
    scrollToTop();
  }

  function onDetailDone() {
    setFetching(false);
  }

  function ondeleteSuccess(response) {
    setTopAlertStatus("success");
    setTopAlertMessage("training program deleted");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
    setForceURL("/admin/training-program");
  }

  function ondeleteError(apiErr) {
    setErrors(apiErr);
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
    scrollToTop();
  }

  function onDeleteDone() {
    setFetching(false);
  }

  function onDone() {
    setFetching(false);
  }
  const handleAddWorkoutClick = (phase) => {
    setSelectedPhase(phase);
  };
  // Helper function to scroll to the top of the page
  const scrollToTop = () => {
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  };

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }
  const getYouTubeVideoId = (url) => {
    const match = url.match(/[?&]v=([^&]+)/);
    return match && match[1];
  };

  function transformResponseForFitnessPlan(phase) {
    let plan = [
      {
        title: phase.name,
        dailyPlans: new Array(),
      },
    ];
    const setInstruction = (exc) => {
      let inst = exc.excercise && exc.excercise.description;
      if (exc.sets) {
        inst =
          inst +
          "\n " +
          "This has to be done for a set of " +
          exc.sets +
          "times.";
      }
      if (exc.restPeriodInSecs) {
        inst =
          inst +
          "\n " +
          "Give a rest period of around  " +
          exc.restPeriodInSecs +
          "seconds.";
      }
      if (exc.targetTimeInSecs) {
        inst =
          inst +
          "\n " +
          "Try to complete in  " +
          exc.targetTimeInSecs +
          "seconds.";
      }

      return inst;
    };
    const getTitle = (trainingDays) => {
      if (!trainingDays) {
        return "";
      }
      if (trainingDays.length) {
        return (
          "Week - " + trainingDays[0].week + ": Day - " + trainingDays[0].day
        );
      } else {
        return "";
      }
    };
    phase.trainingRoutines.map((tr) =>
      plan[0].dailyPlans.push({
        title: getTitle(tr.trainingDays),
        instructions: tr.description,
        planDetails:
          tr.workout.workoutExercises && tr.workout.workoutExercises.length > 0
            ? tr.workout.workoutExercises.map((exc) => ({
                id: exc.id,
                name: exc.exerciseName,
                videoUrl:
                  exc.excercise.videoType == 2
                    ? getYouTubeVideoId(exc.excercise.videoUrl)
                    : exc.excercise.videoUrl,
                thumbnailUrl: exc.excercise.thumbnailUrl,
                description: setInstruction(exc),
                videoType: exc.excercise.videoType,
              }))
            : new Array(),
      })
    );
    return plan;
  }

  return (
    <div className="container">
      <section className="section">
        <nav className="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
          <ul>
            <li className="">
              <Link to="/admin/dashboard" aria-current="page">
                <FontAwesomeIcon className="fas" icon={faGauge} />
                &nbsp;Dashboard
              </Link>
            </li>
            <li className="">
              <Link to="/admin/training-program" aria-current="page">
                <FontAwesomeIcon className="fas" icon={faCalendarPlus} />
                &nbsp;Training Program
              </Link>
            </li>
            <li className="is-active">
              <Link aria-current="page">
                <FontAwesomeIcon className="fas" icon={faEye} />
                &nbsp;Detail
              </Link>
            </li>
          </ul>
        </nav>

        <nav className="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
          <ul>
            <li className="">
              <Link to="/admin/training-program" aria-current="page">
                <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                &nbsp;Back to Training program
              </Link>
            </li>
          </ul>
        </nav>

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
                onClick={(e, ses) => setSelectedWorkoutForDeletion(null)}
              ></button>
            </header>
            <section class="modal-card-body">
              You are about to delete this Training Program and all the data
              associated with it. This action is cannot be undone. Are you sure
              you would like to continue?
            </section>
            <footer class="modal-card-foot">
              <button
                class="button is-success"
                onClick={handleDeleteConfirmButtonClick}
              >
                Confirm
              </button>
              <button
                class="button"
                onClick={(e, ses) => setSelectedWorkoutForDeletion(null)}
              >
                Cancel
              </button>
            </footer>
          </div>
        </div>

        <nav className="box">
          {datum && (
            <div className="columns">
              <div className="column">
                <p className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faCalendarPlus} />
                  &nbsp;Training-program
                </p>
              </div>
              <div className="column has-text-right">
                <Link
                  to={`/admin/training-program/${id}/edit`}
                  className="button is-primary is-small is-fullwidth-mobile mr-2"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faEdit} />
                  &nbsp;Edit
                </Link>
                <Link
                  onClick={(e, s) => {
                    setSelectedWorkoutForDeletion(datum);
                  }}
                  className="button is-danger is-small is-fullwidth-mobile"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faTrash} />
                  &nbsp;Delete
                </Link>
              </div>
            </div>
          )}

          <div className="tabs is-medium is-size-7-mobile">
            <ul className="is-flex is-justify-content-space-between">
              <li className="is-active ">
                <Link className="">
                  <strong>Detail</strong>
                </Link>
              </li>
              <li></li>
            </ul>
          </div>

          {isFetching ? (
            <PageLoadingContent displayMessage={"Please wait..."} />
          ) : (
            <>
              <div className="columns">
                <div className="column is-one-fifth">
                  <PhasePanel
                    phases={datum.trainingPhases}
                    onAddWorkout={handleAddWorkoutClick}
                    setSelectedPhase={setSelectedPhase}
                  />
                </div>
                <div className="column">
                  <h3>Workouts in phase</h3>
                  <p className="label is-small has-text-grey">
                    Click on a phase to know the workouts in each phase
                  </p>
                  {/* Render the available workouts in the selected phase here */}
                  {selectedPhase && (
                    <>
                      <div className="is-flex is-justify-content-space-between is-align-items-center mb-1">
                        <span>Phase: {selectedPhase.name}</span>
                      </div>

                      {selectedPhase.trainingRoutines &&
                      selectedPhase.trainingRoutines.length > 0 ? (
                        selectedPhase.trainingRoutines.map((tr, i) => (
                          <>
                            <FitnessPlanDisplay
                              label={tr.name}
                              key={tr.name + i}
                              weeklyFitnessPlans={transformResponseForFitnessPlan(
                                selectedPhase
                              )}
                            />
                          </>
                        ))
                      ) : (
                        // Show message if no workouts available
                        <section className="hero is-medium has-background-white-ter mt-1">
                          <div className="hero-body">
                            <p className="title">
                              <FontAwesomeIcon className="fas" icon={faTable} />
                              &nbsp;No Workouts
                            </p>
                            <p className="subtitle">
                              No workouts available in this phase. to get.Click
                              edit workouts to get started creating your first
                              workout.
                            </p>
                          </div>
                        </section>
                      )}
                    </>
                  )}
                </div>
              </div>

              <FormErrorBox errors={errors} />
            </>
          )}
        </nav>
      </section>
    </div>
  );
}

export default AdminTPDetailView;
