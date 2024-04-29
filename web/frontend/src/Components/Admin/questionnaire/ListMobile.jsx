import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEdit, faEye, faTrashCan } from "@fortawesome/free-solid-svg-icons";
import { PAGE_SIZE_OPTIONS } from "../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AdminQuestionnaireListMobile(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onSelectQuestionForDeletion,
  } = props;
  return (
    <>
      {listData &&
        listData.results &&
        listData.results.map(function (datum, i) {
          return (
            <div class="pb-2" key={`mobile_tablet_${datum.id}`}>
              <strong>Title:</strong>&nbsp;{datum.title}
              {datum.subtitle && (
                <>
                  <br />
                  <br />
                  <strong>Subtitle:</strong>&nbsp;
                  {datum.subtitle}
                </>
              )}
              <br />
              <br />
              <strong>Status:</strong>&nbsp;
              {datum.status === 1 ? "Active" : "Archived"}
              <br />
              <br />
              {/* Mobile only */}
              <div class="is-hidden-desktop mb-5" key={`mobile_${datum.id}`}>
                <div className="buttons">
                  <Link
                    to={`/admin/questions/${datum.id}`}
                    className="button is-small is-dark"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faEye} />
                    &nbsp;View
                  </Link>
                  <Link
                    to={`/admin/question/${datum.id}`}
                    className="button is-small is-warning"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faEdit} />
                    &nbsp;Edit
                  </Link>
                  <button
                    onClick={(e, ses) => onSelectQuestionForDeletion(e, datum)}
                    className="button is-small is-danger"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                    &nbsp;Delete
                  </button>
                </div>
              </div>
            </div>
          );
        })}

      <div class="columns pt-4 is-mobile">
        <div class="column is-half">
          <span class="select">
            <select
              class={`input has-text-grey-light`}
              name="pageSize"
              onChange={(e) => setPageSize(parseInt(e.target.value))}
            >
              {PAGE_SIZE_OPTIONS.map(function (option, i) {
                return (
                  <option
                    key={`page_size_option_${i}`}
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
            <button class="button" onClick={onPreviousClicked}>
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
    </>
  );
}

export default AdminQuestionnaireListMobile;
