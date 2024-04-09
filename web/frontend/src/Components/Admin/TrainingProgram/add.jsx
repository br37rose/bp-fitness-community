import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faImage,
  faPlus,
  faArrowLeft,
  faGauge,
  faVideoCamera,
  faEye,
  faDumbbell,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { postVideoCollectionCreateAPI } from "../../../API/VideoCollection";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormAttachmentField from "../../Reusable/FormAttachmentField";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormSelectFieldForVideoCategory from "../../Reusable/FormSelectFieldForVideoCategory";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import { VIDEO_COLLECTION_TYPE_OPTIONS_WITH_EMPTY_OPTION } from "../../../Constants/FieldOptions";
import {
  EXERCISE_THUMBNAIL_TYPE_SIMPLE_STORAGE_SERVER,
  EXERCISE_THUMBNAIL_TYPE_EXTERNAL_URL,
  VIDEO_COLLECTION_TYPE_MANY_VIDEOS,
  TRAINING_PROGRAM_TYPE_PHASED,
  TRAINING_PROGRAM_TYPE_WORKOUTS,
} from "../../../Constants/App";

function AdminTrainingProgramAdd() {
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
  const [thumbnailType, setThumbnailType] = useState(0);
  const [thumbnailURL, setThumbnailURL] = useState("");
  const [thumbnailAttachmentID, setThumbnailAttachmentID] = useState("");
  const [thumbnailAttachmentName, setThumbnailAttachmentName] = useState("");
  const [alternateName, setAlternateName] = useState("");
  const [name, setName] = useState("");
  const [phases, setphases] = useState(4);
  const [weeks, setweeks] = useState(4);
  const [summary, setSummary] = useState("");
  const [description, setDescription] = useState("");
  const [type, setType] = useState(VIDEO_COLLECTION_TYPE_MANY_VIDEOS);
  const [videoCategoryID, setVideoCategoryID] = useState("");
  const [isVideoCategoryOther, setIsVideoCategoryOther] = useState("");
  const [programType, setprogramType] = useState(TRAINING_PROGRAM_TYPE_PHASED);

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Starting...");
    setFetching(true);
    setErrors({});
    postVideoCollectionCreateAPI(
      {
        type: type,
        thumbnail_type: thumbnailType,
        thumbnail_upload: thumbnailAttachmentID,
        thumbnail_url: thumbnailURL,
        name: name,
        summary: summary,
        description: description,
        video_category_id: videoCategoryID,
        is_video_category_other: isVideoCategoryOther,
        category_id: videoCategoryID,
      },
      onAddSuccess,
      onAddError,
      onAddDone
    );
    console.log("onSubmitClick: Finished.");
  };

  ////
  //// API.
  ////

  // --- VideoCollection Create --- //

  function onAddSuccess(response) {
    // For debugging purposes only.
    console.log("onAddSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Video collection created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAddSuccess: Delayed for 2 seconds.");
      console.log(
        "onAddSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the organization to the organization attachments page.
    setForceURL("/admin/video-collection/" + response.id + "");
  }

  function onAddError(apiErr) {
    console.log("onAddError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAddError: Delayed for 2 seconds.");
      console.log(
        "onAddError: topAlertMessage, topAlertStatus:",
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

  function onAddDone() {
    console.log("onAddDone: Starting...");
    setFetching(false);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
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

  const isThumbnailUploaded =
    thumbnailAttachmentID !== "" || thumbnailURL !== "";

  return (
    <>
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/training-program" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Training Program
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPlus} />
                  &nbsp;New
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/training-program" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Training programs
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modal */}
          {/* Nothing ... */}

          {/* Page */}
          <nav class="box">
            <div class="columns">
              <div class="column">
                <p class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faPlus} />
                  &nbsp;Add Training Programs
                </p>
              </div>
              <div class="column has-text-right"></div>
            </div>
            <FormErrorBox errors={errors} />

            <p class="pb-4 mb-5 has-text-grey">
              Please fill out all the required fields before submitting this
              form.
            </p>

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                <div class="container">
                  <FormRadioField
                    label="Program Type"
                    name="programtype"
                    value={programType}
                    opt1Value={TRAINING_PROGRAM_TYPE_PHASED}
                    opt1Label="Phased"
                    opt2Value={TRAINING_PROGRAM_TYPE_WORKOUTS}
                    opt2Label="Workouts"
                    errorText={errors && errors.programType}
                    onChange={(e) => setprogramType(parseInt(e.target.value))}
                    maxWidth="180px"
                    disabled={false}
                  />

                  {programType == TRAINING_PROGRAM_TYPE_PHASED && (
                    <>
                      <p class="subtitle is-6 mt-5">
                        <FontAwesomeIcon className="fas" icon={faEye} />
                        &nbsp;Information
                      </p>
                      <hr />
                      <div className="columns">
                        <div className="column is-justify-content-center">
                          <FormInputField
                            label="Phases"
                            name="phases"
                            placeholder="phases"
                            value={phases}
                            errorText={errors && errors.phases}
                            onChange={(e) => setphases(e.target.value)}
                            isRequired={true}
                            maxWidth="80px"
                          />
                        </div>
                        <div className="column">
                          <FormInputField
                            label="Weeks"
                            name="weeks"
                            placeholder="weeks"
                            value={weeks}
                            errorText={errors && errors.weeks}
                            onChange={(e) => setweeks(e.target.value)}
                            isRequired={true}
                            maxWidth="80px"
                          />
                        </div>
                        <div className="column">
                          <FormInputField
                            label="Duration"
                            name="duration"
                            value={phases * weeks + " weeks"}
                            disabled
                            maxWidth="180px"
                          />
                        </div>
                      </div>
                    </>
                  )}
                  <FormTextareaField
                    label="Name"
                    name="Name"
                    placeholder="name"
                    value={name}
                    errorText={errors && errors.name}
                    onChange={(e) => setName(e.target.value)}
                    isRequired={true}
                    maxWidth="380px"
                  />
                  <FormTextareaField
                    label="Description"
                    name="description"
                    placeholder="Description input"
                    value={description}
                    errorText={errors && errors.description}
                    helpText=""
                    onChange={(e) => setDescription(e.target.value)}
                    isRequired={true}
                    maxWidth="380px"
                  />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        class="button is-fullwidth-mobile"
                        to={`/admin/video-collections`}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Video Collections
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        onClick={onSubmitClick}
                        class="button is-success is-fullwidth-mobile"
                        type="button"
                      >
                        <FontAwesomeIcon className="fas" icon={faPlus} />
                        &nbsp;Submit
                      </button>
                    </div>
                  </div>
                </div>
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminTrainingProgramAdd;
