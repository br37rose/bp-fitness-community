import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTrash,
  faArrowLeft,
  faGauge,
  faVideo,
  faEye,
  faDumbbell,
  faTable,
  faSave,
  faCancel,
  faArrowUpRightFromSquare,
  faEdit,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import { deleteWorkoutAPI, getWorkoutListApi } from "../../../API/workout";
import {
  deleteTrainingProgAPI,
  getTrainingProgDetailAPI,
  patchTrainingProgAPI,
} from "../../../API/trainingProgram";
import PhasePanel from "./phasepanel";
import WorkoutDisplay from "../../Reusable/WorkoutsDisplay";
import Modal from "../../Reusable/modal";
import FormInputField from "../../Reusable/FormInputField";
import Accordion from "../../Reusable/accordion";

function AdminTPDetail() {
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
  const [listData, setListData] = useState([]);
  const [selectedWorkoutForDeletion, setSelectedWorkoutForDeletion] =
    useState(null);
  const [showAddWorkoutModal, setShowAddWorkoutModal] = useState(false);
  const [selectedPhase, setSelectedPhase] = useState(null);
  const [isModified, setIsModified] = useState(false);

  const [selectedWorkoutForRoutine, setselectedWorkoutForRoutine] = useState(
    {}
  );
  const [getSelectedWorkouts, setgetSelectedWorkouts] = useState({});
  // Event handling
  const handleAddWorkoutClick = (phase) => {
    setSelectedPhase(phase);
    setShowAddWorkoutModal(true);
  };

  const handleAddWorkoutModalClose = () => {
    setShowAddWorkoutModal(false);
  };

  const handleInputChange = (index, field, value) => {
    const updatedSelectedWorkouts = { ...selectedWorkoutForRoutine };
    const selectedPhaseId = selectedPhase.id;

    // Find the selected workout by index
    const selectedWorkout = updatedSelectedWorkouts[selectedPhaseId][index];

    // Update the corresponding property
    selectedWorkout[field] = value;

    // Update the state
    setselectedWorkoutForRoutine(updatedSelectedWorkouts);
    setIsModified(true);
  };

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
  // Misc.
  useEffect(() => {
    if (datum && !listData.length) {
      setFetching(true);
      getWorkoutListApi(
        {
          visibility: 1,
          sort_field: "created",
          sort_order: -1,
        },
        onListSuccess,
        onListError,
        onListDone
      );
    }
    return () => {};
  }, [datum]);

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
      setselectedWorkoutForRoutine(updatedWorkoutForRoutine);
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

  function onListSuccess(resp) {
    setListData(resp.results);
  }

  function onListError(er) {
    setErrors(er);
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
  }

  function onListDone() {
    setFetching(false);
  }

  const handleSaveButtonClick = () => {
    let phase = new Array();
    // Iterate over each entry in the frontend object
    for (const phaseId in selectedWorkoutForRoutine) {
      if (selectedWorkoutForRoutine.hasOwnProperty(phaseId)) {
        const routines = selectedWorkoutForRoutine[phaseId];
        const phaseRoutines = routines.map((routine) => ({
          workout_id: routine.id,
          routine_day: parseInt(routine.day),
          routine_week: parseInt(routine.week),
        }));
        let phaseVal =
          datum.trainingPhases &&
          datum.trainingPhases.find((phase) => phase.id === phaseId);
        // Create PhaseRequestIDO object
        phase.push({
          phase_id: phaseId,
          phase: phaseVal && phaseVal.phase,
          routines: phaseRoutines,
        });
      }
    }
    setFetching(true);
    let payload = {
      phases: phase,
    };
    patchTrainingProgAPI(id, payload, onPatchOK, onPatchError, onDone);
    setIsModified(false);
  };
  function onPatchOK(response) {
    setTopAlertStatus("success");
    setTopAlertMessage("Program updated");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
    setForceURL("/admin/training-program");
  }

  function onPatchError(apiErr) {
    setErrors(apiErr);
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed updating");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
    scrollToTop();
  }

  function onDone() {
    setFetching(false);
  }

  const handleSaveSelectedWorkouts = () => {
    const currentWorkouts = selectedWorkoutForRoutine[selectedPhase.id] || [];
    const initialWorkouts = getSelectedWorkouts[selectedPhase.id] || [];

    // Check if the current and initial workouts are different
    const hasChanges =
      JSON.stringify(currentWorkouts) !== JSON.stringify(initialWorkouts);

    setIsModified(hasChanges);
    setselectedWorkoutForRoutine((prevState) => ({
      ...prevState,
      [selectedPhase.id]: getSelectedWorkouts[selectedPhase.id],
    }));
    setShowAddWorkoutModal(false);
  };

  // Helper function to scroll to the top of the page
  const scrollToTop = () => {
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  };

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
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
                <FontAwesomeIcon className="fas" icon={faDumbbell} />
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
                  <FontAwesomeIcon className="fas" icon={faVideo} />
                  &nbsp;Training-program
                </p>
              </div>
              <div className="column has-text-right">
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
              <li>
                {isModified && (
                  <div className="mt-3">
                    <button
                      className="button is-success is-small"
                      onClick={handleSaveButtonClick}
                    >
                      <FontAwesomeIcon icon={faSave} />
                      <span>&nbsp;Save changes</span>
                    </button>
                  </div>
                )}
              </li>
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
                        <button
                          className="button is-primary is-small mr-1 ml-2 "
                          onClick={() => setShowAddWorkoutModal(true)}
                        >
                          <FontAwesomeIcon icon={faEdit} />
                          Add/Edit Workouts
                        </button>
                      </div>

                      {selectedWorkoutForRoutine &&
                      selectedWorkoutForRoutine[selectedPhase.id] &&
                      selectedWorkoutForRoutine[selectedPhase.id].length ? (
                        selectedWorkoutForRoutine[selectedPhase.id].map(
                          (routine, index) => (
                            <div key={index}>
                              <Accordion
                                head={
                                  <span>
                                    {routine.name}
                                    <Link
                                      className="ml-1"
                                      to={"/admin/workouts/" + routine.id}
                                      target="_blank"
                                    >
                                      <FontAwesomeIcon
                                        size="sm"
                                        icon={faArrowUpRightFromSquare}
                                      />
                                    </Link>
                                  </span>
                                }
                                content={
                                  <>
                                    <div className="panel-block">
                                      <p>{routine.description}</p>
                                    </div>
                                    <div className="columns px-2">
                                      <div className="column">
                                        <FormInputField
                                          label={"Week"}
                                          placeholder={"week"}
                                          maxWidth={"120px"}
                                          type="number"
                                          value={routine.week}
                                          onChange={(e) =>
                                            handleInputChange(
                                              index,
                                              "week",
                                              e.target.value
                                            )
                                          }
                                        />
                                      </div>
                                      <div className="column">
                                        <FormInputField
                                          label={"Day"}
                                          placeholder={"day"}
                                          maxWidth={"120px"}
                                          type="number"
                                          value={routine.day}
                                          onChange={(e) =>
                                            handleInputChange(
                                              index,
                                              "day",
                                              e.target.value
                                            )
                                          }
                                        />
                                      </div>
                                    </div>
                                    {routine.workoutExercises &&
                                      routine.workoutExercises.map((we, i) => (
                                        <>
                                          <div className="box">
                                            <span>
                                              {i + 1}. {we.exerciseName}
                                            </span>
                                            <span className="label is-inline is-small ml-1">
                                              {we.set || 0} reps -{" "}
                                              {we.restPeriodInSecs} sec rest
                                            </span>
                                          </div>
                                        </>
                                      ))}
                                  </>
                                }
                                isOpenByDefault={index === 0}
                              />
                            </div>
                          )
                        )
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
                              add workouts started creating your first workout.
                            </p>
                          </div>
                        </section>
                      )}
                    </>
                  )}
                </div>
              </div>

              <Modal
                isOpen={showAddWorkoutModal}
                onClose={handleAddWorkoutModalClose}
                contentLabel="Add Workout Modal"
              >
                <h2>Add Workout</h2>
                {selectedPhase && (
                  <div>
                    <div className="is-flex is-justify-content-space-between">
                      <span>Phase: {selectedPhase.name}</span>
                      <div>
                        <button
                          className="button is-small is-success"
                          onClick={handleSaveSelectedWorkouts}
                        >
                          <FontAwesomeIcon icon={faSave} className="mr-1" />
                          Save Selecteds
                        </button>
                        <button
                          className="button is-small is-primary mr-1 ml-1"
                          onClick={() => setShowAddWorkoutModal(false)}
                        >
                          <FontAwesomeIcon icon={faCancel} className="mr-1" />
                          Cancel
                        </button>
                      </div>
                    </div>
                    <WorkoutDisplay
                      workouts={listData}
                      initialState={
                        selectedWorkoutForRoutine[selectedPhase.id] || []
                      }
                      getSelectedWorkouts={(wr) => {
                        setgetSelectedWorkouts({
                          ...getSelectedWorkouts,
                          [selectedPhase.id]: wr,
                        });
                      }}
                    />
                  </div>
                )}
              </Modal>
              <FormErrorBox errors={errors} />
            </>
          )}
        </nav>
      </section>
    </div>
  );
}

export default AdminTPDetail;
