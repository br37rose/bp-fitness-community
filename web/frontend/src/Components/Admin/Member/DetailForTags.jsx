import React, { useState, useEffect } from "react";
import { Link, useParams, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faUsers,
  faTrash,
  faEye,
  faArrowRight,
  faTable,
  faRepeat,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faKey,
  faTag,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getMemberDetailAPI, putMemberUpdateAPI } from "../../../API/member";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import {
  topAlertMessageState,
  topAlertStatusState,
  datumState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
  SUBSCRIPTION_STATUS_WITH_EMPTY_OPTIONS,
  SUBSCRIPTION_TINTERVAL_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import FormTextRow from "../../Reusable/FormTextRow";
import FormTextTagRow from "../../Reusable/FormTextTagRow";
import FormTextYesNoRow from "../../Reusable/FormTextYesNoRow";
import FormTextOptionRow from "../../Reusable/FormTextOptionRow";

function AdminMemberTagList() {
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
  const [datum, setDatum] = useState({});

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
    const modifiedMember = { ...datum };
    if (
      !modifiedMember.hasOwnProperty("tags") ||
      modifiedMember.tags === undefined ||
      modifiedMember.tags === null ||
      modifiedMember.tags === ""
    ) {
      modifiedMember.tags = [];
    }

    const tag = {
      text: text,
      description: description,
    };
    modifiedMember["tags"].push(tag);

    setFetching(true);
    setErrors({});

    // Make API call.
    onUpdateMember(
      modifiedMember,
      onAccountUpdateSuccess,
      onAccountUpdateError,
      onAccountUpdateDone
    );
  };

  // Event to fire when user picks the tag from the row to load up the `Edit Modal`.
  const onSetSelectedTagForEdit = (tag) => {
    setSelectedTagForEdit(tag);
    setText(tag.text);
    setDescription(tag.description);
  };

  // Event to fire when user closes the 'Edit Modal'.
  const onDesetSelectedTagForEdit = () => {
    setSelectedTagForEdit(null);
    setText("");
    setDescription("");
  };

  const onEditButtonClick = () => {
    console.log("onEditButtonClick: Beginning...");

    // Create a copy of our current logged in user and add in `tags` field
    // if it was not previously created.
    let modifiedMember = { ...datum };

    // Use the ES6 `map` function which will iterate through all the tags
    // and search for the `selectedTagForEdit` tag and if found then we
    // will return an edited version of tag, else we will return an non-
    // modified version.
    const updatedTags = modifiedMember.tags.map((tag) => {
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
    modifiedMember.tags = updatedTags;

    // For debugging purposes only.
    console.log("Modified Current User:", modifiedMember);

    // Make API call.
    onUpdateMember(
      modifiedMember,
      onAccountUpdateSuccess,
      onAccountUpdateError,
      onAccountUpdateDone
    );
  };

  const onDeleteConfirmButtonClick = () => {
    // Create a copy of our current logged in user and add in `tags` field
    // if it was not previously created.
    let modifiedMember = { ...datum };

    // Use the ES6 `filter` function to create a new array that contains
    // all tags except the one to be deleted.
    const updatedTags = modifiedMember.tags.filter(
      (tag) => tag.id !== selectedTagForDelete.id
    );

    // Update the user account.
    modifiedMember.tags = updatedTags;

    // For debugging purposes only.
    console.log("Modified Current User:", modifiedMember);

    // Make API call.
    onUpdateMember(
      modifiedMember,
      onAccountUpdateSuccess,
      onAccountUpdateError,
      onAccountUpdateDone
    );
  };

  const onUpdateMember = (member) => {
    const decamelizedData = {
      id: member.id,
      organization_id: member.organizationId,
      first_name: member.firstName,
      last_name: member.lastName,
      email: member.email,
      phone: member.phone,
      postal_code: member.postalCode,
      address_line_1: member.addressLine1,
      address_line_2: member.addressLine2,
      city: member.city,
      region: member.region,
      country: member.country,
      status: member.status,
      password: member.password,
      password_repeated: member.passwordRepeated,
      how_did_you_hear_about_us: member.howDidYouHearAboutUs,
      how_did_you_hear_about_us_other: member.howDidYouHearAboutUsOther,
      agree_promotions_email: member.agreePromotionsEmail,
      tags: member.tags,
    };
    putMemberUpdateAPI(
      decamelizedData,
      onAccountUpdateSuccess,
      onAccountUpdateError,
      onAccountUpdateDone
    );
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onAccountDetailSuccess(response) {
    console.log("onAccountDetailSuccess: Starting...");
    setDatum(response);
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
      console.log(
        "onAccountUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
      setTopAlertMessage("");
    }, 2000);

    // Update our current user.
    setDatum(response);

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
      console.log(
        "onAccountUpdateError: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
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
      window.scrollTo(0, 0); // Start the page at the top of the page.
      setFetching(true);
      setErrors({});
      getMemberDetailAPI(
        id,
        onAccountDetailSuccess,
        onAccountDetailError,
        onAccountDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, []);

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
              <li className="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/members" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Members
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/member/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
              <li className="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faTag} />
                  &nbsp;Tags
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/members" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Members
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <nav class="">
            {/* Create modal */}
            <div class={`modal ${showAddModal ? "is-active" : ""}`}>
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">New Tag</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={(e, ses) => setShowAddModal(false)}
                  ></button>
                </header>
                <section class="modal-card-body">
                  <FormErrorBox errors={errors} />
                  <p class="pb-4">
                    Please fill out all the required fields before submitting
                    this form.
                  </p>
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
                  <button class="button is-success" onClick={onAddButtonClick}>
                    Save Tag
                  </button>
                  <button
                    class="button"
                    onClick={(e, ses) => setShowAddModal(false)}
                  >
                    Cancel
                  </button>
                </footer>
              </div>
            </div>
            {/* Update modal */}
            <div
              class={`modal ${selectedTagForEdit !== null ? "is-active" : ""}`}
            >
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">Update Tag</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={onDesetSelectedTagForEdit}
                  ></button>
                </header>
                <section class="modal-card-body">
                  <FormErrorBox errors={errors} />
                  <p class="pb-4">
                    Please fill out all the required fields before submitting
                    this form.
                  </p>
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
                  <button class="button is-success" onClick={onEditButtonClick}>
                    Save Tag
                  </button>
                  <button class="button" onClick={onDesetSelectedTagForEdit}>
                    Cancel
                  </button>
                </footer>
              </div>
            </div>
            {/* Delete modal */}
            <div
              class={`modal ${
                selectedTagForDelete !== null ? "is-active" : ""
              }`}
            >
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">Are you sure?</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={(e, ses) => setSelectedTagForDelete(null)}
                  ></button>
                </header>
                <section class="modal-card-body">
                  You are about to delete this tag and all the data associated
                  with it. This action is cannot be undone. Are you sure you
                  would like to continue?
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
                    onClick={(e, ses) => setSelectedTagForDelete(null)}
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
                    <FontAwesomeIcon className="fas" icon={faUsers} />
                    &nbsp;Member
                  </p>
                </div>
                <div class="column has-text-right">
                  <Link
                    onClick={(e) => setShowAddModal(true)}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;New Tag
                  </Link>
                </div>
              </div>
            )}

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                {datum && (
                  <div class="container">
                    {/* Tab Navigation */}
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/admin/member/${id}`}>Detail</Link>
                        </li>
                        <li class="is-active">
                          <Link>
                            <strong>Tags</strong>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/member/${datum.id}/profile`}>
                            Profile
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/member/${datum.id}/fitness-plans`}>
                            Fitness Plans
                          </Link>
                        </li>
                        <li>
                          <Link
                            to={`/admin/member/${datum.id}/nutrition-plans`}
                          >
                            Nutrition Plans
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="title is-6">
                      <FontAwesomeIcon className="fas" icon={faTable} />
                      &nbsp;List
                    </p>
                    <hr />

                    {datum.tags !== undefined &&
                    datum.tags !== null &&
                    datum.tags !== "" &&
                    datum.tags.length > 0 ? (
                      <>
                        {/* Non-Empty List */}
                        {datum.tags.map((object, i) => (
                          <TagRow
                            obj={object}
                            index={i}
                            setSelectedTagForDelete={setSelectedTagForDelete}
                            onSetSelectedTagForEdit={onSetSelectedTagForEdit}
                          />
                        ))}
                      </>
                    ) : (
                      <>
                        {/* Empty list */}
                        <section className="hero is-medium has-background-white-ter">
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
                      </>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-medium is-fullwidth-mobile"
                          to={"/dashboard"}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Dashboard
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          onClick={(e) => setShowAddModal(true)}
                          class="button is-medium is-success is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;New Tag
                        </Link>
                      </div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

function TagRow(props) {
  const { obj, index, setSelectedTagForDelete, onSetSelectedTagForEdit } =
    props;
  if (obj === undefined) {
    return null;
  }
  const { text, description } = obj;
  return (
    <div class="box">
      {/*
                ##################################################################
                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                ##################################################################
            */}
      <div class="is-hidden-touch">
        <div className="columns">
          <div className="column is-half">
            <div class="is-pulled-left">
              <p class="title is-5">{text}</p>
              <p class="subtitle is-6">{description}</p>
            </div>
          </div>
          <div className="column is-half">
            <div class="is-pulled-right">
              <Link
                onClick={(e, o) => onSetSelectedTagForEdit(obj)}
                class="button is-warning"
              >
                <FontAwesomeIcon className="mdi" icon={faPencil} />
                &nbsp;Edit
              </Link>
              &nbsp;&nbsp;
              <Link
                onClick={(e, o) => setSelectedTagForDelete(obj)}
                class="button is-danger"
              >
                <FontAwesomeIcon className="mdi" icon={faTrash} />
                &nbsp;Delete
              </Link>
              &nbsp;&nbsp;
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
            <Link
              onClick={(e, o) => onSetSelectedTagForEdit(obj)}
              class="button is-warning is-small is-fullwidth-mobile"
            >
              <FontAwesomeIcon className="mdi" icon={faPencil} />
              &nbsp;Edit
            </Link>
            &nbsp;&nbsp;
          </div>
        </div>
        {/* Mobile only */}
        <div class="is-hidden-tablet pt-2">
          <div class="columns is-mobile pt-2">
            <div class="column">
              <Link
                onClick={(e, o) => onSetSelectedTagForEdit(obj)}
                class="button is-warning is-small is-fullwidth-mobile"
              >
                <FontAwesomeIcon className="mdi" icon={faPencil} />
                &nbsp;Edit
              </Link>
              &nbsp;&nbsp;
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AdminMemberTagList;
