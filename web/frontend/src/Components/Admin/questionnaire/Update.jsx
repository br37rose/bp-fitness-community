import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlus,
  faGauge,
  faPencil,
  faDumbbell,
  faEye,
  faTimesCircle,
  faCheckCircle,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import {
  getQuestionnaireDetailAPI,
  putQuestionnaireUpdateAPI,
} from "../../../API/questionnaire";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";

function AdminQuestionnaireUpdate() {
  const { id } = useParams();

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [showCancelWarning, setShowCancelWarning] = useState(false);

  const [title, setTitle] = useState("");
  const [subtitle, setSubTitle] = useState("");
  const [isMultiselect, setIsMultiselect] = useState(false);
  const [status, setStatus] = useState(true);
  const [options, setOptions] = useState([]);

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    setFetching(true);
    setErrors({});

    // To Snake-case for API from camel-case in React.
    const decamelizedData = {
      id: id,
      title: title,
      subtitle: subtitle,
      is_multiselect: isMultiselect,
      status: status,
      options: options,
    };
    putQuestionnaireUpdateAPI(
      id,
      decamelizedData,
      OnQuestionnaireUpdateSuccess,
      OnQuestionnaireUpdateError,
      OnQuestionnaireUpdateDone
    );
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function OnQuestionnaireDetailSuccess(response) {
    setTitle(response.title);
    setSubTitle(response.subtitle);
    setStatus(response.status);
    setIsMultiselect(response.isMultiselect);
    setOptions(response.options);
  }

  function OnQuestionnaireDetailError(apiErr) {
    setErrors(apiErr);
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function OnQuestionnaireDetailDone() {
    setFetching(false);
  }

  // --- Update --- //

  function OnQuestionnaireUpdateSuccess(response) {
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Question updated successfully");
    setTopAlertStatus("success");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/admin/questions/" + response.id);
  }

  function OnQuestionnaireUpdateError(apiErr) {
    setErrors(apiErr);

    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function OnQuestionnaireUpdateDone() {
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
      getQuestionnaireDetailAPI(
        id,
        OnQuestionnaireDetailSuccess,
        OnQuestionnaireDetailError,
        OnQuestionnaireDetailDone
      );
    }

    return () => {
      mounted = false;
    };
  }, [id]);
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
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/questions" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faDumbbell} />
                  &nbsp;Questions
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/questions/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Edit
                </Link>
              </li>
            </ul>
          </nav>
          <nav class="box">
            <div class={`modal ${showCancelWarning ? "is-active" : ""}`}>
              <div class="modal-background"></div>
              <div class="modal-card">
                <header class="modal-card-head">
                  <p class="modal-card-title">Are you sure?</p>
                  <button
                    class="delete"
                    aria-label="close"
                    onClick={(e) => setShowCancelWarning(false)}
                  ></button>
                </header>
                <section class="modal-card-body">
                  Your record will be cancelled and your work will be lost. This
                  cannot be undone. Do you want to continue?
                </section>
                <footer class="modal-card-foot">
                  <Link
                    class="button is-medium is-success"
                    to={`/admin/questions`}
                  >
                    Yes
                  </Link>
                  <button
                    class="button is-medium"
                    onClick={(e) => setShowCancelWarning(false)}
                  >
                    No
                  </button>
                </footer>
              </div>
            </div>

            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp; Question
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching && (
              <PageLoadingContent displayMessage={"Please wait..."} />
            )}

            <div class="container">
              <p class="subtitle is-6">
                <FontAwesomeIcon className="fas" icon={faEye} />
                &nbsp;Detail
              </p>
              <hr />

              <FormTextareaField
                label="Title"
                name="tilte"
                placeholder="title"
                value={title}
                errorText={errors && errors.title}
                onChange={(e) => setTitle(e.target.value)}
                isRequired={true}
              />

              <FormTextareaField
                label="Subtitle (Optional)"
                name="Subtitle"
                placeholder="subtitle"
                value={subtitle}
                errorText={errors && errors.subtitle}
                onChange={(e) => setSubTitle(e.target.value)}
                isRequired={true}
              />

              <FormRadioField
                label="Is Multiselect?"
                name="is_multiselect"
                value={isMultiselect}
                opt1Value={true}
                opt1Label="Yes"
                opt2Value={false}
                opt2Label="No"
                onChange={() => setIsMultiselect(!isMultiselect)}
              />
              <OptionsComponent
                options={options}
                setOptions={setOptions}
                label="Option(s)"
                errorText={errors && errors.options}
                helpText="click on the button to start adding options"
              />

              <FormRadioField
                label="Status"
                name="status"
                value={status}
                opt1Value={true}
                opt1Label="Active"
                opt2Value={false}
                opt2Label="Archived"
                onChange={() => setStatus(!status)}
              />

              <div class="columns pt-5">
                <div class="column is-half">
                  <button
                    class="button is-medium is-fullwidth-mobile"
                    onClick={(e) => setShowCancelWarning(true)}
                  >
                    <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                    &nbsp;Cancel
                  </button>
                </div>
                <div class="column is-half has-text-right">
                  <button
                    class="button is-medium is-primary is-fullwidth-mobile"
                    onClick={onSubmitClick}
                  >
                    <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                    &nbsp;Submit
                  </button>
                </div>
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

const OptionsComponent = ({
  options,
  setOptions,
  label,
  helpText,
  errorText,
}) => {
  const handleAddOption = () => {
    setOptions([...options, ""]);
  };

  const handleOptionChange = (index, value) => {
    const updatedOptions = [...options];
    updatedOptions[index] = value;
    setOptions(updatedOptions);
  };

  return (
    <div className="options-container">
      <div className="is-flex is-align-items-center mb-4">
        <label className="label">{label}</label>
        <button
          className="ml-3 px-4 py-3 bg-1 button-primary"
          onClick={handleAddOption}
        >
          <FontAwesomeIcon icon={faPlus} className="fas mr-1" />
          Add
        </button>
      </div>
      {helpText && <p class="help">{helpText}</p>}
      {options.map((option, index) => (
        <div key={index} className="space-between">
          <FormInputField
            type="text"
            className="input-field"
            placeholder="Enter option"
            value={option}
            onChange={(e) => handleOptionChange(index, e.target.value)}
            maxWidth="300px"
          />
        </div>
      ))}
      {errorText && <p class="help is-danger">{errorText}</p>}
    </div>
  );
};

export default AdminQuestionnaireUpdate;
