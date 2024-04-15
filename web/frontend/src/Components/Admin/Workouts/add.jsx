import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlus, faArrowLeft } from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormTextareaField from "../../Reusable/FormTextareaField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import DropZone from "../../Reusable/dropzone";
import { getExerciseListAPI } from "../../../API/Exercise";
import ExerciseDisplay from "../../Reusable/ExerciseDisplay";
import { postWorkoutCreateAPI } from "../../../API/workout";
import DragSortListForSelectedWorkouts from "../../Reusable/draglistforSelectWorkouts";

function AdminWorkoutAdd() {
  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [selectedWorkouts, setSelectedWorkouts] = useState([]);
  const [listdata, setlistdata] = useState([]);
  const [selectableExcercises, setselectableExcercises] = useState(listdata);
  const [forceURL, setForceURL] = useState("");

  const onSubmitClick = () => {
    // Logic to submit data
    let payload = {
      name: name,
      description: description,
      visibility: 1, //1. visible to all 2. personal
    };
    let workoutExcercises = new Array();
    selectedWorkouts.map((w, index) =>
      workoutExcercises.push({
        exercise_id: w.isRest ? null : w.id,
        exercise_name: w.isRest ? "REST" : w.name,
        is_rest: w.isRest === true,
        order_number: index + 1,
        sets: w.reps ? parseInt(w.reps) : 0,
        type: w.type ? parseInt(w.type) : 0,
        rest_period_in_secs: parseInt(w.restPeriod),
        target_time_in_secs: w.targetTime ? parseInt(w.targetTime) : 0,
        target_text: w?.targetText,
      })
    );
    payload.workout_exercises = workoutExcercises;
    postWorkoutCreateAPI(payload, onAddSuccess, onAddError, onAddDone);
  };

  function onAddSuccess(response) {
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Exercise created");
    setTopAlertStatus("success");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect the organization to the organization attachments page.
    setForceURL("/admin/workouts/" + response.id + "");
  }

  function onAddError(apiErr) {
    setErrors(apiErr);
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAddDone() {
    setFetching(false);
  }

  const getAllExcericses = () => {
    let params = new Map();
    params.set("page_size", 1000000);
    params.set("sort_field", "created");
    params.set("sort_order", "-1");
    getExerciseListAPI(
      params,
      onExerciseListSuccess,
      onExerciseListError,
      onExerciseListDone
    );
  };

  function onExerciseListSuccess(response) {
    if (response.results !== null) {
      setlistdata(response.results);
      setselectableExcercises(response.results);
      if (response.hasNextPage) {
        // setNextCursor(response.nextCursor);
      }
    } else {
      setlistdata([]);
      // setNextCursor("");
    }
  }

  function onExerciseListError(apiErr) {
    setErrors(apiErr);
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onExerciseListDone() {
    setFetching(false);
  }

  useEffect(() => {
    getAllExcericses();
    window.scrollTo(0, 0);
  }, []);

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  const onDrop = (item) => {
    const exercise = listdata.find((ex) => ex.id === item.id);
    const newWorkout = {
      ...exercise,
      reps: "",
      restPeriod: "",
      targetText: "",
      targetType: "",
    };

    setselectableExcercises((prevExercises) =>
      prevExercises.filter((e) => e.id !== exercise.id)
    );

    setSelectedWorkouts((prevWorkouts) => [...prevWorkouts, newWorkout]);
  };

  const handleInputChange = (e, exerciseId, field) => {
    const { value } = e.target;
    setSelectedWorkouts((prevWorkouts) =>
      prevWorkouts.map((workout) => {
        if (workout.id === exerciseId) {
          return { ...workout, [field]: value };
        }
        return workout;
      })
    );
  };

  const onRemove = (cancelledItem) => {
    // Move the cancelled item back to the exercises column
    setSelectedWorkouts((prevWorkouts) =>
      prevWorkouts.filter((workout) => workout.id !== cancelledItem.id)
    );
    if (!cancelledItem.isRest) {
      const exercise = listdata.find((ex) => ex.id === cancelledItem.id);
      setselectableExcercises((e) => [...e, exercise]);
    }
  };

  const handleAddRest = () => {
    const restId = `rest-${Date.now()}`;
    let restWorkout = { id: restId, restPeriod: 60, isRest: true };

    setSelectedWorkouts((prevWorkouts) => [...prevWorkouts, restWorkout]);
  };

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="container">
        <section className="section">
          <div className="box">
            <p className="title is-4">
              <FontAwesomeIcon icon={faPlus} />
              &nbsp;Add Workouts
            </p>
            <FormErrorBox errors={errors} />
            <p className="pb-4 mb-5 has-text-grey">
              Please fill out all the required fields before submitting this
              form.
            </p>

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                <div className="container">
                  <div className="columns">
                    <div className="column">
                      <FormTextareaField
                        rows={1}
                        name="Name"
                        placeholder="Name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        isRequired={true}
                        maxWidth="150px"
                      />
                    </div>
                    <div className="column">
                      <FormTextareaField
                        rows={1}
                        name="Description"
                        placeholder="Description"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        isRequired={true}
                        maxWidth="380px"
                      />
                    </div>
                  </div>
                  <div className="columns">
                    <div className="column">
                      <div className="workout-column">
                        <h2>Selected Workouts</h2>
                        <div className="is-flex is-justify-content-flex-end mr-3">
                          <button
                            className="button is-small is-success is-light ml-3"
                            onClick={handleAddRest}
                          >
                            Add Rest
                          </button>
                        </div>

                        <DropZone
                          className={"p-2 excersizeWrapper"}
                          onDrop={onDrop}
                          placeholder={!selectedWorkouts.length}
                        >
                          <DragSortListForSelectedWorkouts
                            onRemove={onRemove}
                            onSortChange={setSelectedWorkouts}
                            selectedWorkouts={selectedWorkouts}
                            handleInputChange={handleInputChange}
                          />
                        </DropZone>
                      </div>
                    </div>
                    <div className="column">
                      <div className="exercise-column">
                        <span>Exercises </span>{" "}
                        <span className="label is-inline is-small has-text-grey ">
                          Click or drag item to add
                        </span>
                        <ExerciseDisplay
                          wrapperclass={"excersizeWrapper"}
                          exercises={selectableExcercises}
                          isdraggable
                          onAdd={onDrop}
                          showindex={false}
                        />
                      </div>
                    </div>
                  </div>
                  <div className="columns pt-5">
                    <div className="column is-half">
                      <Link
                        className="button is-fullwidth-mobile"
                        to={`/admin/workouts`}
                      >
                        <FontAwesomeIcon icon={faArrowLeft} />
                        &nbsp;Back to workouts
                      </Link>
                    </div>
                    <div className="column is-half has-text-right">
                      <button
                        onClick={onSubmitClick}
                        className="button is-success is-fullwidth-mobile"
                        type="button"
                      >
                        <FontAwesomeIcon icon={faPlus} />
                        &nbsp;Submit
                      </button>
                    </div>
                  </div>
                </div>
              </>
            )}
          </div>
        </section>
      </div>
    </DndProvider>
  );
}

export default AdminWorkoutAdd;
