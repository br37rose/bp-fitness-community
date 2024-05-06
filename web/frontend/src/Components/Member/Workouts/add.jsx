import {useState, useEffect} from "react";
import {Link, Navigate} from "react-router-dom";
import Scroll from "react-scroll";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
	faPlus,
	faArrowLeft,
	faFilter,
	faSearch,
	faSave,
	faFilterCircleXmark,
	faClose,
} from "@fortawesome/free-solid-svg-icons";
import {useRecoilState} from "recoil";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormTextareaField from "../../Reusable/FormTextareaField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
	currentUserState,
	topAlertMessageState,
	topAlertStatusState,
} from "../../../AppState";
import {DndProvider} from "react-dnd";
import {HTML5Backend} from "react-dnd-html5-backend";
import DropZone from "../../Reusable/dropzone";
import {getExerciseListAPI} from "../../../API/Exercise";
import ExerciseDisplay from "../../Reusable/ExerciseDisplay";
import {postWorkoutCreateAPI} from "../../../API/workout";
import DragSortListForSelectedWorkouts from "../../Reusable/draglistforSelectWorkouts";
import Modal from "../../Reusable/modal";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectField from "../../Reusable/FormSelectField";
import {
	EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION,
	EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION,
	EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION,
	EXERCISE_STATUS_OPTIONS_WITH_EMPTY_OPTION,
	EXERCISE_VIDEO_FILE_TYPE_OPTIONS_WITH_EMPTY_OPTION,
} from "../../../Constants/FieldOptions";
import FormMultiSelectFieldForTags from "../../Reusable/FormMultiSelectFieldForTags";

function MemberWorkoutAdd() {
	const [topAlertMessage, setTopAlertMessage] =
		useRecoilState(topAlertMessageState);
	const [topAlertStatus, setTopAlertStatus] =
		useRecoilState(topAlertStatusState);

	const [currentUser] = useRecoilState(currentUserState);

	const [errors, setErrors] = useState({});
	const [isFetching, setFetching] = useState(false);
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [selectedWorkouts, setSelectedWorkouts] = useState([]);
	const [listdata, setlistdata] = useState([]);
	const [selectableExcercises, setselectableExcercises] = useState(listdata);
	const [forceURL, setForceURL] = useState("");
	const [exersizeLoading, setexersizeLoading] = useState(true);
	const [showExerciseFilter, setshowExerciseFilter] = useState(false);
	const [temporarySearchText, setTemporarySearchText] = useState("");
	const [actualSearchText, setActualSearchText] = useState("");
	const [category, setCategory] = useState("");
	const [movementType, setMovementType] = useState("");
	const [status, setStatus] = useState("");
	const [gender, setGender] = useState("");
	const [videoType, setVideoType] = useState("");
	const [tags, setTags] = useState([]);
	const onSubmitClick = () => {
		// Logic to submit data
		let payload = {
			name: name,
			description: description,
			visibility: 2, //1. visible to all 2. personal
			user_id: currentUser.id,
			user_name: currentUser.name,
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
		setForceURL("/workouts/" + response.id + "");
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

	const getAllExcericses = (clear = false, search = "") => {
		setexersizeLoading(true);
		let params = new Map();
		params.set("page_size", 1000000);
		params.set("sort_field", "created");
		params.set("sort_order", "-1");
		if ((!clear && actualSearchText !== "") || search != "") {
			if (search) {
				params.set("search", search);
			} else {
				params.set("search", actualSearchText);
			}
		}
		if (!clear && category !== "") {
			params.set("category", category);
		}
		if (!clear && movementType !== "") {
			params.set("movement_type", movementType);
		}
		if (!clear && status !== "") {
			params.set("status", status);
		}
		if (!clear && gender !== "") {
			params.set("gender", gender);
		}
		if (!clear && videoType !== "") {
			params.set("video_type", videoType);
		}
		if (tags.length > 0) {
			params.set("tags", tags);
		}
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
		setexersizeLoading(false);
	}

	const onSearchButtonClick = (e) => {
		setActualSearchText(temporarySearchText);
		getAllExcericses(false, temporarySearchText);
		setshowExerciseFilter(false);
	};

	function ApplyFilter() {
		getAllExcericses();
		setshowExerciseFilter(false);
	}
	const onClearFilterClick = (e) => {
		setshowExerciseFilter(false);
		setActualSearchText("");
		setTemporarySearchText("");
		setStatus("");
		setCategory("");
		setMovementType("");
		setVideoType("");
		setGender("");
		getAllExcericses(true);
	};

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
		const {value} = e.target;
		setSelectedWorkouts((prevWorkouts) =>
			prevWorkouts.map((workout) => {
				if (workout.id === exerciseId) {
					return {...workout, [field]: value};
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
		let restWorkout = {id: restId, restPeriod: 60, isRest: true};

		setSelectedWorkouts((prevWorkouts) => [...prevWorkouts, restWorkout]);
	};

	return (
		<DndProvider backend={HTML5Backend}>
			<div className="container is-fluid">
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
								<div>
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
														onClick={handleAddRest}>
														Add Rest
													</button>
												</div>

												<DropZone
													className={"p-2 excersizeWrapper"}
													onDrop={onDrop}
													placeholder={!selectedWorkouts.length}>
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
												<div className="is-flex is-justify-content-space-between">
													<div>
														<span>Exercises </span>{" "}
														<span className="label is-inline is-small has-text-grey ">
															Click or drag item to add
														</span>
													</div>
													<button
														className="button is-small is-primary mr-2 is-light"
														onClick={() => setshowExerciseFilter(true)}>
														<FontAwesomeIcon icon={faFilter} />
														&nbsp;Filter
													</button>
													<Modal
														isOpen={showExerciseFilter}
														onClose={() => setshowExerciseFilter(false)}>
														<div className="box">
															<div className="heading is-size-5">Filter</div>
															<div class="columns mt-1">
																<div class="column">
																	<FormInputFieldWithButton
																		label={"Search"}
																		name="temporarySearchText"
																		type="text"
																		placeholder="Search by name"
																		value={temporarySearchText}
																		helpText=""
																		onChange={(e) =>
																			setTemporarySearchText(e.target.value)
																		}
																		isRequired={true}
																		maxWidth="100%"
																		buttonLabel={
																			<>
																				<FontAwesomeIcon
																					className="fas"
																					icon={faSearch}
																				/>
																			</>
																		}
																		onButtonClick={onSearchButtonClick}
																	/>
																</div>
																<div class="column">
																	<FormSelectField
																		label="Category"
																		name="category"
																		placeholder="Pick"
																		selectedValue={category}
																		errorText={errors && errors.category}
																		helpText=""
																		onChange={(e) =>
																			setCategory(parseInt(e.target.value))
																		}
																		options={
																			EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION
																		}
																	/>
																</div>
																<div class="column">
																	<FormSelectField
																		label="Movement Type"
																		name="movementType"
																		placeholder="Pick"
																		selectedValue={movementType}
																		errorText={errors && errors.movementType}
																		helpText=""
																		onChange={(e) =>
																			setMovementType(parseInt(e.target.value))
																		}
																		options={
																			EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION
																		}
																	/>
																</div>

																<div class="column">
																	<FormSelectField
																		label="Gender"
																		name="gender"
																		placeholder="Pick"
																		selectedValue={gender}
																		errorText={errors && errors.gender}
																		helpText=""
																		onChange={(e) => setGender(e.target.value)}
																		options={
																			EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION
																		}
																	/>
																</div>
																<div class="column">
																	<FormSelectField
																		label="Video Type"
																		name="videoType"
																		placeholder="Pick"
																		selectedValue={videoType}
																		errorText={errors && errors.videoType}
																		helpText=""
																		onChange={(e) =>
																			setVideoType(e.target.value)
																		}
																		options={
																			EXERCISE_VIDEO_FILE_TYPE_OPTIONS_WITH_EMPTY_OPTION
																		}
																	/>
																</div>

																<div class="column">
																	<FormMultiSelectFieldForTags
																		label="Tags"
																		name="tags"
																		placeholder="Pick tags"
																		tags={tags}
																		setTags={setTags}
																		errorText={errors && errors.tags}
																		helpText=""
																		isRequired={true}
																		maxWidth="320px"
																	/>
																</div>
															</div>
															<div className="is-flex is-justify-content-flex-end">
																<button
																	className="button is-small is-success"
																	onClick={ApplyFilter}>
																	<FontAwesomeIcon
																		icon={faSave}
																		className="mr-1"
																	/>
																	Apply
																</button>
																<button
																	className="button is-small is-primary ml-2"
																	onClick={() => setshowExerciseFilter(false)}>
																	<FontAwesomeIcon
																		icon={faClose}
																		className="mr-1"
																	/>
																	Close
																</button>
																<button
																	className="button is-small is-secondary ml-2 is-light"
																	onClick={onClearFilterClick}>
																	<FontAwesomeIcon
																		icon={faFilterCircleXmark}
																		className="mr-1"
																	/>
																	Clear filter
																</button>
															</div>
														</div>
													</Modal>
												</div>
												{exersizeLoading ? (
													<PageLoadingContent
														displayMessage={"Please wait..."}
													/>
												) : (
													<ExerciseDisplay
														wrapperclass={"excersizeWrapper"}
														exercises={selectableExcercises}
														isdraggable
														onAdd={onDrop}
														showindex={false}
														showDescription={false}
													/>
												)}
											</div>
										</div>
									</div>
									<div className="columns pt-5">
										<div className="column is-half">
											<Link class="button is-hidden-touch" to={`/workouts`}>
												<FontAwesomeIcon icon={faArrowLeft} />
												&nbsp;Back to workouts
											</Link>
											<Link
												class="button is-fullwidth is-hidden-desktop"
												to={`/workouts`}>
												<FontAwesomeIcon icon={faArrowLeft} />
												&nbsp;Back to workouts
											</Link>
										</div>
										<div className="column is-half has-text-right">
											<Link
												class="button is-success is-hidden-touch"
												onClick={onSubmitClick}
												disabled={
													!(name && description && selectedWorkouts.length)
												}>
												<FontAwesomeIcon icon={faPlus} />
												&nbsp;Submit
											</Link>
											<Link
												class="button is-success is-fullwidth is-hidden-desktop"
												onClick={onSubmitClick}
												disabled={
													!(name && description && selectedWorkouts.length)
												}>
												<FontAwesomeIcon icon={faPlus} />
												&nbsp;Submit
											</Link>
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

export default MemberWorkoutAdd;
