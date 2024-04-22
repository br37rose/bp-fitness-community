import { useState } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlus,
  faArrowLeft,
  faGauge,
  faEye,
  faQuestionCircle,
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

import { postQuestionnaireCreateAPI } from "../../../API/questionnaire";
import FormRadioField from "../../Reusable/FormRadioField";

function AdminQuestionnaireAdd() {
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

  const [title, setTitle] = useState("");
  const [subtitle, setSubTitle] = useState("");
  const [isMultiselect, setIsMultiselect] = useState(false);
  const [status, setStatus] = useState(true);
  const [options, setOptions] = useState([]);

  const onSubmitClick = (e) => {
    setFetching(true);
    setErrors({});
    postQuestionnaireCreateAPI(
      {
        title: title,
        subtitle: subtitle,
        is_multiselect: isMultiselect,
        status: status,
        options: options,
      },
      onAddSuccess,
      onAddError,
      onAddDone
    );
  };

  function onAddSuccess(response) {
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("question  created");
    setTopAlertStatus("success");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect the organization to the organization attachments page.
    setForceURL("/admin/questions");
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
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/questions" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faQuestionCircle} />
                  &nbsp;Questions
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
                <Link to="/admin/questions" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Questions
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <div class="columns">
              <div class="column">
                <p class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faPlus} />
                  &nbsp;Add Question
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
                  </>

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
                      <Link
                        class="button is-fullwidth-mobile"
                        to={`/admin/questions`}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Questions
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        onClick={onSubmitClick}
                        class="button is-success is-fullwidth-mobile"
                        type="button"
                        disabled={false}
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

export default AdminQuestionnaireAdd;
