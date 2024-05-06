import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTasks, faTachometer, faPlus, faTimesCircle, faCheckCircle, faUserCircle, faGauge, faPencil, faTrophy, faIdCard, faAddressBook, faMessage, faChartPie, faCogs, faEye, faArrowLeft } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getFitnessPlanDetailAPI, putFitnessPlanUpdateAPI } from "../../../API/FitnessPlan";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormRadioField from "../../Reusable/FormRadioField";
import FormInputField from "../../Reusable/FormInputField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormAlternateDateField from "../../Reusable/FormAlternateDateField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormDuelSelectField from "../../Reusable/FormDuelSelectField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../AppState";
import {
    HOME_GYM_EQUIPMENT_OPTIONS,
    FEET_WITH_EMPTY_OPTIONS,
    INCHES_WITH_EMPTY_OPTIONS,
    GENDER_WITH_EMPTY_OPTIONS,
    PHYSICAL_ACTIVITY_WITH_EMPTY_OPTIONS,
    WORKOUT_INTENSITY_WITH_EMPTY_OPTIONS,
    DAYS_PER_WEEK_WITH_EMPTY_OPTIONS,
    TIME_PER_DAY_WITH_EMPTY_OPTIONS,
    MAX_WEEK_WITH_EMPTY_OPTIONS,
    FITNESS_GOAL_OPTIONS,
    WORKOUT_PREFERENCE_OPTIONS
} from "../../../Constants/FieldOptions";
import {
    GENDER_OTHER, GENDER_MALE, GENDER_FEMALE,
    PHYSICAL_ACTIVITY_SEDENTARY, PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE, PHYSICAL_ACTIVITY_MODERATELY_ACTIVE, PHYSICAL_ACTIVITY_VERY_ACTIVE,
    WORKOUT_INTENSITY_LOW, WORKOUT_INTENSITY_MEDIUM, WORKOUT_INTENSITY_HIGH,
} from "../../../Constants/App";
import Layout from "../../Menu/Layout";


function AdminFitnessPlanUpdate() {
    ////
    //// URL Parameters.
    ////

    const { id } = useParams()

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    // Page related states.
    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [showCancelWarning, setShowCancelWarning] = useState(false);

    // Form data related states.
    const [name, setName] = useState("");
    const [equipmentAccess, setEquipmentAccess] = useState(0);
    const [homeGymEquipment, setHomeGymEquipment] = useState([]);
    const [hasWorkoutsAtHome, setHasWorkoutsAtHome] = useState(0);
    const [birthday, setBirthday] = useState(null);
    const [heightFeet, setHeightFeet] = useState(-1);
    const [heightInches, setHeightInches] = useState(-1);
    const [weight, setWeight] = useState(0);
    const [gender, setGender] = useState(0);
    const [genderOther, setGenderOther] = useState("");
    const [idealWeight, setIdealWeight] = useState(0);
    const [physicalActivity, setPhysicalActivity] = useState(0);
    const [workoutIntensity, setWorkoutIntensity] = useState(0);
    const [daysPerWeek, setDaysPerWeek] = useState(0);
    const [timePerDay, setTimePerDay] = useState(0);
    const [maxWeeks, setMaxWeeks] = useState(0);
    const [goals, setGoals] = useState([]);
    const [workoutPreferences, setWorkoutPreferences] = useState([]);

    ////
    //// Event handling.
    ////

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Beginning...");
        setFetching(true);
        setErrors({});

        // To Snake-case for API from camel-case in React.
        const decamelizedData = {
            id: id,
            name: name,
            equipment_access: equipmentAccess,
            home_gym_equipment: homeGymEquipment,
            has_workouts_at_home: hasWorkoutsAtHome,
            birthday: birthday,
            height_feet: heightFeet,
            height_inches: heightInches,
            weight: weight,
            gender: gender,
            gender_other: genderOther,
            ideal_weight: idealWeight,
            physical_activity: physicalActivity,
            workout_intensity: workoutIntensity,
            days_per_week: daysPerWeek,
            time_per_day: timePerDay,
            max_weeks: maxWeeks,
            goals: goals,
            workout_preferences: workoutPreferences,
        };
        console.log("onSubmitClick, decamelizedData:", decamelizedData);
        putFitnessPlanUpdateAPI(
            decamelizedData,
            onAdminFitnessPlanUpdateSuccess,
            onAdminFitnessPlanUpdateError,
            onAdminFitnessPlanUpdateDone);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onFitnessPlanDetailSuccess(response) {
        console.log("onFitnessPlanDetailSuccess: Starting...");
        setName(response.name);
        setEquipmentAccess(response.equipmentAccess);
        setHomeGymEquipment(response.homeGymEquipment);
        setHasWorkoutsAtHome(response.hasWorkoutsAtHome);
        setBirthday(response.birthday);
        setHeightFeet(response.heightFeet);
        setHeightInches(response.heightInches);
        setWeight(response.weight);
        setGender(response.gender);
        setGenderOther(response.genderOther);
        setIdealWeight(response.idealWeight);
        setPhysicalActivity(response.physicalActivity);
        setWorkoutIntensity(response.workoutIntensity);
        setDaysPerWeek(response.daysPerWeek);
        setTimePerDay(response.timePerDay);
        setMaxWeeks(response.maxWeeks);
        setGoals(response.goals);
        setWorkoutPreferences(response.workoutPreferences);
    }

    function onFitnessPlanDetailError(apiErr) {
        console.log("onFitnessPlanDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onFitnessPlanDetailDone() {
        console.log("onFitnessPlanDetailDone: Starting...");
        setFetching(false);
    }

    // --- Update --- //

    function onAdminFitnessPlanUpdateSuccess(response) {
        // For debugging purposes only.
        console.log("onAdminFitnessPlanUpdateSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Fitness plan update");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAdminFitnessPlanUpdateSuccess: Delayed for 2 seconds.");
            console.log("onAdminFitnessPlanUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to a new page.
        setForceURL("/fitness-plan/" + response.id);
    }

    function onAdminFitnessPlanUpdateError(apiErr) {
        console.log("onAdminFitnessPlanUpdateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAdminFitnessPlanUpdateError: Delayed for 2 seconds.");
            console.log("onAdminFitnessPlanUpdateError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAdminFitnessPlanUpdateDone() {
        console.log("onAdminFitnessPlanUpdateDone: Starting...");
        setFetching(false);
    }

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Fitness Plans', link: '/fitness-plans', icon: faTrophy, isActive: false },
            { text: 'Detail', link: `/fitness-plan/${id}`, icon: faEye, isActive: false },
            { text: 'Edit & Re-request', link: '#', icon: faPencil, isActive: true }
        ],
        mobileBackLinkItems: {
            link: `/fitness-plan/${id}`,
            text: "Back to Detail",
            icon: faArrowLeft
        }
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.
            setFetching(true);
            getFitnessPlanDetailAPI(
                id,
                onFitnessPlanDetailSuccess,
                onFitnessPlanDetailError,
                onFitnessPlanDetailDone
            );
        }

        return () => { mounted = false; }
    }, [id]);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div class="box">
                <div class={`modal ${showCancelWarning ? 'is-active' : ''}`}>
                    <div class="modal-background"></div>
                    <div class="modal-card">
                        <header class="modal-card-head">
                            <p class="modal-card-title">Are you sure?</p>
                            <button class="delete" aria-label="close" onClick={(e) => setShowCancelWarning(false)}></button>
                        </header>
                        <section class="modal-card-body">
                            Your form input will be lost and nothing will be updated. This cannot be undone. Do you want to continue?
                        </section>
                        <footer class="modal-card-foot">
                            <Link class="button is-medium is-success" to={`/fitness-plans`}>Yes</Link>
                            <button class="button is-medium" onClick={(e) => setShowCancelWarning(false)}>No</button>
                        </footer>
                    </div>
                </div>

                <p class="title is-4"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit & Re-request an Updated Fitness Plan</p>
                <FormErrorBox errors={errors} />

                {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

                {isFetching && <PageLoadingContent displayMessage={"Please wait..."} />}

                <div>

                    <p class="title is-6">META</p>
                    <hr />

                    <FormInputField
                        label="Name (Optional)"
                        name="name"
                        placeholder="Text input"
                        value={name}
                        errorText={errors && errors.name}
                        helpText=""
                        onChange={(e) => setName(e.target.value)}
                        isRequired={true}
                        maxWidth="380px"
                    />

                    <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;EQUIPMENT ACCESS</p>
                    <hr />

                    <FormRadioField
                        label="What equipment do you have access to?"
                        name="equipmentAccess"
                        placeholder="Pick"
                        value={equipmentAccess}
                        opt1Value={1}
                        opt1Label="No Equipment (calistanic/outdoor options)"
                        opt2Value={2}
                        opt2Label="Full Gym Access"
                        opt3Value={3}
                        opt3Label="Home Gym"
                        errorText={errors && errors.equipmentAccess}
                        onChange={(e) => setEquipmentAccess(parseInt(e.target.value))}
                        maxWidth="180px"
                        disabled={false}
                    />

                    <FormMultiSelectField
                        label="Please select all the home gym equipment that you have (Optional)"
                        name="homeGymEquipment"
                        placeholder="Text input"
                        options={HOME_GYM_EQUIPMENT_OPTIONS}
                        selectedValues={homeGymEquipment}
                        onChange={(e) => {
                            let values = [];
                            for (let option of e) {
                                values.push(option.value);
                            }
                            setHomeGymEquipment(values);
                        }}
                        errorText={errors && errors.homeGymEquipment}
                        helpText=""
                        isRequired={true}
                        maxWidth="640px"
                    />

                    <FormRadioField
                        label="Do you workout at home?"
                        name="hasWorkoutsAtHome"
                        placeholder="Pick"
                        value={hasWorkoutsAtHome}
                        opt1Value={1}
                        opt1Label="Yes"
                        opt2Value={2}
                        opt2Label="No"
                        errorText={errors && errors.hasWorkoutsAtHome}
                        onChange={(e) => setHasWorkoutsAtHome(parseInt(e.target.value))}
                        maxWidth="180px"
                        disabled={false}
                    />

                    <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;PERSONAL DETAILS</p>
                    <hr />

                    <FormAlternateDateField
                        label="Birthday"
                        name="birthday"
                        placeholder="Text input"
                        value={birthday}
                        helpText=""
                        onChange={(date) => setBirthday(date)}
                        errorText={errors && errors.birthday}
                        isRequired={true}
                        maxWidth="180px"
                    />

                    <FormDuelSelectField
                        label="Please enter your height?"
                        oneName="heightFeet"
                        onePlaceholder="Pick"
                        oneSelectedValue={heightFeet}
                        oneErrorText={errors && errors.heightFeet}
                        oneOnChange={(e) => setHeightFeet(parseInt(e.target.value))}
                        oneOptions={FEET_WITH_EMPTY_OPTIONS}
                        oneDisabled={false}
                        oneMaxWidth={{ maxWidth: "100px" }}
                        twoLabel="Height"
                        twoName="heightInches"
                        twoPlaceholder="Pick"
                        twoSelectedValue={heightInches}
                        twoErrorText={errors && errors.heightInches}
                        twoOnChange={(e) => setHeightInches(parseInt(e.target.value))}
                        twoOptions={INCHES_WITH_EMPTY_OPTIONS}
                        twoDisabled={false}
                        twoMaxWidth={{ maxWidth: "100px" }}
                        helpText={heightFeet > -1 && heightInches > -1 && <>(Your height is {heightFeet} ft and {heightInches} inches)</>}
                    />

                    <FormInputField
                        label="What is your current weight (lbs)?"
                        type="number"
                        name="weight"
                        placeholder="Text input"
                        value={weight}
                        errorText={errors && errors.weight}
                        helpText="lbs"
                        onChange={(e) => setWeight(parseFloat(e.target.value))}
                        isRequired={true}
                        maxWidth="80px"
                    />

                    <FormRadioField
                        label="What is your gender?"
                        name="gender"
                        placeholder="Pick"
                        value={gender}
                        opt1Value={GENDER_MALE}
                        opt1Label="Male"
                        opt2Value={GENDER_FEMALE}
                        opt2Label="Female"
                        opt3Value={GENDER_OTHER}
                        opt3Label="Other"
                        errorText={errors && errors.gender}
                        onChange={(e) => setGender(parseInt(e.target.value))}
                        maxWidth="180px"
                        disabled={false}
                    />
                    {gender === GENDER_OTHER &&
                        <FormInputField
                            label="Gender (Other)"
                            name="genderOther"
                            placeholder="Text input"
                            value={genderOther}
                            errorText={errors && errors.genderOther}
                            helpText=""
                            onChange={(e) => setGenderOther(e.target.value)}
                            isRequired={true}
                            maxWidth="380px"
                        />
                    }

                    <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;CURRENT PHYSICAL ACTIVITY</p>
                    <hr />

                    <FormInputField
                        label="What is your ideal weight for your fitness goal?"
                        type="number"
                        name="idealWeight"
                        placeholder="lbs"
                        value={idealWeight}
                        errorText={errors && errors.idealWeight}
                        helpText="lbs"
                        onChange={(e) => setIdealWeight(parseFloat(e.target.value))}
                        isRequired={true}
                        maxWidth="100px"
                    />

                    <FormRadioField
                        label="My current level of physical activity is"
                        name="physicalActivity"
                        placeholder="Pick"
                        value={physicalActivity}
                        opt1Value={PHYSICAL_ACTIVITY_SEDENTARY}
                        opt1Label="Sedentary"
                        opt2Value={PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE}
                        opt2Label="Lightly Active"
                        opt3Value={PHYSICAL_ACTIVITY_MODERATELY_ACTIVE}
                        opt3Label="Moderately Active"
                        opt4Value={PHYSICAL_ACTIVITY_VERY_ACTIVE}
                        opt4Label="Very Active"
                        errorText={errors && errors.physicalActivity}
                        onChange={(e) => setPhysicalActivity(parseInt(e.target.value))}
                        maxWidth="180px"
                        disabled={false}
                    />

                    <FormRadioField
                        label="My current intensity in my exercise routine is"
                        name="workoutIntensity"
                        placeholder="Pick"
                        value={workoutIntensity}
                        opt1Value={WORKOUT_INTENSITY_LOW}
                        opt1Label="Low"
                        opt2Value={WORKOUT_INTENSITY_MEDIUM}
                        opt2Label="Medium"
                        opt3Value={WORKOUT_INTENSITY_HIGH}
                        opt3Label="High"
                        errorText={errors && errors.workoutIntensity}
                        onChange={(e) => setWorkoutIntensity(parseInt(e.target.value))}
                        maxWidth="180px"
                        disabled={false}
                    />

                    <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;GOAL(S) FOR FITNESS PLAN</p>
                    <hr />

                    <FormSelectField
                        label="Enter the number of days per week that you can train"
                        name="daysPerWeek"
                        placeholder="Pick"
                        selectedValue={daysPerWeek}
                        errorText={errors && errors.daysPerWeek}
                        helpText=""
                        onChange={(e) => setDaysPerWeek(parseInt(e.target.value))}
                        options={DAYS_PER_WEEK_WITH_EMPTY_OPTIONS}
                        disabled={false}
                    />

                    <FormRadioField
                        label="Enter the length of time per day that you can train"
                        name="timePerDay"
                        placeholder="Pick"
                        value={timePerDay}
                        opt1Value={30}
                        opt1Label="30 mins"
                        opt2Value={60}
                        opt2Label="60 mins"
                        opt3Value={90}
                        opt3Label="90 mins"
                        errorText={errors && errors.timePerDay}
                        onChange={(e) => setTimePerDay(parseInt(e.target.value))}
                        maxWidth="180px"
                        disabled={false}
                    />

                    <FormSelectField
                        label="Enter the number of weeks that you would like your training plan to last"
                        name="maxWeeks"
                        placeholder="Pick"
                        selectedValue={maxWeeks}
                        errorText={errors && errors.maxWeeks}
                        helpText=""
                        onChange={(e) => setMaxWeeks(parseInt(e.target.value))}
                        options={MAX_WEEK_WITH_EMPTY_OPTIONS}
                        disabled={false}
                    />

                    <FormMultiSelectField
                        label="Enter your fitness goals"
                        name="goals"
                        placeholder="Text input"
                        options={FITNESS_GOAL_OPTIONS}
                        selectedValues={goals}
                        onChange={(e) => {
                            let values = [];
                            for (let option of e) {
                                values.push(option.value);
                            }
                            setGoals(values);
                        }}
                        errorText={errors && errors.goals}
                        helpText=""
                        isRequired={true}
                        maxWidth="640px"
                    />

                    <FormMultiSelectField
                        label="Enter your workout preferences"
                        name="workoutPreferences"
                        placeholder="Text input"
                        options={WORKOUT_PREFERENCE_OPTIONS}
                        selectedValues={workoutPreferences}
                        onChange={(e) => {
                            let values = [];
                            for (let option of e) {
                                values.push(option.value);
                            }
                            setWorkoutPreferences(values);
                        }}
                        errorText={errors && errors.workoutPreferences}
                        helpText=""
                        isRequired={true}
                        maxWidth="640px"
                    />

                    <div class="columns pt-5">
                        <div class="column is-half">
                        <Link
											class="button is-hidden-touch"
											onClick={(e) => setShowCancelWarning(true)}><FontAwesomeIcon className="fas" icon={faTimesCircle} />&nbsp;Cancel
										</Link>
										<Link
											class="button is-fullwidth is-hidden-desktop"
											onClick={(e) => setShowCancelWarning(true)}><FontAwesomeIcon className="fas" icon={faTimesCircle} />&nbsp;Cancel
										</Link>
                        </div>
                        <div class="column is-half has-text-right">
                        <Link
											class="button is-success is-hidden-touch"
											onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save & Submit to Team
										</Link>
										<Link
											class="button is-success is-fullwidth is-hidden-desktop"
											onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save & Submit to Team
										</Link>
                        </div>
                    </div>

                    

                </div>
            </div>
        </Layout>
    );
}

export default AdminFitnessPlanUpdate;
