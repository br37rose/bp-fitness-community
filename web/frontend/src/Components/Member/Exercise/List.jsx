import React, { useState, useEffect } from 'react';
import Scroll from "react-scroll";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faDroplet, faDumbbell, faEye, faFilter, faFire, faGauge, faGenderless, faRefresh, faSearch, faTable } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';
import { useRecoilState } from "recoil";

import Layout from '../../Menu/Layout';
import FormErrorBox from "../../Reusable/FormErrorBox";
import { getExerciseListAPI } from '../../../API/Exercise';
import {
    topAlertMessageState,
    topAlertStatusState,
    currentUserState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectField from "../../Reusable/FormSelectField";
import FormMultiSelectFieldForTags from "../../Reusable/FormMultiSelectFieldForTags";
import {
    PAGE_SIZE_OPTIONS,
    EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION,
    EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION
} from "../../../Constants/FieldOptions";
import MemberExerciseListDesktop from './ListDesktop';
import MemberExerciseListMobile from './ListMobile';


////
//// Custom Component
////

const ExerciseComponent = ({ description,
    name,
    thumbnailUrl,
    thumbnailObjectUrl,
    thumbnailType,
    gender,
    id
}) => {
    return (
        <>
            {/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
            <div class="is-hidden-touch" >
                <MemberExerciseListDesktop name={name} description={description} thumbnailUrl={thumbnailUrl} thumbnailObjectUrl={thumbnailObjectUrl} thumbnailType={thumbnailType} gender={gender} id={id} />
            </div>

            {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
            <div class="is-fullwidth is-hidden-desktop">
                <MemberExerciseListMobile name={name} description={description} thumbnailUrl={thumbnailUrl} thumbnailObjectUrl={thumbnailObjectUrl} thumbnailType={thumbnailType} gender={gender} id={id} />
            </div>


        </>
    );
}


const Modal = ({ isActive, title, children, footer, onClose }) => (
    <div className={`modal ${isActive ? 'is-active' : ''}`}>
        <div className="modal-background"></div>
        <div className="modal-card">
            <header className="modal-card-head">
                <p className="modal-card-title has-text-weight-bold">{title}</p>
                <button onClick={onClose} className="delete" aria-label="close"></button>
            </header>
            <section className="modal-card-body">
                {children}
            </section>
            {footer && <footer className="modal-card-foot">{footer}</footer>}
        </div>
    </div>
);

const MemberExerciseList = () => {

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [listData, setListData] = useState("");
    const [selectedExerciseForDeletion, setSelectedExerciseForDeletion] = useState("");
    const [isFetching, setFetching] = useState(false);
    const [pageSize, setPageSize] = useState(10); // Pagination
    const [previousCursors, setPreviousCursors] = useState([]); // Pagination
    const [nextCursor, setNextCursor] = useState(""); // Pagination
    const [currentCursor, setCurrentCursor] = useState(""); // Pagination
    const [showFilter, setShowFilter] = useState(false); // Filtering + Searching
    const [sortField, setSortField] = useState("created"); // Sorting
    const [temporarySearchText, setTemporarySearchText] = useState(""); // Searching - The search field value as your writes their query.
    const [actualSearchText, setActualSearchText] = useState(""); // Searching - The actual search query value to submit to the API.
    const [category, setCategory] = useState("");
    const [movementType, setMovementType] = useState("");
    const [status, setStatus] = useState("");
    const [gender, setGender] = useState("");
    const [videoType, setVideoType] = useState("");
    const [tags, setTags] = useState([]);

    ////
    //// BREADCRUMB
    ////
    const breadcrumbItems = {
        items: [
            { text: 'Dashboard', link: '/dashboard', isActive: false, icon: faGauge },
            { text: 'Exercises', link: '#', icon: faDumbbell, isActive: true }
        ],
        mobileBackLinkItems: {
            link: "/dashboard",
            text: "Back to Dashboard",
            icon: faArrowLeft
        }
    }

    ////
    //// API.
    ////

    const onExerciseListSuccess = (response) => {
        console.log("onExerciseListSuccess: Starting...");
        if (response.results !== null) {
            setListData(response);
            if (response.hasNextPage) {
                setNextCursor(response.nextCursor); // For pagination purposes.
            }
        } else {
            setListData([]);
            setNextCursor("");
        }
    }

    const onExerciseListError = (apiErr) => {
        console.log("onExerciseListError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    const onExerciseListDone = () => {
        console.log("onExerciseListDone: Starting...");
        setFetching(false);
    }


    ////
    //// Event handling.
    ////

    // Note: currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, tags
    const fetchList = (cur, limit, keywords, cat, mt, st, g, vt, t) => {
        setFetching(true);
        setErrors({});

        let params = new Map();
        params.set("page_size", limit); // Pagination
        params.set("sort_field", "created"); // Sorting
        params.set("sort_order", -1)         // Sorting - descending, meaning most recent start date to oldest start date.

        if (cur !== "") {
            // Pagination
            params.set("cursor", cur);
        }

        // Filtering
        if (keywords !== undefined && keywords !== null && keywords !== "") {
            // Searhcing
            params.set("search", keywords);
        }
        if (cat !== undefined && cat !== null && cat !== "") {
            params.set("category", cat);
        }
        if (mt !== undefined && mt !== null && mt !== "") {
            params.set("movement_type", mt);
        }
        if (st !== undefined && st !== null && st !== "") {
            params.set("status", st);
        }
        if (g !== undefined && g !== null && g !== "") {
            params.set("gender", g);
        }
        if (vt !== undefined && vt !== null && vt !== "") {
            params.set("video_type", vt);
        }
        if (t !== undefined && t !== null && t !== "") {
            params.set("tags", t);
        }

        getExerciseListAPI(
            params,
            onExerciseListSuccess,
            onExerciseListError,
            onExerciseListDone
        );
    };

    const onNextClicked = (e) => {
        console.log("onNextClicked");
        let arr = [...previousCursors];
        arr.push(currentCursor);
        setPreviousCursors(arr);
        setCurrentCursor(nextCursor);
    };

    const onPreviousClicked = (e) => {
        console.log("onPreviousClicked");
        let arr = [...previousCursors];
        const previousCursor = arr.pop();
        setPreviousCursors(arr);
        setCurrentCursor(previousCursor);
    };

    const onSearchButtonClick = (e) => {
        // Searching
        console.log("Search button clicked...");
        setActualSearchText(temporarySearchText);
    };

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, tags);
        }

        return () => {
            mounted = false;
        };
    }, [currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, tags]);

    ////
    //// Component rendering.
    ////

    return (
        <Layout breadcrumbItems={breadcrumbItems}>
            <div class="box">
                <div className="columns">
                    <div className="column">
                        <h1 className="title is-4">
                            <FontAwesomeIcon className="fas" icon={faDumbbell} />
                            &nbsp;Exercises
                        </h1>
                    </div>
                    <div className="column has-text-right">
                        <button onClick={() => fetchList(currentCursor, pageSize, actualSearchText, category, movementType, status, gender, videoType, tags)} class="is-fullwidth-mobile button is-link is-small" type="button">
                            <FontAwesomeIcon className="mdi" icon={faRefresh} />&nbsp;Refresh
                        </button>
                        &nbsp;
                        <button onClick={(e) => setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                            <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                        </button>
                    </div>
                </div>

                {showFilter && (
                    <div class="columns has-background-white-bis" style={{ borderRadius: "15px", padding: "20px" }}>
                        <div class="column">
                            <FormInputFieldWithButton
                                label={"Search"}
                                name="temporarySearchText"
                                type="text"
                                placeholder="Search by name"
                                value={temporarySearchText}
                                helpText=""
                                onChange={(e) => setTemporarySearchText(e.target.value)}
                                isRequired={true}
                                maxWidth="100%"
                                buttonLabel={
                                    <>
                                        <FontAwesomeIcon className="fas" icon={faSearch} />
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
                                onChange={(e) => setCategory(parseInt(e.target.value))}
                                options={EXERCISE_CATEGORY_OPTIONS_WITH_EMPTY_OPTION}
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
                                onChange={(e) => setMovementType(parseInt(e.target.value))}
                                options={EXERCISE_MOMENT_TYPE_OPTIONS_WITH_EMPTY_OPTION}
                            />
                        </div>
                        <div class="column is-size-6">
                            <FormSelectField
                                label="Gender"
                                name="gender"
                                placeholder="Pick"
                                selectedValue={gender}
                                errorText={errors && errors.gender}
                                helpText=""
                                onChange={(e) => setGender(e.target.value)}
                                options={EXERCISE_GENDER_OPTIONS_WITH_EMPTY_OPTION}
                            />
                        </div>
                        <div class="column is-size-6">
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
                )}

                {isFetching ? (
                    <PageLoadingContent displayMessage={"Please wait..."} />
                ) : (
                    <>
                        <FormErrorBox errors={errors} />
                        {listData &&
                            listData.results &&
                            (listData.results.length > 0 || previousCursors.length > 0) ? (
                            <div>
                                {listData.results.map((excercise) => (
                                    <ExerciseComponent
                                        id={excercise.id}
                                        description={excercise.description}
                                        name={excercise.name}
                                        thumbnailUrl={excercise.thumbnailUrl}
                                        thumbnailObjectUrl={excercise.thumbnailObjectUrl}
                                        thumbnailType={excercise.thumbnailType}
                                        gender={excercise.gender}
                                    />
                                ))}
                                <div class="columns">
                                    <div class="column is-half">
                                        <span class="select">
                                            <select
                                                class={`input has-text-grey-light`}
                                                name="pageSize"
                                                onChange={(e) =>
                                                    setPageSize(parseInt(e.target.value))
                                                }
                                            >
                                                {PAGE_SIZE_OPTIONS.map(function (option, i) {
                                                    return (
                                                        <option
                                                            selected={pageSize === option.value}
                                                            value={option.value}
                                                        >
                                                            {option.label}
                                                        </option>
                                                    );
                                                })}
                                            </select>
                                        </span>
                                    </div>
                                    <div class="column is-half has-text-right">
                                        {previousCursors.length > 0 && (
                                            <button
                                                class="button"
                                                onClick={onPreviousClicked}
                                            >
                                                Previous
                                            </button>
                                        )}
                                        {listData.hasNextPage && (
                                            <>
                                                <button class="button" onClick={onNextClicked}>
                                                    Next
                                                </button>
                                            </>
                                        )}
                                    </div>
                                </div>
                            </div>


                        ) : (
                            <section className="hero is-medium has-background-white-ter">
                                <div className="hero-body">
                                    <p className="title">
                                        <FontAwesomeIcon className="fas" icon={faTable} />
                                        &nbsp;No Exercises
                                    </p>
                                    <p className="subtitle">
                                        No exercises found at the moment. Please check back later!
                                    </p>
                                </div>
                            </section>
                        )}
                    </>
                )}
            </div>
        </Layout >
    );
}

export default MemberExerciseList;
