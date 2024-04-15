import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlus,
  faArrowLeft,
  faGauge,
  faEye,
  faDumbbell,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
  currentUserState,
  topAlertMessageState,
  topAlertStatusState,
} from "../../../AppState";
import { postTrainingProgCreateAPI } from "../../../API/trainingProgram";
import { getMemberListOptionsAPI } from "../../../API/member";
import FormSelectField from "../../Reusable/FormSelectField";
import DataDisplayRowText from "../../Reusable/DataDisplayRowText";

function MemberTrainingProgramAdd() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentuser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [thumbnailURL, setThumbnailURL] = useState("");
  const [thumbnailAttachmentID, setThumbnailAttachmentID] = useState("");
  const [memoptions, setmemoptions] = useState([]);

  const [name, setName] = useState("");
  const [phases, setphases] = useState(4);
  const [weeks, setweeks] = useState(4);
  const [description, setDescription] = useState("");

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    setFetching(true);
    setErrors({});
    postTrainingProgCreateAPI(
      {
        name: name,
        description: description,
        phases: parseInt(phases),
        weeks: parseInt(weeks),
        organization_id: currentuser.organizationId,
        user_id: currentuser.id,
      },
      onAddSuccess,
      onAddError,
      onAddDone
    );
  };

  ////
  //// API.
  ////

  // --- VideoCollection Create --- //

  function onAddSuccess(response) {
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("program  created");
    setTopAlertStatus("success");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect the organization to the organization attachments page.
    setForceURL("/training-program/" + response.id + "");
  }

  function onAddError(apiErr) {
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAddDone() {
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
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/training-program" aria-current="page">
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
                <Link to="/training-program" aria-current="page">
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
                  <>
                    <p class="subtitle is-6 mt-5">
                      <FontAwesomeIcon className="fas" icon={faEye} />
                      &nbsp;Information
                    </p>
                    <hr />
                    <div className="columns">
                      <div className="column">
                        <DataDisplayRowText
                          label={" User"}
                          value={currentuser.name}
                        />
                      </div>
                    </div>
                    <div className="columns">
                      <div className="column is-justify-content-center">
                        <FormInputField
                          type="number"
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
                          type="number"
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
                        to={`/training-program`}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to training programs
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        onClick={onSubmitClick}
                        class="button is-success is-fullwidth-mobile"
                        type="button"
                        disabled={
                          !(
                            name &&
                            description &&
                            phases &&
                            weeks &&
                            currentuser.id
                          )
                        }
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

export default MemberTrainingProgramAdd;
