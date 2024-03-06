import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from 'react-scroll';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEllipsis, faTrash, faEye, faArrowRight, faTable, faRepeat, faTasks, faTachometer, faPlus, faArrowLeft, faCheckCircle, faUserCircle, faGauge, faPencil, faIdCard, faAddressBook, faContactCard, faChartPie, faKey } from '@fortawesome/free-solid-svg-icons'
import { useRecoilState } from 'recoil';

import deepClone from "../../Helpers/deepCloneUtility";
import { getAccountDetailAPI, putAccountUpdateAPI } from "../../API/Account";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormTextareaField from "../Reusable/FormTextareaField";
import FormRadioField from "../Reusable/FormRadioField";
import FormMultiSelectField from "../Reusable/FormMultiSelectField";
import FormSelectField from "../Reusable/FormSelectField";
import FormCheckboxField from "../Reusable/FormCheckboxField";
import FormCountryField from "../Reusable/FormCountryField";
import FormRegionField from "../Reusable/FormRegionField";
import { topAlertMessageState, topAlertStatusState, currentUserState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS, SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS } from "../../Constants/FieldOptions";
import FormTextRow from "../Reusable/FormTextRow";
import FormTextTagRow from "../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../Reusable/FormTextOptionRow";
import Layout from "../Menu/Layout";


function AccountTagList() {
    ////
    ////
    ////

    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);
    const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [showAddModal, setShowAddModal] = useState(false);
    const [text, setText] = useState();
    const [description, setDescription] = useState();
    const [selectedTagForDelete, setSelectedTagForDelete] = useState(null);
    const [selectedTagForEdit, setSelectedTagForEdit] = useState(null);

    ////
    //// Event handling.
    ////

    const onAddButtonClick = () => {
        // Create a copy of our current logged in user and add in `tags` field
        // if it was not previously created.
        const modifiedCurrentUser = deepClone(currentUser);
        if (!modifiedCurrentUser.hasOwnProperty('tags') || modifiedCurrentUser.tags === undefined || modifiedCurrentUser.tags === null || modifiedCurrentUser.tags === "") {
            modifiedCurrentUser.tags = [];
        }

        const tag = {
            text: text,
            description: description,
        };
        modifiedCurrentUser["tags"].push(tag);

        setFetching(true);
        setErrors({});

        putAccountUpdateAPI(modifiedCurrentUser, onAccountUpdateSuccess, onAccountUpdateError, onAccountUpdateDone);
    }

    // Event to fire when user picks the tag from the row to load up the `Edit Modal`.
    const onSetSelectedTagForEdit = (tag) => {
        setSelectedTagForEdit(tag);
        setText(tag.text);
        setDescription(tag.description);
    }

    // Event to fire when user closes the 'Edit Modal'.
    const onDesetSelectedTagForEdit = () => {
        setSelectedTagForEdit(null);
        setText("");
        setDescription("");
    }

    const onEditButtonClick = () => {
        console.log("onEditButtonClick: Beginning...");

        // Create a copy of our current logged in user and add in `tags` field
        // if it was not previously created.
        let modifiedCurrentUser = { ...currentUser };

        // Use the ES6 `map` function which will iterate through all the tags
        // and search for the `selectedTagForEdit` tag and if found then we
        // will return an edited version of tag, else we will return an non-
        // modified version.
        const updatedTags = modifiedCurrentUser.tags.map(tag => {
            if (selectedTagForEdit.id === tag.id) {
                return {
                    ...tag,
                    text: text,
                    description: description,
                };
            }
            return tag;
        });

        // Update the user account.
        modifiedCurrentUser.tags = updatedTags;

        // For debugging purposes only.
        console.log("Modified Current User:", modifiedCurrentUser);

        // Make API call.
        putAccountUpdateAPI(modifiedCurrentUser, onAccountUpdateSuccess, onAccountUpdateError, onAccountUpdateDone);
    }

    const onDeleteConfirmButtonClick = () => {
        // Create a copy of our current logged in user and add in `tags` field
        // if it was not previously created.
        let modifiedCurrentUser = { ...currentUser };

        // Use the ES6 `filter` function to create a new array that contains
        // all tags except the one to be deleted.
        const updatedTags = modifiedCurrentUser.tags.filter(tag => tag.id !== selectedTagForDelete.id);

        // Update the user account.
        modifiedCurrentUser.tags = updatedTags;

        // For debugging purposes only.
        console.log("Modified Current User:", modifiedCurrentUser);

        // Make API call.
        putAccountUpdateAPI(modifiedCurrentUser, onAccountUpdateSuccess, onAccountUpdateError, onAccountUpdateDone);
    }

    ////
    //// API.
    ////

    // --- Detail --- //

    function onAccountDetailSuccess(response) {
        console.log("onAccountDetailSuccess: Starting...");
        setCurrentUser(response);
    }

    function onAccountDetailError(apiErr) {
        console.log("onAccountDetailError: Starting...");
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAccountDetailDone() {
        console.log("onAccountDetailDone: Starting...");
        setFetching(false);
    }

    // --- Update --- //

    function onAccountUpdateSuccess(response) {
        // For debugging purposes only.
        console.log("onAccountUpdateSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Account updated");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAccountUpdateSuccess: Delayed for 2 seconds.");
            console.log("onAccountUpdateSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // Update our current user.
        setCurrentUser(response);

        // Close all modals.
        setShowAddModal(false);
        onDesetSelectedTagForEdit();
        setSelectedTagForDelete(null);
    }

    function onAccountUpdateError(apiErr) {
        console.log("onAccountUpdateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAccountUpdateError: Delayed for 2 seconds.");
            console.log("onAccountUpdateError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onAccountUpdateDone() {
        console.log("onAccountUpdateDone: Starting...");
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
            setErrors({});
            getAccountDetailAPI(
                onAccountDetailSuccess,
                onAccountDetailError,
                onAccountDetailDone
            );
        }

        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />
    }

    return (
        <div>
            {/* Modals */}
            <nav class="">
                {/* Create modal */}
                <div class={`modal ${showAddModal ? 'is-active' : ''}`}>
                    <div class="modal-background"></div>
                    <div class="modal-card">
                        <header class="modal-card-head">
                            <p class="modal-card-title">New Tag</p>
                            <button class="delete" aria-label="close" onClick={(e, ses) => setShowAddModal(false)}></button>
                        </header>
                        <section class="modal-card-body">

                            <FormErrorBox errors={errors} />
                            <p class="pb-4">Please fill out all the required fields before submitting this form.</p>
                            <FormInputField
                                label="Text"
                                name="text"
                                type="text"
                                placeholder="Input text"
                                value={text}
                                errorText={errors && errors.text}
                                helpText=""
                                onChange={(e) => setText(e.target.value)}
                                isRequired={true}
                                maxWidth="275px"
                            />

                            <FormTextareaField
                                label="Description (Optional)"
                                name="description"
                                placeholder="Description input"
                                value={description}
                                errorText={errors && errors.description}
                                helpText=""
                                onChange={(e) => setDescription(e.target.value)}
                                isRequired={true}
                                maxWidth="100%"
                                rows={2}
                            />
                        </section>
                        <footer class="modal-card-foot">
                            <button class="button is-success" onClick={onAddButtonClick}>Save Tag</button>
                            <button class="button" onClick={(e, ses) => setShowAddModal(false)}>Cancel</button>
                        </footer>
                    </div>
                </div>
                {/* Update modal */}
                <div class={`modal ${selectedTagForEdit !== null ? 'is-active' : ''}`}>
                    <div class="modal-background"></div>
                    <div class="modal-card">
                        <header class="modal-card-head">
                            <p class="modal-card-title">Update Tag</p>
                            <button class="delete" aria-label="close" onClick={onDesetSelectedTagForEdit}></button>
                        </header>
                        <section class="modal-card-body">

                            <FormErrorBox errors={errors} />
                            <p class="pb-4">Please fill out all the required fields before submitting this form.</p>
                            <FormInputField
                                label="Text"
                                name="text"
                                type="text"
                                placeholder="Input text"
                                value={text}
                                errorText={errors && errors.text}
                                helpText=""
                                onChange={(e) => setText(e.target.value)}
                                isRequired={true}
                                maxWidth="275px"
                            />

                            <FormTextareaField
                                label="Description (Optional)"
                                name="description"
                                placeholder="Description input"
                                value={description}
                                errorText={errors && errors.description}
                                helpText=""
                                onChange={(e) => setDescription(e.target.value)}
                                isRequired={true}
                                maxWidth="100%"
                                rows={2}
                            />
                        </section>
                        <footer class="modal-card-foot">
                            <button class="button is-success" onClick={onEditButtonClick}>Save Tag</button>
                            <button class="button" onClick={onDesetSelectedTagForEdit}>Cancel</button>
                        </footer>
                    </div>
                </div>
                {/* Delete modal */}
                <div class={`modal ${selectedTagForDelete !== null ? 'is-active' : ''}`}>
                    <div class="modal-background"></div>
                    <div class="modal-card">
                        <header class="modal-card-head">
                            <p class="modal-card-title">Are you sure?</p>
                            <button class="delete" aria-label="close" onClick={(e, ses) => setSelectedTagForDelete(null)}></button>
                        </header>
                        <section class="modal-card-body">
                            You are about to delete this tag and all the data associated with it. This action is cannot be undone. Are you sure you would like to continue?
                        </section>
                        <footer class="modal-card-foot">
                            <button class="button is-success" onClick={onDeleteConfirmButtonClick}>Confirm</button>
                            <button class="button" onClick={(e, ses) => setSelectedTagForDelete(null)}>Cancel</button>
                        </footer>
                    </div>
                </div>
            </nav>
            {/* Title + Options */}
            {currentUser && <div class="columns">
                <div class="column">

                </div>
                <div class="column has-text-right">
                    <Link onClick={(e) => setShowAddModal(true)} class="button is-medium is-success is-fullwidth-mobile" type="button">
                        <FontAwesomeIcon className="mdi" icon={faPlus} />&nbsp;New Tag
                    </Link>
                </div>
            </div>}
            {currentUser && <div class="container">

                {currentUser.tags !== undefined && currentUser.tags !== null && currentUser.tags !== "" && currentUser.tags.length > 0
                    ?
                    <>
                        {/* Non-Empty List */}
                        {currentUser.tags.map((object, i) => <TagRow obj={object} index={i} setSelectedTagForDelete={setSelectedTagForDelete} onSetSelectedTagForEdit={onSetSelectedTagForEdit} />)}
                    </>
                    :
                    <>
                        {/* Empty list */}
                        <section className="hero has-background-white-ter">
                            <div className="hero-body">
                                <p className="title">
                                    <FontAwesomeIcon className="fas" icon={faTable} />
                                    &nbsp;No Tags
                                </p>
                                <p className="subtitle">
                                    No tags for your account.{" "}
                                    <b>
                                        <Link onClick={(e) => setShowAddModal(true)}>
                                            Click here&nbsp;
                                            <FontAwesomeIcon
                                                className="mdi"
                                                icon={faArrowRight}
                                            />
                                        </Link>
                                    </b>{" "}
                                    to get started creating your first tag.
                                </p>
                            </div>
                        </section>
                    </>}
            </div>}
        </div>


    );
};

export default AccountTagList;

function TagRow(props) {
    const { obj, index, setSelectedTagForDelete, onSetSelectedTagForEdit } = props;
    if (obj === undefined) {
        return null;
    }
    const { text, description } = obj;
    return (
        <div className="box">
            {/*
                ##################################################################
                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                ##################################################################
            */}
            <div class="is-hidden-touch" >
                <div className="columns">
                    <div className="column is-half">
                        <div class="is-pulled-left">
                            <p class="title is-5">{text}</p>
                            <p class="subtitle is-6">{description}</p>
                        </div>
                    </div>
                    <div className="column is-half">
                        <div class="is-pulled-right">
                            <Link onClick={(e, o) => onSetSelectedTagForEdit(obj)} class="button is-warning"><FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit</Link>&nbsp;&nbsp;
                            <Link onClick={(e, o) => setSelectedTagForDelete(obj)} class="button is-danger"><FontAwesomeIcon className="mdi" icon={faTrash} />&nbsp;Delete</Link>&nbsp;&nbsp;
                        </div>
                    </div>
                </div>
            </div>

            {/*
                ###########################################################################
                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                ###########################################################################
            */}
            <div class="is-hidden-desktop mb-5">
                <strong>Text:</strong>&nbsp;{text}
                <br />
                <br />
                <strong>Description:</strong>&nbsp;{description}
                <br />
                <br />

                {/* Tablet only */}
                <div class="is-hidden-mobile pt-2">
                    <div className="buttons">
                        <Link onClick={(e, o) => onSetSelectedTagForEdit(obj)} class="button is-warning is-small is-fullwidth-mobile"><FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit</Link>&nbsp;&nbsp;
                    </div>
                </div>

                {/* Mobile only */}
                <div class="is-hidden-tablet pt-2">
                    <div class="columns is-mobile pt-2">
                        <div class="column">
                            <Link onClick={(e, o) => onSetSelectedTagForEdit(obj)} class="button is-warning is-small is-fullwidth-mobile"><FontAwesomeIcon className="mdi" icon={faPencil} />&nbsp;Edit</Link>&nbsp;&nbsp;
                        </div>
                    </div>
                </div>
            </div>
        </div>

    );
}