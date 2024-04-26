import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlus,
  faArrowLeft,
  faCalendar,
  faPenRuler,
  faSave,
  faArrowRight,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormTextareaField from "../../Reusable/FormTextareaField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
  currentUserState,
  topAlertMessageState,
  topAlertStatusState,
} from "../../../AppState";
import FormSelectField from "../../Reusable/FormSelectField";
import {
  FITNESS_CHALLENGE,
  MAX_WEEK_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import FormInputField from "../../Reusable/FormInputField";
import FormDateTimeField from "../../Reusable/FormDateTimeField";
import Modal from "../../Reusable/modal";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import { getMemberListOptionsAPI } from "../../../API/member";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import { postfitnessChallengeCreateAPI } from "../../../API/FitnessChallenge";

function AdminFitnessChallengeAdd() {
  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  const [currentUser] = useRecoilState(currentUserState);

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [forceURL, setForceURL] = useState("");
  const [Starton, setStarton] = useState("");
  const [duration, setduration] = useState("");
  const [showRuleModal, setshowRuleModal] = useState(false);
  const [selectedRules, setselectedRules] = useState([]);
  const [memoptions, setmemoptions] = useState([]);
  const [users, setusers] = useState([]);

  const IsRuleSelected = (ruleId) => {
    const index = selectedRules.indexOf(ruleId);
    if (index !== -1) {
      return true;
    }
    return false;
  };
  const ToggleRuleCheckBox = (ruleId) => {
    console.log("ruleid", ruleId);
    if (IsRuleSelected(ruleId)) {
      // If ruleId is present, remove it from the selectedRules array
      setselectedRules((rules) => rules.filter((id) => id !== ruleId));
    } else {
      // If ruleId is not present, add it to the selectedRules array
      setselectedRules((rules) => [...rules, ruleId]);
    }
  };

  const onSubmitClick = () => {
    // Logic to submit data
    setFetching(true);
    let payload = {
      name: name,
      description: description,
      duration: parseInt(duration),
      organization_id: currentUser.organizationId,
      rules: selectedRules,
      users: users,
      start_on: Starton,
    };

    postfitnessChallengeCreateAPI(payload, onAddSuccess, onAddError, onAddDone);
  };

  function onAddSuccess(response) {
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("challenge created");
    setTopAlertStatus("success");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);

    // Redirect the organization to the organization attachments page.
    setForceURL("/admin/fitness-challenge/" + response.id + "");
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

  function onListOK(resp) {
    setFetching(false);
    if (resp?.length) {
      setmemoptions(resp);
    }
  }

  function onListNotOK(resp) {
    setErrors(resp);
    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed gettiing list");
    setTopAlertStatus("danger");
    setTimeout(() => {
      setTopAlertMessage("");
    }, 2000);
  }
  function onDone() {
    setFetching(false);
  }

  useEffect(() => {
    window.scrollTo(0, 0);
    getMemberListOptionsAPI(
      currentUser.organizationId,
      onListOK,
      onListNotOK,
      onDone
    );
  }, []);

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="container">
      <section className="section">
        <div className="box">
          <p className="title is-4">
            <FontAwesomeIcon icon={faPlus} />
            &nbsp;Add Fitness challenge
          </p>
          <FormErrorBox errors={errors} />
          <p className="pb-4 mb-5 has-text-grey">
            Please fill out all the required fields before submitting this form.
          </p>

          {isFetching ? (
            <PageLoadingContent displayMessage={"Please wait..."} />
          ) : (
            <>
              <div className="container">
                <FormInputField
                  name="Name"
                  placeholder="Name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                  label={"Challenge Name : "}
                  errorText={errors && errors.name}
                />

                <FormMultiSelectField
                  label="Add users to the challenge"
                  name="users"
                  placeholder="Add users"
                  options={memoptions}
                  selectedValues={users}
                  onChange={(e) => {
                    let values = [];
                    for (let option of e) {
                      values.push(option.value);
                    }
                    setusers(values);
                  }}
                  errorText={errors && errors.users}
                  helpText=""
                  isRequired={false}
                  // maxWidth="640px"
                />

                <FormTextareaField
                  rows={2}
                  name="Description"
                  placeholder="Description"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                  label={"Challenge Description : "}
                  errorText={errors && errors.description}
                />
                <div className="columns">
                  <div className="column">
                    <div className="is-flex is-align-items-center">
                      <FontAwesomeIcon
                        icon={faCalendar}
                        className="mr-2 mt-3"
                      />
                      <FormDateTimeField
                        label={"Start on :"}
                        onChange={(date) => setStarton(date)}
                        value={Starton}
                        placeholder={"start on"}
                        maxWidth={"240px"}
                        name={"Starton"}
                        errorText={errors && errors.startOn}
                      />
                    </div>
                  </div>
                  <div className="column">
                    <FormSelectField
                      options={MAX_WEEK_WITH_EMPTY_OPTIONS}
                      label={"Duration: "}
                      placeholder={"duration"}
                      selectedValue={duration}
                      onChange={(e) => setduration(e.target.value)}
                      errorText={errors && errors.duration}
                    />
                  </div>
                  <div className="column"></div>
                </div>
                <button
                  className="button is-primary "
                  onClick={() => setshowRuleModal(true)}
                >
                  <FontAwesomeIcon icon={faPenRuler} />
                  &nbsp;Set Rules
                </button>
                <div>
                  <p class="subtitle is-6 mt-6">
                    <FontAwesomeIcon className="fas" icon={faPenRuler} />
                    &nbsp;Rules Added
                  </p>
                  <hr />
                  {selectedRules.map((r) => (
                    <p>
                      <FontAwesomeIcon icon={faArrowRight} />
                      &nbsp;
                      {FITNESS_CHALLENGE.filter((o) => o.value === r).map(
                        (f) => f.label
                      )}
                    </p>
                  ))}
                </div>
              </div>
              <Modal
                isOpen={showRuleModal}
                onClose={() => setshowRuleModal(false)}
              >
                <div className="modal-card-title">Rules</div>
                <div className="p-3 mt-5">
                  {FITNESS_CHALLENGE.map((c) => (
                    <>
                      <FormCheckboxField
                        label={c.label}
                        checked={IsRuleSelected(c.value)}
                        onChange={() => ToggleRuleCheckBox(c.value)}
                      />
                    </>
                  ))}
                </div>

                <button
                  className="button tp-modal-close-btn is-small is-success"
                  onClick={() => setshowRuleModal(false)}
                >
                  <FontAwesomeIcon icon={faSave} className="mr-1 is-bold" />
                  <b>Apply</b>
                </button>
              </Modal>

              <div className="columns pt-5">
                <div className="column is-half">
                  <Link
                    className="button is-fullwidth-mobile"
                    to={`/admin/fitness-challenge`}
                  >
                    <FontAwesomeIcon icon={faArrowLeft} />
                    &nbsp;Back to challenges
                  </Link>
                </div>
                <div className="column is-half has-text-right">
                  <button
                    onClick={onSubmitClick}
                    className="button is-success is-fullwidth-mobile"
                    type="button"
                    disabled={
                      !(
                        name &&
                        description &&
                        Starton &&
                        duration &&
                        selectedRules.length > 0
                      )
                    }
                  >
                    <FontAwesomeIcon icon={faPlus} />
                    &nbsp;Submit
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </section>
    </div>
  );
}

export default AdminFitnessChallengeAdd;
