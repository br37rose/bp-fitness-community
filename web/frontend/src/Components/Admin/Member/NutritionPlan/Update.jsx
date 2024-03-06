import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTasks, faTachometer, faPlus, faTimesCircle, faCheckCircle, faUserCircle, faGauge, faPencil, faLeaf, faIdCard, faAddressBook, faMessage, faChartPie, faCogs, faEye } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import { getNutritionPlanDetailAPI, putNutritionPlanUpdateAPI } from "../../../../API/NutritionPlan";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormInputField from "../../../Reusable/FormInputField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormAlternateDateField from "../../../Reusable/FormAlternateDateField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormDuelSelectField from "../../../Reusable/FormDuelSelectField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../../../AppState";
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
    WORKOUT_PREFERENCE_OPTIONS,
    MEALS_PER_DAY_WITH_EMPTY_OPTIONS,
    CONSUME_FREQUENCY_WITH_EMPTY_OPTIONS,
    NUTRITIONAL_GOAL_WITH_EMPTY_OPTIONS
} from "../../../../Constants/FieldOptions";
import {
    GENDER_OTHER, GENDER_MALE, GENDER_FEMALE,
    PHYSICAL_ACTIVITY_SEDENTARY, PHYSICAL_ACTIVITY_LIGHTLY_ACTIVE, PHYSICAL_ACTIVITY_MODERATELY_ACTIVE, PHYSICAL_ACTIVITY_VERY_ACTIVE,
    WORKOUT_INTENSITY_LOW, WORKOUT_INTENSITY_MEDIUM, WORKOUT_INTENSITY_HIGH,
} from "../../../../Constants/App";


function AdminNutritionPlanUpdate() {
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
    const [hasAllergies, setHasAllergies] = useState(0);
    const [allergies, setAllergies] = useState("");
    const [birthday, setBirthday] = useState(null);
    const [heightFeet, setHeightFeet] = useState(-1);
    const [heightInches, setHeightInches] = useState(-1);
    const [weight, setWeight] = useState(0);
    const [gender, setGender] = useState(0);
    const [genderOther, setGenderOther] = useState("");
    const [idealWeight, setIdealWeight] = useState(0);
    const [physicalActivity, setPhysicalActivity] = useState(0);
    const [workoutIntensity, setWorkoutIntensity] = useState(0);
    const [mealsPerDay, setMealsPerDay] = useState(0);
    const [consumeJunkFood, setConsumeJunkFood] = useState(0);
    const [consumeFruitsAndVegetables, setConsumeFruitsAndVegetables] = useState(0);
    const [hasIntermittentFasting, setHasIntermittentFasting] = useState(0);
    const [maxWeeks, setMaxWeeks] = useState(0);
    const [goals, setGoals] = useState([]);

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
            has_allergies: hasAllergies,
            allergies: allergies,
            meals_per_day: mealsPerDay,
            consume_junk_food: consumeJunkFood,
            consume_fruits_and_vegetables: consumeFruitsAndVegetables,
            birthday: birthday,
            height_feet: heightFeet,
            height_inches: heightInches,
            weight: weight,
            gender: gender,
            gender_other: genderOther,
            ideal_weight: idealWeight,
            physical_activity: physicalActivity,
            workout_intensity: workoutIntensity,
            time_per_day: hasIntermittentFasting,
            max_weeks: maxWeeks,
            goals: goals,
            has_intermittent_fasting: hasIntermittentFasting,
        };
        console.log("onSubmitClick, decamelizedData:", decamelizedData);
        putNutritionPlanUpdateAPI(
            decamelizedData,
            onAdminNutritionPlanUpdateSuccess,
            onAdminNutritionPlanUpdateError,
            onAdminNutritionPlanUpdateDone);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onNutritionPlanDetailSuccess(response){
        console.log("onNutritionPlanDetailSuccess: Starting...");
        setName(response.name);
        setHasAllergies(response.hasAllergies);
        setAllergies(response.allergies);
        setBirthday(response.birthday);
        setHeightFeet(response.heightFeet);
        setHeightInches(response.heightInches);
        setWeight(response.weight);
        setGender(response.gender);
        setGenderOther(response.genderOther);
        setIdealWeight(response.idealWeight);
        setPhysicalActivity(response.physicalActivity);
        setWorkoutIntensity(response.workoutIntensity);
        setMealsPerDay(response.mealsPerDay);
        setConsumeJunkFood(response.consumeJunkFood);
        setConsumeFruitsAndVegetables(response.consumeFruitsAndVegetables);
        setHasIntermittentFasting(response.hasIntermittentFasting);
        setMaxWeeks(response.maxWeeks);
        setGoals(response.goals);
    }

    function onNutritionPlanDetailError(apiErr) {
        console.log("onNutritionPlanDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onNutritionPlanDetailDone() {
        console.log("onNutritionPlanDetailDone: Starting...");
        setFetching(false);
    }

    // --- Update --- //

    function onAdminNutritionPlanUpdateSuccess(response){
        // For debugging purposes only.
        console.log("onAdminNutritionPlanUpdateSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("AdminNutrition plan update");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAdminNutritionPlanUpdateSuccess: Delayed for 2 seconds.");
            console.log("onAdminNutritionPlanUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Redirect the user to a new page.
        setForceURL("/nutrition-plan/"+response.id);
    }

    function onAdminNutritionPlanUpdateError(apiErr) {
        console.log("onAdminNutritionPlanUpdateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAdminNutritionPlanUpdateError: Delayed for 2 seconds.");
            console.log("onAdminNutritionPlanUpdateError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAdminNutritionPlanUpdateDone() {
        console.log("onAdminNutritionPlanUpdateDone: Starting...");
        setFetching(false);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);  // Start the page at the top of the page.
            setFetching(true);
            getNutritionPlanDetailAPI(
                id,
                onNutritionPlanDetailSuccess,
                onNutritionPlanDetailError,
                onNutritionPlanDetailDone
            );
        }

        return () => { mounted = false; }
    }, [id]);
    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL}  />
    }

    return (
        <>
            <div class="container">
                <section class="section">
                    <nav class="breadcrumb" aria-label="breadcrumbs">
                        <ul>
                            <li class=""><Link to="/dashboard" aria-current="page"><FontAwesomeIcon className="fas" icon={faGauge} />&nbsp;Dashboard</Link></li>
                            <li class=""><Link to="/nutrition-plans" aria-current="page"><FontAwesomeIcon className="fas" icon={faLeaf} />&nbsp;AdminNutrition Plans</Link></li>
                            <li class=""><Link to={`/nutrition-plan/${id}`} aria-current="page"><FontAwesomeIcon className="fas" icon={faEye} />&nbsp;Detail</Link></li>
                            <li class="is-active"><Link aria-current="page"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit & Re-request</Link></li>
                        </ul>
                    </nav>
                    <nav class="box">
                        <div class={`modal ${showCancelWarning ? 'is-active' : ''}`}>
                            <div class="modal-background"></div>
                            <div class="modal-card">
                                <header class="modal-card-head">
                                    <p class="modal-card-title">Are you sure?</p>
                                    <button class="delete" aria-label="close" onClick={(e)=>setShowCancelWarning(false)}></button>
                                </header>
                                <section class="modal-card-body">
                                    Your form input will be lost and nothing will be updated. This cannot be undone. Do you want to continue?
                                </section>
                                <footer class="modal-card-foot">
                                    <Link class="button is-medium is-success" to={`/nutrition-plans`}>Yes</Link>
                                    <button class="button is-medium" onClick={(e)=>setShowCancelWarning(false)}>No</button>
                                </footer>
                            </div>
                        </div>

                        <p class="title is-4"><FontAwesomeIcon className="fas" icon={faPencil} />&nbsp;Edit & Re-request an Updated AdminNutrition Plan</p>
                        <FormErrorBox errors={errors} />

                        {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

                        {isFetching && <PageLoadingContent displayMessage={"Please wait..."} />}

                        <div class="container">

                            <p class="title is-6">META</p>
                            <hr />

                            <FormInputField
                                label="Name (Optional)"
                                name="name"
                                placeholder="Text input"
                                value={name}
                                errorText={errors && errors.name}
                                helpText="Give this nutrition plan a name you can use to keep track for your own purposes. Ex: `My Cardio-Plan`."
                                onChange={(e)=>setName(e.target.value)}
                                isRequired={true}
                                maxWidth="380px"
                            />

                            <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;PERSONAL DETAILS</p>
                            <hr />

                            <FormAlternateDateField
                                label="Birthday"
                                name="birthday"
                                placeholder="Text input"
                                value={birthday}
                                helpText=""
                                onChange={(date)=>setBirthday(date)}
                                errorText={errors && errors.birthday}
                                isRequired={true}
                                maxWidth="180px"
                                maxDate={new Date()}
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
                                oneMaxWidth={{maxWidth:"100px"}}
                                twoLabel="Height"
                                twoName="heightInches"
                                twoPlaceholder="Pick"
                                twoSelectedValue={heightInches}
                                twoErrorText={errors && errors.heightInches}
                                twoOnChange={(e) => setHeightInches(parseInt(e.target.value))}
                                twoOptions={INCHES_WITH_EMPTY_OPTIONS}
                                twoDisabled={false}
                                twoMaxWidth={{maxWidth:"100px"}}
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
                                onChange={(e)=>setWeight(parseFloat(e.target.value))}
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
                                    onChange={(e)=>setGenderOther(e.target.value)}
                                    isRequired={true}
                                    maxWidth="380px"
                                />
                            }

                            <FormInputField
                                label="What is your ideal weight for your nutrition goal?"
                                type="number"
                                name="idealWeight"
                                placeholder="lbs"
                                value={idealWeight}
                                errorText={errors && errors.idealWeight}
                                helpText="lbs"
                                onChange={(e)=>setIdealWeight(parseFloat(e.target.value))}
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

                            <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;ALLERGIES</p>
                            <hr />

                            <FormRadioField
                                label="Do you have any allergies?"
                                name="hasAllergies"
                                placeholder="Pick"
                                value={hasAllergies}
                                opt1Value={1}
                                opt1Label="Yes"
                                opt2Value={2}
                                opt2Label="No"
                                errorText={errors && errors.hasAllergies}
                                onChange={(e) => setHasAllergies(parseInt(e.target.value))}
                                maxWidth="180px"
                                disabled={false}
                            />

                            {hasAllergies === 1 &&
                                <>
                                    <FormTextareaField
                                        label="If yes, what are your allergies?"
                                        type="allergies"
                                        name="allergies"
                                        placeholder=""
                                        value={allergies}
                                        errorText={errors && errors.allergies}
                                        helpText=""
                                        onChange={(e)=>setAllergies(e.target.value)}
                                        isRequired={true}
                                        maxWidth="100px"
                                    />
                                </>
                            }


                            <p class="title is-6 pt-5"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;GOAL(S) FOR NUTRITION PLAN</p>
                            <hr />

                            <FormSelectField
                                label="How many meals do you typically eat in a day?"
                                name="mealsPerDay"
                                placeholder="Pick"
                                selectedValue={mealsPerDay}
                                errorText={errors && errors.mealsPerDay}
                                helpText=""
                                onChange={(e) => setMealsPerDay(parseInt(e.target.value))}
                                options={MEALS_PER_DAY_WITH_EMPTY_OPTIONS}
                                disabled={false}
                            />

                            <FormSelectField
                                label="How often do you consume fast food or junk food?"
                                name="consumeJunkFood"
                                placeholder="Pick"
                                selectedValue={consumeJunkFood}
                                errorText={errors && errors.consumeJunkFood}
                                helpText=""
                                onChange={(e) => setConsumeJunkFood(parseInt(e.target.value))}
                                options={CONSUME_FREQUENCY_WITH_EMPTY_OPTIONS}
                                disabled={false}
                            />

                            <FormSelectField
                                label="How often do you consume fruits and/or vegetables?"
                                name="consumeFruitsAndVegetables"
                                placeholder="Pick"
                                selectedValue={consumeFruitsAndVegetables}
                                errorText={errors && errors.consumeFruitsAndVegetables}
                                helpText=""
                                onChange={(e) => setConsumeFruitsAndVegetables(parseInt(e.target.value))}
                                options={CONSUME_FREQUENCY_WITH_EMPTY_OPTIONS}
                                disabled={false}
                            />

                            <FormMultiSelectField
                                label="Enter your nutritional goal(s)"
                                name="goals"
                                placeholder="Text input"
                                options={NUTRITIONAL_GOAL_WITH_EMPTY_OPTIONS}
                                selectedValues={goals}
                                onChange={(e)=>{
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

                            <FormSelectField
                                label="How long would you like to commit to this nutritional plan?"
                                name="maxWeeks"
                                placeholder="Pick"
                                selectedValue={maxWeeks}
                                errorText={errors && errors.maxWeeks}
                                helpText=""
                                onChange={(e) => setMaxWeeks(parseInt(e.target.value))}
                                options={MAX_WEEK_WITH_EMPTY_OPTIONS}
                                disabled={false}
                            />

                            <FormRadioField
                                label="Has Intermittent Fasting"
                                name="hasIntermittentFasting"
                                placeholder="Pick"
                                value={hasIntermittentFasting}
                                opt1Value={1}
                                opt1Label="Yes"
                                opt2Value={2}
                                opt2Label="No"
                                errorText={errors && errors.hasIntermittentFasting}
                                onChange={(e) => setHasIntermittentFasting(parseInt(e.target.value))}
                                maxWidth="180px"
                                disabled={false}
                            />

                            <div class="columns pt-5">
                                <div class="column is-half">
                                    <button class="button is-medium is-fullwidth-mobile" onClick={(e)=>setShowCancelWarning(true)}><FontAwesomeIcon className="fas" icon={faTimesCircle} />&nbsp;Cancel</button>
                                </div>
                                <div class="column is-half has-text-right">
                                    <button class="button is-medium is-primary is-fullwidth-mobile" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save & Submit to Team</button>
                                </div>
                            </div>

                        </div>
                    </nav>
                </section>
            </div>
        </>
    );
}

export default AdminNutritionPlanUpdate;
